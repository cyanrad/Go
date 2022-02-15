package TCP

import (
	"context"
	"io"
	"time"
)

const defaultPingInterval = 30 * time.Second

// >> write pings at regular intervals
// context arg for termination and leakage prevention
// the channel is to signal timer reset
func pinger(ctx context.Context, w io.Writer, reset <-chan time.Duration) {
	// the interval time value
	var interval time.Duration

	// >> getting interval time
	// we put the initial interval duration in the reset channel
	select {
	case <-ctx.Done(): //terminating
		return
	case interval = <-reset: // pulled initial interval off reset channel
	default:
	}

	// >> interval duration handling
	if interval <= 0 {
		interval = defaultPingInterval
	}

	// >> creating ping timer
	timer := time.NewTimer(interval)
	defer func() { // defer call to drain timer channel to avoide leakage
		if !timer.Stop() {
			<-timer.C
		}
	}()

	// >> pinging loop
	// keep track of time-outs by passing the ctx's cancle func
	// and call it here if @ max concecutive timeouts
	for {
		select {
		case <-ctx.Done(): //terminate
			return
		case newInterval := <-reset: // resetting the timer (if data recieved)
			if !timer.Stop() {
				<-timer.C // Blocking wait until it finishes
			}
			if newInterval > 0 {
				interval = newInterval
			}
		case <-timer.C: // ping (timer expires) @ timer expire
			if _, err := w.Write([]byte("ping")); err != nil {
				// track and act on consecutive timeouts here
				return
			}
		}
		_ = timer.Reset(interval) // reset timer
	}
}
