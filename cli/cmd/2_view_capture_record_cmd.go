package cmd

import (
	"comradequinn/hflow/cli/tty"
	"comradequinn/hflow/proxy/store"
	"fmt"
	"io"
	"strconv"
	"strings"
)

var ViewCaptureCommand command

func init() {
	ViewCaptureCommand = command{
		"V",
		2,
		"view the full contents of a captured traffic record",
		func(ttyIO tty.IO, replIn io.Reader, replOut io.Writer) bool {
			if store.Len() == 0 {
				replOut.Write([]byte("\nthere are no traffic capture records to view. hit enter to continue...\n"))
				_, _ = replIn.Read(make([]byte, 1))
				return false
			}

			id, _ := strconv.Atoi(readInput(replIn, replOut, "\nenter the record number: ", numericValidationFunc))

			record, ok := store.Get(id)

			if !ok {
				replOut.Write([]byte("\nthere is no traffic capture records with that id. hit enter to continue...\n"))
				_, _ = replIn.Read(make([]byte, 1))
				return false
			}

			sb := strings.Builder{}

			sb.WriteString(border)
			fmt.Fprintf(&sb, "capture detail for record %v\n", id)
			sb.WriteString(divider)
			sb.WriteString(record.Data)
			sb.WriteString(border)
			sb.WriteString("\nhit enter to continue...\n")

			replOut.Write([]byte(sb.String()))

			_, _ = replIn.Read(make([]byte, 1))

			return true
		},
	}
}
