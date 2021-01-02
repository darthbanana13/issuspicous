package certParser

import (
	"crypto/tls"
	"net"
	"time"
	"fmt"
)

const defaultTimeoutSeconds = 3
const insecureSkipVerify = false

func Get(site, port string) {
	d := &net.Dialer{
			Timeout: time.Duration(defaultTimeoutSeconds) * time.Second,
	}

	conn, _ := tls.DialWithDialer(d, "tcp", site + ":" + port, &tls.Config{
		InsecureSkipVerify: insecureSkipVerify,
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
