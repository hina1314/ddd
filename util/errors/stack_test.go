package errors

import (
	"testing"
	"time"
)

func TestCaptureStackReturnsAfterSkippedFinalFrame(t *testing.T) {
	done := make(chan *StackTrace, 1)
	go func() {
		done <- CaptureStack(0)
	}()

	select {
	case stack := <-done:
		if stack == nil || len(stack.Frames) == 0 {
			t.Fatal("expected at least one application frame")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("CaptureStack did not return")
	}
}
