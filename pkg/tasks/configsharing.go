package tasks

import (
	"github.com/openshift/cluster-monitoring-operator/pkg/client"
	"github.com/openshift/cluster-monitoring-operator/pkg/manifests"
	"github.com/pkg/errors"
)

type ConfigSharingTask struct {
	client	*client.Client
	factory	*manifests.Factory
}

func NewConfigSharingTask(client *client.Client, factory *manifests.Factory) *ConfigSharingTask {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &ConfigSharingTask{client: client, factory: factory}
}
func (t *ConfigSharingTask) Run() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	promRoute, err := t.factory.PrometheusK8sRoute()
	if err != nil {
		return errors.Wrap(err, "initializing Prometheus Route failed")
	}
	promURL, err := t.client.GetRouteURL(promRoute)
	if err != nil {
		return errors.Wrap(err, "failed to retrieve Prometheus host")
	}
	amRoute, err := t.factory.AlertmanagerRoute()
	if err != nil {
		return errors.Wrap(err, "initializing Alertmanager Route failed")
	}
	amURL, err := t.client.GetRouteURL(amRoute)
	if err != nil {
		return errors.Wrap(err, "failed to retrieve Alertmanager host")
	}
	grafanaRoute, err := t.factory.GrafanaRoute()
	if err != nil {
		return errors.Wrap(err, "initializing Grafana Route failed")
	}
	grafanaURL, err := t.client.GetRouteURL(grafanaRoute)
	if err != nil {
		return errors.Wrap(err, "failed to retrieve Grafana host")
	}
	cm := t.factory.SharingConfig(promURL, amURL, grafanaURL)
	err = t.client.CreateOrUpdateConfigMap(cm)
	if err != nil {
		return errors.Wrap(err, "reconciling Sharing Config ConfigMap failed")
	}
	return nil
}
