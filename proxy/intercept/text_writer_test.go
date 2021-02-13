package intercept

import (
	"net/http"
	"strings"
	"testing"
)

func TestTextWriter(t *testing.T) {
	writeCount := 2 // the number of times write is called on the writer passed to TextWriter
	test := func(rsContentType string, expectBinaryResponse bool, c *Context) {
		rqHdrK, rqHdrV, rqbody, rsHdrK, rsHdrV, rsbody, tb := "Rq-Hk", "Rq-Hv", "rq-body-data", "Rs-Hk", "Rs-Hv", "rs-body-data", &TestBuffer{Wrote: make(chan struct{}, writeCount)}

		rq, _ := http.NewRequest(http.MethodPost, "http://www.test.com/echo/?data=some-qs-data", strings.NewReader(rqbody))
		rq.Header.Set(rqHdrK, rqHdrV)

		rqd, _ := readRequestData(rq)

		i, rsd := TextWriter("test-text-writer",
			MatchAllRequests,
			MatchAllResponses,
			false, -1, tb), &ResponseData{Status: "200 OK", Header: http.Header{rsHdrK: []string{rsHdrV}, "Content-Type": []string{rsContentType}}, Body: []byte(rsbody), Request: rqd}

		if err := i.request(c, &rqd); err != nil {
			t.Fatalf("expected no error processing request, got [%v]", err)
		}

		if err := i.response(c, rsd); err != nil {
			t.Fatalf("expected no error processing response, got [%v]", err)
		}

		for i := 0; i < writeCount; i++ {
			<-tb.Wrote
		}

		if !strings.Contains(tb.Buffer.String(), rq.URL.String()) {
			t.Fatalf("expected output to contain url [%v], got [%v]", rq.URL.String(), tb.Buffer.String())
		}

		if !strings.Contains(tb.Buffer.String(), rqHdrK) || !strings.Contains(tb.Buffer.String(), rqHdrV) {
			t.Fatalf("expected output to contain request header [%v:%v], got [%v]", rqHdrK, rqHdrV, tb.Buffer.String())
		}

		if !strings.Contains(tb.Buffer.String(), rq.Method) {
			t.Fatalf("expected output to contain method [%v], got [%v]", rq.Method, tb.Buffer.String())
		}

		if !strings.Contains(tb.Buffer.String(), rqbody) {
			t.Fatalf("expected output to contain request body [%v], got [%v]", rqbody, tb.Buffer.String())
		}

		if !strings.Contains(tb.Buffer.String(), rsd.Status) {
			t.Fatalf("expected output to contain response status [%v], got [%v]", rqbody, tb.Buffer.String())
		}

		if !strings.Contains(tb.Buffer.String(), rsHdrK) || !strings.Contains(tb.Buffer.String(), rsHdrV) {
			t.Fatalf("expected output to contain response header [%v:%v], got [%v]", rsHdrK, rsHdrV, tb.Buffer.String())
		}

		if expectBinaryResponse && !strings.Contains(tb.Buffer.String(), "[binary data]") {
			t.Fatalf("expected output to contain response body [%q], got [%v]", "[binary data]", tb.Buffer.String())
		}

		if !expectBinaryResponse && !strings.Contains(tb.Buffer.String(), rsbody) {
			t.Fatalf("expected output to contain response body [%v], got [%v]", rsbody, tb.Buffer.String())
		}

		assertExpectedExistence := func(s, text string, expectExists bool) {
			exists := strings.Contains(s, text)

			if expectExists != exists {
				t.Fatalf("expected %v exists in %v to be %v. but was %v", text, s, expectExists, exists)
			}
		}

		assertExpectedExistence(tb.Buffer.String(), ">>> *edited", c.RequestEdited)
		assertExpectedExistence(tb.Buffer.String(), "*rerouted", c.RequestRerouted)
		assertExpectedExistence(tb.Buffer.String(), "<<< *edited", c.ResponseEdited)
	}

	test("text/plain", false, &Context{RequestRerouted: true, RequestEdited: true, ResponseEdited: true})
	test("image/gif", true, &Context{RequestRerouted: false, RequestEdited: false, ResponseEdited: false})
}
