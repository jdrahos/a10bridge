package processor_test

import (
	"a10bridge/mocks"
	"a10bridge/model"
	"a10bridge/processor"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ServiceGroupProcessorTestSuite struct {
	suite.Suite
	helper *processor.TestHelper
	client *mocks.Client
}

func (suite *ServiceGroupProcessorTestSuite) SetupTest() {
	suite.client = new(mocks.Client)
}

func TestServiceGroupProcessor(t *testing.T) {
	tests := new(ServiceGroupProcessorTestSuite)
	tests.helper = new(processor.TestHelper)
	suite.Run(t, tests)
}

func (suite *ServiceGroupProcessorTestSuite) TestProcessServiceGroup_noMembers() {
	client := suite.client
	processor := suite.helper.BuildServiceGroupProcessor(client)
	serviceGroup := serviceGroup()
	failedNodeNames := []string{serviceGroup.IngressControllers[0].Nodes[0].A10Server}

	err := processor.ProcessServiceGroup(serviceGroup, failedNodeNames)
	suite.Assert().NotNil(err)
	client.AssertExpectations(suite.T())
}

func (suite *ServiceGroupProcessorTestSuite) TestProcessServiceGroup_notChanged() {
	client := suite.client
	processor := suite.helper.BuildServiceGroupProcessor(client)
	serviceGroup := serviceGroup()
	failedNodeNames := []string{"server_down"}

	client.On("GetServiceGroup", serviceGroup.Name).Once().Return(serviceGroup, nil)
	err := processor.ProcessServiceGroup(serviceGroup, failedNodeNames)
	suite.Assert().Nil(err)
	client.AssertExpectations(suite.T())
}

func (suite *ServiceGroupProcessorTestSuite) TestProcessServiceGroup_GetServiceGroupFails() {
	a10error := new(mocks.A10Error)
	client := suite.client
	processor := suite.helper.BuildServiceGroupProcessor(client)
	serviceGroup := serviceGroup()
	failedNodeNames := []string{"server_down"}

	client.On("GetServiceGroup", serviceGroup.Name).Once().Return(nil, a10error)
	client.On("IsServiceGroupNotFound", a10error).Once().Return(false)
	err := processor.ProcessServiceGroup(serviceGroup, failedNodeNames)
	suite.Assert().NotNil(err)
	client.AssertExpectations(suite.T())
}

func (suite *ServiceGroupProcessorTestSuite) TestProcessServiceGroup_notFound() {
	a10error := new(mocks.A10Error)
	client := suite.client
	processor := suite.helper.BuildServiceGroupProcessor(client)
	serviceGroup := serviceGroup()
	failedNodeNames := []string{"server_down"}

	client.On("GetServiceGroup", serviceGroup.Name).Once().Return(nil, a10error)
	client.On("IsServiceGroupNotFound", a10error).Once().Return(true)
	client.On("CreateServiceGroup", serviceGroup).Once().Return(nil)
	err := processor.ProcessServiceGroup(serviceGroup, failedNodeNames)
	suite.Assert().Nil(err)
	client.AssertExpectations(suite.T())
}

func (suite *ServiceGroupProcessorTestSuite) TestProcessServiceGroup_createSserviceGroupFails() {
	a10error := new(mocks.A10Error)
	client := suite.client
	processor := suite.helper.BuildServiceGroupProcessor(client)
	serviceGroup := serviceGroup()
	failedNodeNames := []string{"server_down"}

	client.On("GetServiceGroup", serviceGroup.Name).Once().Return(nil, a10error)
	client.On("IsServiceGroupNotFound", a10error).Once().Return(true)
	client.On("CreateServiceGroup", serviceGroup).Once().Return(a10error)
	err := processor.ProcessServiceGroup(serviceGroup, failedNodeNames)
	suite.Assert().NotNil(err)
	client.AssertExpectations(suite.T())
}

func (suite *ServiceGroupProcessorTestSuite) TestProcessServiceGroup_healthCheckNameChanged() {
	client := suite.client
	processor := suite.helper.BuildServiceGroupProcessor(client)
	serviceGroup := serviceGroup()
	existing := *serviceGroup
	healthCheck := *serviceGroup.Health
	healthCheck.Name = "changed"
	existing.Health = &healthCheck
	failedNodeNames := []string{"server_down"}
	existing.Members = []*model.Member{
		&model.Member{
			ServerName: "server",
			Port:       8080,
		},
	}

	client.On("GetServiceGroup", serviceGroup.Name).Once().Return(&existing, nil)
	client.On("UpdateServiceGroup", serviceGroup).Once().Return(nil)
	err := processor.ProcessServiceGroup(serviceGroup, failedNodeNames)
	suite.Assert().Nil(err)
	client.AssertExpectations(suite.T())
}

func (suite *ServiceGroupProcessorTestSuite) TestProcessServiceGroup_updateFails() {
	a10error := new(mocks.A10Error)
	client := suite.client
	processor := suite.helper.BuildServiceGroupProcessor(client)
	serviceGroup := serviceGroup()
	existing := *serviceGroup
	healthCheck := *serviceGroup.Health
	healthCheck.Name = "changed"
	existing.Health = &healthCheck
	failedNodeNames := []string{"server_down"}

	client.On("GetServiceGroup", serviceGroup.Name).Once().Return(&existing, nil)
	client.On("UpdateServiceGroup", serviceGroup).Once().Return(a10error)
	err := processor.ProcessServiceGroup(serviceGroup, failedNodeNames)
	suite.Assert().NotNil(err)
	client.AssertExpectations(suite.T())
}

func (suite *ServiceGroupProcessorTestSuite) TestProcessServiceGroup_missingMember() {
	client := suite.client
	processor := suite.helper.BuildServiceGroupProcessor(client)
	serviceGroup := serviceGroup()
	existing := *serviceGroup
	healthCheck := *serviceGroup.Health
	healthCheck.Name = "changed"
	existing.Health = &healthCheck
	failedNodeNames := []string{"server_down"}
	existing.Members = []*model.Member{}

	client.On("GetServiceGroup", serviceGroup.Name).Once().Return(&existing, nil)
	client.On("UpdateServiceGroup", serviceGroup).Once().Return(nil)
	client.On("CreateMember", &model.Member{
		Port:             8080,
		ServerName:       "server",
		ServiceGroupName: "service group",
	}).Once().Return(nil)
	err := processor.ProcessServiceGroup(serviceGroup, failedNodeNames)
	suite.Assert().Nil(err)
	client.AssertExpectations(suite.T())
}

func (suite *ServiceGroupProcessorTestSuite) TestProcessServiceGroup_createMemberFails() {
	a10error := new(mocks.A10Error)
	client := suite.client
	processor := suite.helper.BuildServiceGroupProcessor(client)
	serviceGroup := serviceGroup()
	serviceGroup.IngressControllers[0].Nodes = append(serviceGroup.IngressControllers[0].Nodes, &model.Node{
		Name:      "server2",
		A10Server: "server2",
	})
	existing := *serviceGroup
	healthCheck := *serviceGroup.Health
	healthCheck.Name = "changed"
	existing.Health = &healthCheck
	failedNodeNames := []string{"server_down"}
	existing.Members = []*model.Member{}

	client.On("GetServiceGroup", serviceGroup.Name).Once().Return(&existing, nil)
	client.On("UpdateServiceGroup", serviceGroup).Once().Return(nil)
	client.On("CreateMember", &model.Member{
		Port:             8080,
		ServerName:       "server",
		ServiceGroupName: "service group",
	}).Once().Return(a10error)
	client.On("IsMemberAlreadyExists", a10error).Once().Return(false)
	client.On("CreateMember", &model.Member{
		Port:             8080,
		ServerName:       "server2",
		ServiceGroupName: "service group",
	}).Once().Return(nil)
	err := processor.ProcessServiceGroup(serviceGroup, failedNodeNames)
	suite.Assert().NotNil(err)
	client.AssertExpectations(suite.T())
}

func (suite *ServiceGroupProcessorTestSuite) TestProcessServiceGroup_memberCreatedDuringUpdate() {
	a10error := new(mocks.A10Error)
	client := suite.client
	processor := suite.helper.BuildServiceGroupProcessor(client)
	serviceGroup := serviceGroup()
	existing := *serviceGroup
	healthCheck := *serviceGroup.Health
	healthCheck.Name = "changed"
	existing.Health = &healthCheck
	failedNodeNames := []string{"server_down"}
	existing.Members = []*model.Member{}

	client.On("GetServiceGroup", serviceGroup.Name).Once().Return(&existing, nil)
	client.On("UpdateServiceGroup", serviceGroup).Once().Return(nil)
	client.On("CreateMember", &model.Member{
		Port:             8080,
		ServerName:       "server",
		ServiceGroupName: "service group",
	}).Once().Return(a10error)
	client.On("IsMemberAlreadyExists", a10error).Once().Return(true)
	err := processor.ProcessServiceGroup(serviceGroup, failedNodeNames)
	suite.Assert().Nil(err)
	client.AssertExpectations(suite.T())
}

func (suite *ServiceGroupProcessorTestSuite) TestProcessServiceGroup_extraMember() {
	client := suite.client
	processor := suite.helper.BuildServiceGroupProcessor(client)
	serviceGroup := serviceGroup()
	existing := *serviceGroup
	healthCheck := *serviceGroup.Health
	healthCheck.Name = "changed"
	existing.Health = &healthCheck
	failedNodeNames := []string{"server_down"}
	extraServer1 := &model.Member{
		ServerName: "server2",
		Port:       8080,
	}
	extraServer2 := &model.Member{
		ServerName: "server",
		Port:       8081,
	}
	existing.Members = []*model.Member{
		&model.Member{
			ServerName: "server",
			Port:       8080,
		},
		extraServer1,
		extraServer2,
	}

	client.On("GetServiceGroup", serviceGroup.Name).Once().Return(&existing, nil)
	client.On("UpdateServiceGroup", serviceGroup).Once().Return(nil)
	client.On("DeleteMember", extraServer1).Once().Return(nil)
	client.On("DeleteMember", extraServer2).Once().Return(nil)
	err := processor.ProcessServiceGroup(serviceGroup, failedNodeNames)
	suite.Assert().Nil(err)
	client.AssertExpectations(suite.T())
}

func (suite *ServiceGroupProcessorTestSuite) TestProcessServiceGroup_deleteMemberFails() {
	a10error := new(mocks.A10Error)
	client := suite.client
	processor := suite.helper.BuildServiceGroupProcessor(client)
	serviceGroup := serviceGroup()
	existing := *serviceGroup
	healthCheck := *serviceGroup.Health
	healthCheck.Name = "changed"
	existing.Health = &healthCheck
	failedNodeNames := []string{"server_down"}
	extraServer := &model.Member{
		ServerName: "server2",
		Port:       8080,
	}
	existing.Members = []*model.Member{
		&model.Member{
			ServerName: "server",
			Port:       8080,
		},
		extraServer,
	}

	client.On("GetServiceGroup", serviceGroup.Name).Once().Return(&existing, nil)
	client.On("UpdateServiceGroup", serviceGroup).Once().Return(nil)
	client.On("DeleteMember", extraServer).Once().Return(a10error)
	err := processor.ProcessServiceGroup(serviceGroup, failedNodeNames)
	suite.Assert().NotNil(err)
	client.AssertExpectations(suite.T())
}

func serviceGroup() *model.ServiceGroup {
	return &model.ServiceGroup{
		Health: &model.HealthCheck{
			Name: "test",
		},
		Name: "service group",
		IngressControllers: []*model.IngressController{
			&model.IngressController{
				Health: &model.HealthCheck{
					Name: "test",
				},
				Name: "ingress 1",
				Nodes: []*model.Node{
					&model.Node{
						Name:      "server",
						A10Server: "server",
					},
				},
				Port: 8080,
			},
		},
	}
}
