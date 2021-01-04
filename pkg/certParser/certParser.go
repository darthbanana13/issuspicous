package certParser

import (
	"crypto/tls"
	"net"
	"time"
	"fmt"
	"bytes"
	"text/template"
	"sync"
	"github.com/darthrevan13/issuspicous/pkg/addrParser"
)

type Cert struct {
	Addr		addrParser.Site
	Subject		string
	Issuer		string
	SANs		[]string
	Validity	time.Time

}

const certInfoTempl = `Subject:		{{.Subject}}
Issuer:			{{.Issuer}}
SANs:			{{.SANs}}
Validity:		{{.Validity}}`

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

func NewCert(addr string) (Cert, error) {
	// TODO: Handle errors
	site, _ := addrParser.NewSite(addr)
	return NewCertFromSite(site)
}

func NewCertFromSite(site addrParser.Site) (Cert, error) {
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

	return	Cert{
				Addr:		site,
				Subject:	cert.Subject.CommonName,
				Issuer:		cert.Issuer.CommonName,
				SANs:		cert.DNSNames,
				Validity:	cert.NotAfter.In(time.Local),
			},
			nil
}

// TODO: Refactor common approach to concurency between addrParser and certParser
func NewCerts(addrs []string) <-chan Cert {
	sitesChan := addrParser.NewSites(addrs)
	var done = make(chan Cert)
	var wg sync.WaitGroup
	for site := range sitesChan {
		wg.Add(1)
		go newCertFromSiteConcur(site, done, &wg)
	}
	go func() {
		wg.Wait()
		close(done)
	}()
	return done
}

func newCertFromSiteConcur(site addrParser.Site, done chan<- Cert, wg *sync.WaitGroup) {
	// TODO: Send errors to separate channel
	cert, err := NewCertFromSite(site)
	if err != nil {
		fmt.Println(err)
		wg.Done()
		return
	}
	done<- cert
	wg.Done()
}

func (c Cert) CertificateInfo() string {
	var b bytes.Buffer

	t := template.Must(template.New("default").Parse(certInfoTempl))
	if err := t.Execute(&b, c); err != nil {
		panic(err)
	}
	return b.String()
}
