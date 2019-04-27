package manifests

import (
	"testing"
)

func TestImageParsing(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	imageCases := []struct {
		str	string
		image	image
	}{{str: "quay.io/test/image:tag", image: image{repo: "quay.io/test/image", tag: "tag"}}, {str: "image:tag", image: image{repo: "image", tag: "tag"}}, {str: "quay.io:443/test/image:tag", image: image{repo: "quay.io:443/test/image", tag: "tag"}}}
	for _, imageCase := range imageCases {
		image, err := imageFromString(imageCase.str)
		if err != nil {
			t.Errorf("error parsing image string %s : %v", imageCase.str, err)
			continue
		}
		if imageCase.image != *image {
			t.Errorf("parsed image %+v does not match expected image %+v", *image, imageCase.image)
			continue
		}
		if imageCase.str != image.String() {
			t.Errorf("parsed image string %s does not match expected image string %s", image.String(), imageCase.str)
		}
	}
}
