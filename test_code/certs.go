package main

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

func testCert(caCertContent string) error {
	block, _ := pem.Decode([]byte(caCertContent))
	if block == nil {
		return fmt.Errorf("failed to parse certificate PEM")
	}
	_, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse certificate: %v", err)
	}

	return nil
}
