package copy

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWriter(t *testing.T) {
	data := "some test data"
	r := httptest.NewRequest("GET", "http://somewhere.outhere", strings.NewReader(data))

	s, err := CloserToString(&r.Body)

	if err != nil {
		t.Fatalf("expected no error but got [%v]", err)
	}

	if s != data {
		t.Fatalf("expected [%v] in response but got [%v]", data, s)
	}

	var b []byte

	b, err = io.ReadAll(r.Body)

	s = string(b)

	if s != data {
		t.Fatalf("expected [%v] in original closer but got [%v]", data, s)
	}
}
