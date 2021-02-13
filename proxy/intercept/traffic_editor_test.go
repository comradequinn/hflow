package intercept

import (
	"comradequinn/hflow/proxy/internal/copy"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"
)

func TestTrafficEditorForGetRq(t *testing.T) {
	newQSData, oldQSData := "new-qs-data", "old-qs-data"

	expectedHTTPRq := httptest.NewRequest(http.MethodGet, "http://www.test.com/echo/?data="+newQSData, nil)

	inTerminal = func() bool { return true }
	runCmd = func(cmd *exec.Cmd) error { // simulate the behaviour of the editor by writing changes to the passed file
		data, _ := readRequestData(expectedHTTPRq)

		if err := os.WriteFile(cmd.Args[1], []byte(data.String()), 0666); err != nil {
			t.Fatalf("expected no error. got %v", err)
		}

		return nil
	}

	rqd, _ := readRequestData(httptest.NewRequest(http.MethodGet, "http://www.test.com/echo/?data="+oldQSData, nil))

	if err := TrafficEditor("test-traffic-editor", MatchAllRequests, matchNoResponses, func() (release func()) { return func() {} }, NoOpFunc).request(NewContext(), &rqd); err != nil {
		t.Fatalf("expected no error during request processing. got %v", err)
	}

	if string(rqd.URL.Query().Get("data")) != newQSData {
		t.Fatalf("request body was not updated. expected: %v, got: %v", newQSData, string(rqd.Body))
	}
}

func TestTrafficEditorForPostRqAndRs(t *testing.T) {
	stdIOCaptured, stdIOReleased, newBodyData, oldBodyData := false, false, "new-body-data", "old-body-data"

	expectedHTTPRq := httptest.NewRequest(http.MethodPost, "http://www.test.com/echo/?data=some-qs-data", strings.NewReader(newBodyData))
	expectedHTTPRq.Header.Add("Content-Length", strconv.Itoa(len(newBodyData)))

	expectedHTTPRs := &http.Response{Header: http.Header{}}
	expectedHTTPRs.Request = expectedHTTPRq
	expectedHTTPRs.Status = "200 OK"
	expectedHTTPRs.StatusCode = 200
	expectedHTTPRs.Header.Add("Content-Length", strconv.Itoa(len(newBodyData)))
	expectedHTTPRs.Body = copy.BytesToCloser([]byte(newBodyData))

	ctx := NewContext()

	runCmdForRequest := true

	inTerminal = func() bool { return true }
	runCmd = func(cmd *exec.Cmd) error { // simulate the behaviour of the editor by writing changes to the passed file
		if !stdIOCaptured {
			t.Fatalf("stdio was not captured during editing process")
		}

		var data fmt.Stringer

		if runCmdForRequest {
			data, _ = readRequestData(expectedHTTPRq)
		} else {
			data, _ = readResponseData(expectedHTTPRs)
		}

		if err := os.WriteFile(cmd.Args[1], []byte(data.String()), 0666); err != nil {
			t.Fatalf("expected no error. got %v", err)
		}

		runCmdForRequest = !runCmdForRequest // this gets called first for a request, then again for the response; this bool tracks that state
		return nil
	}

	rqd, _ := readRequestData(expectedHTTPRq)
	rqd.Body = []byte(oldBodyData)

	ted := TrafficEditor("test-traffic-editor", MatchAllRequests, MatchAllResponses, func() (release func()) { stdIOCaptured = true; return func() { stdIOReleased = true } }, NoOpFunc)

	err := ted.request(ctx, &rqd)

	if err != nil {
		t.Fatalf("expected no error during request processing. got %v", err)
	}

	if !stdIOCaptured || !stdIOReleased {
		t.Fatalf("stdio was not captured and released during request processing. expected: true, true. got: %v, %v", stdIOCaptured, stdIOReleased)
	}

	if string(rqd.Body) != newBodyData {
		t.Fatalf("request body was not updated. expected: %v, got: %v", newBodyData, string(rqd.Body))
	}

	if !ctx.RequestEdited {
		t.Fatalf("context was not set as request edited")
	}

	rsd, _ := readResponseData(expectedHTTPRs)
	rsd.Body = []byte(oldBodyData)
	stdIOCaptured, stdIOReleased = false, false // reset these as they will true from the request test

	err = ted.response(ctx, &rsd)

	if err != nil {
		t.Fatalf("expected no error during request processing. got %v", err)
	}

	if !stdIOCaptured || !stdIOReleased {
		t.Fatalf("stdio was not captured and released during request processing. expected: true, true. got: %v, %v", stdIOCaptured, stdIOReleased)
	}

	if string(rqd.Body) != newBodyData {
		t.Fatalf("response body was not updated. expected: %v, got: %v", newBodyData, string(rqd.Body))
	}

	if !ctx.ResponseEdited {
		t.Fatalf("context was not set as response edited")
	}
}
