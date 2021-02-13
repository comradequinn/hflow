package main

import (
	"comradequinn/hflow/log"
	"fmt"
	"io"
	"net/http"
)

func handler(rw http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf(0, "recovered from panic: %v", err)
		}
	}()

	rw.Header().Set("X-From", "hflow-stub")

	if r.URL.Path != "/echo/" {
		rw.WriteHeader(http.StatusNotFound)
		log.Printf(0, "x<< rejected request for unsupported path [%v]", r.URL.String())
		return
	}

	log.Printf(1, ">>> received request for [%v] at host [%v]", r.URL.String(), r.Host)
	fmt.Fprintf(rw, "qs-data: %v\n", r.URL.RawQuery)

	var (
		b   []byte
		err error
	)

	if b, err = io.ReadAll(r.Body); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Printf(0, ">>x unable to read body of request for [%v]: [%v]", r.URL.String(), err)
		return
	}

	fmt.Fprintf(rw, "body-data: %v\n", string(b))
	log.Printf(2, "<<< wrote [%v] bytes in response to request for [%v]", len(b), r.URL.String())
}
