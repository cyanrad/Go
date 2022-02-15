package TCP
import (
    "context"
    "net"
    "syscall"
    "testing"
    "time"
)

func TestDialContext(t *testing.T) {
	// >> creating deadline for 5 sec into the future
	// After which the context will automatically cancle
    dl := time.Now().Add(2 * time.Second)	
	
	// >> Creating the context
	// context.WithCancle(context.Background()) can work if you don't want a deadline
    ctx, cancel := context.WithDeadline(context.Background(), dl)
    defer cancel()	// good practice, to make sure ctx is garbage collected
	
	go func(ctx context.Context, t *testing.T){
		<-ctx.Done()
		//t.Fatal("context timedout")
		//t.Fail()
	}(ctx, t)
	
	
	// >> Dialer setup
    var d net.Dialer // DialContext is a method on a Dialer
    d.Control = func(_, _ string, _ syscall.RawConn) error {
        // Sleep long enough to reach the context's deadline.
        time.Sleep(5*time.Second + time.Millisecond)
        return nil
    }
	// Dialing with a context
	
    conn, err := d.DialContext(ctx, "tcp", "192.168.0.189:80")
    if err == nil {
        conn.Close()
        t.Fatal("connection did not time out")
    }
    nErr, ok := err.(net.Error)
    if !ok {
        t.Error(err)
    } else {
        if !nErr.Timeout() {
            t.Errorf("error is not a timeout: %v", err)
        }
    }
	
	// makes sure that the error comes from ctx exeeding
	if ctx.Err() != context.DeadlineExceeded {
		t.Fatal("expected deadline exceeded; actual: ", ctx.Err())
    }
	
}