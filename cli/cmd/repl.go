package cmd

import (
	"bufio"
	"io"
	"strconv"
)

var (
	border                = "_________________________________________________________________________________________\n\n"
	divider               = "_ _ _ _ _ _ _ _ _ _ _ _ _ _ _ _ _ _ _ _ _ _ _ _ _ _ _ _ _ _ _ _ _ _ _ _ _ _ _ _ _ _ _ _ _ \n\n"
	noValidationFunc      = func(_ string) (bool, string) { return true, "" }
	numericValidationFunc = func(s string) (bool, string) {
		_, err := strconv.Atoi(s)
		return err == nil, "value must be numeric"
	}
)

func readInput(replIn io.Reader, replOut io.Writer, prompt string, validate func(string) (bool, string)) string {
	replOut.Write([]byte(prompt))

	scanner, input := bufio.NewScanner(replIn), ""

	for scanner.Scan() {
		input = scanner.Text()

		if valid, invalidMsg := validate(input); !valid {
			replOut.Write([]byte(invalidMsg + "\n"))
			replOut.Write([]byte(prompt))
			continue
		}
		break
	}

	return input
}
