package encryption

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
)

func ParsePublic(f string) (*rsa.PublicKey, error) {

	file, err := os.ReadFile(f)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("[ERROR] Unable to read certificate: %s", err.Error()))
	}

	block, _ := pem.Decode(file)
	if block == nil {
		return nil, errors.New(fmt.Sprintf("[ERROR] Unable to decode certificate: %s", err.Error()))
	}

	var parsed interface{}
	if parsed, err = x509.ParsePKIXPublicKey(block.Bytes); err != nil {
		return nil, errors.New(fmt.Sprintf("[ERROR] Unable to parse RSA public key: %s", err.Error()))
	}

	var public *rsa.PublicKey

	ok := false
	if public, ok = parsed.(*rsa.PublicKey); !ok {
		return nil, errors.New(fmt.Sprintf("[ERROR] Unable to parse RSA public key: %s", err.Error()))
	}

	return public, nil
}

func ParsePrivate(f string) (*rsa.PrivateKey, error) {

	file, err := os.ReadFile(f)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("[ERROR] Unable to read certificate: %s", err.Error()))
	}

	block, _ := pem.Decode(file)
	if block == nil {
		return nil, errors.New(fmt.Sprintf("[ERROR] Unable to decode certificate: %s", err.Error()))
	}

	var parsed interface{}
	if parsed, err = x509.ParsePKCS8PrivateKey(block.Bytes); err != nil {
		return nil, errors.New(fmt.Sprintf("[ERROR] Unable to parse RSA private key: %s", err.Error()))
	}

	var private *rsa.PrivateKey

	ok := false
	if private, ok = parsed.(*rsa.PrivateKey); !ok {
		return nil, errors.New(fmt.Sprintf("[ERROR] Unable to parse RSA private key: %s", err.Error()))
	}

	return private, nil
}
