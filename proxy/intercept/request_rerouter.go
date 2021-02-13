package intercept

import (
	"comradequinn/hflow/log"
)

// RequestRerouter switches the destination of matching requests to the specified host
func RequestRerouter(label string, mrq MatchRequestFunc, targetHost string) Intercept {
	return New(label, mrq, nil,
		func(c *Context, rq *RequestData) error {
			rq.URL.Host = targetHost
			c.RequestRerouted = true

			log.Printf(1, "intercept labelled [%v] updated target host for [%v] to [%v]", label, rq.URL.String(), targetHost)
			return nil
		}, nil,
	)
}
