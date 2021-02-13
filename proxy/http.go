package proxy

import (
	"comradequinn/hflow/log"
	"comradequinn/hflow/proxy/intercept"
	"comradequinn/hflow/proxy/internal/copy"
	"net/http"
	"runtime/debug"
)

// HTTPHandler is is a http.HandlerFunc that acts as HTTP Proxy
func HTTPHandler() http.HandlerFunc {
	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse },
	}

	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf(0, "panic recovered while proxying request: [%v] [%s]", err, debug.Stack())
				return
			}
		}()

		log.Printf(1, "<<< received proxy request for [%v] on host [%v]", r.URL.String(), r.Host)

		var (
			ctx *intercept.Context
			err error
		)

		ctx, r, err = intercept.Request(r, TrafficSummaryWriter(), CaptureFileWriter(), RequestRerouter(), TrafficEditor())

		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			log.Printf(0, "error intercepting request: [%v]", err)
			return
		}

		log.Printf(2, ">>> requesting [%v] from host [%v]", r.URL.String(), r.Host)

		rs, err := client.Do(r)

		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			log.Printf(0, "error proxying request for [%v] on host [%v]: [%v]", r.URL.String(), r.Host, err)
			return
		}

		log.Printf(2, "<<< received [%v] in response to [%v] on [%v]", rs.StatusCode, r.URL.String(), r.Host)

		rs, err = intercept.Response(ctx, r, rs, TrafficSummaryWriter(), CaptureFileWriter(), TrafficEditor())

		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			log.Printf(0, "error intercepting response to [%v] on host [%v]: [%v]", r.URL.String(), r.Host, err)
			return
		}

		copy.Header(rs.Header, w.Header())

		b, err := copy.CloserToBytes(&rs.Body)

		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			log.Printf(0, "error reading response body from [%v]: [%v]", r.URL.String(), err)
			return
		}

		if _, err = w.Write(b); err != nil {
			log.Printf(0, "error writing response body from [%v] on [%v] to hflow client: [%v]", r.URL.String(), r.Host, err)
			return
		}

		log.Printf(2, ">>> wrote proxy response for [%v]", r.URL.String())
	}
}
