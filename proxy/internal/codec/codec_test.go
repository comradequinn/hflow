package codec

import "testing"

func testHandler(t *testing.T, scheme string, c codec) {
	data := "some test\ndata"

	var (
		b   []byte
		err error
	)

	b, err = c.enc([]byte(data))

	if err != nil {
		t.Fatalf("unable to encode data with [%v]: [%v]", scheme, err)
	}

	b, err = c.dec(b)

	if err != nil {
		t.Fatalf("unable to decode data with [%v]: [%v]", scheme, err)
	}

	if string(b) != data {
		t.Fatalf("data decoded with [%v] was incorrect. Got [%v] expected [%v]", scheme, string(b), data)
	}
}

func TestGzip(t *testing.T) {
	testHandler(t, "gzip", gzip{})
}

func TestBrotli(t *testing.T) {
	testHandler(t, "brotli", brotli{})
}
