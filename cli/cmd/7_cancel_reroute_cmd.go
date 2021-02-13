package cmd

import (
	"comradequinn/hflow/cli/tty"
	"comradequinn/hflow/proxy"
	"comradequinn/hflow/proxy/intercept"
	"io"
)

var CancelRerouteCommand command

func init() {
	CancelRerouteCommand = command{"\\",
		7,
		"cancel rerouting requests to a different host",
		func(ttyIO tty.IO, replIn io.Reader, replOut io.Writer) bool {
			proxy.SetRequestRerouter(intercept.Unset)

			delete(activeSet, CancelRerouteCommand.accessKey)
			activeSet[RerouteCommand.accessKey] = RerouteCommand

			replOut.Write([]byte("\nactive request rerouting cancelled\n"))
			replOut.Write([]byte("\nto display the menu, press the 'm' key and hit enter\n\n"))

			return false
		},
	}
}
