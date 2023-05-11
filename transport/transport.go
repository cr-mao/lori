package transport

import (
	"context"
	"net/url"
)

// Server is transport server.
type Server interface {
	Start(context.Context) error
	Stop(context.Context) error
}

// Endpointer is registry endpoint.
type Endpointer interface {
	Endpoint() (*url.URL, error)
}

// Transporter is transport context value interface.
type Transporter interface {
	// Kind transporter
	// grpc
	// http
	// tcp
	Kind() Kind
	// Endpoint return server or client endpoint
	// Server Transport: grpc://127.0.0.1:9000
	// Client Transport: discovery:///provider-demo
	Endpoint() string
}

// Kind defines the type of Transport
type Kind string

func (k Kind) String() string { return string(k) }

// Defines a set of transport kind
const (
	KindGRPC Kind = "grpc"
	KindHTTP Kind = "http"
	KindTCP  Kind = "tcp"
)

type (
	serverTransportKey struct{}
	clientTransportKey struct{}
)

// NewServerContext returns a new Context that carries value.
func NewServerContext(ctx context.Context, tr Transporter) context.Context {
	return context.WithValue(ctx, serverTransportKey{}, tr)
}

// FromServerContext returns the Transport value stored in ctx, if any.
func FromServerContext(ctx context.Context) (tr Transporter, ok bool) {
	tr, ok = ctx.Value(serverTransportKey{}).(Transporter)
	return
}

// NewClientContext returns a new Context that carries value.
func NewClientContext(ctx context.Context, tr Transporter) context.Context {
	return context.WithValue(ctx, clientTransportKey{}, tr)
}

// FromClientContext returns the Transport value stored in ctx, if any.
func FromClientContext(ctx context.Context) (tr Transporter, ok bool) {
	tr, ok = ctx.Value(clientTransportKey{}).(Transporter)
	return
}
