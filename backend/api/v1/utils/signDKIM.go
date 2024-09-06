package utils

import (
	"bytes"
	"fmt"
	"github.com/toorop/go-dkim"
)

func SignEmail(email *[]byte, domain, selector, privateKey string) ([]byte, error) {
	options := dkim.NewSigOptions()
	options.PrivateKey = []byte(privateKey)
	options.Domain = domain
	options.Selector = selector
	options.SignatureExpireIn = 3600
	options.BodyLength = 50
	options.Headers = []string{"from", "to", "subject"}
	options.AddSignatureTimestamp = true
	options.Canonicalization = "relaxed/relaxed"

	var buffer bytes.Buffer
	err := dkim.Sign(email, options)
	if err != nil {
		return nil, fmt.Errorf("error signing email: %w", err)
	}

	return buffer.Bytes(), nil
}
