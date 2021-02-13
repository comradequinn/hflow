package cmd

import (
	"comradequinn/hflow/cli/tty"
	"io"
)

type (
	// CommandFunc represents a func that executes in response
	// to CLI input.
	//
	// The REPL IO parameters can be used in this REPL context to read and to to the
	// users terminal
	//
	// The ttyIO parameter refers to the terminal in normal usage, outside of the current
	// REPL context function executes in. This may be passed other functions created in the
	// REPL context but invoked outside of it, that therefore cannot use REPL IO
	CommandFunc func(ttyIO tty.IO, replIn io.Reader, replOut io.Writer) (renderMenu bool)
	command     struct {
		accessKey   string
		order       int
		description string
		handler     CommandFunc
	}
)

var (
	activeSet = map[string]command{}
)

func init() {
	activeSet[StatusCommand.accessKey] = StatusCommand
	activeSet[CaptureCommand.accessKey] = CaptureCommand
	activeSet[BreakpointCommand.accessKey] = BreakpointCommand
	activeSet[RerouteCommand.accessKey] = RerouteCommand
	activeSet[AboutCommand.accessKey] = AboutCommand
	activeSet[ExitMenuCommand.accessKey] = ExitMenuCommand
}

// Returns any matching CommandFunc for the specified accessKey.
// If found true is returned as the second return value, otherwise false.
func For(accessKey string) (CommandFunc, bool) {
	c, ok := activeSet[accessKey]
	return c.handler, ok
}
