package addressParser

import (
	"fmt"
	"strings"
	"net"
	"net/url"
	"sync"
)

const defaultPort = "443"
const defaultProtocol = "https"

var wg sync.WaitGroup

func Parse(addresses []string) {
	for _, address := range addresses {
		wg.Add(1)
		go parseSingle(address)
	}
	wg.Wait()
}

func parseSingle(address string) {
	var site string
	var protocol string
	if strings.HasPrefix(address, "http") {
		u, _ := url.Parse(address)
		site = u.Host
		protocol = u.Scheme
	} else {
		site = address
		protocol = defaultProtocol
	}

	if !strings.Contains(site, ":") {
		fmt.Println("Protocol:" + protocol + " Site:" + site + " Port:" + defaultPort)
		wg.Done()
		return
	}
	site, port, err := net.SplitHostPort(site)
	if err == nil {
		fmt.Println("Protocol:" + protocol + " Site:" + site + " Port:" + port)
		wg.Done()
		return
	}
	fmt.Println(err)
	wg.Done()
}
