package grpc

import (
	"testing"
	"time"

	"context"
)

func TestServer(t *testing.T) {
	go func() {
		s := NewServer(WithAddress("0.0.0.0:9000"))
		ctx := context.Background()
		err := s.Start(ctx)
		t.Log(err)
	}()
	time.Sleep(time.Second * 5)
}
