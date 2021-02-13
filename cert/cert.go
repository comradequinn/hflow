package cert

import (
	"bytes"
	"comradequinn/hflow/log"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"
	"net"
	"sync"
	"time"
)

type (
	dynamicCertQuery struct {
		subject string
		result  chan<- *dynamicCert
	}
	dynamicCert struct {
		val *tls.Certificate
		mx  sync.Mutex
	}
)

var (
	caCertificate      *x509.Certificate
	caPrivateKey       *rsa.PrivateKey
	hostIP             net.IP
	dynamicCertQueries = make(chan dynamicCertQuery, 1)
)

func init() {
	log.Printf(1, "initialising dynamic certificate generation")

	conn, err := net.Dial("udp", "192.0.0.0:9999") // Doesn't matter that nothing is listening, we just need to get the IP it selects to connect with

	if err != nil {
		log.Fatalf(0, "unable to ascertain host ip [%v]", err)
	}

	defer conn.Close()

	hostIP = conn.LocalAddr().(*net.UDPAddr).IP
	caCertificate, caPrivateKey = parseCA()

	dynamicCerts := make(map[string]*dynamicCert, 100) // set this to an arbitary, but high default value to avoid the map repeatedly resizing under burst when the proxy starts

	go func() {
		log.Printf(1, "started dynamic certificate worker")

		for {
			dCertRq := <-dynamicCertQueries
			dCert, exists := dynamicCerts[dCertRq.subject]

			if !exists {
				dCert = &dynamicCert{mx: sync.Mutex{}}
				dynamicCerts[dCertRq.subject] = dCert

				log.Printf(2, "added uninitialised dynamic certificate record for [%v]", dCertRq.subject)
			}

			dCertRq.result <- dCert
		}
	}()
}

// For returns a dynamically generated certificate, signed by the HFLOW CA, that matches the domain or ip that was requested
func For(chi *tls.ClientHelloInfo) (*tls.Certificate, error) {
	subject, dCertQryResult := chi.ServerName, make(chan *dynamicCert, 1)

	if subject == "" {
		log.Printf(3, "using ip as subject of dynamic certificate due to absence of server name in client hello")
		subject = hostIP.String()
	}

	dynamicCertQueries <- dynamicCertQuery{
		subject: subject,
		result:  dCertQryResult,
	}

	dCert := <-dCertQryResult

	dCert.mx.Lock()
	defer dCert.mx.Unlock()

	if dCert.val == nil {
		var err error

		if dCert.val, err = createCert(subject, hostIP, caCertificate, caPrivateKey); err != nil {
			dCert.val = nil
			return nil, fmt.Errorf("unable to create dynamic certificate for [%v]: [%v]", subject, err)
		}

		log.Printf(2, "initialised dynamic certificate record for [%v] with new certificate ", subject)
	}

	log.Printf(1, "returned dynamic certificate for [%v]", subject)

	return dCert.val, nil
}

// WriteCA writes the HFLOW CA X509 certificate in PEM format to the specified io.Writer
func WriteCA(w io.Writer) error {
	_, err := w.Write(hflowCACertPEM)
	return err
}

// parseCA extracts the X509 certificate and private key of the HFLOW CA from local PEM files
func parseCA() (*x509.Certificate, *rsa.PrivateKey) {
	readPEM := func(bytes []byte, pemType string) []byte {
		pemBlock, _ := pem.Decode(bytes)

		if pemBlock == nil {
			log.Fatalf(0, "expected pem [%v] data but data is not in pem format", pemType)
		}

		if pemBlock.Type != pemType {
			log.Fatalf(0, "expected pem data for type [%v] but found [%v]", pemType, pemBlock.Type)
		}

		return pemBlock.Bytes
	}

	var err error
	var privateKey *rsa.PrivateKey

	if privateKey, err = x509.ParsePKCS1PrivateKey(readPEM(hflowCAPrivateKeyPEM, "RSA PRIVATE KEY")); err != nil {
		log.Fatalf(0, "unable to parse hflow ca private key pem as pkcs1 private key [%v]. ", err)
	}

	var pkixPublicKey any

	if pkixPublicKey, err = x509.ParsePKIXPublicKey(readPEM(hflowCAPublicKeyPEM, "PUBLIC KEY")); err != nil {
		log.Fatalf(0, "unable to parse hflow ca public key pem as pkix public key [%v]", err)
	}

	var rsaPublicKey *rsa.PublicKey
	var ok bool

	if rsaPublicKey, ok = pkixPublicKey.(*rsa.PublicKey); !ok {
		log.Fatalf(0, "unable to parse hflow ca public key pkix data as rsa public key [%v]", err)
	}

	privateKey.PublicKey = *rsaPublicKey

	log.Printf(3, "loaded hflow ca key public and private keys")

	var certificate *x509.Certificate

	if certificate, err = x509.ParseCertificate(readPEM(hflowCACertPEM, "CERTIFICATE")); err != nil {
		log.Fatalf(0, "unable to parse hflow ca cert pem data as x509 certificate [%v]", err)
	}

	log.Printf(3, "loaded hflow ca certificate")

	return certificate, privateKey
}

func createCert(subject string, ip net.IP, caCertificate *x509.Certificate, caPrivateKey *rsa.PrivateKey) (*tls.Certificate, error) {
	log.Printf(3, "generating dynamic certificate for subject [%v]", subject)

	now := time.Now()

	template := x509.Certificate{
		Subject: pkix.Name{
			Organization:  []string{"HFLOW Dynamic Cert"},
			StreetAddress: []string{"1 Virtual Avenue"},
			PostalCode:    []string{"12345"},
			Province:      []string{"Ether"},
			Locality:      []string{"Net"},
			Country:       []string{"UK"},
			CommonName:    subject,
		},
		NotBefore:             now,
		NotAfter:              now.Add(87658 * time.Hour),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCRLSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		IsCA:                  false,
		BasicConstraintsValid: true,
		DNSNames:              []string{subject},
		IPAddresses:           []net.IP{ip},
		SerialNumber:          big.NewInt(1),
	}

	var (
		err error
		pk  *rsa.PrivateKey
	)

	if pk, err = rsa.GenerateKey(rand.Reader, 4096); err != nil {
		return nil, fmt.Errorf("failed to generate private key for dynamic certificate with subject [%v]. [%v]", subject, err)
	}

	log.Printf(3, "generated private key for dynamic certificate with subject [%v]", subject)

	var der []byte

	if der, err = x509.CreateCertificate(rand.Reader, &template, caCertificate, &pk.PublicKey, caPrivateKey); err != nil {
		return nil, fmt.Errorf("failed to generate der encoded dynamic certificate with subject [%v]. [%v]", subject, err)
	}

	certPEM, keyPEM := bytes.Buffer{}, bytes.Buffer{}

	if err = pem.Encode(&certPEM, &pem.Block{Type: "CERTIFICATE", Bytes: der}); err != nil {
		return nil, fmt.Errorf("failed to transcode dynamic certificate with subject [%v] from der to pem [%v]", subject, err)
	}

	if err = pem.Encode(&keyPEM, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(pk)}); err != nil {
		return nil, fmt.Errorf("failed to encode private key for dynamic certificate with subject [%v] to pem [%v]", subject, err)
	}

	var cert tls.Certificate

	if cert, err = tls.X509KeyPair(certPEM.Bytes(), keyPEM.Bytes()); err != nil {
		return nil, fmt.Errorf("failed to generate dynamic certificate with subject [%v] from pem [%v]", subject, err)
	}

	log.Printf(3, "generated dynamic certificate with subject [%v]", subject)

	return &cert, nil
}
