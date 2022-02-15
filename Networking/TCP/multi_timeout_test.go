package TCP
import (
    "context"
    "net"
    "sync"
    "testing"
    "time"
)
func TestDialContextCancelFanOut(t *testing.T) {
	//>> creating the context (10s)
    ctx, cancel := context.WithDeadline(
        context.Background(),
        time.Now().Add(10*time.Second),
    )
	
	//>> Creating a listener
    listener, err := net.Listen("tcp", "127.0.0.1:")
    if err != nil {
        t.Fatal(err)
    }
    defer listener.Close()
	
	//>> Accepting Connections & Closing it after success
    go func() {
            // Only accepting a single connection.
            conn, err := listener.Accept()
            if err == nil {
                conn.Close()
            }
    }()
	
	//>> creating Dialers
    dial := func(ctx context.Context, address string, response chan int,
        id int, wg *sync.WaitGroup) {
        defer wg.Done()
        var d net.Dialer
		
        c, err := d.DialContext(ctx, "tcp", address)
		//@ connection Fail:
		// Simply return
		//@ connection success:
		// Close connection, and place the Dialer id in the response channel
        if err != nil {
            return
        }
        c.Close()
		
        select {
        case <-ctx.Done():
        case response <- id:
        }
    }
	
	//>> Calling the Dialers
    res := make(chan int)
    var wg sync.WaitGroup
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go dial(ctx, listener.Addr().String(), res, i+1, &wg)
    }
	
	//>> canceling ctx, and closing channel
    response := <-res
    cancel()
    wg.Wait()
    close(res)
	
	//>> checking err type
    if ctx.Err() != context.Canceled {
        t.Errorf("expected canceled context; actual: %s",
            ctx.Err(),
        )
    }
	
    t.Logf("dialer %d retrieved the resource", response)
}