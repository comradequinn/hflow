package intercept

import (
	"comradequinn/hflow/log"
	"fmt"
	"io"
)

// TextWriter writes
// * request traffic to the specified io.Writer where the mrq matches the request
// * response traffic to the specified io.Writer where mrs matches the response
//
// Unless binary is set to true, only text-based mime-type bodies are written to
// If limit is greater than or equal to 0, then text response body writes are capped at that number of bytes
func TextWriter(label string, mrq MatchRequestFunc, mrs MatchResponseFunc, binary bool, limit int, w io.Writer) Intercept {
	return New(label, mrq, mrs,
		func(c *Context, r *RequestData) error {
			ctxText := ""

			if c.RequestEdited {
				ctxText = " *edited"
			}

			if c.RequestRerouted {
				ctxText += " *rerouted"
			}

			if _, err := w.Write([]byte(">>>" + ctxText + "\n" + r.string(true, binary, limit))); err != nil {
				return fmt.Errorf("error in intercept labelled [%v] for request [%v]. unable to write request data: [%v]", label, r.URL.String(), err)
			}

			log.Printf(2, "intercept labelled [%v] wrote request data for [%v]", label, r.URL.String())

			return nil
		},
		func(c *Context, r *ResponseData) error {
			ctxText := ""

			if c.ResponseEdited {
				ctxText = " *edited"
			}

			if _, err := w.Write([]byte(fmt.Sprintf("\n<<<%v source-request: %v %v\n", ctxText, r.Request.Method, r.Request.URL.String()) + r.string(binary, limit))); err != nil {
				return fmt.Errorf("error in intercept labelled [%v] for request [%v]. unable to write response data: [%v]", label, r.Request.URL.String(), err)
			}

			log.Printf(2, "intercept labelled [%v] wrote response data for [%v]", label, r.Request.URL.String())

			return nil
		},
	)
}
