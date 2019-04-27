package manifests

import (
	"errors"
	"strings"
)

type image struct {
	repo	string
	tag	string
}

var InvalidImage = errors.New("image string invalid")

func imageFromString(s string) (*image, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	parts := strings.Split(s, ":")
	if len(parts) < 2 {
		return nil, InvalidImage
	}
	return &image{repo: strings.Join(parts[0:len(parts)-1], ":"), tag: parts[len(parts)-1]}, nil
}
func (i *image) String() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return i.repo + ":" + i.tag
}
func (i *image) SetTagIfNotEmpty(tag string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if tag != "" {
		i.tag = tag
	}
}
