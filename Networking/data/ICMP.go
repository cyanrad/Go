package data

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"
)

var ( //CLI options, that provide some of the ping functionality
	Count    = flag.Int("c", 3, "number of pings: <= 0 means forever")
	Interval = flag.Duration("i", time.Second, "interval between pings")
	Timeout  = flag.Duration("W", 5*time.Second, "time to wait for a reply")
)

func init() {
	flag.Usage = func() {
		fmt.Printf("Usage: %s [options] host:port\nOptions:\n", os.Args[0])
		flag.PrintDefaults()
	}
}

func temp() {
	// flag err handling
	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Print("host:port is required\n\n")
		flag.Usage()
		os.Exit(1)
	}

	target := flag.Arg(0)
	fmt.Println("PING", target)

	if *Count <= 0 { // 0 pings
		fmt.Println("CTRL+C to stop.")
	}

	msg := 0
	for (*Count <= 0) || (msg < *Count) {
		msg++
		fmt.Print(msg, " ")

		start := time.Now()
		c, err := net.DialTimeout("tcp", target, *Timeout) // connecting to target
		dur := time.Since(start)

		if err != nil {
			fmt.Printf("fail in %s: %v\n", dur, err)
			if nErr, ok := err.(net.Error); !ok || !nErr.Temporary() {
				os.Exit(1)
			}
		} else {
			_ = c.Close()
			fmt.Println(dur)
		}
		time.Sleep(*Interval)
	}

}

func PingTarget(address string) {
	if *Count <= 0 { // 0 pings
		fmt.Println("CTRL+C to stop.")
	}

	msg := 0
	for (*Count <= 0) || (msg < *Count) {
		msg++
		fmt.Print(msg, " ")

		start := time.Now()
		c, err := net.DialTimeout("tcp", address, *Timeout) // connecting to target
		dur := time.Since(start)

		if err != nil {
			fmt.Printf("fail in %s: %v\n", dur, err)
			if nErr, ok := err.(net.Error); !ok || !nErr.Temporary() {
				os.Exit(1)
			}
		} else {
			_ = c.Close()
			fmt.Println(dur)
		}
		time.Sleep(*Interval)
	}

}
