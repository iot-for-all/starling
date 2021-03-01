package util

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

// ComputeHmac computes a SHA256 hmac of the data based on the provided secret.
func ComputeHmac(secret string, data string) (string, error) {
	s, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		return "", err
	}

	h := hmac.New(sha256.New, s)
	h.Write([]byte(data))
	dest := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return dest, nil
}
