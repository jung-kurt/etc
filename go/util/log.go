package util

import (
	"bufio"
	"encoding/binary"
	"encoding/hex"
	"io"
	"log"
	"strings"
)

// LogWriter returns the log instance wrapped as a WriteCloser suitable for
// logging shell command output. The return value's Close() method should be
// called after use. This code was adapted from
// https://play.golang.org/p/ruQZM8Bhf- In an environment in which p is used by
// multiple goroutines (such as the standard library log), this can result in
// interleaved output lines. If p is nil, log.Print() from the standard library
// is used instead.
func LogWriter(p Printer) io.WriteCloser {
	pr, pw := io.Pipe()
	br := bufio.NewReader(pr)
	done := make(chan bool)
	go func() {
		var err error
		for err == nil {
			var line string
			line, err = br.ReadString('\n')
			if err == nil {
				line = strings.TrimSpace(line)
				if line != "" {
					if p == nil {
						log.Print(line)
					} else {
						p.Print(line)
					}
				}
			}
		}
		done <- true
	}()
	return &waitWriteCloser{pw, done}
}

// waitWriteCloser implements an io.WriteCloser. It includes a channel named
// Done that will emit true when closed.
type waitWriteCloser struct {
	io.WriteCloser
	done chan bool
}

// Printer supports basic output of all types. It is supported by the
// log.Logger type.
type Printer interface {
	Print(v ...interface{})
}

// Close implements the io.Closer interface.
func (c *waitWriteCloser) Close() error {
	err := c.WriteCloser.Close()
	<-c.done
	return err
}

// BinaryWrite packs the values specified by args to the writer specified by w.
// Little endian format is used.
func BinaryWrite(w io.Writer, args ...interface{}) {
	for _, arg := range args {
		binary.Write(w, binary.LittleEndian, arg)
	}
}

// BinaryRead unpacks the bytes from r into the values pointed to by the
// arguments specified by args. Little endian format is used. An error is
// returned if a problem occurs.
func BinaryRead(r io.Reader, args ...interface{}) (err error) {
	ln := len(args)
	for j := 0; j < ln && err == nil; j++ {
		err = binary.Read(r, binary.LittleEndian, args[j])
	}
	return
}

// BinaryLog logs the binary data specified by d in combined hexadecimal/ascii
// lines of up to sixteen bytes.
func BinaryLog(d []byte) {
	remLen := len(d)
	for remLen > 0 {
		showLen := remLen
		if showLen > 16 {
			showLen = 16
		}
		log.Printf("%s", hex.Dump(d[:showLen]))
		d = d[showLen:]
		remLen -= showLen
	}
}
