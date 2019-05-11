package manifests

import (
	"crypto/rand"
	"encoding/base64"
)

func GeneratePassword(n int) (string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), err
}
