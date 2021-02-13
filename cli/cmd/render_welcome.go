package cmd

import (
	"fmt"
	"io"
)

var RenderWelcome = func(w io.Writer, httpPort, httpsPort int) {
	fmt.Fprintf(w, "\nhflow is listening for http traffic on port %v and https traffic on port %v.\n\npress the 'm' key and hit enter to display the menu...\n\n", httpPort, httpsPort)
}
