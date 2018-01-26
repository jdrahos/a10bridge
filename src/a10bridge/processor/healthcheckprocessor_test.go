package processor_test

import (
	"a10bridge/mocks"
	"a10bridge/model"
	"a10bridge/processor"
	"testing"

	"github.com/stretchr/testify/suite"
)

type HealthCheckProcessorTestSuite struct {
	suite.Suite
	helper *processor.TestHelper
	client *mocks.Client
}

func (suite *HealthCheckProcessorTestSuite) SetupTest() {
	suite.client = new(mocks.Client)
}

func TestHealthCheckProcessor(t *testing.T) {
	tests := new(HealthCheckProcessorTestSuite)
	tests.helper = new(processor.TestHelper)
	suite.Run(t, tests)
}

func (suite *HealthCheckProcessorTestSuite) TestProcessHealthCheck_notChanged() {
	client := suite.client
	processor := suite.helper.BuildHealthcheckProcessor(client)
	healthCheck := healthCheck()

	client.On("GetHealthMonitor", healthCheck.Name).Once().Return(healthCheck, nil)
	err := processor.ProcessHealthCheck(healthCheck)
	suite.Assert().Nil(err)
	client.AssertExpectations(suite.T())
}

func (suite *HealthCheckProcessorTestSuite) TestProcessHealthCheck_notFound() {
	a10error := new(mocks.A10Error)
	client := suite.client
	processor := suite.helper.BuildHealthcheckProcessor(client)
	healthCheck := healthCheck()

	client.On("GetHealthMonitor", healthCheck.Name).Once().Return(nil, a10error)
	client.On("IsHealthMonitorNotFound", a10error).Once().Return(true)
	client.On("CreateHealthMonitor", healthCheck).Once().Return(nil)
	err := processor.ProcessHealthCheck(healthCheck)
	suite.Assert().Nil(err)
	client.AssertExpectations(suite.T())
}

func (suite *HealthCheckProcessorTestSuite) TestProcessHealthCheck_getFails() {
	a10error := new(mocks.A10Error)
	client := suite.client
	processor := suite.helper.BuildHealthcheckProcessor(client)
	healthCheck := healthCheck()

	client.On("GetHealthMonitor", healthCheck.Name).Once().Return(nil, a10error)
	client.On("IsHealthMonitorNotFound", a10error).Once().Return(false)
	err := processor.ProcessHealthCheck(healthCheck)
	suite.Assert().NotNil(err)
	client.AssertExpectations(suite.T())
}

func (suite *HealthCheckProcessorTestSuite) TestProcessHealthCheck_createFails() {
	a10error := new(mocks.A10Error)
	client := suite.client
	processor := suite.helper.BuildHealthcheckProcessor(client)
	healthCheck := healthCheck()

	client.On("GetHealthMonitor", healthCheck.Name).Once().Return(nil, a10error)
	client.On("IsHealthMonitorNotFound", a10error).Once().Return(true)
	client.On("CreateHealthMonitor", healthCheck).Once().Return(a10error)
	err := processor.ProcessHealthCheck(healthCheck)
	suite.Assert().NotNil(err)
	client.AssertExpectations(suite.T())
}

func (suite *HealthCheckProcessorTestSuite) TestProcessHealthCheck_endpointChanged() {
	client := suite.client
	processor := suite.helper.BuildHealthcheckProcessor(client)

	healthCheck := healthCheck()
	existing := *healthCheck
	existing.Endpoint = "/ws"

	client.On("GetHealthMonitor", healthCheck.Name).Once().Return(&existing, nil)
	client.On("UpdateHealthMonitor", healthCheck).Once().Return(nil)
	err := processor.ProcessHealthCheck(healthCheck)
	suite.Assert().Nil(err)
	client.AssertExpectations(suite.T())
}

func (suite *HealthCheckProcessorTestSuite) TestProcessHealthCheck_expectedCodeChanged() {
	client := suite.client
	processor := suite.helper.BuildHealthcheckProcessor(client)

	healthCheck := healthCheck()
	existing := *healthCheck
	existing.ExpectCode = "505"

	client.On("GetHealthMonitor", healthCheck.Name).Once().Return(&existing, nil)
	client.On("UpdateHealthMonitor", healthCheck).Once().Return(nil)
	err := processor.ProcessHealthCheck(healthCheck)
	suite.Assert().Nil(err)
	client.AssertExpectations(suite.T())
}

func (suite *HealthCheckProcessorTestSuite) TestProcessHealthCheck_intervalChanged() {
	client := suite.client
	processor := suite.helper.BuildHealthcheckProcessor(client)

	healthCheck := healthCheck()
	existing := *healthCheck
	existing.Interval = 505

	client.On("GetHealthMonitor", healthCheck.Name).Once().Return(&existing, nil)
	client.On("UpdateHealthMonitor", healthCheck).Once().Return(nil)
	err := processor.ProcessHealthCheck(healthCheck)
	suite.Assert().Nil(err)
	client.AssertExpectations(suite.T())
}

func (suite *HealthCheckProcessorTestSuite) TestProcessHealthCheck_portChanged() {
	client := suite.client
	processor := suite.helper.BuildHealthcheckProcessor(client)

	healthCheck := healthCheck()
	existing := *healthCheck
	existing.Port = 808080

	client.On("GetHealthMonitor", healthCheck.Name).Once().Return(&existing, nil)
	client.On("UpdateHealthMonitor", healthCheck).Once().Return(nil)
	err := processor.ProcessHealthCheck(healthCheck)
	suite.Assert().Nil(err)
	client.AssertExpectations(suite.T())
}

func (suite *HealthCheckProcessorTestSuite) TestProcessHealthCheck_requiredConsecutivePassesChanged() {
	client := suite.client
	processor := suite.helper.BuildHealthcheckProcessor(client)

	healthCheck := healthCheck()
	existing := *healthCheck
	existing.RequiredConsecutivePasses = 10

	client.On("GetHealthMonitor", healthCheck.Name).Once().Return(&existing, nil)
	client.On("UpdateHealthMonitor", healthCheck).Once().Return(nil)
	err := processor.ProcessHealthCheck(healthCheck)
	suite.Assert().Nil(err)
	client.AssertExpectations(suite.T())
}

func (suite *HealthCheckProcessorTestSuite) TestProcessHealthCheck_retryCountChanged() {
	client := suite.client
	processor := suite.helper.BuildHealthcheckProcessor(client)

	healthCheck := healthCheck()
	existing := *healthCheck
	existing.RetryCount = 10

	client.On("GetHealthMonitor", healthCheck.Name).Once().Return(&existing, nil)
	client.On("UpdateHealthMonitor", healthCheck).Once().Return(nil)
	err := processor.ProcessHealthCheck(healthCheck)
	suite.Assert().Nil(err)
	client.AssertExpectations(suite.T())
}

func (suite *HealthCheckProcessorTestSuite) TestProcessHealthCheck_timeoutChanged() {
	client := suite.client
	processor := suite.helper.BuildHealthcheckProcessor(client)

	healthCheck := healthCheck()
	existing := *healthCheck
	existing.Timeout = 50

	client.On("GetHealthMonitor", healthCheck.Name).Once().Return(&existing, nil)
	client.On("UpdateHealthMonitor", healthCheck).Once().Return(nil)
	err := processor.ProcessHealthCheck(healthCheck)
	suite.Assert().Nil(err)
	client.AssertExpectations(suite.T())
}

func (suite *HealthCheckProcessorTestSuite) TestProcessHealthCheck_updateFails() {
	a10error := new(mocks.A10Error)
	client := suite.client
	processor := suite.helper.BuildHealthcheckProcessor(client)

	healthCheck := healthCheck()
	existing := *healthCheck
	existing.Endpoint = "/ws"

	client.On("GetHealthMonitor", healthCheck.Name).Once().Return(&existing, nil)
	client.On("UpdateHealthMonitor", healthCheck).Once().Return(a10error)
	err := processor.ProcessHealthCheck(healthCheck)
	suite.Assert().NotNil(err)
	client.AssertExpectations(suite.T())
}

func healthCheck() *model.HealthCheck {
	return &model.HealthCheck{
		RequiredConsecutivePasses: 2,
		Endpoint:                  "test",
		ExpectCode:                "200",
		Interval:                  10,
		Name:                      "test",
		Port:                      80,
		RetryCount:                4,
		Timeout:                   10,
	}
}
