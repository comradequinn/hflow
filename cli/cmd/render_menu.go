package cmd

import (
	"io"
	"sort"
	"strings"
)

var RenderMenu = func(w io.Writer) {
	sb := strings.Builder{}

	sb.WriteString("\n_________________________________________________________________________________________\n")
	sb.WriteString("\nhflow menu\n")
	sb.WriteString("_________________________________________________________________________________________\n\n")

	options, i := make([]command, len(activeSet)), 0

	for _, c := range activeSet {
		options[i] = c
		i++
	}

	sort.Slice(options, func(i, j int) bool { return options[i].order < options[j].order })

	for _, c := range options {
		sb.WriteString(c.accessKey + " - " + c.description + "\n")
	}

	sb.WriteString("_________________________________________________________________________________________\n")

	sb.WriteString("\nenter the required option: ")

	w.Write([]byte(sb.String()))
}
