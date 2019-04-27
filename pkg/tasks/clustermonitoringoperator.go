package tasks

import (
	"github.com/openshift/cluster-monitoring-operator/pkg/client"
	"github.com/openshift/cluster-monitoring-operator/pkg/manifests"
	"github.com/pkg/errors"
)

type ClusterMonitoringOperatorTask struct {
	client	*client.Client
	factory	*manifests.Factory
}

func NewClusterMonitoringOperatorTask(client *client.Client, factory *manifests.Factory) *ClusterMonitoringOperatorTask {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &ClusterMonitoringOperatorTask{client: client, factory: factory}
}
func (t *ClusterMonitoringOperatorTask) Run() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	svc, err := t.factory.ClusterMonitoringOperatorService()
	if err != nil {
		return errors.Wrap(err, "initializing Cluster Monitoring Operator Service failed")
	}
	err = t.client.CreateOrUpdateService(svc)
	if err != nil {
		return errors.Wrap(err, "reconciling Cluster Monitoring Operator Service failed")
	}
	cr, err := t.factory.ClusterMonitoringClusterRole()
	if err != nil {
		return errors.Wrap(err, "initializing cluster-monitoring ClusterRole failed")
	}
	err = t.client.CreateOrUpdateClusterRole(cr)
	if err != nil {
		return errors.Wrap(err, "reconciling cluster-monitoring ClusterRole failed")
	}
	smcmo, err := t.factory.ClusterMonitoringOperatorServiceMonitor()
	if err != nil {
		return errors.Wrap(err, "initializing Cluster Monitoring Operator ServiceMonitor failed")
	}
	err = t.client.CreateOrUpdateServiceMonitor(smcmo)
	return errors.Wrap(err, "reconciling Cluster Monitoring Operator ServiceMonitor failed")
}
