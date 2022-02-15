package TCP

import (
	"context"
	"io"
	"net"
	"testing"
	"time"
)

func TestPingerAdvanceDeadline(t *testing.T) {
	// >> Listener
	done := make(chan struct{})
	listener, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}

	// >> pinger
	// deadline: 5 seconds
	// ping interval: every second
	// from client side: recieve 4 pings then an io.EOF
	begin := time.Now()
	go func() {
		defer func() { close(done) }() // closing the done channel when process is complete

		// >> accepting conneciton
		conn, err := listener.Accept()
		if err != nil {
			t.Log(err)
			return
		}

		// pinger control context
		ctx, cancel := context.WithCancel(context.Background())
		defer func() {
			cancel()
			conn.Close()
		}()

		// buffered timer reset channel
		resetTimer := make(chan time.Duration, 1)
		resetTimer <- time.Second

		// the actual pinger
		go pinger(ctx, conn, resetTimer)

		// setting connection deadline
		err = conn.SetDeadline(time.Now().Add(5 * time.Second))
		if err != nil {
			t.Error(err)
			return
		}

		buf := make([]byte, 1024)
		for {
			// writing to buffer client data
			n, err := conn.Read(buf) //blocks until data is recieved i think
			if err != nil {
				t.Log("read err: ", err)
				return
			}
			t.Logf("[%s] %s", //logging read data
				time.Since(begin).Truncate(time.Second), buf[:n])

			// resetting pinger
			resetTimer <- 0

			// resetting deadline
			err = conn.SetDeadline(time.Now().Add(5 * time.Second))
			if err != nil {
				t.Error(err)
				return
			}
		}
	}()

	// >> client init
	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	// >> Reading pings
	buf := make([]byte, 1024)
	for i := 0; i < 4; i++ { // read up to four pings
		n, err := conn.Read(buf)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("[%s] %s", time.Since(begin).Truncate(time.Second), buf[:n])
	}

	// writing to server
	_, err = conn.Write([]byte("PONG!!!")) // should reset the ping timer
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 4; i++ { // read up to four more pings
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				t.Fatal(err)
			}
			break
		}
		t.Logf("[%s] %s", time.Since(begin).Truncate(time.Second), buf[:n])
	}

	// @ finish
	<-done
	end := time.Since(begin).Truncate(time.Second)
	t.Logf("[%s] done", end)
	if end != 9*time.Second {
		t.Fatalf("expected EOF at 9 seconds; actual %s", end)
	}
}
