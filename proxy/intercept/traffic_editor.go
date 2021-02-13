package intercept

import (
	"bufio"
	"bytes"
	"comradequinn/hflow/log"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"
)

// CaptureStdIO represents a func that, when invoked, prevents any other routine writing to stdio
// until the release function returned from the initial invocation is, itself, invoked
type CaptureStdIO func() (release func())

// aliases to allow shimming in tests
var (
	runCmd     = func(cmd *exec.Cmd) error { return cmd.Run() }
	inTerminal = func() func() bool { // figure out whether the stdout is connected to a terminal, this only needs doing once to close around the result
		fileInfo, _ := os.Stdout.Stat()
		terminal := true

		if (fileInfo.Mode() & os.ModeCharDevice) != os.ModeCharDevice {
			terminal = false
		}

		return func() bool { return terminal }
	}()
)

// TrafficEditor opens the default text editor showing the contents of the request and/or response where it matches the matchfunc.
// The data can then be edited as required by the user, when the editor terminates, the updated data
// is forwarded
//
// As the editor requires the use of the process's stdio, the captureStdio func must be provided which, when invoked,
// ensures no other routine will act on them
func TrafficEditor(label string, mrq MatchRequestFunc, mrs MatchResponseFunc, captureStdIO CaptureStdIO, onApplyFunc func()) Intercept {
	randomInt := func() func() int64 {
		rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
		return func() int64 { return rnd.Int63() }
	}()

	editor := os.Getenv("EDITOR")

	if editor == "" {
		editor = "vim"
	}

	log.Printf(3, "intercept labelled [%v] assigned [%v] as the editor", label, editor)

	edit := func(context string, data []byte) ([]byte, error) {
		if !inTerminal() {
			return nil, fmt.Errorf("error in intercept labelled [%v] for %v. stdio of process must be a terminal in order to support editing the request", label, context)
		}

		file := "." + strconv.FormatInt(randomInt(), 10) + ".hflow.edit"

		if err := os.WriteFile(file, data, 0666); err != nil {
			return nil, fmt.Errorf("error in intercept labelled [%v] for %v. unable to create an edit file: %v", label, context, err)
		}

		log.Printf(3, "intercept labelled [%v] wrote edit file for %v named [%v]", label, context, file)

		defer func() {
			if err := os.Remove(file); err != nil {
				log.Printf(0, "error in intercept labelled [%v] for %v. unable to remove edit file named [%v]: %v", label, context, file, err)
			}
		}()

		editorCmd := exec.Command(editor, file)
		editorCmd.Stdin, editorCmd.Stdout, editorCmd.Stderr = os.Stdin, os.Stdout, os.Stderr

		log.Printf(3, "intercept labelled [%v] prepared editor command for %v as [%v]", label, context, editorCmd.String())

		releaseStdIO := captureStdIO()
		err := runCmd(editorCmd)
		releaseStdIO()

		if err != nil {
			return nil, fmt.Errorf("error in intercept labelled [%v] for %v. error code returned from editor used to process edit file. editor cmd was [%v]: %v", label, context, editorCmd.String(), err)
		}

		log.Printf(3, "intercept labelled [%v] executed editor command for %v successfully", label, context)

		return os.ReadFile(file)
	}

	return New(label, mrq, mrs,
		func(c *Context, origRq *RequestData) error {
			onApplyFunc()

			editedData, err := edit(fmt.Sprintf("request [%v]", origRq.URL.String()), []byte(origRq.String()))

			if err != nil {
				return fmt.Errorf("error in intercept labelled [%v] for request [%v]. unable to read edit file after editing: %v", label, origRq.URL.String(), err)
			}

			parsedHTTP, err := http.ReadRequest(bufio.NewReader(bytes.NewReader(editedData)))

			if err != nil {
				return fmt.Errorf("error in intercept labelled [%v] for request [%v]. unable to parse the contents of the edit file as a http request after editing: %v", label, origRq.URL.String(), err)
			}

			editedRq, err := readRequestData(parsedHTTP)

			if err != nil {
				return fmt.Errorf("error in intercept labelled [%v] for request [%v]. unable to read the http request data after editing: %v", label, origRq.URL.String(), err)
			}

			editedRq.InitialHost, editedRq.URL.Scheme, editedRq.URL.Host = origRq.InitialHost, origRq.URL.Scheme, origRq.URL.Host // some data is lost during serialization, so restore it here

			*origRq = editedRq
			c.RequestEdited = true

			log.Printf(2, "intercept labelled [%v] updated data for [%v]", label, origRq.URL.String())

			return nil
		},
		func(c *Context, origRs *ResponseData) error {
			onApplyFunc()

			editedData, err := edit(fmt.Sprintf("response to [%v]", origRs.Request.URL.String()), []byte(origRs.String()))

			if err != nil {
				return fmt.Errorf("error in intercept labelled [%v] for response to [%v]. unable to read edit file after editing: %v", label, origRs.Request.URL.String(), err)
			}

			parsedHTTP, err := http.ReadResponse(bufio.NewReader(bytes.NewReader(editedData)), nil)

			if err != nil {
				return fmt.Errorf("error in intercept labelled [%v] for response to [%v]. unable to parse the contents of the edit file as a http response after editing: %v", label, origRs.Request.URL.String(), err)
			}

			parsedHTTP.Request, err = origRs.Request.toHTTP()

			if err != nil {
				return fmt.Errorf("error in intercept labelled [%v] for response to [%v]. unable to assign the original request to the http response after editing: %v", label, origRs.Request.URL.String(), err)
			}

			editedRs, err := readResponseData(parsedHTTP)

			if err != nil {
				return fmt.Errorf("error in intercept labelled [%v] for response to [%v]. unable to read the http response data after editing: %v", label, origRs.Request.URL.String(), err)
			}

			*origRs = editedRs
			c.ResponseEdited = true

			log.Printf(2, "intercept labelled [%v] updated data for response to [%v]", label, origRs.Request.URL.String())

			return nil
		},
	)
}
