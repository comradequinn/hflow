package main

import (
	"comradequinn/hflow/log"
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
)

// Start begins the stub listening on the specified port
func main() {
	httpPort, httpsPort := 0, 0

	flag.IntVar(&httpPort, "port", 8081, "the port for the stub server to listen for http on")
	flag.IntVar(&httpsPort, "tls", 4444, "the port for the stub server to listen for https on")

	v := flag.Int("v", 0, "the verbosity of the log output")

	flag.Parse()

	log.SetVerbosity(*v)

	var err error
	tlsCert, err := tls.X509KeyPair(certPEM, privateKeyPEM)

	if err != nil {
		log.Fatalf(0, "error reading stub https server certificate or key [%v]", err)
	}

	go func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%v", httpPort), http.HandlerFunc(handler)); err != nil {
			log.Fatalf(0, "error starting stub http server [%v]", err)
		}
	}()

	svr := &http.Server{
		Addr:    fmt.Sprintf(":%v", httpsPort),
		Handler: http.HandlerFunc(handler),
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{tlsCert},
		},
	}

	go func() {
		if err := svr.ListenAndServeTLS("", ""); err != nil {
			log.Fatalf(0, "error starting stub https server [%v]", err)
		}
	}()

	log.Printf(0, "stub server started on port [%v] for http and [%v] for https", httpPort, httpsPort)

	select {}
}
