package integration

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bf2fc6cc711aee1a0c2a/kas-fleet-manager/pkg/api"
	publicOpenapi "github.com/bf2fc6cc711aee1a0c2a/kas-fleet-manager/pkg/api/openapi"
	"github.com/bf2fc6cc711aee1a0c2a/kas-fleet-manager/pkg/api/private/openapi"
	"github.com/bf2fc6cc711aee1a0c2a/kas-fleet-manager/pkg/constants"
	"github.com/bf2fc6cc711aee1a0c2a/kas-fleet-manager/test"
	utils "github.com/bf2fc6cc711aee1a0c2a/kas-fleet-manager/test/common"
	"github.com/bf2fc6cc711aee1a0c2a/kas-fleet-manager/test/mocks"
	"github.com/dgrijalva/jwt-go"
	. "github.com/onsi/gomega"
	v1 "github.com/openshift-online/ocm-sdk-go/accountsmgmt/v1"
	"gopkg.in/resty.v1"
)

type TestServer struct {
	OcmServer     *httptest.Server
	TearDown      func()
	ClusterID     string
	Token         string
	Client        *publicOpenapi.APIClient
	PrivateClient *openapi.APIClient
	Helper        *test.Helper
	Ctx           context.Context
}

type claimsFunc func(account *v1.Account, clusterId string, h *test.Helper) jwt.MapClaims

func setup(t *testing.T, claims claimsFunc) TestServer {
	ocmServer := mocks.NewMockConfigurableServerBuilder().Build()
	h, client, tearDown := test.RegisterIntegration(t, ocmServer)

	clusterId, getClusterErr := utils.GetOSDClusterID(h, t, nil)
	if getClusterErr != nil {
		t.Fatalf("Failed to retrieve cluster details: %v", getClusterErr)
	}

	account := h.NewAllowedServiceAccount()
	ctx := h.NewAuthenticatedContext(account, claims(account, clusterId, h))
	token := h.CreateJWTStringWithClaim(account, claims(account, clusterId, h))

	config := openapi.NewConfiguration()
	config.BasePath = fmt.Sprintf("http://%s", h.AppConfig.Server.BindAddress)
	config.DefaultHeader = map[string]string{
		"Authorization": "Bearer " + token,
	}
	privateClient := openapi.NewAPIClient(config)

	return TestServer{
		OcmServer:     ocmServer,
		Client:        client,
		PrivateClient: privateClient,
		TearDown: func() {
			ocmServer.Close()
			tearDown()
		},
		ClusterID: clusterId,
		Token:     token,
		Helper:    h,
		Ctx:       ctx,
	}
}

func TestDataPlaneEndpoints_AuthzSuccess(t *testing.T) {
	clusterId := "test-cluster-id"
	testServer := setup(t, func(account *v1.Account, cid string, h *test.Helper) jwt.MapClaims {
		return jwt.MapClaims{
			"iss": h.AppConfig.Keycloak.KafkaRealm.ValidIssuerURI,
			"realm_access": map[string][]string{
				"roles": {"kas_fleetshard_operator"},
			},
			"kas-fleetshard-operator-cluster-id": clusterId,
		}
	})

	defer testServer.TearDown()

	body := map[string]openapi.DataPlaneKafkaStatus{
		testServer.ClusterID: {},
	}
	restyResp, err := resty.R().
		SetHeader("Content-Type", "application/json").
		SetAuthToken(testServer.Token).
		SetBody(body).
		Put(testServer.Helper.RestURL("/agent_clusters/" + clusterId + "/kafkas/status"))

	Expect(err).NotTo(HaveOccurred())
	Expect(restyResp.StatusCode()).To(Equal(http.StatusBadRequest)) //the clusterId is not valid

	clusterStatusUpdateRequest := openapi.DataPlaneClusterUpdateStatusRequest{}
	restyResp, err = resty.R().
		SetHeader("Content-Type", "application/json").
		SetAuthToken(testServer.Token).
		SetBody(clusterStatusUpdateRequest).
		Put(testServer.Helper.RestURL("/agent_clusters/" + clusterId + "/status"))

	Expect(err).NotTo(HaveOccurred())
	Expect(restyResp.StatusCode()).To(Equal(http.StatusBadRequest)) //the clusterId is not valid
}

//TODO: this test is added here to verify the "/agent_clusters" endpoint is backward compatible with "/agent-clusters". It should be removed once the backward compatibility is removed.
func TestDataPlaneEndpoints_AuthzSuccess_Old_Path(t *testing.T) {
	clusterId := "test-cluster-id"
	testServer := setup(t, func(account *v1.Account, cid string, h *test.Helper) jwt.MapClaims {
		return jwt.MapClaims{
			"iss": h.AppConfig.Keycloak.KafkaRealm.ValidIssuerURI,
			"realm_access": map[string][]string{
				"roles": {"kas_fleetshard_operator"},
			},
			"kas-fleetshard-operator-cluster-id": clusterId,
		}
	})

	defer testServer.TearDown()

	body := map[string]openapi.DataPlaneKafkaStatus{
		testServer.ClusterID: {},
	}
	restyResp, err := resty.R().
		SetHeader("Content-Type", "application/json").
		SetAuthToken(testServer.Token).
		SetBody(body).
		Put(testServer.Helper.RestURL("/agent-clusters/" + clusterId + "/kafkas/status"))

	Expect(err).NotTo(HaveOccurred())
	Expect(restyResp.StatusCode()).To(Equal(http.StatusBadRequest)) //the clusterId is not valid
}

func TestDataPlaneEndpoints_AuthzFailWhenNoRealmRole(t *testing.T) {
	testServer := setup(t, func(account *v1.Account, cid string, h *test.Helper) jwt.MapClaims {
		return jwt.MapClaims{
			"iss":                                h.AppConfig.Keycloak.KafkaRealm.ValidIssuerURI,
			"kas-fleetshard-operator-cluster-id": "test-cluster-id",
		}
	})

	defer testServer.TearDown()

	body := map[string]openapi.DataPlaneKafkaStatus{
		"test-cluster-id": {},
	}
	restyResp, err := resty.R().
		SetHeader("Content-Type", "application/json").
		SetAuthToken(testServer.Token).
		SetBody(body).
		Put(testServer.Helper.RestURL("/agent_clusters/" + testServer.ClusterID + "/kafkas/status"))

	Expect(err).NotTo(HaveOccurred())
	Expect(restyResp.StatusCode()).To(Equal(http.StatusNotFound))

	clusterStatusUpdateRequest := openapi.DataPlaneClusterUpdateStatusRequest{}
	restyResp, err = resty.R().
		SetHeader("Content-Type", "application/json").
		SetAuthToken(testServer.Token).
		SetBody(clusterStatusUpdateRequest).
		Put(testServer.Helper.RestURL("/agent_clusters/" + testServer.ClusterID + "/status"))

	Expect(err).NotTo(HaveOccurred())
	Expect(restyResp.StatusCode()).To(Equal(http.StatusNotFound))
}

func TestDataPlaneEndpoints_AuthzFailWhenClusterIdNotMatch(t *testing.T) {
	testServer := setup(t, func(account *v1.Account, cid string, h *test.Helper) jwt.MapClaims {
		return jwt.MapClaims{
			"iss": h.AppConfig.Keycloak.KafkaRealm.ValidIssuerURI,
			"realm_access": map[string][]string{
				"roles": {"kas_fleetshard_operator"},
			},
			"kas-fleetshard-operator-cluster-id": "different-cluster-id",
		}
	})
	defer testServer.TearDown()

	body := map[string]openapi.DataPlaneKafkaStatus{
		"test-cluster-id": {},
	}
	restyResp, err := resty.R().
		SetHeader("Content-Type", "application/json").
		SetAuthToken(testServer.Token).
		SetBody(body).
		Put(testServer.Helper.RestURL("/agent_clusters/" + testServer.ClusterID + "/kafkas/status"))

	Expect(err).NotTo(HaveOccurred())
	Expect(restyResp.StatusCode()).To(Equal(http.StatusNotFound))

	clusterStatusUpdateRequest := openapi.DataPlaneClusterUpdateStatusRequest{}
	restyResp, err = resty.R().
		SetHeader("Content-Type", "application/json").
		SetAuthToken(testServer.Token).
		SetBody(clusterStatusUpdateRequest).
		Put(testServer.Helper.RestURL("/agent_clusters/" + testServer.ClusterID + "/status"))

	Expect(err).NotTo(HaveOccurred())
	Expect(restyResp.StatusCode()).To(Equal(http.StatusNotFound))
}

func TestDataPlaneEndpoints_GetAndUpdateManagedKafkas(t *testing.T) {
	testServer := setup(t, func(account *v1.Account, cid string, h *test.Helper) jwt.MapClaims {
		username, _ := account.GetUsername()
		return jwt.MapClaims{
			"username": username,
			"iss":      h.AppConfig.Keycloak.KafkaRealm.ValidIssuerURI,
			"realm_access": map[string][]string{
				"roles": {"kas_fleetshard_operator"},
			},
			"kas-fleetshard-operator-cluster-id": cid,
		}
	})
	defer testServer.TearDown()
	bootstrapServerHost := "some-bootstrap⁻host"
	ssoClientID := "some-sso-client-id"
	ssoSecret := "some-sso-secret"

	var testKafkas = []*api.KafkaRequest{
		{
			ClusterID:           testServer.ClusterID,
			MultiAZ:             false,
			Name:                mockKafkaName1,
			Status:              constants.KafkaRequestStatusDeprovision.String(),
			BootstrapServerHost: bootstrapServerHost,
			SsoClientID:         ssoClientID,
			SsoClientSecret:     ssoSecret,
			Version:             "2.7.0",
		},
		{
			ClusterID:           testServer.ClusterID,
			MultiAZ:             false,
			Name:                mockKafkaName2,
			Status:              constants.KafkaRequestStatusProvisioning.String(),
			BootstrapServerHost: bootstrapServerHost,
			SsoClientID:         ssoClientID,
			SsoClientSecret:     ssoSecret,
			Version:             "2.6.0",
		},
		{
			ClusterID:           testServer.ClusterID,
			MultiAZ:             false,
			Name:                mockKafkaName3,
			Status:              constants.KafkaRequestStatusPreparing.String(),
			BootstrapServerHost: bootstrapServerHost,
			SsoClientID:         ssoClientID,
			SsoClientSecret:     ssoSecret,
			Version:             "2.7.1",
		},
		{
			ClusterID:           testServer.ClusterID,
			MultiAZ:             false,
			Name:                mockKafkaName4,
			Status:              constants.KafkaRequestStatusReady.String(),
			BootstrapServerHost: bootstrapServerHost,
			SsoClientID:         ssoClientID,
			SsoClientSecret:     ssoSecret,
			Version:             "2.7.2",
		},
		{
			ClusterID:           testServer.ClusterID,
			MultiAZ:             false,
			Name:                mockKafkaName4,
			Status:              constants.KafkaRequestStatusFailed.String(),
			BootstrapServerHost: bootstrapServerHost,
			SsoClientID:         ssoClientID,
			SsoClientSecret:     ssoSecret,
			Version:             "2.7.2",
		},
	}

	db := testServer.Helper.Env().DBFactory.New()

	// create dummy kafkas
	if err := db.Create(&testKafkas).Error; err != nil {
		Expect(err).NotTo(HaveOccurred())
		return
	}

	// create an additional kafka in failed state without "ssoSecret", "ssoClientID" and bootstrapServerHost. This indicates that the
	// kafka failed in preparing state and should not be returned in the list
	additionalKafka := &api.KafkaRequest{
		ClusterID: testServer.ClusterID,
		MultiAZ:   false,
		Name:      mockKafkaName4,
		Status:    constants.KafkaRequestStatusFailed.String(),
		Version:   "2.7.2",
	}

	if err := db.Save(additionalKafka).Error; err != nil {
		Expect(err).NotTo(HaveOccurred())
		return
	}

	list, resp, err := testServer.PrivateClient.AgentClustersApi.GetKafkas(testServer.Ctx, testServer.ClusterID)
	Expect(err).NotTo(HaveOccurred())
	Expect(resp.StatusCode).To(Equal(http.StatusOK))
	Expect(len(list.Items)).To(Equal(4)) // only count valid Managed Kafka CR

	find := func(slice []openapi.ManagedKafka, match func(kafka openapi.ManagedKafka) bool) *openapi.ManagedKafka {
		for _, item := range slice {
			if match(item) {
				return &item
			}
		}
		return nil
	}

	for _, k := range testKafkas {
		if k.Status != constants.KafkaRequestStatusPreparing.String() {
			if mk := find(list.Items, func(item openapi.ManagedKafka) bool { return item.Metadata.Annotations.Id == k.ID }); mk != nil {
				Expect(mk.Metadata.Name).To(Equal(k.Name))
				Expect(mk.Metadata.Annotations.DeprecatedBf2OrgPlacementId).To(Equal(k.PlacementId))
				Expect(mk.Metadata.Annotations.DeprecatedBf2OrgId).To(Equal(k.ID))
				Expect(mk.Metadata.Annotations.PlacementId).To(Equal(k.PlacementId))
				Expect(mk.Metadata.Annotations.Id).To(Equal(k.ID))
				Expect(mk.Metadata.Namespace).NotTo(BeEmpty())
				Expect(mk.Spec.Deleted).To(Equal(k.Status == constants.KafkaRequestStatusDeprovision.String()))
				Expect(mk.Spec.Versions.Kafka).To(Equal(k.Version))
				Expect(mk.Spec.Endpoint.Tls).To(BeNil())
			} else {
				t.Error("failed matching managedkafka id with kafkarequest id")
				break
			}
		}
	}

	var readyClusters, deletedClusters []string
	updates := map[string]openapi.DataPlaneKafkaStatus{}
	for _, item := range list.Items {
		if !item.Spec.Deleted {
			updates[item.Metadata.Annotations.Id] = openapi.DataPlaneKafkaStatus{
				Conditions: []openapi.DataPlaneClusterUpdateStatusRequestConditions{{
					Type:   "Ready",
					Status: "True",
				}},
			}
			readyClusters = append(readyClusters, item.Metadata.Annotations.Id)
		} else {
			updates[item.Metadata.Annotations.Id] = openapi.DataPlaneKafkaStatus{
				Conditions: []openapi.DataPlaneClusterUpdateStatusRequestConditions{{
					Type:   "Ready",
					Status: "False",
					Reason: "Deleted",
				}},
			}
			deletedClusters = append(deletedClusters, item.Metadata.Annotations.Id)
		}
	}

	_, err = testServer.PrivateClient.AgentClustersApi.UpdateKafkaClusterStatus(testServer.Ctx, testServer.ClusterID, updates)
	Expect(err).NotTo(HaveOccurred())

	for _, cid := range readyClusters {
		c := &api.KafkaRequest{}
		if err := db.First(c, "id = ?", cid).Error; err != nil {
			t.Errorf("failed to find kafka cluster with id %s due to error: %v", cid, err)
		}
		Expect(c.Status).To(Equal(constants.KafkaRequestStatusReady.String()))
	}

	for _, cid := range deletedClusters {
		c := &api.KafkaRequest{}
		// need to use Unscoped here as there is a chance the entry is soft deleted already
		if err := db.Unscoped().Where("id = ?", cid).First(c).Error; err != nil {
			t.Errorf("failed to find kafka cluster with id %s due to error: %v", cid, err)
		}
		Expect(c.Status).To(Equal(constants.KafkaRequestStatusDeleting.String()))
	}
}
