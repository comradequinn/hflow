// Package log provides a basic logging, modelled on stdlib's log package, with configurable verbosity
package log

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

var (
	verbosity int = 0
)

func init() {
	log.SetPrefix("hflow ")
}

// SetWriter specifies the writer to send log data to
func SetWriter(w io.Writer) {
	log.SetOutput(w)
	Printf(1, "log writer set to new writer")
}

// SetVerbosity specifies the maximum verbosity of logs written
func SetVerbosity(v int) {
	verbosity = v
	Printf(1, "log verbosity set to [%v]", v)
}

// Printf writes s to the log formatted with args if the configured verbosity is <= to v
func Printf(v int, s string, args ...any) {
	if v <= verbosity {
		log.Printf("lv="+strconv.Itoa(v)+" "+s+"\n", args...)
	}
}

// Fatalf calls Printf then exits the program with a return code of 1
func Fatalf(v int, s string, args ...any) {
	Printf(v, s, args...)
	os.Exit(1)
}

// Panicf calls Printf then panics with the same information used for Printf
func Panicf(v int, s string, args ...any) {
	Printf(v, s, args...)
	panic(fmt.Sprintf(s, args...))
}
