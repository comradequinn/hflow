package cmd

import (
	"comradequinn/hflow/cli/tty"
	"comradequinn/hflow/proxy"
	"comradequinn/hflow/proxy/intercept"
	"fmt"
	"io"
)

var CaptureCommand command

func init() {
	CaptureCommand = command{
		"C",
		1,
		"write captured traffic to the terminal (optional traffic filters can be applied)",
		func(ttyIO tty.IO, replIn io.Reader, replOut io.Writer) bool {
			rqFilter := readInput(replIn, replOut, "\nfilter traffic for request urls containing (hit enter for no filter): ", noValidationFunc)
			rsFilter := readInput(replIn, replOut, "filter traffic for response statuses containing (hit enter for no filter): ", noValidationFunc)

			mrq := intercept.MatchRequestURL(rqFilter)

			fmtFilter := func(s string) string {
				if s == "" {
					return "any"
				}

				return s
			}

			proxy.SetTrafficSummaryWriter(intercept.SummaryTextWriter(fmt.Sprintf("write to terminal when request urls match [%v] and the associated response statuses match [%v]", fmtFilter(rqFilter), fmtFilter(rsFilter)), mrq, intercept.MatchRequestAndResponseStatus(rsFilter, mrq), ttyIO.Out()))

			delete(activeSet, CaptureCommand.accessKey)
			activeSet[StopCaptureCommand.accessKey] = StopCaptureCommand
			activeSet[ViewCaptureCommand.accessKey] = ViewCaptureCommand

			replOut.Write([]byte("\nto display the menu at any time during the capture, press the 'm' key and hit enter. press enter to start capturing\n"))
			_, _ = replIn.Read(make([]byte, 1))

			replOut.Write([]byte("capture started\n\n"))

			return false
		},
	}
}
