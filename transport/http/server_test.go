package http

import (
	"context"
	"testing"
)

func TestServer(t *testing.T) {
	s := NewServer(WithAddress("0.0.0.0:8080"))
	_ = s.Start(context.Background())
}
