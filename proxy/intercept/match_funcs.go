package intercept

import "strings"

// MatchRequestFunc describes a func which should return (true, nil)
// when passed a *Request in order for an action to be applied
type MatchRequestFunc func(*RequestData) (bool, error)

// MatchResponseFunc describes a func which should return (true, nil) when passed
// the *Request & *Response in order for an action to be applied
type MatchResponseFunc func(*ResponseData) (bool, error)

// RequestFunc describes a func that performs an action on the specified request.
type RequestFunc func(*Context, *RequestData) error

// ResponseFunc describes a func that performs an action on the specified response.
type ResponseFunc func(*Context, *ResponseData) error

var (
	// MatchAllRequests returns a MatchRequestFunc matches any request
	MatchAllRequests MatchRequestFunc = func(r *RequestData) (bool, error) { return true, nil }

	// MatchNoRequests returns a MatchRequestFunc that matches no requests
	MatchNoRequests MatchRequestFunc = func(r *RequestData) (bool, error) { return false, nil }

	// MatchAllResponses returns a MatchResponseFunc that matches any response
	MatchAllResponses MatchResponseFunc = func(rs *ResponseData) (bool, error) { return true, nil }

	// MatchNoResponses returns a MatchResponseFunc that matches no responses
	MatchNoResponses MatchResponseFunc = func(rs *ResponseData) (bool, error) { return false, nil }

	// MatchRequestURL returns a MatchRequestFunc that matches all requests with a URL that contains s, or all requests if s is empty
	MatchRequestURL = func(s string) MatchRequestFunc {
		if s == "" {
			return MatchAllRequests
		}

		return func(r *RequestData) (bool, error) {
			return strings.Contains(r.URL.String(), s), nil
		}
	}

	// MatchRequestHost returns a MatchRequestFunc that matches all requests with a host of s
	MatchRequestHost = func(s string) MatchRequestFunc {
		return func(r *RequestData) (bool, error) {
			return strings.EqualFold(r.URL.Host, s), nil
		}
	}

	// MatchSourceRequestURL returns a MatchResponseFunc that matches all responses
	// with a source request that matches MatchRequestFunc
	MatchSourceRequestURL = func(mrq MatchRequestFunc) MatchResponseFunc {
		return func(r *ResponseData) (bool, error) {
			return mrq(&r.Request)
		}
	}

	// MatchRequestAndResponseStatus returns a MatchResponseFunc that matches all responses with a status
	// that contains s and a source request that matches MatchRequestFunc
	MatchRequestAndResponseStatus = func(s string, mrq MatchRequestFunc) MatchResponseFunc {
		rsMatchFunc := func(string, *ResponseData) bool { return true }

		if s != "" {
			rsMatchFunc = func(s string, rs *ResponseData) bool { return strings.Contains(rs.Status, s) }
		}

		return func(r *ResponseData) (bool, error) {
			rqMatch, err := mrq(&r.Request)

			if !rqMatch || err != nil {
				return rqMatch, err
			}

			return rsMatchFunc(s, r), nil
		}
	}
)
