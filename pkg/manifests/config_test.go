package manifests

import (
	"errors"
	"fmt"
	"os"
	"testing"
	configv1 "github.com/openshift/api/config/v1"
)

func TestConfigParsing(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	f, err := os.Open("../../examples/config/config.yaml")
	if err != nil {
		t.Fatal(err)
	}
	_, err = NewConfig(f)
	if err != nil {
		t.Fatal(err)
	}
}
func TestEmptyConfigIsValid(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_, err := NewConfigFromString("")
	if err != nil {
		t.Fatal(err)
	}
}
func TestTelemeterClientConfig(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	truev, falsev := true, false
	tcs := []struct {
		enabled	bool
		cfg	*TelemeterClientConfig
	}{{cfg: nil, enabled: false}, {cfg: &TelemeterClientConfig{}, enabled: false}, {cfg: &TelemeterClientConfig{Enabled: &truev}, enabled: false}, {cfg: &TelemeterClientConfig{Enabled: &falsev}, enabled: false}, {cfg: &TelemeterClientConfig{ClusterID: "test"}, enabled: false}, {cfg: &TelemeterClientConfig{ClusterID: "test", Enabled: &falsev}, enabled: false}, {cfg: &TelemeterClientConfig{ClusterID: "test", Enabled: &truev}, enabled: false}, {cfg: &TelemeterClientConfig{Token: "test"}, enabled: false}, {cfg: &TelemeterClientConfig{Token: "test", Enabled: &falsev}, enabled: false}, {cfg: &TelemeterClientConfig{Token: "test", Enabled: &truev}, enabled: false}, {cfg: &TelemeterClientConfig{ClusterID: "test", Token: "test"}, enabled: true}, {cfg: &TelemeterClientConfig{ClusterID: "test", Token: "test", Enabled: &truev}, enabled: true}, {cfg: &TelemeterClientConfig{ClusterID: "test", Token: "test", Enabled: &falsev}, enabled: false}}
	for i, tc := range tcs {
		if got := tc.cfg.IsEnabled(); got != tc.enabled {
			t.Errorf("testcase %d: expected enabled %t, got %t", i, tc.enabled, got)
		}
	}
}
func TestEtcdDefaultsToDisabled(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c, err := NewConfigFromString("")
	if err != nil {
		t.Fatal(err)
	}
	if c.EtcdConfig.IsEnabled() {
		t.Error("an empty configuration should have etcd disabled")
	}
	c, err = NewConfigFromString(`{"etcd":{}}`)
	if err != nil {
		t.Fatal(err)
	}
	if c.EtcdConfig.IsEnabled() {
		t.Error("an empty etcd configuration should have etcd disabled")
	}
}

type configCheckFunc func(*Config, error) error

func configChecks(fs ...configCheckFunc) configCheckFunc {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return configCheckFunc(func(c *Config, err error) error {
		for _, f := range fs {
			if e := f(c, err); e != nil {
				return e
			}
		}
		return nil
	})
}
func hasError(expected bool) configCheckFunc {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return configCheckFunc(func(_ *Config, err error) error {
		if got := err != nil; got != expected {
			return fmt.Errorf("expected error %t, got %t", expected, got)
		}
		return nil
	})
}
func TestLoadProxy(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	hasHTTPProxy := func(expected string) configCheckFunc {
		return configCheckFunc(func(c *Config, _ error) error {
			if got := c.HTTPConfig.HTTPProxy; got != expected {
				return fmt.Errorf("want http proxy %v, got %v", expected, got)
			}
			return nil
		})
	}
	hasHTTPSProxy := func(expected string) configCheckFunc {
		return configCheckFunc(func(c *Config, _ error) error {
			if got := c.HTTPConfig.HTTPSProxy; got != expected {
				return fmt.Errorf("want https proxy %v, got %v", expected, got)
			}
			return nil
		})
	}
	hasNoProxy := func(expected string) configCheckFunc {
		return configCheckFunc(func(c *Config, _ error) error {
			if got := c.HTTPConfig.NoProxy; got != expected {
				return fmt.Errorf("want noproxy %v, got %v", expected, got)
			}
			return nil
		})
	}
	for _, tc := range []struct {
		name	string
		load	func() (*configv1.Proxy, error)
		check	configCheckFunc
	}{{name: "error loading proxy", load: func() (*configv1.Proxy, error) {
		return nil, errors.New("failure")
	}, check: configChecks(hasHTTPProxy(""), hasHTTPSProxy(""), hasNoProxy(""), hasError(true))}, {name: "empty spec", load: func() (*configv1.Proxy, error) {
		return &configv1.Proxy{}, nil
	}, check: configChecks(hasHTTPProxy(""), hasHTTPSProxy(""), hasNoProxy(""), hasError(false))}, {name: "proxies", load: func() (*configv1.Proxy, error) {
		return &configv1.Proxy{Spec: configv1.ProxySpec{HTTPProxy: "http://proxy", HTTPSProxy: "https://proxy", NoProxy: "localhost,svc.cluster"}}, nil
	}, check: configChecks(hasHTTPProxy("http://proxy"), hasHTTPSProxy("https://proxy"), hasNoProxy("localhost,svc.cluster"), hasError(false))}} {
		t.Run(tc.name, func(t *testing.T) {
			c := NewDefaultConfig()
			err := c.LoadProxy(tc.load)
			if err := tc.check(c, err); err != nil {
				t.Error(err)
			}
		})
	}
}
