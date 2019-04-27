package tasks

import (
	"github.com/openshift/cluster-monitoring-operator/pkg/client"
	"github.com/openshift/cluster-monitoring-operator/pkg/manifests"
	"github.com/pkg/errors"
)

type NodeExporterTask struct {
	client	*client.Client
	factory	*manifests.Factory
}

func NewNodeExporterTask(client *client.Client, factory *manifests.Factory) *NodeExporterTask {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &NodeExporterTask{client: client, factory: factory}
}
func (t *NodeExporterTask) Run() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	scc, err := t.factory.NodeExporterSecurityContextConstraints()
	if err != nil {
		return errors.Wrap(err, "initializing node-exporter SecurityContextConstraints failed")
	}
	err = t.client.CreateOrUpdateSecurityContextConstraints(scc)
	if err != nil {
		return errors.Wrap(err, "reconciling node-exporter SecurityContextConstraints failed")
	}
	sa, err := t.factory.NodeExporterServiceAccount()
	if err != nil {
		return errors.Wrap(err, "initializing node-exporter Service failed")
	}
	err = t.client.CreateOrUpdateServiceAccount(sa)
	if err != nil {
		return errors.Wrap(err, "reconciling node-exporter ServiceAccount failed")
	}
	cr, err := t.factory.NodeExporterClusterRole()
	if err != nil {
		return errors.Wrap(err, "initializing node-exporter ClusterRole failed")
	}
	err = t.client.CreateOrUpdateClusterRole(cr)
	if err != nil {
		return errors.Wrap(err, "reconciling node-exporter ClusterRole failed")
	}
	crb, err := t.factory.NodeExporterClusterRoleBinding()
	if err != nil {
		return errors.Wrap(err, "initializing node-exporter ClusterRoleBinding failed")
	}
	err = t.client.CreateOrUpdateClusterRoleBinding(crb)
	if err != nil {
		return errors.Wrap(err, "reconciling node-exporter ClusterRoleBinding failed")
	}
	svc, err := t.factory.NodeExporterService()
	if err != nil {
		return errors.Wrap(err, "initializing node-exporter Service failed")
	}
	err = t.client.CreateOrUpdateService(svc)
	if err != nil {
		return errors.Wrap(err, "reconciling node-exporter Service failed")
	}
	ds, err := t.factory.NodeExporterDaemonSet()
	if err != nil {
		return errors.Wrap(err, "initializing node-exporter DaemonSet failed")
	}
	err = t.client.CreateOrUpdateDaemonSet(ds)
	if err != nil {
		return errors.Wrap(err, "reconciling node-exporter DaemonSet failed")
	}
	smn, err := t.factory.NodeExporterServiceMonitor()
	if err != nil {
		return errors.Wrap(err, "initializing node-exporter ServiceMonitor failed")
	}
	err = t.client.CreateOrUpdateServiceMonitor(smn)
	return errors.Wrap(err, "reconciling node-exporter ServiceMonitor failed")
}
