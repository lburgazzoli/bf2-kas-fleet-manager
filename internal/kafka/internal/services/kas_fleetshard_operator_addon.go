package services

import (
	"github.com/bf2fc6cc711aee1a0c2a/kas-fleet-manager/pkg/api"
	"github.com/bf2fc6cc711aee1a0c2a/kas-fleet-manager/pkg/clusters"
	"github.com/bf2fc6cc711aee1a0c2a/kas-fleet-manager/pkg/clusters/ocm"
	"github.com/bf2fc6cc711aee1a0c2a/kas-fleet-manager/pkg/clusters/types"
	"github.com/bf2fc6cc711aee1a0c2a/kas-fleet-manager/pkg/config"
	"github.com/bf2fc6cc711aee1a0c2a/kas-fleet-manager/pkg/errors"
	"github.com/bf2fc6cc711aee1a0c2a/kas-fleet-manager/pkg/services"
	"github.com/goava/di"
	"github.com/golang/glog"
)

const (
	KasFleetshardOperatorRoleName = "kas_fleetshard_operator"

	//parameter names for the kas-fleetshard-operator service account
	kasFleetshardOperatorParamMasSSOBaseUrl        = "sso-auth-server-url"
	kasFleetshardOperatorParamServiceAccountId     = "sso-client-id"
	kasFleetshardOperatorParamServiceAccountSecret = "sso-secret"
	// parameter names for the cluster id
	kasFleetshardOperatorParamClusterId = "cluster-id"
	// parameter names for the control plane url
	kasFleetshardOperatorParamControlPlaneBaseURL = "control-plane-url"
	//parameter names for fleetshardoperator synchronizer
	kasFleetshardOperatorParamPollinterval   = "poll-interval"
	kasFleetshardOperatorParamResyncInterval = "resync-interval"
)

//go:generate moq -out kas_fleetshard_operator_addon_moq.go . KasFleetshardOperatorAddon
type KasFleetshardOperatorAddon interface {
	Provision(cluster api.Cluster) (bool, *errors.ServiceError)
	ReconcileParameters(cluster api.Cluster) *errors.ServiceError
	RemoveServiceAccount(cluster api.Cluster) *errors.ServiceError
}

func NewKasFleetshardOperatorAddon(o kasFleetshardOperatorAddon) KasFleetshardOperatorAddon {
	return &o
}

type kasFleetshardOperatorAddon struct {
	di.Inject
	SsoService          services.KafkaKeycloakService
	ProviderFactory     clusters.ProviderFactory
	ServerConfig        *config.ServerConfig
	KasFleetShardConfig *config.KasFleetshardConfig
	OCMConfig           *config.OCMConfig
	KeycloakConfig      *config.KeycloakConfig
}

func (o *kasFleetshardOperatorAddon) Provision(cluster api.Cluster) (bool, *errors.ServiceError) {
	kasFleetshardAddonID := o.OCMConfig.KasFleetshardAddonID
	params, paramsErr := o.getAddonParams(cluster)
	if paramsErr != nil {
		return false, paramsErr
	}
	p, err := o.ProviderFactory.GetAddonProvider(cluster.ProviderType)
	if err != nil {
		return false, errors.NewWithCause(errors.ErrorGeneral, err, "failed to get provider implementation")
	}
	glog.V(5).Infof("Provision addon %s for cluster %s", kasFleetshardAddonID, cluster.ClusterID)
	spec := &types.ClusterSpec{
		InternalID:     cluster.ClusterID,
		ExternalID:     cluster.ExternalID,
		Status:         cluster.Status,
		AdditionalInfo: cluster.ClusterSpec,
	}
	if ready, err := p.InstallAddonWithParams(spec, kasFleetshardAddonID, params); err != nil {
		return false, errors.NewWithCause(errors.ErrorGeneral, err, "failed to install addon %s for cluster %s", kasFleetshardAddonID, cluster.ClusterID)
	} else {
		return ready, nil
	}
}

func (o *kasFleetshardOperatorAddon) ReconcileParameters(cluster api.Cluster) *errors.ServiceError {
	kasFleetshardAddonID := o.OCMConfig.KasFleetshardAddonID
	params, paramsErr := o.getAddonParams(cluster)
	if paramsErr != nil {
		return paramsErr
	}
	p, err := o.ProviderFactory.GetAddonProvider(cluster.ProviderType)
	if err != nil {
		return errors.NewWithCause(errors.ErrorGeneral, err, "failed to get provider implementation")
	}

	glog.V(5).Infof("Reconcile parameters for addon %s on cluster %s", kasFleetshardAddonID, cluster.ClusterID)
	spec := &types.ClusterSpec{
		InternalID:     cluster.ClusterID,
		ExternalID:     cluster.ExternalID,
		Status:         cluster.Status,
		AdditionalInfo: cluster.ClusterSpec,
	}
	if updated, err := p.InstallAddonWithParams(spec, kasFleetshardAddonID, params); err != nil {
		return errors.NewWithCause(errors.ErrorGeneral, err, "failed to update parameters for addon %s for cluster %s", kasFleetshardAddonID, cluster.ClusterID)
	} else if updated {
		glog.V(5).Infof("Addon parameters for addon %s on cluster %s are updated", kasFleetshardAddonID, cluster.ClusterID)
		return nil
	} else {
		glog.V(5).Infof("Addon parameters for addon %s on cluster %s are not updated", kasFleetshardAddonID, cluster.ClusterID)
		return nil
	}
}

func (o *kasFleetshardOperatorAddon) getAddonParams(cluster api.Cluster) ([]ocm.AddonParameter, *errors.ServiceError) {
	acc, pErr := o.provisionServiceAccount(cluster.ClusterID)
	if pErr != nil {
		return nil, errors.GeneralError("failed to create service account for cluster %s due to error: %v", cluster.ClusterID, pErr)
	}
	params := o.buildAddonParams(acc, cluster.ClusterID)
	return params, nil
}

func (o *kasFleetshardOperatorAddon) provisionServiceAccount(clusterId string) (*api.ServiceAccount, *errors.ServiceError) {
	glog.V(5).Infof("Provisioning service account for cluster %s", clusterId)
	return o.SsoService.RegisterKasFleetshardOperatorServiceAccount(clusterId, KasFleetshardOperatorRoleName)
}

func (o *kasFleetshardOperatorAddon) buildAddonParams(serviceAccount *api.ServiceAccount, clusterId string) []ocm.AddonParameter {
	p := []ocm.AddonParameter{

		{
			Id:    kasFleetshardOperatorParamMasSSOBaseUrl,
			Value: o.KeycloakConfig.KafkaRealm.ValidIssuerURI,
		},
		{
			Id:    kasFleetshardOperatorParamServiceAccountId,
			Value: serviceAccount.ClientID,
		},
		{
			Id:    kasFleetshardOperatorParamServiceAccountSecret,
			Value: serviceAccount.ClientSecret,
		},
		{
			Id:    kasFleetshardOperatorParamControlPlaneBaseURL,
			Value: o.ServerConfig.PublicHostURL,
		},
		{
			Id:    kasFleetshardOperatorParamClusterId,
			Value: clusterId,
		},
		{
			Id:    kasFleetshardOperatorParamPollinterval,
			Value: o.KasFleetShardConfig.PollInterval,
		},
		{
			Id:    kasFleetshardOperatorParamResyncInterval,
			Value: o.KasFleetShardConfig.ResyncInterval,
		},
	}
	return p
}

func (o *kasFleetshardOperatorAddon) RemoveServiceAccount(cluster api.Cluster) *errors.ServiceError {
	glog.V(5).Infof("Removing kas-fleetshard-operator service account for cluster %s", cluster.ClusterID)
	return o.SsoService.DeRegisterKasFleetshardOperatorServiceAccount(cluster.ClusterID)
}