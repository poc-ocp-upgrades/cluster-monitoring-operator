package client

import (
	"fmt"
	"reflect"
	"sort"
	"testing"
	configv1 "github.com/openshift/api/config/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestConditions(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	type checkFunc func(*conditions) error
	hasConditions := func(want []configv1.ClusterOperatorStatusCondition) checkFunc {
		return func(cs *conditions) error {
			got := cs.entries()
			sort.Sort(byType(got))
			sort.Sort(byType(want))
			if !reflect.DeepEqual(got, want) {
				return fmt.Errorf("got conditions\n%+v\nwant\n%+v", got, want)
			}
			return nil
		}
	}
	allUnknown := hasConditions([]configv1.ClusterOperatorStatusCondition{{Type: configv1.OperatorProgressing, Status: configv1.ConditionUnknown, LastTransitionTime: v1.Time{}, Message: ""}, {Type: configv1.OperatorAvailable, Status: configv1.ConditionUnknown, LastTransitionTime: v1.Time{}, Message: ""}, {Type: configv1.OperatorFailing, Status: configv1.ConditionUnknown, LastTransitionTime: v1.Time{}, Message: ""}})
	for _, tc := range []struct {
		name		string
		conditions	func() *conditions
		check		checkFunc
	}{{name: "initial nil conditions", conditions: func() *conditions {
		return newConditions(configv1.ClusterOperatorStatus{}, "", v1.Time{})
	}, check: allUnknown}, {name: "initial empty conditions", conditions: func() *conditions {
		return newConditions(configv1.ClusterOperatorStatus{Conditions: []configv1.ClusterOperatorStatusCondition{}}, "", v1.Time{})
	}, check: allUnknown}, {name: "initial failing condition", conditions: func() *conditions {
		return newConditions(configv1.ClusterOperatorStatus{Conditions: []configv1.ClusterOperatorStatusCondition{{Type: configv1.OperatorFailing, Status: configv1.ConditionTrue}}}, "", v1.Time{})
	}, check: hasConditions([]configv1.ClusterOperatorStatusCondition{{Type: configv1.OperatorProgressing, Status: configv1.ConditionUnknown, LastTransitionTime: v1.Time{}, Message: ""}, {Type: configv1.OperatorAvailable, Status: configv1.ConditionUnknown, LastTransitionTime: v1.Time{}, Message: ""}, {Type: configv1.OperatorFailing, Status: configv1.ConditionTrue, LastTransitionTime: v1.Time{}, Message: ""}})}, {name: "progressing, previously unknown availability", conditions: func() *conditions {
		cs := newConditions(configv1.ClusterOperatorStatus{}, "", v1.Time{})
		cs.setCondition(configv1.OperatorProgressing, configv1.ConditionTrue, "", v1.Unix(0, 0))
		return cs
	}, check: hasConditions([]configv1.ClusterOperatorStatusCondition{{Type: configv1.OperatorProgressing, Status: configv1.ConditionTrue, LastTransitionTime: v1.Unix(0, 0), Message: ""}, {Type: configv1.OperatorAvailable, Status: configv1.ConditionUnknown, LastTransitionTime: v1.Time{}, Message: ""}, {Type: configv1.OperatorFailing, Status: configv1.ConditionUnknown, LastTransitionTime: v1.Time{}, Message: ""}})}, {name: "progressing, previously unavailable", conditions: func() *conditions {
		cs := newConditions(configv1.ClusterOperatorStatus{}, "", v1.Time{})
		cs.setCondition(configv1.OperatorAvailable, configv1.ConditionFalse, "", v1.Time{})
		cs.setCondition(configv1.OperatorFailing, configv1.ConditionFalse, "", v1.Time{})
		cs.setCondition(configv1.OperatorProgressing, configv1.ConditionTrue, "", v1.Unix(0, 0))
		return cs
	}, check: hasConditions([]configv1.ClusterOperatorStatusCondition{{Type: configv1.OperatorProgressing, Status: configv1.ConditionTrue, LastTransitionTime: v1.Unix(0, 0), Message: ""}, {Type: configv1.OperatorAvailable, Status: configv1.ConditionFalse, LastTransitionTime: v1.Time{}, Message: ""}, {Type: configv1.OperatorFailing, Status: configv1.ConditionFalse, LastTransitionTime: v1.Time{}, Message: ""}})}, {name: "not progressing, previously available", conditions: func() *conditions {
		cs := newConditions(configv1.ClusterOperatorStatus{Conditions: []configv1.ClusterOperatorStatusCondition{{Type: configv1.OperatorFailing, Status: configv1.ConditionFalse}, {Type: configv1.OperatorAvailable, Status: configv1.ConditionTrue}, {Type: configv1.OperatorProgressing, Status: configv1.ConditionFalse}}}, "", v1.Time{})
		cs.setCondition(configv1.OperatorProgressing, configv1.ConditionTrue, "", v1.Unix(0, 0))
		return cs
	}, check: hasConditions([]configv1.ClusterOperatorStatusCondition{{Type: configv1.OperatorProgressing, Status: configv1.ConditionFalse, LastTransitionTime: v1.Time{}, Message: ""}, {Type: configv1.OperatorAvailable, Status: configv1.ConditionTrue, LastTransitionTime: v1.Time{}, Message: ""}, {Type: configv1.OperatorFailing, Status: configv1.ConditionFalse, LastTransitionTime: v1.Time{}, Message: ""}})}, {name: "not progressing, previously available, same version", conditions: func() *conditions {
		cs := newConditions(configv1.ClusterOperatorStatus{Versions: []configv1.OperandVersion{{Version: "1.0"}}, Conditions: []configv1.ClusterOperatorStatusCondition{{Type: configv1.OperatorFailing, Status: configv1.ConditionFalse}, {Type: configv1.OperatorAvailable, Status: configv1.ConditionTrue}, {Type: configv1.OperatorProgressing, Status: configv1.ConditionFalse}}}, "1.0", v1.Time{})
		cs.setCondition(configv1.OperatorProgressing, configv1.ConditionTrue, "", v1.Unix(0, 0))
		return cs
	}, check: hasConditions([]configv1.ClusterOperatorStatusCondition{{Type: configv1.OperatorProgressing, Status: configv1.ConditionFalse, LastTransitionTime: v1.Time{}, Message: ""}, {Type: configv1.OperatorAvailable, Status: configv1.ConditionTrue, LastTransitionTime: v1.Time{}, Message: ""}, {Type: configv1.OperatorFailing, Status: configv1.ConditionFalse, LastTransitionTime: v1.Time{}, Message: ""}})}, {name: "progressing, previously available, different version", conditions: func() *conditions {
		cs := newConditions(configv1.ClusterOperatorStatus{Versions: []configv1.OperandVersion{{Version: "1.0"}}, Conditions: []configv1.ClusterOperatorStatusCondition{{Type: configv1.OperatorFailing, Status: configv1.ConditionFalse}, {Type: configv1.OperatorAvailable, Status: configv1.ConditionTrue}, {Type: configv1.OperatorProgressing, Status: configv1.ConditionFalse}}}, "1.1", v1.Time{})
		cs.setCondition(configv1.OperatorProgressing, configv1.ConditionTrue, "", v1.Unix(0, 0))
		return cs
	}, check: hasConditions([]configv1.ClusterOperatorStatusCondition{{Type: configv1.OperatorProgressing, Status: configv1.ConditionTrue, LastTransitionTime: v1.Unix(0, 0), Message: ""}, {Type: configv1.OperatorAvailable, Status: configv1.ConditionTrue, LastTransitionTime: v1.Time{}, Message: ""}, {Type: configv1.OperatorFailing, Status: configv1.ConditionFalse, LastTransitionTime: v1.Time{}, Message: ""}})}, {name: "progressing, previously unavailable, different version", conditions: func() *conditions {
		cs := newConditions(configv1.ClusterOperatorStatus{Versions: []configv1.OperandVersion{{Version: "1.0"}}, Conditions: []configv1.ClusterOperatorStatusCondition{{Type: configv1.OperatorFailing, Status: configv1.ConditionFalse}, {Type: configv1.OperatorAvailable, Status: configv1.ConditionFalse}, {Type: configv1.OperatorProgressing, Status: configv1.ConditionFalse}}}, "1.1", v1.Time{})
		cs.setCondition(configv1.OperatorProgressing, configv1.ConditionTrue, "", v1.Unix(0, 0))
		return cs
	}, check: hasConditions([]configv1.ClusterOperatorStatusCondition{{Type: configv1.OperatorProgressing, Status: configv1.ConditionTrue, LastTransitionTime: v1.Unix(0, 0), Message: ""}, {Type: configv1.OperatorAvailable, Status: configv1.ConditionFalse, LastTransitionTime: v1.Time{}, Message: ""}, {Type: configv1.OperatorFailing, Status: configv1.ConditionFalse, LastTransitionTime: v1.Time{}, Message: ""}})}, {name: "progressing, previously unavailable, same version", conditions: func() *conditions {
		cs := newConditions(configv1.ClusterOperatorStatus{Versions: []configv1.OperandVersion{{Version: "1.0"}}, Conditions: []configv1.ClusterOperatorStatusCondition{{Type: configv1.OperatorFailing, Status: configv1.ConditionFalse}, {Type: configv1.OperatorAvailable, Status: configv1.ConditionFalse}, {Type: configv1.OperatorProgressing, Status: configv1.ConditionFalse}}}, "1.0", v1.Time{})
		cs.setCondition(configv1.OperatorProgressing, configv1.ConditionTrue, "", v1.Unix(0, 0))
		return cs
	}, check: hasConditions([]configv1.ClusterOperatorStatusCondition{{Type: configv1.OperatorProgressing, Status: configv1.ConditionTrue, LastTransitionTime: v1.Unix(0, 0), Message: ""}, {Type: configv1.OperatorAvailable, Status: configv1.ConditionFalse, LastTransitionTime: v1.Time{}, Message: ""}, {Type: configv1.OperatorFailing, Status: configv1.ConditionFalse, LastTransitionTime: v1.Time{}, Message: ""}})}, {name: "change due to message change", conditions: func() *conditions {
		cs := newConditions(configv1.ClusterOperatorStatus{Conditions: []configv1.ClusterOperatorStatusCondition{{Type: configv1.OperatorAvailable, Status: configv1.ConditionTrue, Message: "foo", LastTransitionTime: v1.Time{}}}}, "", v1.Time{})
		cs.setCondition(configv1.OperatorAvailable, configv1.ConditionTrue, "bar", v1.Unix(0, 0))
		return cs
	}, check: hasConditions([]configv1.ClusterOperatorStatusCondition{{Type: configv1.OperatorAvailable, Status: configv1.ConditionTrue, LastTransitionTime: v1.Unix(0, 0), Message: "bar"}, {Type: configv1.OperatorProgressing, Status: configv1.ConditionUnknown, LastTransitionTime: v1.Time{}, Message: ""}, {Type: configv1.OperatorFailing, Status: configv1.ConditionUnknown, LastTransitionTime: v1.Time{}, Message: ""}})}, {name: "change due to status change", conditions: func() *conditions {
		cs := newConditions(configv1.ClusterOperatorStatus{Conditions: []configv1.ClusterOperatorStatusCondition{{Type: configv1.OperatorAvailable, Status: configv1.ConditionTrue, Message: "foo", LastTransitionTime: v1.Time{}}}}, "", v1.Time{})
		cs.setCondition(configv1.OperatorAvailable, configv1.ConditionFalse, "foo", v1.Unix(0, 0))
		return cs
	}, check: hasConditions([]configv1.ClusterOperatorStatusCondition{{Type: configv1.OperatorAvailable, Status: configv1.ConditionFalse, LastTransitionTime: v1.Unix(0, 0), Message: "foo"}, {Type: configv1.OperatorProgressing, Status: configv1.ConditionUnknown, LastTransitionTime: v1.Time{}, Message: ""}, {Type: configv1.OperatorFailing, Status: configv1.ConditionUnknown, LastTransitionTime: v1.Time{}, Message: ""}})}, {name: "no change due to no message/status change", conditions: func() *conditions {
		cs := newConditions(configv1.ClusterOperatorStatus{Conditions: []configv1.ClusterOperatorStatusCondition{{Type: configv1.OperatorAvailable, Status: configv1.ConditionTrue, Message: "foo", LastTransitionTime: v1.Time{}}}}, "", v1.Time{})
		cs.setCondition(configv1.OperatorAvailable, configv1.ConditionTrue, "foo", v1.Unix(0, 0))
		return cs
	}, check: hasConditions([]configv1.ClusterOperatorStatusCondition{{Type: configv1.OperatorAvailable, Status: configv1.ConditionTrue, LastTransitionTime: v1.Time{}, Message: "foo"}, {Type: configv1.OperatorProgressing, Status: configv1.ConditionUnknown, LastTransitionTime: v1.Time{}, Message: ""}, {Type: configv1.OperatorFailing, Status: configv1.ConditionUnknown, LastTransitionTime: v1.Time{}, Message: ""}})}} {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.check(tc.conditions()); err != nil {
				t.Error(err)
			}
		})
	}
}
