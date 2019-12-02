package jsonipc

import "testing"

type ipcArgs struct {
	Val string
}

type ipcReply struct {
	Val   string
	Error interface{}
}

func TestIpc(t *testing.T) {
	cl, err := NewJSONIpc("./fakeipc/fakeipc")
	if err != nil {
		t.Fatal("NewJSONIpc failed", err)
	}

	defer func() {
		if err := cl.Close(); err != nil {
			t.Error("unable to close client", err)
		}
	}()

	reply := &ipcReply{}
	msg := "test arg"
	err = cl.Call(ipcArgs{Val: msg}, reply)
	if err != nil {
		t.Fatal("Call failed", err)
	}
	if reply.Val != msg {
		t.Fatalf("invalid Result, expected: %s received: %s", msg, reply.Val)
	}
}

func TestMissingBinary(t *testing.T) {
	_, err := NewJSONIpc("./nonexistent-binary")
	if err == nil {
		t.Fatal("expected failure for missing binary")
	}
}

func TestFailingBinary(t *testing.T) {
	cl, err := NewJSONIpc("/usr/bin/false")
	if err != nil {
		t.Fatal("NewJSONIpc failed", err)
	}
	reply := &ipcReply{}
	msg := "test arg"
	err = cl.Call(ipcArgs{Val: msg}, reply)
	if err == nil {
		t.Fatal("Call should have failed with failing binary")
	}
}
