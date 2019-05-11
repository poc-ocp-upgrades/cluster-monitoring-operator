package client

import (
	v1 "github.com/openshift/api/config/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type conditions struct {
	entryMap						map[v1.ClusterStatusConditionType]v1.ClusterOperatorStatusCondition
	currentVersion, targetVersion	string
}

func newConditions(cos v1.ClusterOperatorStatus, targetVersion string, time metav1.Time) *conditions {
	_logClusterCodePath()
	defer _logClusterCodePath()
	entries := map[v1.ClusterStatusConditionType]v1.ClusterOperatorStatusCondition{v1.OperatorAvailable: {Type: v1.OperatorAvailable, Status: v1.ConditionUnknown, LastTransitionTime: time}, v1.OperatorProgressing: {Type: v1.OperatorProgressing, Status: v1.ConditionUnknown, LastTransitionTime: time}, v1.OperatorFailing: {Type: v1.OperatorFailing, Status: v1.ConditionUnknown, LastTransitionTime: time}}
	for _, c := range cos.Conditions {
		entries[c.Type] = c
	}
	cs := &conditions{entryMap: entries}
	if len(cos.Versions) > 0 {
		cs.currentVersion = cos.Versions[0].Version
	}
	cs.targetVersion = targetVersion
	return cs
}
func (cs *conditions) setCondition(condition v1.ClusterStatusConditionType, status v1.ConditionStatus, message string, time metav1.Time) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	entries := make(map[v1.ClusterStatusConditionType]v1.ClusterOperatorStatusCondition)
	for k, v := range cs.entryMap {
		entries[k] = v
	}
	c, ok := cs.entryMap[condition]
	if !ok || c.Status != status || c.Message != message {
		entries[condition] = v1.ClusterOperatorStatusCondition{Type: condition, Status: status, LastTransitionTime: time, Message: message}
	}
	wantsProgressing := condition == v1.OperatorProgressing && status == v1.ConditionTrue
	available, hasAvailable := cs.entryMap[v1.OperatorAvailable]
	abort := wantsProgressing && hasAvailable && available.Status == v1.ConditionTrue
	abort = abort && cs.targetVersion == cs.currentVersion
	if abort {
		return
	}
	cs.entryMap = entries
}
func (cs *conditions) entries() []v1.ClusterOperatorStatusCondition {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var res []v1.ClusterOperatorStatusCondition
	for _, v := range cs.entryMap {
		res = append(res, v)
	}
	return res
}

type byType []v1.ClusterOperatorStatusCondition

func (b byType) Len() int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return len(b)
}
func (b byType) Swap(i, j int) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	b[i], b[j] = b[j], b[i]
}
func (b byType) Less(i, j int) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return b[i].Type < b[j].Type
}
