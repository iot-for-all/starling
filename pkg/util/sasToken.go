package util

import (
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// CreateSasToken creates a shared access signature token.
func CreateSasToken(
	key string,
	resource string,
	keyName string,
	expiresAfter time.Duration) (string, error) {

	sr := url.QueryEscape(resource)
	skn := url.QueryEscape(keyName)
	se := strconv.FormatInt(time.Now().Add(expiresAfter).Unix(), 10)
	sig, err := ComputeHmac(key, sr+"\n"+se)
	if err != nil {
		return "", err
	}

	sig = url.QueryEscape(sig)
	return fmt.Sprintf("SharedAccessSignature sr=%s&sig=%s&se=%s&skn=%s", sr, sig, se, skn), nil
}
