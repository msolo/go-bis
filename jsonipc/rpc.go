package jsonipc

import (
	"io"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
	"os/exec"
	"sync"
)

type jsonRpc struct {
	cmdArgs []string

	mu  sync.Mutex
	cmd *exec.Cmd
	cl  *rpc.Client
}

type stdioRWC struct {
	r io.ReadCloser
	w io.WriteCloser
}

func (rwc stdioRWC) Read(p []byte) (n int, err error) {
	return rwc.r.Read(p)
}

func (rwc stdioRWC) Write(p []byte) (n int, err error) {
	return rwc.w.Write(p)
}

func (rwc stdioRWC) Close() error {
	rErr := rwc.r.Close()
	wErr := rwc.w.Close()
	if rErr != nil {
		return rErr
	}
	return wErr
}

func NewJSONRPC(cmdArgs ...string) (*jsonRpc, error) {
	jsr := &jsonRpc{cmdArgs: cmdArgs}
	if err := jsr.start(); err != nil {
		return nil, err
	}
	return jsr, nil
}

func (jsr *jsonRpc) start() error {
	cmd := exec.Command(jsr.cmdArgs[0], jsr.cmdArgs[1:]...)
	wr, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	rd, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	cmd.Stderr = os.Stderr

	rwc := stdioRWC{r: rd, w: wr}

	jsr.cmd = cmd
	jsr.cl = jsonrpc.NewClient(rwc)

	return cmd.Start()
}

func (jsr *jsonRpc) Close() error {
	jsr.mu.Lock()
	defer jsr.mu.Unlock()
	if err := jsr.cl.Close(); err != nil {
		// If this doesn't close cleanly, escalate to termination.
		if killErr := jsr.cmd.Process.Kill(); killErr != nil {
			return killErr
		}
	}
	return jsr.cmd.Wait()
}

func (jsr *jsonRpc) Call(method string, args interface{}, reply interface{}) error {
	jsr.mu.Lock()
	defer jsr.mu.Unlock()
	return jsr.cl.Call(method, args, reply)
}
