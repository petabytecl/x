package cmdx

import (
	"fmt"
	"os"
)

// Must fatals with the optional message if err is not nil.
func Must(err error, message string, args ...interface{}) {
	if err == nil {
		return
	}

	Fatalf(message, args...)
}

// Fatalf prints to os.Stderr and exists with code 1.
func Fatalf(message string, args ...interface{}) {
	if len(args) > 0 {
		_, _ = fmt.Fprintf(os.Stderr, message+"\n", args...)
	} else {
		_, _ = fmt.Fprintln(os.Stderr, message)
	}
	os.Exit(1)
}
