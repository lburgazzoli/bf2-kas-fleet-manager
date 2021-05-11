package kafka_mgrs

import (
	"sync"

	"github.com/bf2fc6cc711aee1a0c2a/kas-fleet-manager/pkg/services/syncsetresources"
	"github.com/bf2fc6cc711aee1a0c2a/kas-fleet-manager/pkg/workers"
	"github.com/pkg/errors"

	"github.com/bf2fc6cc711aee1a0c2a/kas-fleet-manager/pkg/metrics"

	"github.com/bf2fc6cc711aee1a0c2a/kas-fleet-manager/pkg/api"

	"github.com/bf2fc6cc711aee1a0c2a/kas-fleet-manager/pkg/constants"
	"github.com/bf2fc6cc711aee1a0c2a/kas-fleet-manager/pkg/services"
	"github.com/golang/glog"
)

// ReadyKafkaManager represents a kafka manager that periodically reconciles kafka requests
type ReadyKafkaManager struct {
	id              string
	workerType      string
	isRunning       bool
	kafkaService    services.KafkaService
	keycloakService services.KeycloakService
	configService   services.ConfigService
	imStop          chan struct{}
	syncTeardown    sync.WaitGroup
	reconciler      workers.Reconciler
}

// NewReadyKafkaManager creates a new kafka manager
func NewReadyKafkaManager(kafkaService services.KafkaService, id string, keycloakService services.KeycloakService, configService services.ConfigService) *ReadyKafkaManager {
	return &ReadyKafkaManager{
		id:              id,
		workerType:      "ready_kafka",
		kafkaService:    kafkaService,
		keycloakService: keycloakService,
		configService:   configService,
	}
}

func (k *ReadyKafkaManager) GetStopChan() *chan struct{} {
	return &k.imStop
}

func (k *ReadyKafkaManager) GetSyncGroup() *sync.WaitGroup {
	return &k.syncTeardown
}

func (k *ReadyKafkaManager) GetID() string {
	return k.id
}

func (c *ReadyKafkaManager) GetWorkerType() string {
	return c.workerType
}

// Start initializes the kafka manager to reconcile kafka requests
func (k *ReadyKafkaManager) Start() {
	metrics.SetLeaderWorkerMetric(k.workerType, true)
	k.reconciler.Start(k)
}

// Stop causes the process for reconciling kafka requests to stop.
func (k *ReadyKafkaManager) Stop() {
	k.reconciler.Stop(k)
	metrics.ResetMetricsForKafkaManagers()
	metrics.SetLeaderWorkerMetric(k.workerType, false)
}

func (c *ReadyKafkaManager) IsRunning() bool {
	return c.isRunning
}

func (c *ReadyKafkaManager) SetIsRunning(val bool) {
	c.isRunning = val
}

func (k *ReadyKafkaManager) Reconcile() []error {
	glog.Infoln("reconciling ready kafkas")
	if !k.configService.GetConfig().Keycloak.EnableAuthenticationOnKafka {
		return nil
	}

	var errors []error

	readyKafkas, serviceErr := k.kafkaService.ListByStatus(constants.KafkaRequestStatusReady)
	if serviceErr != nil {
		glog.Errorf("failed to list ready kafkas: %s", serviceErr.Error())
		errors = append(errors, serviceErr)
	} else {
		glog.Infof("ready kafkas count = %d", len(readyKafkas))
	}

	for _, kafka := range readyKafkas {
		glog.V(10).Infof("ready kafka id = %s", kafka.ID)
		if err := k.reconcileSsoClientIDAndSecret(kafka); err != nil {
			glog.Errorf("failed to get provisioning kafkas sso client%s: %s", kafka.SsoClientID, err.Error())
			errors = append(errors, err)
		}
	}

	return errors
}

func (k *ReadyKafkaManager) reconcileSsoClientIDAndSecret(kafkaRequest *api.KafkaRequest) error {
	if kafkaRequest.SsoClientID == "" && kafkaRequest.SsoClientSecret == "" {
		kafkaRequest.SsoClientID = syncsetresources.BuildKeycloakClientNameIdentifier(kafkaRequest.ID)
		secret, err := k.keycloakService.GetKafkaClientSecret(kafkaRequest.SsoClientID)
		if err != nil {
			return errors.Wrapf(err, "failed to get sso client id & secret for kafka cluster: %s", kafkaRequest.SsoClientID)
		}
		kafkaRequest.SsoClientSecret = secret
		if err = k.kafkaService.Update(kafkaRequest); err != nil {
			return errors.Wrapf(err, "failed to update kafka %s with cluster details", kafkaRequest.ID)
		}
	}
	return nil
}