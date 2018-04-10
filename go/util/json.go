package util

import (
	"encoding/json"
	"io"
	"os"
)

// JSONPut writes the value specified by val to the writer specified by w
func JSONPut(w io.Writer, val interface{}) (err error) {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	err = enc.Encode(val)
	return
}

// JSONPutFile writes the value specified by val to the file specified by
// fileStr. The fie is overwritten if it already exists.
func JSONPutFile(fileStr string, val interface{}) (err error) {
	var fl *os.File

	fl, err = os.Create(fileStr)
	if err == nil {
		err = JSONPut(fl, val)
		fl.Sync()
		fl.Close()
	}
	return
}

// JSONGet reads a JSON-encoded structure from the reader specified by r and
// places it into the buffer pointed to by valPtr
func JSONGet(r io.Reader, valPtr interface{}) (err error) {
	dec := json.NewDecoder(r)
	err = dec.Decode(valPtr)
	return
}

// JSONGetFile reads a JSON-encoded structure from the file specified by
// fileStr and places it into the buffer pointed to by valPtr
func JSONGetFile(fileStr string, valPtr interface{}) (err error) {
	var fl *os.File

	fl, err = os.Open(fileStr)
	if err == nil {
		err = JSONGet(fl, valPtr)
		fl.Close()
	}
	return
}
