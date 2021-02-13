package intercept

import (
	"net/http"
	"strings"
	"testing"
)

func TestRequestString(t *testing.T) {
	httpRq, _ := http.NewRequest(http.MethodPost, "http://www.test.com/echo/?data=some-qs-data", strings.NewReader("body-data"))
	httpRq.Header.Add("h1", "v1")
	rqd, err := readRequestData(httpRq)

	if err != nil {
		t.Fatalf("expected no error. got %v", err)
	}

	expected := `POST /echo/?data=some-qs-data HTTP/1.1
H1: v1

body-data
`
	if string(rqd.String()) != expected {
		t.Fatalf("expected:\n%v\n. got:\n%v\n", expected, rqd.String())
	}
}

func TestResponseString(t *testing.T) {
	rsd := ResponseData{
		Header:     http.Header{"k1": []string{"v1"}},
		Body:       []byte("body-data"),
		Status:     "200 OK",
		StatusCode: 200,
		Request:    RequestData{Protocol: "HTTP/1.1"},
	}

	expected := `HTTP/1.1 200 OK
k1: v1

body-data
`
	if string(rsd.String()) != expected {
		t.Fatalf("expected:\n%v\n. got:\n%v\n", expected, rsd.String())
	}
}
