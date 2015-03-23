package gofigure

import (
	"fmt"
	"os"
	"sync"
	"syscall"
	"testing"
	"time"
)

func TestReloader(t *testing.T) {

	r := NewSignalMonitor()
	if r == nil {
		t.Fatal("Could not create reloader")
	}

	n := 0
	wg := sync.WaitGroup{}
	wg.Add(1)

	r.Monitor(ReloadFunc(func() {
		fmt.Println("Received update")
		n++

		wg.Done()
	}))

	time.Sleep(10 * time.Millisecond)

	pid := os.Getpid()
	fmt.Println(pid)

	p := &os.Process{Pid: pid}
	p.Signal(syscall.SIGHUP)

	wg.Wait()
	r.Stop()
	if n != 1 {
		t.Errorf("No updates")
	}

}
