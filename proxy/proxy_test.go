package proxy

import (
	"comradequinn/hflow/proxy/intercept"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestProxy(t *testing.T) {
	test := func(t *testing.T, icpt bool, clientTLS *tls.Config, proxyHandler http.HandlerFunc, newStubSvrFunc func(http.Handler) *httptest.Server) {
		expectedRsCode, expectedRsHeaderKey, expectedRsHeaderVal, expectedInterceptedRsHeaderKey, expectedRsBody := http.StatusOK, "rsqHdrK", "rsHdrV", "intRsHdrK", "rsBdy"

		expectedRqMethod, expectedRqPath, expectedRqQSKey, expectedRqQSVal, expectedRqHeaderKey, expectedRqHeaderVal, expectedRqBody, expectedInterceptedRqHeaderKey :=
			http.MethodPost, "test-path", "rqQSK", "rqQSV", "rqHdrK", "rqHdrV", "rqBdy", "intRqHdrK"

		proxy, client := httptest.NewServer(proxyHandler), http.Client{}
		proxyURL, _ := url.Parse(proxy.URL)

		client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyURL), TLSClientConfig: clientTLS}

		defer proxy.Close()

		var actualRqMethod, actualRqPath, actualRqQSVal, actualRqHeaderVal, actualRqBody, actualInterceptedRqHeaderVal string

		stub := newStubSvrFunc(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			b, _ := io.ReadAll(r.Body)
			actualRqMethod, actualRqPath, actualRqQSVal, actualRqHeaderVal, actualInterceptedRqHeaderVal, actualRqBody = r.Method, r.URL.Path, r.URL.Query().Get(expectedRqQSKey), r.Header.Get(expectedRqHeaderKey), r.Header.Get(expectedInterceptedRqHeaderKey), string(b)

			w.Header().Add(expectedRsHeaderKey, expectedRsHeaderVal)
			w.WriteHeader(expectedRsCode)
			w.Write([]byte(expectedRsBody))
		}))

		defer stub.Close()

		rq, _ := http.NewRequest(expectedRqMethod, fmt.Sprintf("%v/%v/?%v=%v", stub.URL, expectedRqPath, expectedRqQSKey, expectedRqQSVal), strings.NewReader(expectedRqBody))
		rq.Header.Set(expectedRqHeaderKey, expectedRqHeaderVal)

		if icpt {
			newIcpt := func(t string) intercept.Intercept {
				return intercept.New("test-intercept",
					intercept.MatchAllRequests,
					intercept.MatchAllResponses,
					func(_ *intercept.Context, r *intercept.RequestData) error {
						r.Header.Set(expectedInterceptedRqHeaderKey, r.Header.Get(expectedInterceptedRqHeaderKey)+" "+t)
						return nil
					},
					func(_ *intercept.Context, r *intercept.ResponseData) error {
						r.Header.Set(expectedInterceptedRsHeaderKey, r.Header.Get(expectedInterceptedRsHeaderKey)+" "+t)
						return nil
					},
				)
			}

			SetTrafficSummaryWriter(newIcpt("traffic"))
			SetCaptureFileWriter(newIcpt("capture"))
			SetRequestRerouter(newIcpt("rerouter"))
			SetTrafficEditor(newIcpt("editor"))

			defer SetTrafficSummaryWriter(intercept.Unset)
			defer SetRequestRerouter(intercept.Unset)
			defer SetTrafficEditor(intercept.Unset)
		}

		rs, err := client.Do(rq)

		if err != nil {
			t.Fatalf("expected no error proxying request, got [%v]", err)
		}

		assert := func(attr, exp, got string) {
			if got != exp {
				t.Fatalf("expected to %v [%v], got [%v]", attr, exp, got)
			}
		}

		assert("receive request method of", expectedRqMethod, actualRqMethod)
		assert("receive request path of", "/"+expectedRqPath+"/", actualRqPath)
		assert("receive request querystring value of", expectedRqQSVal, actualRqQSVal)
		assert("receive request header value of", expectedRqHeaderVal, actualRqHeaderVal)
		assert("receive request body of", expectedRqBody, actualRqBody)
		assert("receive response header value of", expectedRsHeaderVal, rs.Header.Get(expectedRsHeaderKey))

		b, _ := io.ReadAll(rs.Body)

		assert("receive response body of", expectedRsBody, string(b))

		if icpt {
			if !strings.Contains(actualInterceptedRqHeaderVal, "traffic") || !strings.Contains(actualInterceptedRqHeaderVal, "capture") || !strings.Contains(actualInterceptedRqHeaderVal, "editor") || !strings.Contains(actualInterceptedRqHeaderVal, "rerouter") {
				t.Fatalf("expected to find all request intercepts have marked header. got %v", actualInterceptedRqHeaderVal)
			}

			rsHeader := rs.Header.Get(expectedInterceptedRsHeaderKey)

			if !strings.Contains(rsHeader, "traffic") || !strings.Contains(rsHeader, "capture") || !strings.Contains(rsHeader, "editor") {
				t.Fatalf("expected to find all response intercepts have marked header. got %v", rsHeader)
			}
		}
	}

	tlsCfg := &tls.Config{InsecureSkipVerify: true}

	t.Run("HTTP", func(t *testing.T) { test(t, false, nil, HTTPHandler(), httptest.NewServer) })
	t.Run("InterceptedHTTP", func(t *testing.T) { test(t, true, nil, HTTPHandler(), httptest.NewServer) })
	t.Run("HTTPS", func(t *testing.T) { test(t, false, tlsCfg, HTTPSHandler(), httptest.NewTLSServer) })
	t.Run("InterceptedHTTPS", func(t *testing.T) { test(t, true, tlsCfg, HTTPSHandler(), httptest.NewTLSServer) })
}
