package certificates

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"time"

	"go.uber.org/zap"
)

func GenerateSelfSignedCerts(logger zap.Logger) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		logger.Error("failed to generate private key [%w]", zap.Error(err))
	}

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		logger.Error("failed to generate serial number: [%w]", zap.Error(err))
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName:         "POC",
			OrganizationalUnit: []string{"FIBER"},
			Country:            []string{"FR"},
			Organization:       []string{"GO"},
		},
		DNSNames:              []string{"localhost"},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
	}

	// Create self-signed certificate.
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		logger.Error("failed to create certificate: [%w]", zap.Error(err))
	}

	pemCert := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	if pemCert == nil {
		logger.Error("failed to encode certificate to PEM: [%w]", zap.Error(err))
	}
	if err := os.WriteFile("cert.pem", pemCert, 0644); err != nil {
		logger.Error("failed to write certificate: [%w]", zap.Error(err))
	}
	logger.Info("wrote cert.pem")

	privBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		logger.Error("unable to marshal private key: [%w]", zap.Error(err))
	}
	pemKey := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privBytes})
	if pemKey == nil {
		logger.Error("failed to encode pem to memory")
	}
	if err := os.WriteFile("key.pem", pemKey, 0600); err != nil {
		logger.Error("unable to write key.pem: [%w]", zap.Error(err))
	}
	logger.Info("wrote key.pem")
}
