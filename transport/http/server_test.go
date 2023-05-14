package http

import (
	"context"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	go func() {
		s := NewServer(WithAddress("0.0.0.0:8080"))
		_ = s.Start(context.Background())
	}()

	time.Sleep(time.Second * 5)
}
