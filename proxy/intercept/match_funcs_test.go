package intercept

import (
	"net/url"
	"testing"
)

var (
	matchNoRequests  MatchRequestFunc  = func(_ *RequestData) (bool, error) { return false, nil }
	matchNoResponses MatchResponseFunc = func(_ *ResponseData) (bool, error) { return false, nil }
)

func TestMatchRequestURL(t *testing.T) {
	u, _ := url.Parse("http://some.domain.com:8080/somepath/?someqs=someval")
	rq := RequestData{URL: *u}

	test := func(t *testing.T, pattern string, expected bool) {
		matchFunc := MatchRequestURL(pattern)

		match, err := matchFunc(&rq)

		if match != expected || err != nil {
			t.Fatalf("expected matchrequesturl to return [%v] for pattern [%v], got [%v]", expected, pattern, match)
		}
	}

	t.Run("MatchesEmptyString", func(t *testing.T) { test(t, "", true) })
	t.Run("MatchesURL", func(t *testing.T) { test(t, "http://some.domain.com:8080/somepath/?someqs=someval", true) })
	t.Run("MatchesDomain", func(t *testing.T) { test(t, "some.domain.com", true) })
	t.Run("MatchesPort", func(t *testing.T) { test(t, ":8080", true) })
	t.Run("MatchesPath", func(t *testing.T) { test(t, "somepath", true) })
	t.Run("MatchesQueryStringKey", func(t *testing.T) { test(t, "someqs", true) })
	t.Run("MatchesQueryStringValue", func(t *testing.T) { test(t, "someval", true) })

	t.Run("NoMatchWhiteSpace", func(t *testing.T) { test(t, " ", false) })
	t.Run("NoMatchURL", func(t *testing.T) { test(t, "http://missing.domain.com:8080/somepath/?someqs=someval", false) })
	t.Run("NoMatchDomain", func(t *testing.T) { test(t, "missing.domain.com", false) })
	t.Run("NoMatchPort", func(t *testing.T) { test(t, ":8081", false) })
	t.Run("NoMatchPath", func(t *testing.T) { test(t, "missingpath", false) })
	t.Run("NoMatchQueryStringKey", func(t *testing.T) { test(t, "missingqs", false) })
	t.Run("NoMatchQueryStringValue", func(t *testing.T) { test(t, "missingval", false) })
}

func TestMatchResponseCode(t *testing.T) {
	rs := ResponseData{Status: "200 OK"}

	test := func(t *testing.T, pattern string, matchRequest bool, expected bool) {
		matchRqFunc := MatchAllRequests

		if !matchRequest {
			matchRqFunc = matchNoRequests
		}

		matchFunc := MatchRequestAndResponseStatus(pattern, matchRqFunc)

		match, err := matchFunc(&rs)

		if match != expected || err != nil {
			t.Fatalf("expected matchresponsestatus to return [%v] for pattern [%v] with matchrequest set to [%v], got [%v]", expected, pattern, matchRequest, match)
		}
	}

	t.Run("MatchEmptyString", func(t *testing.T) { test(t, "", true, true) })
	t.Run("MatchCode", func(t *testing.T) { test(t, "200", true, true) })
	t.Run("MatchPartialCode", func(t *testing.T) { test(t, "2", true, true) })
	t.Run("MatchStatus", func(t *testing.T) { test(t, "OK", true, true) })
	t.Run("MatchPartialStatus", func(t *testing.T) { test(t, "K", true, true) })

	t.Run("NoMatchCode", func(t *testing.T) { test(t, "202", true, false) })
	t.Run("NoMatchPartialCode", func(t *testing.T) { test(t, "4", true, false) })
	t.Run("NoMatchStatus", func(t *testing.T) { test(t, "BAD REQUEST", true, false) })
	t.Run("NoMatchPartialStatus", func(t *testing.T) { test(t, "BAD", true, false) })

	t.Run("NoMatchRequest", func(t *testing.T) { test(t, "", false, false) })
}
