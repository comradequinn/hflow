// Package copy provides copy related utility functions
package copy

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"
)

// CloserToString returns the contents of the r as a string and resets r so it can be read again
func CloserToString(r *io.ReadCloser) (string, error) {
	b, err := CloserToBytes(r)

	return string(b), err
}

// CloserToBytes returns the contents of the r as a []byte and resets r so it can be read again
func CloserToBytes(r *io.ReadCloser) ([]byte, error) {
	b, err := io.ReadAll(*r)

	if err != nil {
		return nil, err
	}

	*r = BytesToCloser(b)

	return b, nil
}

// BytesToCloser sets b as the contents of the returned io.ReadCloser leaving it ready to be read
func BytesToCloser(b []byte) io.ReadCloser {
	return io.NopCloser(bytes.NewReader(b))
}

// Header copies the contents of src to dest
func Header(src http.Header, dest http.Header) {
	for k, v := range src {
		hv := strings.Join(v, " ")
		dest.Set(k, hv)
	}
}

// HTTPRequest returns a new *http.Request based on src
func HTTPRequest(src *http.Request) *http.Request {
	dest := src.Clone(context.Background())
	dest.RequestURI = ""

	return dest
}
