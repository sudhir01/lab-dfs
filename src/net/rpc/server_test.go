package rpc

import "testing"

func TestIsRegistered(test *testing.T) {
    server := rpc.NewServer()

    if server.IsRegistered("testing") {
        t.Fatalf("got true when default is false\n")
    }
}
