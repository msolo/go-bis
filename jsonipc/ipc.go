// A simple JSON-over-stdio API for communicating with an external process.
//
// JSON objects are serially encoded to stdin for the child process
// and replies are decorde from stdout.
package jsonipc

import (
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"sync"
)

type jsonIpc struct {
	args []string

	mu   sync.Mutex
	cmd  *exec.Cmd
	wr   io.WriteCloser
	jsWr *json.Encoder
	jsRd *json.Decoder
}

// A minimal IPC mechanism that communicates over stdio.  The child
// process will be started immediately. The child process should exit
// gracefully when stdin is closed. If not, it will received SIGTERM.
func NewJSONIpc(cmdArgs ...string) (*jsonIpc, error) {
	jsi := &jsonIpc{args: cmdArgs}
	if err := jsi.start(); err != nil {
		return nil, err
	}
	return jsi, nil
}

func (jsi *jsonIpc) start() error {
	cmd := exec.Command(jsi.args[0], jsi.args[1:]...)
	wr, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	rd, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	cmd.Stderr = os.Stderr

	jsi.cmd = cmd
	jsi.wr = wr
	jsi.jsWr = json.NewEncoder(wr)
	jsi.jsRd = json.NewDecoder(rd)

	return cmd.Start()
}

func (jsi *jsonIpc) Close() error {
	jsi.mu.Lock()
	defer jsi.mu.Unlock()
	if err := jsi.wr.Close(); err != nil {
		// If this doesn't close cleanly, escalate to termination.
		// FIXME(msolo) do we need to SIGKILL?
		if killErr := jsi.cmd.Process.Kill(); killErr != nil {
			return killErr
		}
	}
	return jsi.cmd.Wait()
}

// Each Call is a serialized request-response pair. There is no
// concurrent access to the underlying process.
//
// This API matches the RPC transport, but communicating application
// level errors is left as a part of the protocol. Only I/O and
// encoding errors will be returned.
func (jsi *jsonIpc) Call(args interface{}, reply interface{}) error {
	jsi.mu.Lock()
	defer jsi.mu.Unlock()
	if err := jsi.jsWr.Encode(args); err != nil {
		return err
	}
	return jsi.jsRd.Decode(reply)
}
