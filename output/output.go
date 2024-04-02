package output

import (
	"fmt"
	"os"
)

// ErrorExit reports an error to the console, if not nil, and performs os.Exit(1) if error found.
func ErrorExit(domain string, err error) {
	if err == nil {
		return
	}
	Error(domain, err)
	os.Exit(1)
}

// Error reports an error to the console.
func Error(domain string, err error) {
	if err == nil {
		return
	}
	fmt.Printf("%v: \x1b[31;1m%s\x1b[0m\n", domain, fmt.Sprintf("error: %s", err))
}

// Info reports an info to the console.
func Info(prefix, format string, args ...interface{}) {
	fmt.Printf("\x1b[34;1m%s\x1b[0m: ", prefix)
	fmt.Printf("%s\n", fmt.Sprintf(format, args...))
}

// Warning reports a warning to the console.
func Warning(format string, args ...interface{}) {
	fmt.Printf("\x1b[36;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}
