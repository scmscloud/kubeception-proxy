package encryption

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"log"
)

func Sign(private *rsa.PrivateKey, endpoint string) (*[]byte, error) {

	hashed := sha256.Sum256([]byte(endpoint))
	signature, err := rsa.SignPKCS1v15(rand.Reader, private, crypto.SHA256, hashed[:])
	if err != nil {
		log.Fatal(err)
	}

	return &signature, nil
}

func Verify(public *rsa.PublicKey, endpoint string, signature string) (bool, error) {

	digest := sha256.Sum256([]byte(endpoint))

	sign, err := hex.DecodeString(signature)
	if err != nil {
		return false, err
	}

	if err := rsa.VerifyPKCS1v15(public, crypto.SHA256, digest[:], sign); err != nil {
		return false, err
	}

	return true, nil
}
