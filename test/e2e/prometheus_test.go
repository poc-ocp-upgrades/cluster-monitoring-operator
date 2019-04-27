package e2e

import (
	"log"
	"testing"
	"time"
	"github.com/pkg/errors"
	"k8s.io/api/apps/v1beta2"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

func TestPrometheusVolumeClaim(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	err := f.OperatorClient.WaitForStatefulsetRollout(&v1beta2.StatefulSet{ObjectMeta: metav1.ObjectMeta{Name: "prometheus-k8s", Namespace: f.Ns}})
	if err != nil {
		log.Fatal(err)
	}
	cm := &v1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: "cluster-monitoring-config", Namespace: f.Ns}, Data: map[string]string{"config.yaml": `prometheusK8s:
  volumeClaimTemplate:
    spec:
      storageClassName: gp2
      resources:
        requests:
          storage: 2Gi
`}}
	if err := f.OperatorClient.CreateOrUpdateConfigMap(cm); err != nil {
		log.Fatal(err)
	}
	var lastErr error
	err = wait.Poll(time.Second, 5*time.Minute, func() (bool, error) {
		_, err := f.KubeClient.CoreV1().PersistentVolumeClaims(f.Ns).Get("prometheus-k8s-db-prometheus-k8s-0", metav1.GetOptions{})
		lastErr = errors.Wrap(err, "getting prometheus persistent volume claim failed")
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
	err = f.OperatorClient.WaitForStatefulsetRollout(&v1beta2.StatefulSet{ObjectMeta: metav1.ObjectMeta{Name: "prometheus-k8s", Namespace: f.Ns}})
	if err != nil {
		log.Fatal(err)
	}
}
