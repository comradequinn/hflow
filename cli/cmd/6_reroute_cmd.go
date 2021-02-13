package cmd

import (
	"comradequinn/hflow/cli/tty"
	"comradequinn/hflow/proxy"
	"comradequinn/hflow/proxy/intercept"
	"fmt"
	"io"
)

var RerouteCommand command

func init() {
	RerouteCommand = command{
		"R",
		6,
		"reroute requests to a different host",
		func(ttyIO tty.IO, replIn io.Reader, replOut io.Writer) bool {
			fromHost := readInput(replIn, replOut, "enter the host to reroute traffic from: ", func(i string) (bool, string) {
				if len(i) < 2 {
					return false, "invalid host input"
				}

				return true, ""
			})

			toHost := readInput(replIn, replOut, "enter the host to reroute traffic to: ", func(i string) (bool, string) {
				if len(i) < 2 {
					return false, "invalid host input"
				}

				return true, ""
			})

			proxy.SetRequestRerouter(intercept.RequestRerouter(fmt.Sprintf("from [%v] to [%v]", fromHost, toHost), intercept.MatchRequestHost(fromHost), toHost))

			delete(activeSet, RerouteCommand.accessKey)
			activeSet[CancelRerouteCommand.accessKey] = CancelRerouteCommand

			replOut.Write([]byte("\rrerouting configuration ready. press enter to apply....\n"))
			_, _ = replIn.Read(make([]byte, 1))

			replOut.Write([]byte("rerouting configuration applied. to display the menu at any time, press the 'm' key and hit enter.\n\n"))

			return false
		},
	}
}
