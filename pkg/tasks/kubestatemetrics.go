package tasks

import (
	"github.com/openshift/cluster-monitoring-operator/pkg/client"
	"github.com/openshift/cluster-monitoring-operator/pkg/manifests"
	"github.com/pkg/errors"
)

type KubeStateMetricsTask struct {
	client	*client.Client
	factory	*manifests.Factory
}

func NewKubeStateMetricsTask(client *client.Client, factory *manifests.Factory) *KubeStateMetricsTask {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &KubeStateMetricsTask{client: client, factory: factory}
}
func (t *KubeStateMetricsTask) Run() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	sa, err := t.factory.KubeStateMetricsServiceAccount()
	if err != nil {
		return errors.Wrap(err, "initializing kube-state-metrics Service failed")
	}
	err = t.client.CreateOrUpdateServiceAccount(sa)
	if err != nil {
		return errors.Wrap(err, "reconciling kube-state-metrics ServiceAccount failed")
	}
	cr, err := t.factory.KubeStateMetricsClusterRole()
	if err != nil {
		return errors.Wrap(err, "initializing kube-state-metrics ClusterRole failed")
	}
	err = t.client.CreateOrUpdateClusterRole(cr)
	if err != nil {
		return errors.Wrap(err, "reconciling kube-state-metrics ClusterRole failed")
	}
	crb, err := t.factory.KubeStateMetricsClusterRoleBinding()
	if err != nil {
		return errors.Wrap(err, "initializing kube-state-metrics ClusterRoleBinding failed")
	}
	err = t.client.CreateOrUpdateClusterRoleBinding(crb)
	if err != nil {
		return errors.Wrap(err, "reconciling kube-state-metrics ClusterRoleBinding failed")
	}
	svc, err := t.factory.KubeStateMetricsService()
	if err != nil {
		return errors.Wrap(err, "initializing kube-state-metrics Service failed")
	}
	err = t.client.CreateOrUpdateService(svc)
	if err != nil {
		return errors.Wrap(err, "reconciling kube-state-metrics Service failed")
	}
	dep, err := t.factory.KubeStateMetricsDeployment()
	if err != nil {
		return errors.Wrap(err, "initializing kube-state-metrics Deployment failed")
	}
	err = t.client.CreateOrUpdateDeployment(dep)
	if err != nil {
		return errors.Wrap(err, "reconciling kube-state-metrics Deployment failed")
	}
	sm, err := t.factory.KubeStateMetricsServiceMonitor()
	if err != nil {
		return errors.Wrap(err, "initializing kube-state-metrics ServiceMonitor failed")
	}
	err = t.client.CreateOrUpdateServiceMonitor(sm)
	return errors.Wrap(err, "reconciling kube-state-metrics ServiceMonitor failed")
}
