package framework

import (
	"testing"
)

func TestGetFirstValueFromPromQuery(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	tests := []struct {
		Name	string
		F	func(t *testing.T)
	}{{Name: "should fail on multiple timeseries", F: func(t *testing.T) {
		body := `
{"status":"success","data":{"resultType":"vector","result":[{"metric":{"__name__":"ALERTS","alertname":"TargetDown","alertstate":"firing","job":"metrics","severity":"warning"},"value":[1551102571.196,"1"]},{"metric":{"__name__":"ALERTS","alertname":"Watchdog","alertstate":"firing","severity":"none"},"value":[1551102571.196,"1"]}]}}
`
		_, err := GetFirstValueFromPromQuery([]byte(body))
		if err == nil || err.Error() != "expected body to contain single timeseries but got 2" {
			t.Fatalf("expected GetFirstValueFromPromQuery to fail on multiple timeseries but got err %q instead", err)
		}
	}}, {Name: "should return first value", F: func(t *testing.T) {
		body := `
{"status":"success","data":{"resultType":"vector","result":[{"metric":{"__name__":"ALERTS","alertname":"Watchdog","alertstate":"firing","severity":"none"},"value":[1551102571.196,"1"]}]}}
`
		v, err := GetFirstValueFromPromQuery([]byte(body))
		if err != nil {
			t.Fatal(err)
		}
		if v != 1 {
			t.Fatalf("expected query to return %v but got %v", 1, v)
		}
	}}}
	for _, test := range tests {
		t.Run(test.Name, test.F)
	}
}
