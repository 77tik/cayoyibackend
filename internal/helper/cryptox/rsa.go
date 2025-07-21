package cryptox

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
)

func GenRSAKey(bits int) (private *rsa.PrivateKey, publicDERBase64 string, err error) {
	// 生成一个长为bits的RSA密钥
	private, err = rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, "", err
	}

	public := private.Public()
	prvate, err := x509.MarshalPKIXPublicKey(public)
	if err != nil {
		return nil, "", err
	}

	return private, base64.StdEncoding.EncodeToString(prvate), nil
}
