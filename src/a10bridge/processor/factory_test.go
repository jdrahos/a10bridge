package processor_test

import (
	"a10bridge/a10/api"
	"a10bridge/apiserver"
	"a10bridge/config"
	"a10bridge/mocks"
	"a10bridge/processor"
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"
)

type FactoryTestSuite struct {
	suite.Suite
	helper processor.TestHelper
}

func TestFactory(t *testing.T) {
	tests := new(FactoryTestSuite)
	tests.helper = processor.TestHelper{}
	suite.Run(t, tests)
}

func (suite *FactoryTestSuite) TestBuildK8sProcessor() {
	k8sClient := new(mocks.K8sClient)
	original := suite.helper.SetApiserverCreateClient(func() (apiserver.K8sClient, error) {
		return k8sClient, nil
	})
	defer suite.helper.SetApiserverCreateClient(original)

	kprocessor, err := processor.BuildK8sProcessor()
	suite.Assert().Nil(err)
	suite.Assert().NotNil(kprocessor)
}

func (suite *FactoryTestSuite) TestBuildK8sProcessor_createClientFailure() {
	original := suite.helper.SetApiserverCreateClient(func() (apiserver.K8sClient, error) {
		return nil, errors.New("test")
	})
	defer suite.helper.SetApiserverCreateClient(original)

	kprocessor, err := processor.BuildK8sProcessor()
	suite.Assert().NotNil(err)
	suite.Assert().Nil(kprocessor)
}

func (suite *FactoryTestSuite) TestBuildA10Processors() {
	a10Client := new(mocks.Client)
	original := suite.helper.SetA10BuildClient(func(a10Instance *config.A10Instance) (api.Client, api.A10Error) {
		return a10Client, nil
	})
	defer suite.helper.SetA10BuildClient(original)

	a10Processors, err := processor.BuildA10Processors(&config.A10Instance{APIVersion: 2})
	suite.Assert().Nil(err)
	suite.Assert().NotNil(a10Processors)
}

func (suite *FactoryTestSuite) TestBuildA10Processors_clientBuildFails() {
	a10Error := new(mocks.A10Error)
	original := suite.helper.SetA10BuildClient(func(a10Instance *config.A10Instance) (api.Client, api.A10Error) {
		return nil, a10Error
	})
	defer suite.helper.SetA10BuildClient(original)

	a10Processors, err := processor.BuildA10Processors(&config.A10Instance{APIVersion: 2})
	suite.Assert().NotNil(err)
	suite.Assert().Nil(a10Processors)
}

func (suite *FactoryTestSuite) TestDestroy() {
	a10Client := new(mocks.Client)
	original := suite.helper.SetA10BuildClient(func(a10Instance *config.A10Instance) (api.Client, api.A10Error) {
		return a10Client, nil
	})
	defer suite.helper.SetA10BuildClient(original)
	a10Processors, _ := processor.BuildA10Processors(&config.A10Instance{APIVersion: 2})

	a10Client.On("Close").Return(nil)
	a10Processors.Destroy()
	a10Client.AssertCalled(suite.T(), "Close")
}
