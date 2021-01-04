package certParser

import (
	"crypto/tls"
	"net"
	"time"
	// "fmt"
	"bytes"
	"text/template"
	"github.com/darthrevan13/issuspicous/pkg/addrParser"
)

type Certificate struct {
	Addr		addrParser.Site
	Subject		string
	Issuer		string
	SANs		[]string
	Validity	time.Time

}

const certInfoTempl = `Subject:	{{.Subject}}
Issuer:		{{.Issuer}}
SANs:		{{.SANs}}
Validity:	{{.Validity}}`

const defaultTimeoutSeconds = 3
const insecureSkipVerify = false

// To add older TLS Ciphers for compatibility reasons consult https://wiki.mozilla.org/Security/Server_Side_TLS
var cipherSuites = []uint16{
	//TLS 1.3
	tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
	tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
	tls.TLS_CHACHA20_POLY1305_SHA256,

	// TLS 1.2
	tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
	tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
	tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
	tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
	tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
	tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
	tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
	tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
}

func NewCertificate(addr string) (Certificate, error) {
	// TODO: Handle errors
	site, _ := addrParser.NewSite(addr)
	d := &net.Dialer{
			Timeout: time.Duration(defaultTimeoutSeconds) * time.Second,
	}

	// TODO: Handle errors
	conn, _ := tls.DialWithDialer(d, "tcp", site.Addr(), &tls.Config{
		InsecureSkipVerify: insecureSkipVerify,
		CipherSuites:       cipherSuites,
		MaxVersion:         tls.VersionTLS13,
	})
	defer conn.Close()

	certChain := conn.ConnectionState().PeerCertificates
	cert := certChain[0]

	return	Certificate{
				Addr:		site,
				Subject:	cert.Subject.CommonName,
				Issuer:		cert.Issuer.CommonName,
				SANs:		cert.DNSNames,
				Validity:	cert.NotAfter.In(time.Local),
			},
			nil
}

func (c Certificate) CertificateInfo() string {
	var b bytes.Buffer

	t := template.Must(template.New("default").Parse(certInfoTempl))
	if err := t.Execute(&b, c); err != nil {
		panic(err)
	}
	return b.String()
}
