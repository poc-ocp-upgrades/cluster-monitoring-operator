package e2e

import (
	"flag"
	"fmt"
	"log"
	"testing"
	"time"
	"github.com/openshift/cluster-monitoring-operator/test/e2e/framework"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/clientcmd"
)

var f *framework.Framework

func TestMain(m *testing.M) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if err := testMain(m); err != nil {
		log.Fatal(err)
	}
}
func testMain(m *testing.M) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	kubeConfigPath := flag.String("kubeconfig", clientcmd.RecommendedHomeFile, "kube config path, default: $HOME/.kube/config")
	flag.Parse()
	var (
		err	error
		cleanUp	func() error
	)
	f, cleanUp, err = framework.New(*kubeConfigPath)
	if cleanUp != nil {
		defer cleanUp()
	}
	if err != nil {
		return err
	}
	err = wait.Poll(time.Second, 5*time.Minute, func() (bool, error) {
		_, err := f.KubeClient.Apps().Deployments(f.Ns).Get("prometheus-operator", metav1.GetOptions{})
		if err != nil {
			return false, nil
		}
		return true, nil
	})
	if err != nil {
		return err
	}
	var loopErr error
	err = wait.Poll(5*time.Second, 1*time.Minute, func() (bool, error) {
		var (
			body	[]byte
			v	int
		)
		body, loopErr = f.PrometheusK8sClient.Query("count(up{job=\"prometheus-k8s\"})")
		if loopErr != nil {
			return false, nil
		}
		v, loopErr = framework.GetFirstValueFromPromQuery(body)
		if loopErr != nil {
			return false, nil
		}
		if v != 2 {
			loopErr = fmt.Errorf("expected 2 Prometheus instances but got: %v", v)
			return false, nil
		}
		return true, nil
	})
	if err != nil {
		return errors.Wrapf(err, "wait for prometheus-k8s: %v", loopErr)
	}
	if m.Run() != 0 {
		return errors.New("tests failed")
	}
	return nil
}
func TestTargetsUp(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	targets := []string{"node-exporter", "kubelet", "apiserver", "kube-state-metrics", "prometheus-k8s", "prometheus-operator", "alertmanager-main", "crio", "telemeter-client", "etcd"}
	for _, target := range targets {
		f.PrometheusK8sClient.WaitForQueryReturnOne(t, time.Minute, "max(up{job=\""+target+"\"})")
	}
}
func TestMemoryUsageRecordingRule(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	f.PrometheusK8sClient.WaitForQueryReturnGreaterEqualOne(t, time.Minute, "count(namespace:container_memory_usage_bytes:sum)")
}
