package certParser

import (
	"crypto/tls"
	"net"
	"time"
	"fmt"
)

const defaultTimeoutSeconds = 3
const insecureSkipVerify = false

// Only TLS 1.3 Cipers for security reasons
// To add older TLS Ciphers for compatibility reasons consult https://wiki.mozilla.org/Security/Server_Side_TLS
var cipherSuites = []uint16{
	tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
	tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
	tls.TLS_CHACHA20_POLY1305_SHA256,
}

func Get(site, port string) {
	d := &net.Dialer{
			Timeout: time.Duration(defaultTimeoutSeconds) * time.Second,
	}

	conn, _ := tls.DialWithDialer(d, "tcp", site + ":" + port, &tls.Config{
		InsecureSkipVerify: insecureSkipVerify,
		CipherSuites:       cipherSuites,
		MaxVersion:         tls.VersionTLS13,
	})
	defer conn.Close()

	certChain := conn.ConnectionState().PeerCertificates
	cert := certChain[0]

	fmt.Println("Subject: " + cert.Subject.CommonName)
	fmt.Println("Issuer: " + cert.Issuer.CommonName)
	fmt.Println("SANs:")
	fmt.Println(cert.DNSNames)
	fmt.Println("Validity: " + cert.NotAfter.In(time.Local).String())
}
