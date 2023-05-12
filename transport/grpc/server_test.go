package grpc

import (
	"context"
	"testing"
)

func TestServer(t *testing.T) {
	s := NewServer(WithAddress("0.0.0.0:9000"))
	ctx := context.Background()
	err := s.Start(ctx)
	t.Log(err)
}
