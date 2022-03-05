package main

import (
	"flag"
	"fmt"
	"net"
	"strconv"
	"time"
)

func scanPort(protocol string, hostname string, port int, maxTimeout time.Duration) bool {
	var address string = hostname + ":" + strconv.Itoa(port)
	conn, err := net.DialTimeout(protocol, address, maxTimeout*time.Second)

	if err != nil {
		return false
	}

	defer conn.Close()
	return true
}

func main() {

	// Defining CLI Flags
	protoPTR := flag.String("protocol", "tcp", "Set the protocol you want, TCP or UDP. DEFAULT: TCP")
	hostPTR := flag.String("hostname", "localhost", "Set the hostname you want to connect to. DEFAULT: localhost")
	portPtr := flag.Int("port", 80, "Set the port to scan. DEFAULT: 80")
	timeoutPtr := flag.Duration("timeout", 60, "Set the timout allowed on ports. DEFAULT: 60")

	// Parse Flags
	flag.Parse()
	open := scanPort(*protoPTR, *hostPTR, *portPtr, *timeoutPtr)
	open2 := scanPort(*protoPTR, *hostPTR, 40, *timeoutPtr)

	fmt.Printf("Connected? %t\n", open)
	fmt.Printf("Connected? %t\n", open2)

}
