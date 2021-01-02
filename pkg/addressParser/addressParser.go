package addressParser

import (
	"fmt"
	"strings"
	"net"
	"sync"
)

const defaultPort = "443"

var wg sync.WaitGroup

func Parse(addresses []string) {
	for _, address := range addresses {
		wg.Add(1)
		go parseSingle(address)
	}
	wg.Wait()
}

func parseSingle(address string) {
	if !strings.Contains(address, ":") {
		fmt.Println("Site:" + address + " Port:" + defaultPort)
		wg.Done()
		return
	}
	site, port, err := net.SplitHostPort(address)
	if err == nil {
		fmt.Println("Site:" + site + " Port:" + port)
		wg.Done()
		return
	}
	fmt.Println(err)
	wg.Done()
}
