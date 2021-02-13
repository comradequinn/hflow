package cmd

import (
	"comradequinn/hflow/cli/tty"
	"comradequinn/hflow/proxy"
	"comradequinn/hflow/proxy/intercept"
	"io"
)

var StopCaptureCommand command

func init() {
	StopCaptureCommand = command{
		"|",
		3,
		"stop writing captured traffic to the terminal ",
		func(ttyIO tty.IO, replIn io.Reader, replOut io.Writer) bool {
			proxy.SetTrafficSummaryWriter(intercept.Unset)
			delete(activeSet, StopCaptureCommand.accessKey)
			activeSet[CaptureCommand.accessKey] = CaptureCommand

			replOut.Write([]byte("\nstopped writing captured traffic to the terminal\n"))
			replOut.Write([]byte("\nto view the menu, press the 'm' key and hit enter \n\n"))

			return false
		},
	}
}
