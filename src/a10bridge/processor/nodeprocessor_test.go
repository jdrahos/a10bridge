package processor_test

import (
	"a10bridge/mocks"
	"a10bridge/model"
	"a10bridge/processor"
	"testing"

	"github.com/stretchr/testify/suite"
)

type NodeProcessorTestSuite struct {
	suite.Suite
	helper *processor.TestHelper
	client *mocks.Client
}

func (suite *NodeProcessorTestSuite) SetupTest() {
	suite.client = new(mocks.Client)
}

func TestNodeProcessor(t *testing.T) {
	tests := new(NodeProcessorTestSuite)
	tests.helper = new(processor.TestHelper)
	suite.Run(t, tests)
}

func (suite *NodeProcessorTestSuite) TestProcessNode_notChanged() {
	client := suite.client
	processor := suite.helper.BuildNodeProcessor(client)
	node := node()

	client.On("GetServer", node.A10Server).Once().Return(node, nil)
	err := processor.ProcessNode(node)
	suite.Assert().Nil(err)
	client.AssertExpectations(suite.T())
}

func (suite *NodeProcessorTestSuite) TestProcessNode_getFails() {
	a10error := new(mocks.A10Error)
	client := suite.client
	processor := suite.helper.BuildNodeProcessor(client)
	node := node()

	client.On("GetServer", node.A10Server).Once().Return(nil, a10error)
	client.On("IsServerNotFound", a10error).Once().Return(false)
	err := processor.ProcessNode(node)
	suite.Assert().NotNil(err)
	client.AssertExpectations(suite.T())
}

func (suite *NodeProcessorTestSuite) TestProcessNode_serverNotFound() {
	a10error := new(mocks.A10Error)
	client := suite.client
	processor := suite.helper.BuildNodeProcessor(client)
	node := node()

	client.On("GetServer", node.A10Server).Once().Return(nil, a10error)
	client.On("IsServerNotFound", a10error).Once().Return(true)
	client.On("CreateServer", node).Once().Return(nil)
	err := processor.ProcessNode(node)
	suite.Assert().Nil(err)
	client.AssertExpectations(suite.T())
}

func (suite *NodeProcessorTestSuite) TestProcessNode_createServerFails() {
	a10error := new(mocks.A10Error)
	client := suite.client
	processor := suite.helper.BuildNodeProcessor(client)
	node := node()

	client.On("GetServer", node.A10Server).Once().Return(nil, a10error)
	client.On("IsServerNotFound", a10error).Once().Return(true)
	client.On("CreateServer", node).Once().Return(a10error)
	err := processor.ProcessNode(node)
	suite.Assert().NotNil(err)
	client.AssertExpectations(suite.T())
}

func (suite *NodeProcessorTestSuite) TestProcessNode_ipChanged() {
	client := suite.client
	processor := suite.helper.BuildNodeProcessor(client)
	node := node()
	existing := *node
	existing.IPAddress = "blah"

	client.On("GetServer", node.A10Server).Once().Return(&existing, nil)
	client.On("UpdateServer", node).Once().Return(nil)
	err := processor.ProcessNode(node)
	suite.Assert().Nil(err)
	client.AssertExpectations(suite.T())
}

func (suite *NodeProcessorTestSuite) TestProcessNode_weightChanged() {
	client := suite.client
	processor := suite.helper.BuildNodeProcessor(client)
	node := node()
	existing := *node
	existing.Weight = "blah"

	client.On("GetServer", node.A10Server).Once().Return(&existing, nil)
	client.On("UpdateServer", node).Once().Return(nil)
	err := processor.ProcessNode(node)
	suite.Assert().Nil(err)
	client.AssertExpectations(suite.T())
}

func (suite *NodeProcessorTestSuite) TestProcessNode_updateFails() {
	a10error := new(mocks.A10Error)
	client := suite.client
	processor := suite.helper.BuildNodeProcessor(client)
	node := node()
	existing := *node
	existing.Weight = "blah"

	client.On("GetServer", node.A10Server).Once().Return(&existing, nil)
	client.On("UpdateServer", node).Once().Return(a10error)
	err := processor.ProcessNode(node)
	suite.Assert().NotNil(err)
	client.AssertExpectations(suite.T())
}

func node() *model.Node {
	return &model.Node{
		A10Server: "a10server",
		IPAddress: "10.10.10.10",
		Labels: map[string]string{
			"test": "value",
		},
		Name:   "server",
		Weight: "1",
	}
}
