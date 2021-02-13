// Package syncio provides syncronised access to io operations
package syncio

import (
	"io"
	"sync"
)

// Writer wraps an io.Writer allowing only syncronous writes to it
type Writer struct {
	mx sync.Mutex
	w  io.Writer
}

// Writer wraps an io.Writer allowing only syncronous writes to it
type Reader struct {
	mx sync.Mutex
	r  io.Reader
}

// NewWriter returns a new Writer based on the passed io.Writer
func NewWriter(w io.Writer) *Writer {
	return &Writer{w: w}
}

// NewReader returns a new Reader based on the passed io.Reader and timeout
// Where timeout expires before the read returns, an empty byte is returned
func NewReader(r io.Reader) *Reader {
	return &Reader{r: r}
}

// Write sends the passed []byte to the underlying io.Writer once it has acquired an implicitly associated exclusive lock
func (w *Writer) Write(p []byte) (n int, err error) {
	w.mx.Lock()
	defer w.mx.Unlock()

	return w.w.Write(p)
}

// Read sends the passed []byte to the underlying io.Reader, once it has acquired an implicitly associated exclusive lock, and returns the result
func (r *Reader) Read(p []byte) (n int, err error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	return r.r.Read(p)
}

// Disable causes any and calls to write to block, until the release func returned is invoked
func (w *Writer) Lock() func() {
	w.mx.Lock()

	return func() { w.mx.Unlock() }
}

// Disable causes any and calls to read to block, until the release func returned is invoked
func (r *Reader) Lock() func() {
	r.mx.Lock()

	return func() { r.mx.Unlock() }
}
