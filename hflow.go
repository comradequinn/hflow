package main

import (
	"comradequinn/hflow/cli"
	"comradequinn/hflow/cli/tty"
	"comradequinn/hflow/log"
	"comradequinn/hflow/proxy"
	"comradequinn/hflow/proxy/store"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	logFile, err := os.Create("hflow.log")

	if err != nil {
		log.Fatalf(0, "error creating log file: %v", err)
	}

	defer logFile.Close()

	log.SetWriter(logFile)

	execMode := flag.String("e", "c", "the execution mode; (e)export-ca or (c)li. eg '-e=c'")
	httpPort := flag.Int("p", 8080, "the port to proxy http over")
	httpsPort := flag.Int("t", 4443, "the port to proxy https/tls over")
	storeCapacity := flag.Int("c", 10000, "the number of traffic records (requests and responses) to hold in memory for viewing on request when traffic captures are being written to the terminal")

	trafficCapture := cli.TrafficCapture{}
	flag.StringVar(&trafficCapture.File, "f", "", "create a traffic capture (f)ile which, unless otherwise filtered, contains all traffic from the session. use -u and -s to apply optional filters. eg '-f=hflow.capture'")
	flag.StringVar(&trafficCapture.URLFilter, "u", "", "only write requests to the capture file that contain the url-pattern specified as this argument. eg '-u=www.example.com'. has no effect is -f is not also specified")
	flag.StringVar(&trafficCapture.StatusFilter, "s", "", "only write responses to the capture file that that contain the status-pattern specified as this argument. . eg '-s=500'. has no effect is -f is not also specified")
	flag.BoolVar(&trafficCapture.Binary, "b", false, "if set, writes non-text response bodies to the capture file. has no effect is -f is not also specified")
	flag.IntVar(&trafficCapture.Limit, "l", -1, "limit text response bodies written to the capture file to the specified byte count when sending to writers, -1 is no limit. eg '-l=250'. has no effect is -f is not also specified")

	verbosity := flag.Int("v", 0, "the verbosity of the log output. 0 is least verbose. 3 is most verbose. eg '-v=1'")

	flag.Parse()

	log.SetVerbosity(*verbosity)

	store.Init(*storeCapacity)

	ttyIO := tty.New(os.Stdin, time.Millisecond*100, os.Stdout, os.Stderr)

	log.Printf(1, "received command line args [%v]", os.Args[1:])

	switch *execMode {
	case "e":
		if err := cli.ExportCA(ttyIO.Out()); err != nil {
			log.Fatalf(0, "stdio starting in ca export execution mode: [%v]", err)
		}
		return
	case "c":
		shutdownFunc, err := cli.Interactive(ttyIO, trafficCapture, *httpPort, *httpsPort)
		defer shutdownFunc()

		if err != nil {
			log.Fatalf(0, "error starting in daemon execution mode: [%v]", err)
		}
	default:
		log.Fatalf(0, "unrecognised command mode [%v]", *execMode)
	}

	startSvr := func(name string, port int, handler http.Handler) {
		svr := http.Server{
			Addr:    fmt.Sprintf(":%v", port),
			Handler: handler,
		}

		go func() {
			if err := svr.ListenAndServe(); err != nil {
				log.Fatalf(0, "error starting proxy server on port [%v]: [%v]", port, err)
			}
		}()

		log.Printf(0, "%v started on port [%v]", name, port)
	}

	startSvr("http proxy server", *httpPort, proxy.HTTPHandler())
	startSvr("https proxy server", *httpsPort, proxy.HTTPSHandler())

	sigInt := make(chan os.Signal, 1)
	signal.Notify(sigInt, os.Interrupt)
	<-sigInt

	log.Printf(0, "shutting down")
}
