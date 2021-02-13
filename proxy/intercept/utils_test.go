package intercept

import (
	"bytes"
	"sync"
)

// TestBuffer embeds a bytes.Buffer to which it ensures synchronous writes and notifies of them via Wrote
type TestBuffer struct {
	mx     sync.Mutex
	Wrote  chan struct{}
	Buffer bytes.Buffer
}

// Write implements io.Writer and signals written when it completes a write
func (t *TestBuffer) Write(p []byte) (n int, err error) {
	t.mx.Lock()
	defer t.mx.Unlock()

	n, err = t.Buffer.Write(p)
	t.Wrote <- struct{}{}

	return n, err
}
