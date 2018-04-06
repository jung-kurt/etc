package util

import (
	"encoding/json"
	"io"
	"os"
)

// JsonPut writes the value specified by val to the writer specified by w
func JsonPut(w io.Writer, val interface{}) (err error) {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	err = enc.Encode(val)
	return
}

// JsonPutFile writes the value specified by val to the file specified by
// fileStr. The fie is overwritten if it already exists.
func JsonPutFile(fileStr string, val interface{}) (err error) {
	var fl *os.File

	fl, err = os.Create(fileStr)
	if err == nil {
		err = JsonPut(fl, val)
		fl.Sync()
		fl.Close()
	}
	return
}

// JsonGet reads a JSON-encoded structure from the reader specified by r and
// places it into the buffer pointed to by valPtr
func JsonGet(r io.Reader, valPtr interface{}) (err error) {
	dec := json.NewDecoder(r)
	err = dec.Decode(valPtr)
	return
}

// JsonGetFile reads a JSON-encoded structure from the file specified by
// fileStr and places it into the buffer pointed to by valPtr
func JsonGetFile(fileStr string, valPtr interface{}) (err error) {
	var fl *os.File

	fl, err = os.Open(fileStr)
	if err == nil {
		err = JsonGet(fl, valPtr)
		fl.Close()
	}
	return
}
