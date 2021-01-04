package addrParser

import (
	"fmt"
	"strings"
	"net"
	"net/url"
	// "sync"
	"strconv"
)

const (
	HTTP	= iota
	HTTPS
)

var Proto = map[string]uint8{
	"http":		HTTP,
	"https":	HTTPS,
}

var ProtoPort = map[uint8]uint16{
	HTTP:	80,
	HTTPS:	443,
}

type Site struct {
	Protocol	uint8
	Host		string
	Port	    uint16
}

const defaultProto = HTTPS

func NewSites(addrs []string) {
	// var wg sync.WaitGroup
	var done = make(chan Site)
	defer close(done)
	for _, addr := range addrs {
		// wg.Add(1)
		go newSiteConcur(addr, done)
		// wg.Done()
	}
	fmt.Println(<-done)
	fmt.Println(<-done)
	fmt.Println(<-done)
}

func newSiteConcur(addr string, done chan<- Site) {
	Site, err := NewSite(addr)
	if err != nil {
		// TODO: Return errors on a separate channel
		fmt.Println(err)
		return
	}
	done<- Site
}

func NewSite(addr string) (Site, error) {
	site, err := parseAddrWithProto(addr)

	if !strings.Contains(site.Host, ":") {
		return	site, nil
	}
	var port string
	site.Host, port, err = net.SplitHostPort(site.Host)
	if err == nil {
		//TODO: Manage port out of range error
		uintport, _ := strconv.ParseUint(port, 10, 16)
		site.Port = uint16(uintport)
		return site, nil
	}
	return Site{}, nil
}

func parseAddrWithProto(addr string) (site Site, err error) {
	if strings.HasPrefix(addr, "http") {
		u, err := url.Parse(addr)
		if err != nil {
			return Site{}, err
		}
		return	Site{
					Protocol: Proto[u.Scheme],
					Host: u.Host,
					Port: ProtoPort[Proto[u.Scheme]],
				},
				nil
	}
	return	Site{
				Protocol: defaultProto,
				Host: addr,
				Port: ProtoPort[defaultProto],
			},
			nil
}

func (s Site) Addr() string {
	return s.Host + ":" + strconv.Itoa(int(s.Port))
}
