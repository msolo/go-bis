// A simple JSON-over-stdio API for communicating with an external process.
package jsonipc

import (
	"encoding/json"
	"os"
	"os/exec"
	"sync"
)

type jsonIpc struct {
	args []string

	mu   sync.Mutex
	cmd  *exec.Cmd
	jsWr *json.Encoder
	jsRd *json.Decoder
}

func newJSONIpc(cmdArgs ...string) *jsonIpc {
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
	jsi.jsWr = json.NewEncoder(wr)
	jsi.jsRd = json.NewDecoder(rd)

	return cmd.Start()
}

func (jsi *jsonRpc) Close() error {
	jsi.mu.Lock()
	defer jsi.mu.Unlock()
	if err := jsi.cl.Close(); err != nil {
		// If this doesn't close cleanly, escalate to termination.
		if killErr := jsi.cmd.Process.Kill(); killErr != nil {
			return killErr
		}
	}
	return jsi.cmd.Wait()
}

func (jsi *jsonIpc) Call(args interface{}, reply interface{}) error {
	jsi.mu.Lock()
	defer jsi.mu.Unlock()
	if err := jsi.jsWr.Encode(args); err != nil {
		return err
	}
	return jsi.jsid.Decode(reply)
}
