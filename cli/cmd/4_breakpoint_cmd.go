package cmd

import (
	"comradequinn/hflow/cli/tty"
	"comradequinn/hflow/proxy"
	"comradequinn/hflow/proxy/intercept"
	"fmt"
	"io"
	"strconv"
)

var BreakpointCommand command

func init() {
	BreakpointCommand = command{
		"B",
		4,
		"set a breakpoint to allow request or response editing",
		func(ttyIO tty.IO, replIn io.Reader, replOut io.Writer) bool {
			breakCondition := readInput(replIn, replOut, "break on traffic where the request matches: ", func(i string) (bool, string) {
				if len(i) < 2 {
					return false, "invalid breakpoint condition"
				}

				return true, ""
			})

			breakMode := readInput(replIn, replOut, "break on request only (0), response only (1), both (2): ", func(i string) (bool, string) {
				v, err := strconv.Atoi(i)

				return err == nil && v >= 0 && v < 3, "invalid break mode"
			})

			mrq, mrs := intercept.MatchRequestURL(breakCondition), intercept.MatchNoResponses

			if breakMode == "1" {
				mrq, mrs = intercept.MatchNoRequests, intercept.MatchSourceRequestURL(mrq)
			}

			if breakMode == "2" {
				mrs = intercept.MatchSourceRequestURL(mrq)
			}

			breakModeLabel := ""

			switch breakMode {
			case "0":
				breakModeLabel = "requests for"
			case "1":
				breakModeLabel = "responses to"
			case "2":
				breakModeLabel = "requests for, and responses to,"
			}

			proxy.SetTrafficEditor(intercept.TrafficEditor(fmt.Sprintf("break on %v [%v]", breakModeLabel, breakCondition), mrq, mrs, ttyIO.Lock, func() { ttyIO.Out().Write([]byte("\rX      - breakpoint condition met. hit enter to edit...")) }))

			delete(activeSet, BreakpointCommand.accessKey)
			activeSet[RemoveBreakpointCommand.accessKey] = RemoveBreakpointCommand

			replOut.Write([]byte("\nbreakpoint configured. hit enter to apply....\n"))
			_, _ = replIn.Read(make([]byte, 1))

			replOut.Write([]byte("breakpoint applied. to display the menu at any time, press the 'm' key and hit enter.\n\n"))

			return false
		},
	}
}
