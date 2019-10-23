package slog

import (
	"bytes"
	"encoding/json"
	"sync"
)

func JsonFmtEntry(e Entry) string {
	data, err := json.MarshalIndent(e, "", "")
	if err == nil {
		buf := &bytes.Buffer{}
		if err := json.Compact(buf, data); err == nil {
			data = buf.Bytes()
		}
	}
	if err != nil {
		// No point in handling this error.
		data, _ = json.MarshalIndent(map[string]string{"JsonErr": err.Error()}, "", "")
	}
	data = append(data, '\n')
	// FIXME(msolo) Yes, this is pointless. Fix crazy copying.
	return string(data)
}

type multiHandler struct {
	mu       sync.Mutex
	handlers []Handler
}

func (mh *multiHandler) WriteEntry(e Entry) error {
	mh.mu.Lock()
	defer mh.mu.Unlock()
	for _, h := range mh.handlers {
		if err := h.WriteEntry(e); err != nil {
			return err
		}
	}
	return nil
}

func NewMultiHandler(handlers ...Handler) Handler {
	return &multiHandler{handlers: handlers}
}
