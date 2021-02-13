// Package codec provides a generic mechanism for reading and writing various content encoding schemes
package codec

import "fmt"

type codec interface {
	enc(data []byte) ([]byte, error)
	dec(data []byte) ([]byte, error)
}

var codecs = map[string]codec{
	"gzip": &gzip{},
	"br":   &brotli{},
}

// Supported returns true if scheme is a supported content encoding
func Supported(scheme string) bool {
	return codecs[scheme] != nil
}

// Decode decodes data using scheme
func Decode(scheme string, data []byte) ([]byte, error) {
	c := codecs[scheme]

	if c == nil {
		return nil, fmt.Errorf("unsupported encoding scheme [%v]", scheme)
	}

	return c.dec(data)
}

// Encode encodes data using scheme
func Encode(scheme string, data []byte) ([]byte, error) {
	c := codecs[scheme]

	if c == nil {
		return nil, fmt.Errorf("unsupported encoding scheme [%v]", scheme)
	}

	return c.enc(data)
}
