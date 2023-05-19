package tcp_test

import (
	"testing"

	"github.com/cr-mao/lori/transport/network"
	"github.com/cr-mao/lori/transport/network/tcp"
)

func TestServer(t *testing.T) {
	t.SkipNow()
	server := tcp.NewServer()
	server.OnStart(func() {
		t.Logf("server is started")
	})
	server.OnConnect(func(conn network.Conn) {
		t.Logf("connection is opened, connection id: %d", conn.ID())
	})
	server.OnDisconnect(func(conn network.Conn) {
		t.Logf("connection is closed, connection id: %d", conn.ID())
	})
	server.OnReceive(func(conn network.Conn, msg []byte, msgType int) {
		t.Logf("receive msg from client, connection id: %d, msg: %s", conn.ID(), string(msg))

		if err := conn.Push([]byte("I'm fine~~")); err != nil {
			t.Error(err)
		}
	})

	if err := server.Start(); err != nil {
		t.Fatal(err)
	}
}
