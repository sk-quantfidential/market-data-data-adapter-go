//go:build integration

package tests

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/quantfidential/trading-ecosystem/market-data-adapter-go/pkg/interfaces"
)

// ServiceDiscoveryBehaviorTestSuite tests the behavior of service discovery operations
type ServiceDiscoveryBehaviorTestSuite struct {
	BehaviorTestSuite
}

// TestServiceDiscoveryBehaviorSuite runs the service discovery behavior test suite
func TestServiceDiscoveryBehaviorSuite(t *testing.T) {
	suite.Run(t, new(ServiceDiscoveryBehaviorTestSuite))
}

// TestServiceRegistration tests service registration and discovery
func (suite *ServiceDiscoveryBehaviorTestSuite) TestServiceRegistration() {
	var serviceID = GenerateTestID("service")

	suite.Given("a service to register", func() {
		// Service defined below
	}).When("registering the service", func() {
		service := suite.CreateTestServiceInfo(serviceID, func(s *interfaces.ServiceInfo) {
			s.ServiceName = "test-market-data-service"
			s.Version = "1.0.0"
		})

		err := suite.adapter.ServiceDiscoveryRepository().Register(suite.ctx, service)
		suite.Require().NoError(err)
		suite.trackCreatedService(serviceID)
	}).Then("the service should be discoverable", func() {
		services, err := suite.adapter.ServiceDiscoveryRepository().Discover(suite.ctx, "test-market-data-service")
		suite.Require().NoError(err)
		suite.GreaterOrEqual(len(services), 1)

		var found bool
		for _, svc := range services {
			if svc.ServiceID == serviceID {
				found = true
				suite.Equal("1.0.0", svc.Version)
				break
			}
		}
		suite.True(found, "Should find registered service")
	})
}

// TestServiceHeartbeat tests service heartbeat updates
func (suite *ServiceDiscoveryBehaviorTestSuite) TestServiceHeartbeat() {
	var serviceID = GenerateTestID("heartbeat-service")

	suite.Given("a registered service", func() {
		service := suite.CreateTestServiceInfo(serviceID, func(s *interfaces.ServiceInfo) {
			s.ServiceName = "heartbeat-test-service"
		})
		err := suite.adapter.ServiceDiscoveryRepository().Register(suite.ctx, service)
		suite.Require().NoError(err)
		suite.trackCreatedService(serviceID)
	}).When("sending a heartbeat", func() {
		err := suite.adapter.ServiceDiscoveryRepository().Heartbeat(suite.ctx, serviceID)
		suite.Require().NoError(err)
	}).Then("the service last heartbeat should be updated", func() {
		service, err := suite.adapter.ServiceDiscoveryRepository().GetServiceInfo(suite.ctx, serviceID)
		suite.Require().NoError(err)
		suite.NotNil(service)
		suite.WithinDuration(time.Now(), service.LastHeartbeat, 5*time.Second)
	})
}

// TestServiceDeregistration tests service deregistration
func (suite *ServiceDiscoveryBehaviorTestSuite) TestServiceDeregistration() {
	var serviceID = GenerateTestID("deregister-service")

	suite.Given("a registered service", func() {
		service := suite.CreateTestServiceInfo(serviceID, func(s *interfaces.ServiceInfo) {
			s.ServiceName = "deregister-test-service"
		})
		err := suite.adapter.ServiceDiscoveryRepository().Register(suite.ctx, service)
		suite.Require().NoError(err)
		suite.trackCreatedService(serviceID)
	}).When("deregistering the service", func() {
		err := suite.adapter.ServiceDiscoveryRepository().Deregister(suite.ctx, serviceID)
		suite.Require().NoError(err)
	}).Then("the service should not be discoverable", func() {
		_, err := suite.adapter.ServiceDiscoveryRepository().GetServiceInfo(suite.ctx, serviceID)
		suite.Error(err, "Should not find deregistered service")
	})
}
