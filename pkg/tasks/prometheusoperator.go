package tasks

import (
	"github.com/openshift/cluster-monitoring-operator/pkg/client"
	"github.com/openshift/cluster-monitoring-operator/pkg/manifests"
	"github.com/pkg/errors"
)

type PrometheusOperatorTask struct {
	client	*client.Client
	factory	*manifests.Factory
}

func NewPrometheusOperatorTask(client *client.Client, factory *manifests.Factory) *PrometheusOperatorTask {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &PrometheusOperatorTask{client: client, factory: factory}
}
func (t *PrometheusOperatorTask) Run() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	sa, err := t.factory.PrometheusOperatorServiceAccount()
	if err != nil {
		return errors.Wrap(err, "initializing Prometheus Operator ServiceAccount failed")
	}
	err = t.client.CreateOrUpdateServiceAccount(sa)
	if err != nil {
		return errors.Wrap(err, "reconciling Prometheus Operator ServiceAccount failed")
	}
	cr, err := t.factory.PrometheusOperatorClusterRole()
	if err != nil {
		return errors.Wrap(err, "initializing Prometheus Operator ClusterRole failed")
	}
	err = t.client.CreateOrUpdateClusterRole(cr)
	if err != nil {
		return errors.Wrap(err, "reconciling Prometheus Operator ClusterRole failed")
	}
	crb, err := t.factory.PrometheusOperatorClusterRoleBinding()
	if err != nil {
		return errors.Wrap(err, "initializing Prometheus Operator ClusterRoleBinding failed")
	}
	err = t.client.CreateOrUpdateClusterRoleBinding(crb)
	if err != nil {
		return errors.Wrap(err, "reconciling Prometheus Operator ClusterRoleBinding failed")
	}
	svc, err := t.factory.PrometheusOperatorService()
	if err != nil {
		return errors.Wrap(err, "initializing Prometheus Operator Service failed")
	}
	err = t.client.CreateOrUpdateService(svc)
	if err != nil {
		return errors.Wrap(err, "reconciling Prometheus Operator Service failed")
	}
	namespaces, err := t.client.NamespacesToMonitor()
	if err != nil {
		return errors.Wrap(err, "listing namespaces to monitor failed")
	}
	d, err := t.factory.PrometheusOperatorDeployment(namespaces)
	if err != nil {
		return errors.Wrap(err, "initializing Prometheus Operator Deployment failed")
	}
	err = t.client.CreateOrUpdateDeployment(d)
	if err != nil {
		return errors.Wrap(err, "reconciling Prometheus Operator Deployment failed")
	}
	err = t.client.WaitForPrometheusOperatorCRDsReady()
	if err != nil {
		return errors.Wrap(err, "waiting for Prometheus CRDs to become available failed")
	}
	smpo, err := t.factory.PrometheusOperatorServiceMonitor()
	if err != nil {
		return errors.Wrap(err, "initializing Prometheus Operator ServiceMonitor failed")
	}
	err = t.client.CreateOrUpdateServiceMonitor(smpo)
	return errors.Wrap(err, "reconciling Prometheus Operator ServiceMonitor failed")
}
