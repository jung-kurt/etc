package util

import (
	"bytes"
	"errors"
	"io"
	"os/exec"
)

// Shell executes the command specified by cmdStr. Standard input to the
// command is streamed from rdr. Standard output from the command is streamed
// to wr. Any content that the command writes to standard error is converted to
// the returned error code. Any practical number of string arguments to the
// command can follow cmdStr.
func Shell(wr io.Writer, rdr io.Reader, cmdStr string, args ...string) (err error) {
	var errBuf bytes.Buffer
	var cmd *exec.Cmd

	cmd = exec.Command(cmdStr, args...)
	cmd.Stdin = rdr
	cmd.Stdout = wr
	cmd.Stderr = &errBuf
	err = cmd.Run()
	if errBuf.Len() > 0 {
		// Overwrite error with more descriptive message
		err = errors.New(errBuf.String())
	}
	return
}

// ShellBuf executes the command specified by cmdStr. Standard input to the
// command is read from inBuf. Standard output from the command is returned as
// outBuf. Any content that the command writes to standard error is converted
// to the returned error code. Any practical number of string arguments to the
// command can follow cmdStr.
func ShellBuf(inBuf []byte, cmdStr string, args ...string) (outBuf []byte, err error) {
	var wr bytes.Buffer
	err = Shell(&wr, bytes.NewReader(inBuf), cmdStr, args...)
	if err == nil {
		outBuf = wr.Bytes()
	}
	return
}
