package codec

import (
	"bytes"
	gz "compress/gzip"
	"fmt"
	"io"
)

type gzip struct{}

func (gzip) dec(data []byte) ([]byte, error) {
	r, err := gz.NewReader(bytes.NewReader(data))

	if err != nil {
		return nil, fmt.Errorf("unable to read compressed data with gzip: [%v]", err)
	}

	b, err := io.ReadAll(r)

	if err != nil {
		return nil, fmt.Errorf("unable to decompress data with gzip: [%v]", err)
	}

	if err = r.Close(); err != nil {
		return nil, fmt.Errorf("unable to close gzip reader: [%v]", err)
	}

	return b, nil
}

func (gzip) enc(data []byte) ([]byte, error) {
	b := bytes.Buffer{}
	w := gz.NewWriter(&b)

	_, err := w.Write(data)

	if err != nil {
		return nil, fmt.Errorf("unable to compress data with gzip: [%v]", err)
	}

	if err = w.Close(); err != nil {
		return nil, fmt.Errorf("unable to close gzip writer: [%v]", err)
	}

	return b.Bytes(), nil
}
