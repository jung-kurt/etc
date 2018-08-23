package util

import (
	"io"
	"os"
	"strings"
)

// CaptureOutput begins buffering the content that is written to fl. It returns
// a retrieval function and an error. If no error occurs when setting up the
// buffering, that is, err is nil, the retrieval function can be called to stop
// the buffering and obtain the populated string buffer.
func CaptureOutput(fl **os.File) (get func() *strings.Builder, err error) {
	var flSave, flRead, flWrite *os.File
	var ch chan *strings.Builder
	var gotten bool

	flSave = *fl
	flRead, flWrite, err = os.Pipe()
	if err == nil {
		ch = make(chan *strings.Builder)
		*fl = flWrite

		go func() {
			var buf strings.Builder
			io.Copy(&buf, flRead)
			ch <- &buf
		}()

		get = func() (buf *strings.Builder) {
			if !gotten {
				flWrite.Close() // This terminates copy routine in goroutine
				*fl = flSave    // Restore original file
				buf = <-ch      // Retrieve buffer that was posted in goroutine
				gotten = true   // Don't try this again
			}
			return
		}

	}
	return
}

// CaptureStdOutAndErr begins buffering the content that is written to stdout
// and stderr. It returns a retrieval function and an error. If no error occurs
// when setting up the buffering, that is, err is nil, the retrieval function
// can be called to stop the buffering and obtain the accumulated output. The
// return values of the retrieval function are the content written to stdout
// and stderr respectively.
func CaptureStdOutAndErr() (get func() (outStr, errStr string), err error) {
	var outGet, errGet func() *strings.Builder

	outGet, err = CaptureOutput(&os.Stdout)
	if err == nil {
		errGet, err = CaptureOutput(&os.Stderr)
		if err == nil {
			var gotten bool
			get = func() (outStr, errStr string) {
				if !gotten {
					outStr = outGet().String()
					errStr = errGet().String()
					gotten = true
				}
				return
			}
		} else {
			outGet() // Stop goroutine
		}
	}
	return
}
