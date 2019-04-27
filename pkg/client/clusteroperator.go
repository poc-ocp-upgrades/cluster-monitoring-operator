package client

import (
	"fmt"
	v1 "github.com/openshift/api/config/v1"
	clientv1 "github.com/openshift/client-go/config/clientset/versioned/typed/config/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type StatusReporter struct {
	client			clientv1.ClusterOperatorInterface
	clusterOperatorName	string
	version			string
}

func NewStatusReporter(client clientv1.ClusterOperatorInterface, name, version string) *StatusReporter {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &StatusReporter{client: client, clusterOperatorName: name, version: version}
}
func (r *StatusReporter) SetDone() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	co, err := r.client.Get(r.clusterOperatorName, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		co = r.newClusterOperator()
		co, err = r.client.Create(co)
	}
	if err != nil && !apierrors.IsNotFound(err) {
		return err
	}
	time := metav1.Now()
	conditions := newConditions(co.Status, r.version, time)
	conditions.setCondition(v1.OperatorAvailable, v1.ConditionTrue, "Successfully rolled out the stack.", time)
	conditions.setCondition(v1.OperatorProgressing, v1.ConditionFalse, "", time)
	conditions.setCondition(v1.OperatorFailing, v1.ConditionFalse, "", time)
	co.Status.Conditions = conditions.entries()
	if len(r.version) > 0 {
		co.Status.Versions = []v1.OperandVersion{{Name: "operator", Version: r.version}}
	} else {
		co.Status.Versions = nil
	}
	_, err = r.client.UpdateStatus(co)
	return err
}
func (r *StatusReporter) SetInProgress() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	co, err := r.client.Get(r.clusterOperatorName, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		co = r.newClusterOperator()
		co, err = r.client.Create(co)
	}
	if err != nil && !apierrors.IsNotFound(err) {
		return err
	}
	time := metav1.Now()
	conditions := newConditions(co.Status, r.version, time)
	conditions.setCondition(v1.OperatorProgressing, v1.ConditionTrue, "Rolling out the stack.", time)
	co.Status.Conditions = conditions.entries()
	_, err = r.client.UpdateStatus(co)
	return err
}
func (r *StatusReporter) SetFailed(statusErr error) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	co, err := r.client.Get(r.clusterOperatorName, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		co = r.newClusterOperator()
		co, err = r.client.Create(co)
	}
	if err != nil && !apierrors.IsNotFound(err) {
		return err
	}
	time := metav1.Now()
	conditions := newConditions(co.Status, r.version, time)
	conditions.setCondition(v1.OperatorAvailable, v1.ConditionFalse, "", time)
	conditions.setCondition(v1.OperatorProgressing, v1.ConditionFalse, "", time)
	conditions.setCondition(v1.OperatorFailing, v1.ConditionTrue, fmt.Sprintf("Failed to rollout the stack. Error: %v", statusErr), time)
	co.Status.Conditions = conditions.entries()
	_, err = r.client.UpdateStatus(co)
	return err
}
func (r *StatusReporter) newClusterOperator() *v1.ClusterOperator {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	time := metav1.Now()
	co := &v1.ClusterOperator{TypeMeta: metav1.TypeMeta{APIVersion: "config.openshift.io/v1", Kind: "ClusterOperator"}, ObjectMeta: metav1.ObjectMeta{Name: r.clusterOperatorName}, Spec: v1.ClusterOperatorSpec{}, Status: v1.ClusterOperatorStatus{}}
	co.Status.Conditions = newConditions(co.Status, r.version, time).entries()
	return co
}
