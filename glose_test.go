package glose_test

import (
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/karim-w/glose"
)

type dummyCloser struct{}

func (d *dummyCloser) Close() error {
	return nil
}

func TestClosable(t *testing.T) {
	os.Setenv("GLOSE_SKIP_EXIT", "true")

	dummy := &dummyCloser{}

	glose.Register(dummy)

	go glose.Watch()

	pid := os.Getpid()
	t.Logf("Current process PID: %d", pid)

	// Send SIGTERM signal to the current process after a short delay.
	go func() {
		time.Sleep(2 * time.Second)
		t.Log("Sending SIGTERM signal...")
		if err := syscall.Kill(pid, syscall.SIGTERM); err != nil {
			t.Fatalf("Error sending SIGTERM: %v", err)
		}
	}()

	time.Sleep(5 * time.Second)
}

func TestPanikEmpty(t *testing.T) {
	os.Setenv("GLOSE_SKIP_EXIT", "true")

	err := os.Chdir("non-existing-dir")

	// defer recover()
	defer func() {
		if r := recover(); r != nil {
			t.Logf("Recovered: %v", r)
		}
	}()

	glose.Panik(err)
}

func TestPanik(t *testing.T) {
	os.Setenv("GLOSE_SKIP_EXIT", "true")

	err := os.Chdir("non-existing-dir")
	glose.Register(&dummyCloser{})

	defer func() {
		if r := recover(); r != nil {
			t.Logf("Recovered: %v", r)
		}
	}()

	glose.Panik(err)
}
