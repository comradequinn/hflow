package cli

import (
	"comradequinn/hflow/cert"
	"comradequinn/hflow/log"
	"fmt"
	"io"
)

// ExportCA writes the flow CA to stdout
func ExportCA(stdout io.Writer) error {
	if err := cert.WriteCA(stdout); err != nil {
		return fmt.Errorf("error writing hflow ca certificate [%v]", err)
	}

	log.Printf(0, "hflow ca certificate written to stdout in pem format")

	return nil
}
