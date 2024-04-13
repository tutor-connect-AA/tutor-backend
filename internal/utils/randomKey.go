package utils

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
)

func GenerateKey() (*ecdsa.PrivateKey, error) {
	var key, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	if err != nil {
		return &ecdsa.PrivateKey{}, err
	}
	return key, nil
}
