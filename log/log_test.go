package log

import (
	sl "log"
	"strings"
	"testing"
)

func TestLogVerbosity(t *testing.T) {
	sb := strings.Builder{}
	sl.SetOutput(&sb)

	SetVerbosity(0)

	Printf(0, "zero")

	if !strings.Contains(sb.String(), "zero") {
		t.Fatal("log message not written when verbosity was enabled")
	}

	sb.Reset()

	Printf(1, "one")

	if strings.Contains(sb.String(), "one") {
		t.Fatal("log message was written when verbosity was not enabled")
	}

	sb.Reset()
	SetVerbosity(1)

	Printf(1, "one")

	if !strings.Contains(sb.String(), "one") {
		t.Fatal("log message was not written when verbosity was enabled")
	}
}

func TestLogContent(t *testing.T) {
	sb := strings.Builder{}

	sl.SetOutput(&sb)

	Printf(0, "log content")

	if !strings.Contains(sb.String(), "log content") || !strings.Contains(sb.String(), "hflow") || !strings.Contains(sb.String(), "lv=0") {
		t.Fatal("log message did contain the correct content")
	}
}
