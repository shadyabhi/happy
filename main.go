package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

func main() {
	var server = flag.String("server", "linkedin.com:443", "Server to query for")
	var timeout = flag.Int("timeout", 300, "Time in milliseconds for timeout")
	flag.Parse()

	log.Printf("Connecting to server: %s with timeout: %dms", *server, *timeout)

	addresses, err := resolveAddress(*server)
	if err != nil {
		log.Fatalf("Error resolving addresses: %s", err)
	}

	var wg sync.WaitGroup
	// connectAddress returns "ip" to signal which one was connected first
	// size handles return from both v4 and v6 goroutine
	result := make(chan net.IP, 2)

	for _, addr := range addresses {
		wg.Add(1)
		go connectAddress(addr, *timeout, result, &wg)
	}
	wg.Add(1)
	go whoWon(result, &wg)

	wg.Wait()

}

func resolveAddress(server string) (addresses [2]*net.TCPAddr, err error) {

	for i, addrType := range []string{"tcp4", "tcp6"} {
		addr, err := net.ResolveTCPAddr(addrType, server)
		if err != nil {
			return [2]*net.TCPAddr{}, fmt.Errorf("Error resolving ip address from server: server=%s, err=%s", server, err)
		}
		addresses[i] = addr
	}
	return addresses, nil
}

func connectAddress(addr *net.TCPAddr, timeout int, result chan net.IP, wg *sync.WaitGroup) error {
	start := time.Now()
	d := net.Dialer{Timeout: time.Duration(timeout) * time.Millisecond}
	conn, err := d.Dial("tcp", addr.String())
	if err != nil {
		log.Printf("Dial failed for address: %s, err: %s", addr.String(), err.Error())
		wg.Done()
		return err
	} else {
		result <- addr.IP
	}

	elasped := time.Since(start)
	log.Printf("Connected to address: %s in %dms", addr.String(), elasped.Nanoseconds()/1000000)
	conn.Close()
	wg.Done()

	return nil
}

func whoWon(result chan net.IP, wg *sync.WaitGroup) {
	r := <-result
	if r.To4() == nil {
		log.Printf("IPv6 won!")
	} else {
		log.Printf("IPv4 won!")
	}
	wg.Done()
}
