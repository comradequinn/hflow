package intercept

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRequestRerouter(t *testing.T) {
	targetHost, originalHost, originalBody, originalMethod := "new.newhost.dev", "www.originalhost.com", "body-data", http.MethodPost
	rq, _ := readRequestData(httptest.NewRequest(originalMethod, "http://"+originalHost+"/api/resource/?qs=1", strings.NewReader(originalBody)))
	rq.Header.Set("Host", originalHost)
	originalURL := rq.URL.String()

	ctx := NewContext()
	if err := RequestRerouter("test-reroute", MatchAllRequests, targetHost).request(ctx, &rq); err != nil {
		t.Fatalf("expected no error. got %v", err)
	}

	if rq.Method != originalMethod {
		t.Fatalf("expected no method change. got %v", rq.Method)
	}

	if string(rq.Body) != originalBody {
		t.Fatalf("expected no body change. got %v", rq.Body)
	}

	if string(rq.URL.String()) != strings.Replace(originalURL, originalHost, targetHost, -1) {
		t.Fatalf("host was not updated to target host. got %v", rq.URL.String())
	}

	if !ctx.RequestRerouted {
		t.Fatalf("context was not set as rerouted")
	}
}
