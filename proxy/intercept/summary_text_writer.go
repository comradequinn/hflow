package intercept

import (
	"comradequinn/hflow/log"
	"comradequinn/hflow/proxy/store"
	"fmt"
	"io"
	"strings"
)

var storeInsert = store.Insert

// SummaryTextWriter writes
// * summarised request traffic to the specified io.Writer where the mrq matches the request
// * summarised response traffic to the specified io.Writer where mrs matches the response
func SummaryTextWriter(label string, mrq MatchRequestFunc, mrs MatchResponseFunc, w io.Writer) Intercept {
	return New(label, mrq, mrs,
		func(c *Context, r *RequestData) error {
			ctxText := ""

			if c.RequestEdited {
				ctxText = " *edited"
			}

			if c.RequestRerouted {
				ctxText += " *rerouted"
			}

			if ctxText != "" {
				ctxText = " (" + strings.TrimSpace(ctxText) + ")"
			}

			id := storeInsert(store.Record{Data: r.String()})

			url := r.URL.String()

			if len(url) > 100 {
				url = url[:100]
			}

			if _, err := w.Write([]byte(fmt.Sprintf("\r%-6d>> %v%v %v\n", id, r.Method, ctxText, url))); err != nil {
				return fmt.Errorf("error in intercept labelled [%v] for request [%v]. unable to write request summary data: [%v]", label, r.URL.String(), err)
			}

			log.Printf(2, "intercept labelled [%v] wrote request summary data for [%v]", label, r.URL.String())

			return nil
		},
		func(c *Context, r *ResponseData) error {
			ctxText := ""

			if c.ResponseEdited {
				ctxText = " (*edited)"
			}

			id := storeInsert(store.Record{Data: r.String()})

			url := r.Request.URL.String()

			if len(url) > 100 {
				url = url[:100]
			}

			if _, err := w.Write([]byte(fmt.Sprintf("\r%-7d< %v%v %v bytes in body (source: %v %v)\n", id, r.Status, ctxText, len(r.Body), r.Request.Method, url))); err != nil {
				return fmt.Errorf("error in intercept labelled [%v] for request [%v]. unable to write request summary data: [%v]", label, r.Request.URL.String(), err)
			}

			log.Printf(2, "intercept labelled [%v] wrote response summary data for [%v]", label, r.Request.URL.String())

			return nil
		},
	)
}
