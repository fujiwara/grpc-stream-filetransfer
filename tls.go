package grpcp

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"time"
)

func genTLS(certFile string, keyFile string) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load key pair: %w", err)
	}
	return &tls.Config{Certificates: []tls.Certificate{cert}}, nil
}

func genSelfSignedTLS() (*tls.Config, error) {
	dir, err := os.MkdirTemp("", "grpcp")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer func() {
		os.RemoveAll(dir)
	}()

	// プライベートキーを生成
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key: %w", err)
	}

	// 証明書のテンプレートを作成
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"My Organization"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	// 証明書を自己署名で作成
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return nil, fmt.Errorf("failed to create certificate: %w", err)
	}

	// PEM形式でプライベートキーを書き出し
	keyOut, err := os.Create(filepath.Join(dir, "server.key"))
	if err != nil {
		return nil, fmt.Errorf("failed to open server.key for writing: %w", err)
	}
	keyBytes, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal private: %w", err)
	}
	pem.Encode(keyOut, &pem.Block{Type: "EC PRIVATE KEY", Bytes: keyBytes})
	keyOut.Close()

	// PEM形式で証明書を書き出し
	certOut, err := os.Create(filepath.Join(dir, "server.crt"))
	if err != nil {
		return nil, fmt.Errorf("failed to open server.crt for writing: %w", err)
	}
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	certOut.Close()

	// TLS設定を使用してサーバーを起動する場合の例
	cert, err := tls.LoadX509KeyPair(
		filepath.Join(dir, "server.crt"),
		filepath.Join(dir, "server.key"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load key pair: %w", err)
	}
	return &tls.Config{Certificates: []tls.Certificate{cert}}, nil
}
