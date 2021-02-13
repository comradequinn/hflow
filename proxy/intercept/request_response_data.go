package intercept

import (
	"comradequinn/hflow/log"
	"comradequinn/hflow/proxy/internal/codec"
	"comradequinn/hflow/proxy/internal/copy"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// RequestData represents a http.Request being currently processed by proxy
type RequestData struct {
	URL         url.URL
	Method      string
	Header      http.Header
	InitialHost string
	Protocol    string
	Body        []byte
}

func readRequestData(hr *http.Request) (RequestData, error) {
	if hr == nil {
		return RequestData{}, fmt.Errorf("cannot read nil request")
	}

	r := RequestData{Header: http.Header{}, Method: hr.Method}

	r.Protocol = hr.Proto
	r.URL = *hr.URL
	r.InitialHost = hr.Host

	for k, v := range hr.Header {
		hv := strings.Join(v, " ")
		r.Header.Set(k, hv)
	}

	var err error

	if r.Body, err = copy.CloserToBytes(&hr.Body); err != nil {
		return RequestData{}, fmt.Errorf("unable to read request body: [%v]", err)
	}

	return r, nil
}

func (r RequestData) toHTTP() (*http.Request, error) {
	nr, err := http.NewRequest(r.Method, r.URL.String(), copy.BytesToCloser(r.Body))

	if err == nil {
		copy.Header(r.Header, nr.Header)
	}

	nr.Host = r.InitialHost

	return nr, err
}

func (r RequestData) String() string {
	return r.string(false, true, 0)
}

func (r RequestData) string(appendHost bool, binary bool, limit int) string {
	sb := strings.Builder{}

	sb.WriteString(fmt.Sprintf("%v %v?%v %v\n", r.Method, r.URL.Path, r.URL.RawQuery, r.Protocol))

	if appendHost && r.Header.Get("Host") == "" {
		sb.WriteString("Host: " + r.InitialHost + "\n")
	}

	writeString(r.Header, r.Body, binary, limit, &sb)

	return sb.String()
}

// ResponseData represents a http.Response being currently processed by proxy
type ResponseData struct {
	Header     http.Header
	Body       []byte
	Status     string
	StatusCode int
	Proto      string
	ProtoMajor int
	ProtoMinor int
	Request    RequestData
	TLS        *tls.ConnectionState
}

func readResponseData(hr *http.Response) (ResponseData, error) {
	r := ResponseData{Header: http.Header{}}

	copy.Header(hr.Header, r.Header)

	r.Header.Set("Via", "hflow")

	var err error

	if r.Request, err = readRequestData(hr.Request); err != nil {
		return ResponseData{}, err
	}

	r.Status, r.StatusCode, r.Proto, r.ProtoMajor, r.ProtoMinor, r.TLS =
		hr.Status, hr.StatusCode, hr.Proto, hr.ProtoMajor, hr.ProtoMinor, hr.TLS

	if r.Body, err = copy.CloserToBytes(&hr.Body); err != nil {
		return ResponseData{}, fmt.Errorf("unable to read response body: [%v]", err)
	}

	ct := r.Header.Get("Content-Encoding")

	if codec.Supported(ct) {
		log.Printf(1, "decoding response body from [%v] using scheme from content-encoding header [%v]", ct, hr.Request.URL.String())

		if r.Body, err = codec.Decode(ct, r.Body); err != nil {
			return ResponseData{}, fmt.Errorf("unable to decode response body using scheme from content-encoding header [%v]: [%v]", ct, err)
		}
	}

	return r, nil
}

func (r ResponseData) toHTTP() (*http.Response, error) {
	hr := http.Response{Header: http.Header{}}

	copy.Header(r.Header, hr.Header)

	var err error

	if hr.Request, err = r.Request.toHTTP(); err != nil {
		return nil, err
	}

	hr.Status, hr.StatusCode, hr.Proto, hr.ProtoMajor, hr.ProtoMinor, hr.TLS =
		r.Status, r.StatusCode, r.Proto, r.ProtoMajor, r.ProtoMinor, r.TLS

	ct := r.Header.Get("Content-Encoding")

	if codec.Supported(ct) {
		log.Printf(1, "encoding response body from [%v] using scheme from content-encoding header [%v]", ct, hr.Request.URL.String())

		if r.Body, err = codec.Encode(ct, r.Body); err != nil {
			return nil, fmt.Errorf("unable to encode response body using scheme from content-encoding header [%v]: [%v]", ct, err)
		}
	}

	hr.Body, hr.ContentLength = copy.BytesToCloser(r.Body), int64(len(r.Body))

	return &hr, nil
}

func (r ResponseData) String() string {
	return r.string(true, 0)
}

func (r ResponseData) string(binary bool, limit int) string {
	sb := strings.Builder{}

	sb.WriteString(fmt.Sprintf("%v %v\n", r.Request.Protocol, r.Status))
	writeString(r.Header, r.Body, binary, limit, &sb)

	return sb.String()
}

func writeString(h http.Header, b []byte, binary bool, limit int, sb *strings.Builder) {
	contentType, textContentTypes := "", []string{"text/", "/json", "xml", "/javascript", "urlencoded"}

	for k, vs := range h {
		if k == "Content-Type" {
			contentType = strings.Split(vs[0], ";")[0]
		}

		sb.WriteString(fmt.Sprintf("%v: ", k))
		sb.WriteString(strings.Join(vs, ","))
		sb.WriteString("\n")
	}

	if limit > 0 && len(b) > limit {
		b = b[:limit]
	}

	body := string(b)

	if contentType != "" {
		text := false

		for _, tct := range textContentTypes {
			if strings.Contains(contentType, tct) {
				text = true
				break
			}
		}

		if !text && !binary {
			body = "[binary data]"
		}
	}

	if len(body) > 0 {
		sb.WriteString("\n" + body + "\n")
	} else {
		sb.WriteString("\n")
	}
}
