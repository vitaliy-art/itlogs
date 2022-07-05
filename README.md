# itlogs
Package itlogs is a wrapper for log standard library.

## Installation
`go get github.com/vitaliy-art/itlogs`

## Usage
```go
package main

import (
	"os"

	"github.com/vitaliy-art/itlogs"
)

func main() {
	// Get the default logger.
	logger := itlogs.GetDefaultLogger()

	flags := []itlogs.LogMsgFlag{itlogs.Lmsgprefix, itlogs.Ldate, itlogs.Ltime}
	// Create a new logger.
	logger = itlogs.NewLogger(os.Stdout, flags, itlogs.Debug)

	// Set the default logger.
	itlogs.SetDefaultLogger(logger)

	// Just log it!
	logger.Info("Hello world!")
}
```
