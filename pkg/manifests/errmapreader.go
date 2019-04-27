package manifests

import (
	"encoding/json"
	"fmt"
)

type errMapReader struct {
	src	map[string]string
	err	error
}

func newErrMapReader(src map[string]string) *errMapReader {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &errMapReader{src: src}
}
func (r *errMapReader) value(key string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if r.err != nil {
		return ""
	}
	result, ok := r.src[key]
	if !ok {
		r.err = fmt.Errorf("key %s is missing", key)
		return ""
	}
	return result
}
func (r *errMapReader) slice(key string) []string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if r.err != nil {
		return nil
	}
	v := r.value(key)
	if r.err != nil {
		return nil
	}
	if len(v) == 0 {
		return nil
	}
	var ret []string
	if err := json.Unmarshal([]byte(v), &ret); err != nil {
		r.err = err
		return nil
	}
	return ret
}
func (r *errMapReader) Error() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return r.err
}
