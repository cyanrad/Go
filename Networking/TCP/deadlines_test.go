package TCP

import (
	"net"
	"testing"
	"time"
)

func TestDeadline(t *testing.T) {
	sync := make(chan struct{})

	// >> Creating listener
	listener, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}

	// >> Accepting connections
	go func() {
		conn, err := listener.Accept()
		if err != nil {
			t.Log(err)
			return
		}
		defer func() {
			conn.Close()
			close(sync) // read from sync shouldn't block due to early return
		}()

		// >> Setting connection deadline
		err = conn.SetDeadline(time.Now().Add(5 * time.Second))
		if err != nil {
			t.Error(err)
			return
		}

		// >> Reading data
		buf := make([]byte, 1)
		// blocka until Data is recived
		_, err = conn.Read(buf) //should return an error after 5 sec
		// in this example it will fails since the clien't write function is stopped by the sync channel
		nErr, ok := err.(net.Error) // >> Verifying that the error is a timeout
		if !ok || nErr.Timeout() {
			t.Errorf("expected timeout error; actual: %v", err)
		}

		// Any future reads will result in a timeout error
		// however read functionality can be restored with pushing the deadline
		sync <- struct{}{}
		err = conn.SetDeadline(time.Now().Add(5 * time.Second))
		if err != nil {
			t.Error(err)
			return
		}

		// this read will succseed since we extended the deadline
		_, err = conn.Read(buf)
		if err != nil {
			t.Error(err)
		}
	}()

	// >> Creating connectino
	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	<-sync //witing for client and server to be done

	// writing to the server
	_, err = conn.Write([]byte("1"))
	if err != nil {
		t.Fatal(err)
	}

	buf := make([]byte, 1)
	data, _ := conn.Read(buf) // blocked until remote node sends data
	//sends data @ closing connection
	t.Log(data)
}
