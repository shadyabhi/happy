package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

type response struct {
	ip   net.IP
	time int64
}

func main() {
	var server = flag.String("server", "linkedin.com:443", "Server to query for")
	var timeout = flag.Int("timeout", 300, "Time in milliseconds for timeout")
	flag.Parse()

	log.Printf("Connecting to server: %s with timeout: %dms", *server, *timeout)

	addresses, err := resolveAddress(*server)
	if err != nil {
		log.Fatalf("Error resolving addresses: %s", err)
	}

	var wgConnect sync.WaitGroup
	// connectAddress returns "ip" to signal which one was connected first
	// size handles return from both v4 and v6 goroutine
	results := make(chan response, 2)

	for _, addr := range addresses {
		wgConnect.Add(1)
		go connectAddress(addr, *timeout, results, &wgConnect)
	}
	wgConnect.Wait()
	close(results)

	whoWon(results, timeout)
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

func connectAddress(addr *net.TCPAddr, timeout int, results chan response, wg *sync.WaitGroup) error {
	start := time.Now()
	d := net.Dialer{Timeout: time.Duration(timeout) * time.Millisecond}
	conn, err := d.Dial("tcp", addr.String())
	if err != nil {
		log.Printf("Dial failed for address: %s, err: %s", addr.String(), err.Error())
		wg.Done()
		return err
	}

	elasped := time.Since(start)
	results <- response{ip: addr.IP, time: elasped.Nanoseconds() / 1000000}
	log.Printf("Connected to address: %s in %dms", addr.String(), elasped.Nanoseconds()/1000000)
	conn.Close()
	wg.Done()

	return nil
}

func whoWon(results chan response, timeout *int) {
	n := len(results)
	r := <-results

	if r.time > int64(*timeout) {
		if r.ip.To4() == nil {
			log.Printf("As per happy eyeballs, IPv6 won!")
		} else {
			log.Printf("As per happy eyeballs, IPv4 won!")
		}
	}
	if r.time < int64(*timeout) && r.ip.To4() != nil {
		// v4 returned before v6
		if n == 1 {
			log.Printf("As per happy eyeballs, IPv4 won!")
		}

		if n == 2 {
			v6 := <-results
			if v6.time < int64(*timeout) {
				log.Printf("As per happy eyeballs, IPv6 won!")
			} else {
				log.Printf("As per happy eyeballs, IPv4 won!")
			}
		}
	} else {
		log.Printf("As per happy eyeballs, IPv6 won!")
	}
}
