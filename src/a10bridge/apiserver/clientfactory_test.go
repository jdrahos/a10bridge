package apiserver_test

import (
	"a10bridge/apiserver"
	"errors"
	"testing"

	"k8s.io/client-go/kubernetes/fake"

	"k8s.io/client-go/kubernetes"

	"github.com/stretchr/testify/suite"
	"k8s.io/client-go/rest"
)

type ClientFactoryTestSuite struct {
	suite.Suite
	helper *apiserver.TestHelper
}

func TestClientFactory(t *testing.T) {
	tests := new(ClientFactoryTestSuite)
	tests.helper = new(apiserver.TestHelper)
	suite.Run(t, tests)
}

func (suite *ClientFactoryTestSuite) TestCreateClient_inCluster() {
	originalRestInClusterConfig := suite.helper.SetRestInClusterConfig(func() (*rest.Config, error) {
		return &rest.Config{}, nil
	})
	defer suite.helper.SetRestInClusterConfig(originalRestInClusterConfig)
	originalKubernetesNewForConfig := suite.helper.SetKubernetesNewForConfigFunc(func(c *rest.Config) (*kubernetes.Clientset, error) {
		return &kubernetes.Clientset{}, nil
	})
	defer suite.helper.SetKubernetesNewForConfigFunc(originalKubernetesNewForConfig)

	client, err := apiserver.CreateClient()
	suite.Assert().Nil(err)
	suite.Assert().NotNil(client)
}

func (suite *ClientFactoryTestSuite) TestCreateClient_inClusterConfigCreationFails() {
	originalRestInClusterConfig := suite.helper.SetRestInClusterConfig(func() (*rest.Config, error) {
		return nil, errors.New("fail")
	})
	defer suite.helper.SetRestInClusterConfig(originalRestInClusterConfig)
	originalClientcmdBuildConfigFromFlags := suite.helper.SetClientcmdBuildConfigFromFlagsFunc(func(masterUrl, kubeconfigPath string) (*rest.Config, error) {
		return &rest.Config{}, nil
	})
	defer suite.helper.SetClientcmdBuildConfigFromFlagsFunc(originalClientcmdBuildConfigFromFlags)
	originalKubernetesNewForConfig := suite.helper.SetKubernetesNewForConfigFunc(func(c *rest.Config) (*kubernetes.Clientset, error) {
		return &kubernetes.Clientset{}, nil
	})
	defer suite.helper.SetKubernetesNewForConfigFunc(originalKubernetesNewForConfig)

	client, err := apiserver.CreateClient()
	suite.Assert().Nil(err)
	suite.Assert().NotNil(client)
}

func (suite *ClientFactoryTestSuite) TestCreateClient_inClusterClientSetCreationFails() {
	originalRestInClusterConfig := suite.helper.SetRestInClusterConfig(func() (*rest.Config, error) {
		return &rest.Config{}, nil
	})
	defer suite.helper.SetRestInClusterConfig(originalRestInClusterConfig)
	originalKubernetesNewForConfig := suite.helper.SetKubernetesNewForConfigFunc(func(c *rest.Config) (*kubernetes.Clientset, error) {
		//second pass
		suite.helper.SetKubernetesNewForConfigFunc(func(c *rest.Config) (*kubernetes.Clientset, error) {
			return &kubernetes.Clientset{}, nil
		})
		//first fail
		return nil, errors.New("fail")
	})
	defer suite.helper.SetKubernetesNewForConfigFunc(originalKubernetesNewForConfig)
	originalClientcmdBuildConfigFromFlags := suite.helper.SetClientcmdBuildConfigFromFlagsFunc(func(masterUrl, kubeconfigPath string) (*rest.Config, error) {
		return &rest.Config{}, nil
	})
	defer suite.helper.SetClientcmdBuildConfigFromFlagsFunc(originalClientcmdBuildConfigFromFlags)

	client, err := apiserver.CreateClient()
	suite.Assert().Nil(err)
	suite.Assert().NotNil(client)
}

func (suite *ClientFactoryTestSuite) TestCreateClient_inClusterFailsKubectlClientSetCreationFails() {
	originalRestInClusterConfig := suite.helper.SetRestInClusterConfig(func() (*rest.Config, error) {
		return nil, errors.New("fail in incluster")
	})
	defer suite.helper.SetRestInClusterConfig(originalRestInClusterConfig)
	originalClientcmdBuildConfigFromFlags := suite.helper.SetClientcmdBuildConfigFromFlagsFunc(func(masterUrl, kubeconfigPath string) (*rest.Config, error) {
		return nil, errors.New("fail in kubectl")
	})
	defer suite.helper.SetClientcmdBuildConfigFromFlagsFunc(originalClientcmdBuildConfigFromFlags)

	client, err := apiserver.CreateClient()
	suite.Assert().NotNil(err)
	suite.Assert().Nil(client)
}

func (suite *ClientFactoryTestSuite) TestCreateClient_inClusterFailsKubectlConfigCreationFails() {
	originalRestInClusterConfig := suite.helper.SetRestInClusterConfig(func() (*rest.Config, error) {
		return nil, errors.New("fail in incluster")
	})
	defer suite.helper.SetRestInClusterConfig(originalRestInClusterConfig)
	originalClientcmdBuildConfigFromFlags := suite.helper.SetClientcmdBuildConfigFromFlagsFunc(func(masterUrl, kubeconfigPath string) (*rest.Config, error) {
		return &rest.Config{}, nil
	})
	defer suite.helper.SetClientcmdBuildConfigFromFlagsFunc(originalClientcmdBuildConfigFromFlags)
	originalKubernetesNewForConfig := suite.helper.SetKubernetesNewForConfigFunc(func(c *rest.Config) (*kubernetes.Clientset, error) {
		return nil, errors.New("fail in kubectl")
	})
	defer suite.helper.SetKubernetesNewForConfigFunc(originalKubernetesNewForConfig)

	client, err := apiserver.CreateClient()
	suite.Assert().NotNil(err)
	suite.Assert().Nil(client)
}

func (suite *ClientFactoryTestSuite) TestFakeClientInjection() {
	fakeClientSet := fake.NewSimpleClientset()
	apiserver.InjectFakeClient(fakeClientSet)
	client, err := apiserver.CreateClient()
	suite.Assert().Nil(err)
	suite.Assert().NotNil(client)
}
