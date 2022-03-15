package main

import (
	"context"
	"fmt"
	"net"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
)

type PortScanner struct {
	hostname string
	lock     *semaphore.Weighted
}

func Ulimit() int64 {
	out, err := exec.Command("/bin/bash", "-c", "ulimit -n").Output()
	if err != nil {
		panic(err)
	}

	s := strings.TrimSpace(string(out))

	i, err := strconv.ParseInt(s, 10, 64)

	if err != nil {
		panic(err)
	}

	return i
}

func ScanPort(hostname string, port int, maxTimeout time.Duration) {
	address := fmt.Sprintf("%s:%d", hostname, port)
	conn, err := net.DialTimeout("tcp", address, maxTimeout)

	if err != nil {
		//fmt.Println(err)
		if strings.Contains(err.Error(), "too many open files") {
			time.Sleep(maxTimeout)
			// If we open too many files, sleep and restart the scan
			ScanPort(hostname, port, maxTimeout)
		} else {
			//fmt.Println(port, "closed")
		}
		return
	}

	conn.Close()
	fmt.Println(address, " is open")
}

func (ps *PortScanner) Start(f, l int, timeout time.Duration) {
	wg := sync.WaitGroup{}
	defer wg.Wait()

	for port := f; port <= l; port++ {
		ps.lock.Acquire(context.TODO(), 1)
		wg.Add(1)
		go func(port int) {
			defer ps.lock.Release(1)
			defer wg.Done()
			ScanPort(ps.hostname, port, timeout)
		}(port)
	}
}

func main() {

	// Defining CLI Flags
	//protoPTR := flag.String("protocol", "tcp", "Set the protocol you want, TCP or UDP.")
	//hostPTR := flag.String("hostname", "localhost", "Set the hostname you want to connect to.")
	//portPtr := flag.Int("port", 80, "Set the port to scan.")
	//timeoutPTR := flag.Duration("timeout", 50, "Set the timout, lower is better but too low would make any port seem closed.")

	// Parse Flags
	//flag.Parse()

	//scanPort(*protoPTR, *hostPTR, *portPtr, *timeoutPTR)

	ps := &PortScanner{
		hostname: "google.com",
		lock:     semaphore.NewWeighted(Ulimit()),
	}
	ps.Start(1, 65535, 500*time.Millisecond)
}
