package cmd

import (
	"comradequinn/hflow/cli/tty"
	"io"
	"strings"
)

var AboutCommand = command{"A", 8, "display information about hflow", func(ttyIO tty.IO, replIn io.Reader, replOut io.Writer) bool {
	sb := strings.Builder{}

	sb.WriteString(border)
	sb.WriteString("hflow is a fast, simple, command-line http(s) debugging proxy\n")
	sb.WriteString(divider)
	sb.WriteString("hflow supports:\n")
	sb.WriteString("- traffic capture: capture all traffic or filter by request url and/or response status\n")
	sb.WriteString("- edit & continue: edit requests and/or responses, in transit, that match a specified url pattern\n")
	sb.WriteString("- request re-routing: route requests destined for one host to another, such as a local or test instance\n")
	sb.WriteString("- tls decryption: decrypt https traffic. optionally, add hflow's root ca into your store to avoid security warnings\n")
	sb.WriteString("- response decoding: automatically decodes gzip and brotli encoded responses\n")
	sb.WriteString(border)

	replOut.Write([]byte(sb.String()))

	return false
}}
