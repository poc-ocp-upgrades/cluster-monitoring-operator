package e2e

import (
	"log"
	"strconv"
	"testing"
	"time"
	monv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func TestMultinamespacePrometheusRule(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	t.Parallel()
	nsName := "openshift-test-prometheus-rules" + strconv.FormatInt(time.Now().Unix(), 36)
	err := f.OperatorClient.CreateOrUpdateNamespace(&v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: nsName, Labels: map[string]string{"openshift.io/cluster-monitoring": "true"}}})
	if err != nil {
		log.Fatal(err)
	}
	defer f.OperatorClient.DeleteIfExists(nsName)
	err = f.OperatorClient.CreateOrUpdatePrometheusRule(&monv1.PrometheusRule{ObjectMeta: metav1.ObjectMeta{Name: "non-monitoring-prometheus-rules", Namespace: nsName}, Spec: monv1.PrometheusRuleSpec{Groups: []monv1.RuleGroup{{Name: "test-group", Rules: []monv1.Rule{{Alert: "AdditionalTestAlertRule", Expr: intstr.FromString("vector(1)")}}}}}})
	if err != nil {
		log.Fatal(err)
	}
	f.PrometheusK8sClient.WaitForQueryReturnOne(t, 10*time.Minute, `count(ALERTS{alertname="AdditionalTestAlertRule"} == 1)`)
}
