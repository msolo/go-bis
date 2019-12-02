package jsonipc

import "testing"

type rpcArgs struct {
	Val string
}

type rpcReply struct {
	Val string
}

func TestRpc(t *testing.T) {
	cl, err := NewJSONRPC("./fakerpc/fakerpc")
	if err != nil {
		t.Fatal("NewJSONRPC failed", err)
	}

	defer func() {
		if err := cl.Close(); err != nil {
			t.Error("unable to close client", err)
		}
	}()

	reply := &rpcReply{}
	msg := "test arg"
	err = cl.Call("Echo.Echo", rpcArgs{Val: msg}, reply)
	if err != nil {
		t.Fatal("Call failed", err)
	}
	if reply.Val != msg {
		t.Fatalf("invalid Result, expected: %s received: %s", msg, reply.Val)
	}
}

func TestRpcMissingBinary(t *testing.T) {
	_, err := NewJSONRPC("./nonexistent-binary")
	if err == nil {
		t.Fatal("expected failure for missing binary")
	}
}

func TestRpcFailingBinary(t *testing.T) {
	cl, err := NewJSONRPC("/usr/bin/false")
	if err != nil {
		t.Fatal("NewJSONRPC failed", err)
	}
	reply := &rpcReply{}
	msg := "test arg"
	err = cl.Call("Echo.Echo", rpcArgs{Val: msg}, reply)
	if err == nil {
		t.Fatal("Call should have failed with failing binary")
	}
}
