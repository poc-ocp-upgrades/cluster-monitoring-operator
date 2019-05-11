package tasks

import (
	"github.com/golang/glog"
	"github.com/openshift/cluster-monitoring-operator/pkg/client"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

type TaskRunner struct {
	client	*client.Client
	tasks	[]*TaskSpec
}

func NewTaskRunner(client *client.Client, tasks []*TaskSpec) *TaskRunner {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &TaskRunner{client: client, tasks: tasks}
}
func (tl *TaskRunner) RunAll() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var g errgroup.Group
	for i, ts := range tl.tasks {
		ts := ts
		i := i
		g.Go(func() error {
			glog.V(4).Infof("running task %d of %d: %v", i+1, len(tl.tasks), ts.Name)
			err := tl.ExecuteTask(ts)
			glog.V(4).Infof("ran task %d of %d: %v", i+1, len(tl.tasks), ts.Name)
			return errors.Wrapf(err, "running task %v failed", ts.Name)
		})
	}
	return g.Wait()
}
func (tl *TaskRunner) ExecuteTask(ts *TaskSpec) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return ts.Task.Run()
}
func NewTaskSpec(name string, task Task) *TaskSpec {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &TaskSpec{Name: name, Task: task}
}

type TaskSpec struct {
	Name	string
	Task	Task
}
type Task interface{ Run() error }
