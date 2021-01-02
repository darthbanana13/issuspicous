package addrParser

import (
	"fmt"
	"strings"
	"net"
	"net/url"
	"sync"
)

const defaultPort = "443"
const defaultProto = "https"

var wg sync.WaitGroup

func Parse(addrs []string) {
	for _, addr := range addrs {
		wg.Add(1)
		go parseSingle(addr)
	}
	wg.Wait()
}

func parseSingle(addr string) {
	var site string
	var proto string
	if strings.HasPrefix(addr, "http") {
		// TODO: Take care of errors
		u, _ := url.Parse(addr)
		site = u.Host
		proto = u.Scheme
	} else {
		site = addr
		proto = defaultProto
	}

	if !strings.Contains(site, ":") {
		fmt.Println("Protocol:" + proto + " Site:" + site + " Port:" + defaultPort)
		wg.Done()
		return
	}
	site, port, err := net.SplitHostPort(site)
	if err == nil {
		fmt.Println("Protocol:" + proto + " Site:" + site + " Port:" + port)
		wg.Done()
		return
	}
	fmt.Println(err)
	wg.Done()
}
