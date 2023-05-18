package netlib

import (
	"testing"
)

func TestParseAddr(t *testing.T) {
	listenAddr, exposeAddr, err := ParseAddr(":0")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(listenAddr, exposeAddr)
}
