package cmd

import (
	"comradequinn/hflow/cli/tty"
	"comradequinn/hflow/proxy"
	"io"
	"strings"
)

var StatusCommand = command{"S", 0, "display the current proxy settings", func(ttyIO tty.IO, replIn io.Reader, replOut io.Writer) bool {
	sb := strings.Builder{}

	sb.WriteString(border)
	sb.WriteString("the proxy has the following configuration:\n")
	sb.WriteString(divider)
	sb.WriteString("capture:\t" + proxy.TrafficSummaryWriter().Label() + "\n")
	sb.WriteString("breakpoint:\t" + proxy.TrafficEditor().Label() + "\n")
	sb.WriteString("reroute:\t" + proxy.RequestRerouter().Label() + "\n")
	sb.WriteString(border)

	replOut.Write([]byte(sb.String()))

	return false
}}
