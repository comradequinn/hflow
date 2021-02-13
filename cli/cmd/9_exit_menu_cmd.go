package cmd

import (
	"comradequinn/hflow/cli/tty"
	"io"
)

var ExitMenuCommand = command{
	"X",
	9,
	"exit the menu without making any changes",
	func(ttyIO tty.IO, replIn io.Reader, replOut io.Writer) bool {
		replOut.Write([]byte("menu exited. to display the menu at any time, press the 'm' key and hit enter\n\n"))
		return false
	},
}
