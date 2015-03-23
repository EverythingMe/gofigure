package gofigure

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

// Reloader is an interface to be implemented by the calling program, that gets called when we need
// to reload our configs
type Reloader interface {
	Reload()
}

// ReloadFunc can be used to make a simple func conform to the Reloader interface
type ReloadFunc func()

func (f ReloadFunc) Reload() {
	f()
}

// ReloadMonitor is an interface for waiting for external notifications that we need to reload our configs.
// The only implementation right now is a SIGHUP listener
type ReloadMonitor interface {
	Watch(Reloader)
	Stop()
}

// SignalMonitor is a monitor that waits for SIGHUP and calls its Reloader
type SignalMonitor struct {
	stopch chan bool
}

// NewSignalMonitor creates a new signal monitor
func NewSignalMonitor() *SignalMonitor {
	return &SignalMonitor{
		make(chan bool),
	}
}

// Monitor waits for SIGHUP and calls the reloader's Reload method
func (m *SignalMonitor) Monitor(r Reloader) {

	go func() {
		sigch := make(chan os.Signal, 1)
		signal.Notify(sigch, syscall.SIGHUP)
		for {
			// Block until a signal is received.

			select {
			case sig := <-sigch:
				log.Printf("Received signal %v", sig)
				// call the reloader to reload its configs
				r.Reload()
			case <-m.stopch:
				log.Printf("Stopping reload listener")
				break
			}
		}
		signal.Stop(sigch)

	}()
}

// Stop stops the signal monitor
func (m *SignalMonitor) Stop() {
	m.stopch <- true
	log.Printf("Stopped signal monitor")
}
