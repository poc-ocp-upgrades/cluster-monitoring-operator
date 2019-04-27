package e2e

import (
	"fmt"
	"log"
	"sort"
	"testing"
	"time"
	"github.com/openshift/cluster-monitoring-operator/test/e2e/framework"
	"github.com/openshift/cluster-monitoring-operator/pkg/manifests"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

func TestAggregatedMetricPermissions(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	present := func(where []string, what string) bool {
		sort.Strings(where)
		i := sort.SearchStrings(where, what)
		return i < len(where) && where[i] == what
	}
	type checkFunc func(clusterRole string) error
	hasRule := func(apiGroup, resource, verb string) checkFunc {
		return func(clusterRole string) error {
			return framework.Poll(time.Second, 5*time.Minute, func() error {
				viewRole, err := f.KubeClient.RbacV1().ClusterRoles().Get(clusterRole, metav1.GetOptions{})
				if err != nil {
					return errors.Wrapf(err, "getting %s cluster role failed", clusterRole)
				}
				for _, rule := range viewRole.Rules {
					if !present(rule.APIGroups, apiGroup) {
						continue
					}
					if !present(rule.Resources, resource) {
						continue
					}
					if !present(rule.Verbs, verb) {
						continue
					}
					return nil
				}
				return fmt.Errorf("could not find metrics in cluster role %s", clusterRole)
			})
		}
	}
	canGetPodMetrics := hasRule("metrics.k8s.io", "pods", "get")
	for _, tc := range []struct {
		clusterRole	string
		check		checkFunc
	}{{clusterRole: "view", check: canGetPodMetrics}, {clusterRole: "edit", check: canGetPodMetrics}, {clusterRole: "admin", check: canGetPodMetrics}} {
		t.Run(tc.clusterRole, func(t *testing.T) {
			if err := tc.check(tc.clusterRole); err != nil {
				t.Error(err)
			}
		})
	}
}
func TestPrometheusAdapterCARotation(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	var lastErr error
	err := wait.Poll(time.Second, 5*time.Minute, func() (bool, error) {
		_, err := f.KubeClient.Apps().Deployments(f.Ns).Get("prometheus-adapter", metav1.GetOptions{})
		lastErr = errors.Wrap(err, "getting prometheus-adapter deployment failed")
		if err != nil {
			return false, nil
		}
		return true, nil
	})
	if err != nil {
		if err == wait.ErrWaitTimeout && lastErr != nil {
			err = lastErr
		}
		log.Fatal(err)
	}
	apiAuth, err := f.KubeClient.CoreV1().ConfigMaps("kube-system").Get("extension-apiserver-authentication", metav1.GetOptions{})
	if err != nil {
		log.Fatal(err)
	}
	tls, err := f.KubeClient.CoreV1().Secrets("openshift-monitoring").Get("prometheus-adapter-tls", metav1.GetOptions{})
	if err != nil {
		log.Fatal(err)
	}
	apiAuth.Data["requestheader-client-ca-file"] = apiAuth.Data["requestheader-client-ca-file"] + "\n"
	apiAuth, err = f.KubeClient.CoreV1().ConfigMaps("kube-system").Update(apiAuth)
	if err != nil {
		log.Fatal(err)
	}
	factory := manifests.NewFactory("openshift-monitoring", nil)
	newSecret, err := factory.PrometheusAdapterSecret(tls, apiAuth)
	if err != nil {
		log.Fatal(err)
	}
	err = wait.Poll(time.Second, 5*time.Minute, func() (bool, error) {
		_, err := f.KubeClient.CoreV1().Secrets(f.Ns).Get(newSecret.Name, metav1.GetOptions{})
		lastErr = errors.Wrap(err, "getting new api auth secret failed")
		if err != nil {
			return false, nil
		}
		return true, nil
	})
	if err != nil {
		if err == wait.ErrWaitTimeout && lastErr != nil {
			err = lastErr
		}
		log.Fatal(err)
	}
	err = wait.Poll(time.Second, 5*time.Minute, func() (bool, error) {
		d, err := f.KubeClient.Apps().Deployments(f.Ns).Get("prometheus-adapter", metav1.GetOptions{})
		lastErr = errors.Wrap(err, "getting new prometheus adapter deployment failed")
		if err != nil {
			return false, nil
		}
		lastErr = fmt.Errorf("waiting for updated replica count=%d to be spec replica count=%d", d.Status.UpdatedReplicas, *d.Spec.Replicas)
		return d.Status.UpdatedReplicas == *d.Spec.Replicas, nil
	})
	if err != nil {
		if err == wait.ErrWaitTimeout && lastErr != nil {
			err = lastErr
		}
		log.Fatal(err)
	}
}
