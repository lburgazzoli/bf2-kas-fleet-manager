// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package services

import (
	"github.com/bf2fc6cc711aee1a0c2a/kas-fleet-manager/pkg/api"
	apiErrors "github.com/bf2fc6cc711aee1a0c2a/kas-fleet-manager/pkg/errors"
	"sync"
)

// Ensure, that QuotaServiceMock does implement QuotaService.
// If this is not the case, regenerate this file with moq.
var _ QuotaService = &QuotaServiceMock{}

// QuotaServiceMock is a mock implementation of QuotaService.
//
// 	func TestSomethingThatUsesQuotaService(t *testing.T) {
//
// 		// make and configure a mocked QuotaService
// 		mockedQuotaService := &QuotaServiceMock{
// 			CheckQuotaFunc: func(kafka *api.KafkaRequest) *apiErrors.ServiceError {
// 				panic("mock out the CheckQuota method")
// 			},
// 			DeleteQuotaFunc: func(subscriptionId string) *apiErrors.ServiceError {
// 				panic("mock out the DeleteQuota method")
// 			},
// 			ReserveQuotaFunc: func(kafka *api.KafkaRequest) (string, *apiErrors.ServiceError) {
// 				panic("mock out the ReserveQuota method")
// 			},
// 		}
//
// 		// use mockedQuotaService in code that requires QuotaService
// 		// and then make assertions.
//
// 	}
type QuotaServiceMock struct {
	// CheckQuotaFunc mocks the CheckQuota method.
	CheckQuotaFunc func(kafka *api.KafkaRequest) *apiErrors.ServiceError

	// DeleteQuotaFunc mocks the DeleteQuota method.
	DeleteQuotaFunc func(subscriptionId string) *apiErrors.ServiceError

	// ReserveQuotaFunc mocks the ReserveQuota method.
	ReserveQuotaFunc func(kafka *api.KafkaRequest) (string, *apiErrors.ServiceError)

	// calls tracks calls to the methods.
	calls struct {
		// CheckQuota holds details about calls to the CheckQuota method.
		CheckQuota []struct {
			// Kafka is the kafka argument value.
			Kafka *api.KafkaRequest
		}
		// DeleteQuota holds details about calls to the DeleteQuota method.
		DeleteQuota []struct {
			// SubscriptionId is the subscriptionId argument value.
			SubscriptionId string
		}
		// ReserveQuota holds details about calls to the ReserveQuota method.
		ReserveQuota []struct {
			// Kafka is the kafka argument value.
			Kafka *api.KafkaRequest
		}
	}
	lockCheckQuota   sync.RWMutex
	lockDeleteQuota  sync.RWMutex
	lockReserveQuota sync.RWMutex
}

// CheckQuota calls CheckQuotaFunc.
func (mock *QuotaServiceMock) CheckQuota(kafka *api.KafkaRequest) *apiErrors.ServiceError {
	if mock.CheckQuotaFunc == nil {
		panic("QuotaServiceMock.CheckQuotaFunc: method is nil but QuotaService.CheckQuota was just called")
	}
	callInfo := struct {
		Kafka *api.KafkaRequest
	}{
		Kafka: kafka,
	}
	mock.lockCheckQuota.Lock()
	mock.calls.CheckQuota = append(mock.calls.CheckQuota, callInfo)
	mock.lockCheckQuota.Unlock()
	return mock.CheckQuotaFunc(kafka)
}

// CheckQuotaCalls gets all the calls that were made to CheckQuota.
// Check the length with:
//     len(mockedQuotaService.CheckQuotaCalls())
func (mock *QuotaServiceMock) CheckQuotaCalls() []struct {
	Kafka *api.KafkaRequest
} {
	var calls []struct {
		Kafka *api.KafkaRequest
	}
	mock.lockCheckQuota.RLock()
	calls = mock.calls.CheckQuota
	mock.lockCheckQuota.RUnlock()
	return calls
}

// DeleteQuota calls DeleteQuotaFunc.
func (mock *QuotaServiceMock) DeleteQuota(subscriptionId string) *apiErrors.ServiceError {
	if mock.DeleteQuotaFunc == nil {
		panic("QuotaServiceMock.DeleteQuotaFunc: method is nil but QuotaService.DeleteQuota was just called")
	}
	callInfo := struct {
		SubscriptionId string
	}{
		SubscriptionId: subscriptionId,
	}
	mock.lockDeleteQuota.Lock()
	mock.calls.DeleteQuota = append(mock.calls.DeleteQuota, callInfo)
	mock.lockDeleteQuota.Unlock()
	return mock.DeleteQuotaFunc(subscriptionId)
}

// DeleteQuotaCalls gets all the calls that were made to DeleteQuota.
// Check the length with:
//     len(mockedQuotaService.DeleteQuotaCalls())
func (mock *QuotaServiceMock) DeleteQuotaCalls() []struct {
	SubscriptionId string
} {
	var calls []struct {
		SubscriptionId string
	}
	mock.lockDeleteQuota.RLock()
	calls = mock.calls.DeleteQuota
	mock.lockDeleteQuota.RUnlock()
	return calls
}

// ReserveQuota calls ReserveQuotaFunc.
func (mock *QuotaServiceMock) ReserveQuota(kafka *api.KafkaRequest) (string, *apiErrors.ServiceError) {
	if mock.ReserveQuotaFunc == nil {
		panic("QuotaServiceMock.ReserveQuotaFunc: method is nil but QuotaService.ReserveQuota was just called")
	}
	callInfo := struct {
		Kafka *api.KafkaRequest
	}{
		Kafka: kafka,
	}
	mock.lockReserveQuota.Lock()
	mock.calls.ReserveQuota = append(mock.calls.ReserveQuota, callInfo)
	mock.lockReserveQuota.Unlock()
	return mock.ReserveQuotaFunc(kafka)
}

// ReserveQuotaCalls gets all the calls that were made to ReserveQuota.
// Check the length with:
//     len(mockedQuotaService.ReserveQuotaCalls())
func (mock *QuotaServiceMock) ReserveQuotaCalls() []struct {
	Kafka *api.KafkaRequest
} {
	var calls []struct {
		Kafka *api.KafkaRequest
	}
	mock.lockReserveQuota.RLock()
	calls = mock.calls.ReserveQuota
	mock.lockReserveQuota.RUnlock()
	return calls
}
