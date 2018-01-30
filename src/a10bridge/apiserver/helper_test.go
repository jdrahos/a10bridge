package apiserver

import (
	"net"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/rest"
)

type TestHelper struct{}
type RestInClusterConfigFunc func() (*rest.Config, error)
type ClientcmdBuildConfigFromFlagsFunc func(masterUrl, kubeconfigPath string) (*rest.Config, error)
type KubernetesNewForConfigFunc func(c *rest.Config) (*kubernetes.Clientset, error)
type NetLookupIPFunc func(host string) ([]net.IP, error)

func (helper *TestHelper) BuildClient(clientset *fake.Clientset) K8sClient {
	return clientImpl{
		corev1Impl:            clientset.CoreV1(),
		extensionsv1beta1Impl: clientset.ExtensionsV1beta1(),
	}
}

func (helper TestHelper) SetRestInClusterConfig(inClusterConfigFunc RestInClusterConfigFunc) RestInClusterConfigFunc {
	old := restInClusterConfig
	restInClusterConfig = inClusterConfigFunc
	return old
}

func (helper TestHelper) SetClientcmdBuildConfigFromFlagsFunc(buildConfigFromFlagsFunc ClientcmdBuildConfigFromFlagsFunc) ClientcmdBuildConfigFromFlagsFunc {
	old := clientcmdBuildConfigFromFlags
	clientcmdBuildConfigFromFlags = buildConfigFromFlagsFunc
	return old
}

func (helper TestHelper) SetKubernetesNewForConfigFunc(newForConfigFunc KubernetesNewForConfigFunc) KubernetesNewForConfigFunc {
	old := kubernetesNewForConfig
	kubernetesNewForConfig = newForConfigFunc
	return old
}

func (helper TestHelper) SetNetLookupIPFunc(netLookupIPFunc NetLookupIPFunc) NetLookupIPFunc {
	old := netLookupIP
	netLookupIP = netLookupIPFunc
	return old
}
