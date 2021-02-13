package proxy

import (
	"bufio"
	"comradequinn/hflow/cert"
	"comradequinn/hflow/log"
	"comradequinn/hflow/proxy/intercept"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"
)

// HTTPSHandler is is a http.HandlerFunc that acts as HTTPS Proxy
func HTTPSHandler() http.HandlerFunc {
	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse },
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	return func(connectRs http.ResponseWriter, connectRq *http.Request) {
		if connectRq.Method != http.MethodConnect {
			connectRs.WriteHeader(http.StatusMethodNotAllowed)
			log.Printf(3, "rejected request for [%v] on [%v] via unsupported method [%v]", connectRq.URL.String(), connectRq.Host, connectRq.Method)
			return
		}

		hj, ok := connectRs.(http.Hijacker)

		if !ok {
			log.Printf(0, "http connect request for [%v] not hijackable", connectRq.Host)
			return
		}

		tcpConn, _, err := hj.Hijack()

		if err != nil {
			log.Printf(0, "error hijacking http connect request for [%v]. [%v]", connectRq.Host, err)
			return
		}

		fmt.Fprintf(tcpConn, "HTTP/1.1 200 Connection Established\r\n\r\n")

		tlsConn := tls.Server(tcpConn, &tls.Config{GetCertificate: cert.For})

		log.Printf(3, "tunneling to [%v] on behalf of [%v]", connectRq.Host, tcpConn.RemoteAddr())

		go func() {
			defer func() {
				tlsConn.Close()
				log.Printf(3, "closed tunnel to [%v] on behalf of [%v]", connectRq.Host, tcpConn.RemoteAddr())

				if err := recover(); err != nil {
					log.Printf(0, "panic while tunneling from remote client [%v] to remote host [%v]. [%+v]", connectRq.RemoteAddr, connectRq.Host, err)
				}
			}()

			eof := func(buffer *bufio.Reader) bool {
				log.Printf(3, "waiting to receive from remote client [%v]", connectRq.RemoteAddr)

				if err := tcpConn.SetReadDeadline(time.Now().Add(time.Second * 60)); err != nil {
					log.Printf(0, "error setting read deadline on connection with remote client [%v]. [%v]", connectRq.RemoteAddr, err)
					return true
				}

				if _, err := buffer.Peek(1); err != nil {
					log.Printf(3, "unable to read from connection with remote client [%v]. [%v]", connectRq.RemoteAddr, err)
					return true
				}

				log.Printf(3, "receiving from remote client [%v]", connectRq.RemoteAddr)

				return false
			}

			buffer := bufio.NewReader(tlsConn)

			for !eof(buffer) {
				rq, err := http.ReadRequest(buffer)

				if err != nil {
					log.Printf(0, "error reading https request from remote client [%v]. [%v]", connectRq.RemoteAddr, err)
					return
				}

				log.Printf(1, "<<< received proxy request for [%v] on host [%v]", rq.URL.String(), rq.Host)

				rq.RequestURI, rq.URL.Scheme, rq.URL.Host = "", "https", connectRq.Host

				var ctx *intercept.Context

				ctx, rq, err = intercept.Request(rq, TrafficSummaryWriter(), CaptureFileWriter(), RequestRerouter(), TrafficEditor())

				if err != nil {
					log.Printf(0, "error intercepting https request from remote client [%v]. [%v]", connectRq.RemoteAddr, err)
					return
				}

				log.Printf(3, ">>> requesting [%v] from host [%v]", rq.URL.String(), rq.Host)

				rs, err := client.Do(rq)

				if err != nil {
					log.Printf(0, "error proxying request for [%v] on host [%v]: [%v]", rq.URL.String(), rq.Host, err)
					return
				}

				log.Printf(3, "<<< received [%v] in response to [%v] on [%v]", rs.StatusCode, rq.URL.String(), rq.Host)

				rs, err = intercept.Response(ctx, rq, rs, TrafficSummaryWriter(), CaptureFileWriter(), TrafficEditor())

				if err != nil {
					log.Printf(0, "error intercepting response to [%v] on host [%v]: [%v]", rq.URL.String(), rq.Host, err)
					return
				}

				rs.Write(tlsConn)

				log.Printf(2, ">>> wrote proxy response for [%v] on [%v]", rq.URL.String(), rq.Host)
			}
		}()
	}
}
