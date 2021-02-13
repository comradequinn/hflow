package cli

import (
	"bufio"
	"comradequinn/hflow/cli/cmd"
	"comradequinn/hflow/cli/tty"
	"comradequinn/hflow/log"
	"comradequinn/hflow/proxy"
	"comradequinn/hflow/proxy/intercept"
	"comradequinn/hflow/syncio"
	"io"
	"os"
	"strings"
)

type TrafficCapture struct {
	File         string
	URLFilter    string
	StatusFilter string
	Binary       bool
	Limit        int
}

// Interactive starts an interactive CLI session. It returns a closeFunc to be executed when the process is about to terminate
func Interactive(ttyIO tty.IO, tc TrafficCapture, httpPort, httpsPort int) (func(), error) {
	closeFunc := func() {}

	if tc.File != "" {
		f, err := os.Create(tc.File)

		if err != nil {
			log.Fatalf(0, "error creating specified capture file [%v]: [%v]", tc.File, err)
		}

		mrq := intercept.MatchRequestURL(tc.URLFilter)
		proxy.SetCaptureFileWriter(intercept.TextWriter("traffic capture file writer", mrq, intercept.MatchRequestAndResponseStatus(tc.StatusFilter, mrq), tc.Binary, tc.Limit, syncio.NewWriter(f)))

		closeFunc = func() { f.Close() }
	}

	cmd.RenderWelcome(ttyIO.Out(), httpPort, httpsPort)

	ttyIO.ReadFunc(109, func(in io.Reader, out io.Writer) {
		cmd.RenderMenu(out)

		scanner := bufio.NewScanner(in)

		for scanner.Scan() {
			if cmdFunc, ok := cmd.For(strings.ToUpper(scanner.Text())); ok {
				if cmdFunc(ttyIO, in, out) {
					cmd.RenderMenu(out)
					continue
				}
				return
			}
		}
	})

	return closeFunc, nil
}
