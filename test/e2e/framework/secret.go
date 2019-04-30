package framework

import (
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes"
	poTestFramework "github.com/coreos/prometheus-operator/test/framework"
	"github.com/pkg/errors"
)

func CreateSecret(kubeClient kubernetes.Interface, namespace string, relativePath string) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	secret, err := parseSecretYaml(relativePath)
	if err != nil {
		return errors.Wrap(err, "parsing secret failed")
	}
	if _, err := kubeClient.CoreV1().Secrets(namespace).Create(secret); err != nil {
		return errors.Wrap(err, "creating secret failed")
	}
	return nil
}
func parseSecretYaml(relativePath string) (*v1.Secret, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	manifest, err := poTestFramework.PathToOSFile(relativePath)
	if err != nil {
		return nil, err
	}
	secret := v1.Secret{}
	if err := yaml.NewYAMLOrJSONDecoder(manifest, 100).Decode(&secret); err != nil {
		return nil, err
	}
	return &secret, nil
}
