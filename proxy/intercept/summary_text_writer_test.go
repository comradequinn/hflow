package intercept

import (
	"comradequinn/hflow/proxy/store"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	id := 0

	storeInsert = func(r store.Record) int {
		id++
		return id
	}

	os.Exit(m.Run())
}

func TestSummaryTextWriter(t *testing.T) {
	bodyData := "body-data"
	url, _ := url.Parse("http://www.test.com/echo/?data=some-qs-data")

	act := func(edited, rerouted bool) {
		rq := RequestData{Method: http.MethodGet, URL: *url, Body: []byte(bodyData)}
		rs := ResponseData{Status: "200 OK", Body: []byte(bodyData), Request: rq}

		sb, ctx := strings.Builder{}, &Context{RequestEdited: edited, RequestRerouted: rerouted, ResponseEdited: edited}

		i := SummaryTextWriter("test-summary-text-writer", MatchAllRequests, MatchAllResponses, &sb)

		i.request(ctx, &rq)
		rqOutput := sb.String()
		sb.Reset()

		if !strings.Contains(rqOutput, rq.Method) {
			t.Fatalf("expected %v in request output. got %v", rq.Method, rqOutput)
		}

		if !strings.Contains(rqOutput, rq.URL.String()) {
			t.Fatalf("expected %v in request output. got %v", rq.URL, rqOutput)
		}

		assertOnBool := func(s, data string, expected bool) {
			if expected && !strings.Contains(s, data) {
				t.Fatalf("expected %v in request output. got %v", data, s)
				return
			}

			if !expected && strings.Contains(rqOutput, data) {
				t.Fatalf("did not expect %v in request output. got %v", data, s)
				return
			}
		}

		assertOnBool(rqOutput, "*rerouted", rerouted)
		assertOnBool(rqOutput, "*edited", edited)

		i.response(ctx, &rs)
		rsOutput := sb.String()
		sb.Reset()

		if !strings.Contains(rsOutput, rs.Status) {
			t.Fatalf("expected %v in request output. got %v", rq.Method, rqOutput)
		}

		if !strings.Contains(rsOutput, rq.Method) {
			t.Fatalf("expected %v in request output. got %v", rq.Method, rqOutput)
		}

		if !strings.Contains(rsOutput, rq.URL.String()) {
			t.Fatalf("expected %v in request output. got %v", rq.URL, rqOutput)
		}

		assertOnBool(rsOutput, "*edited", edited)
	}

	act(true, true)
	act(false, false)
	act(false, true)
	act(true, false)
}
