package cmd

import (
	"comradequinn/hflow/cli/tty"
	"comradequinn/hflow/proxy"
	"comradequinn/hflow/proxy/intercept"
	"io"
)

var RemoveBreakpointCommand command

func init() {
	RemoveBreakpointCommand = command{
		"/",
		5,
		"remove the active breakpoint",
		func(ttyIO tty.IO, replIn io.Reader, replOut io.Writer) bool {
			proxy.SetTrafficEditor(intercept.Unset)

			delete(activeSet, RemoveBreakpointCommand.accessKey)
			activeSet[BreakpointCommand.accessKey] = BreakpointCommand

			replOut.Write([]byte("\nremoved active breakpoint\n"))
			replOut.Write([]byte("\nto view the menu, press the 'm' key and hit enter\n\n"))

			return false
		},
	}
}
