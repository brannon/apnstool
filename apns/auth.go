// Copyright 2019 Brannon Jones. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package apns

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"os"
	"time"

	"github.com/pascaldekloe/jwt"
)

func LoadKeyFromFile(filePath string) (*ecdsa.PrivateKey, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	if block == nil || block.Type != "PRIVATE KEY" {
		return nil, errors.New("failed to decode PEM block containing private key")
	}

	pkcs8Key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	var ecdsaKey *ecdsa.PrivateKey

	ecdsaKey, ok := pkcs8Key.(*ecdsa.PrivateKey)
	if !ok {
		return nil, errors.New("failed to parse key as ECDSA private key")
	}

	return ecdsaKey, nil
}

func GenerateJWTFromKey(key *ecdsa.PrivateKey, keyId string, teamId string, issuedAt time.Time, expiresAfter time.Duration) (string, error) {
	var claims jwt.Claims

	claims.KeyID = keyId
	claims.Issuer = teamId
	claims.Issued = jwt.NewNumericTime(issuedAt)
	claims.Expires = jwt.NewNumericTime(issuedAt.Add(expiresAfter))

	token, err := claims.ECDSASign(jwt.ES256, key)

	return string(token), err
}

func GenerateJWTFromKeyFile(keyFile string, keyId string, teamId string, issuedAt time.Time, expiresAfter time.Duration) (string, error) {
	key, err := LoadKeyFromFile(keyFile)
	if err != nil {
		return "", err
	}

	return GenerateJWTFromKey(key, keyId, teamId, issuedAt, expiresAfter)
}
