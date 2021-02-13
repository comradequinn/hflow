// Package intercept provides functions for creating and configuring intercepts with package proxy
package intercept

import (
	"comradequinn/hflow/log"
	"fmt"
	"net/http"
)

// Intercept describes an action to be taken on a http exchange
// and the conditions to be met in order for that action to be applied
type Intercept struct {
	label    string
	request  RequestFunc
	response ResponseFunc
}

// Context maintains data about the effect the intercepts applied have had on the traffic
type Context struct {
	RequestRerouted bool
	RequestEdited   bool
	ResponseEdited  bool
}

var (
	// Unset is an interceptor that does not act on any request or response
	Unset = New("unset", MatchNoRequests, MatchNoResponses, nil, nil)
	// NoOp func is an alias for a function that takes no action
	NoOpFunc = func() {}
)

func NewContext() *Context { return &Context{} }

// New returns a new Intercept based on the passed arguments
func New(label string, mrq MatchRequestFunc, mrs MatchResponseFunc, rqf RequestFunc, rsf ResponseFunc) Intercept {
	requestFunc := func(c *Context, rd *RequestData) error {
		if mrq == nil || rqf == nil {
			return nil
		}

		match, err := mrq(rd)

		if !match || err != nil {
			return err
		}

		log.Printf(2, "applying intercept labelled [%v] to request for [%v]", label, rd.URL.String())

		return rqf(c, rd)
	}

	responseFunc := func(c *Context, r *ResponseData) error {
		if mrs == nil || rsf == nil {
			return nil
		}

		match, err := mrs(r)

		if !match || err != nil {
			return err
		}

		log.Printf(2, "applying intercept labelled [%v] to response to [%v]", label, r.Request.URL.String())

		return rsf(c, r)
	}

	return Intercept{label: label, request: requestFunc, response: responseFunc}
}

// Label returns the Label of the Intercept
func (i Intercept) Label() string {
	return i.label
}

// Request returns a new *http.Request which is the result of applying any matching intercepts to hr
func Request(hr *http.Request, trafficSummaryWriter, captureFileWriter, requestRerouter, trafficEditor Intercept) (*Context, *http.Request, error) {
	log.Printf(3, "intercepting request for [%v]", hr.URL.String())

	c := NewContext()
	r, err := readRequestData(hr)

	if err != nil {
		return nil, nil, fmt.Errorf("error reading request data from https request to remote client [%v]. [%v]", hr.URL.String(), err)
	}

	if err = trafficEditor.request(c, &r); err != nil {
		return nil, nil, fmt.Errorf("error applying traffic editor labelled [%v] to request for [%v]: [%v]", trafficEditor.label, r.URL.String(), err)
	}

	if err = requestRerouter.request(c, &r); err != nil {
		return nil, nil, fmt.Errorf("error applying request rerouter labelled [%v] to request for [%v]: [%v]", requestRerouter.label, r.URL.String(), err)
	}

	if err = trafficSummaryWriter.request(c, &r); err != nil {
		return nil, nil, fmt.Errorf("error applying traffic summary writer labelled [%v] to request for [%v]: [%v]", trafficSummaryWriter.label, r.URL.String(), err)
	}

	if err = captureFileWriter.request(c, &r); err != nil {
		return nil, nil, fmt.Errorf("error applying capture file writer labelled [%v] to request for [%v]: [%v]", captureFileWriter.label, r.URL.String(), err)
	}

	nhr, err := r.toHTTP()
	return c, nhr, err
}

// Response returns a new *http.Response which is the result of applying any matching intercepts to hrs
func Response(c *Context, hr *http.Request, hrs *http.Response, trafficSummaryWriter, captureFileWriter, trafficEditor Intercept) (*http.Response, error) {
	log.Printf(3, "intercepting response for [%v]", hr.URL.String())

	rs, err := readResponseData(hrs)

	if err != nil {
		return nil, fmt.Errorf("error reading response data from https response to [%v]: [%v]", hr.URL.String(), err)
	}

	if err = trafficEditor.response(c, &rs); err != nil {
		return nil, fmt.Errorf("error applying traffic editor labelled [%v] to response to [%v]: [%v]", trafficEditor.Label(), hr.URL.String(), err)
	}

	if err = trafficSummaryWriter.response(c, &rs); err != nil {
		return nil, fmt.Errorf("error applying traffic summary writer labelled [%v] to response to [%v]: [%v]", trafficSummaryWriter.Label(), hr.URL.String(), err)
	}

	if err = captureFileWriter.response(c, &rs); err != nil {
		return nil, fmt.Errorf("error applying capture file writer labelled [%v] to response to [%v]: [%v]", captureFileWriter.Label(), hr.URL.String(), err)
	}

	return rs.toHTTP()
}
