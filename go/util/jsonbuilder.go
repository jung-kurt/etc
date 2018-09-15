package util

import (
	"encoding/json"
	"fmt"
	"strings"
)

const (
	cnNone     = 0
	cnObject   = 1
	cnArray    = 2
	cnFinished = 3
)

type elementType struct {
	state int
	count int
}

// JSONBuilder facilitates the incremental construction of a JSON string. Its
// zero-value is ready for use.
type JSONBuilder struct {
	bld    strings.Builder
	stk    []elementType
	topPtr *elementType // nil if stack is empty
	count  int          // number of items in base (0 or 1)
	err    error
}

func (j *JSONBuilder) stateStr(st int) string {
	switch st {
	case cnObject:
		return "object"
	case cnArray:
		return "array"
	case cnFinished:
		return "finished"
	}
	return "none"
}

func (j *JSONBuilder) errorf(format string, args ...interface{}) {
	if j.err == nil {
		j.err = fmt.Errorf(format, args...)
	}
}

func (j *JSONBuilder) push(state int) {
	if j.err == nil {
		if state == cnArray || state == cnObject {
			j.stk = append(j.stk, elementType{state: state, count: 0})
			j.topPtr = &j.stk[len(j.stk)-1]
		} else {
			j.errorf("%s state cannot be pushed", j.stateStr(state))
		}
	}
}

func (j *JSONBuilder) pop(mustMatch int) {
	if j.err == nil {
		ln := len(j.stk)
		if ln > 0 {
			if j.topPtr.state == mustMatch {
				ln--
				j.stk = j.stk[:ln]
				if ln == 0 {
					j.topPtr = nil
				} else {
					j.topPtr = &j.stk[ln-1]
				}
			} else {
				j.errorf("cannot pop %s state while in %s state", j.stateStr(mustMatch), j.stateStr(j.topPtr.state))
			}
		} else {
			j.errorf("cannot pop %s state from empty state stack", j.stateStr(mustMatch))
		}
	}
}

// ArrayOpen opens an array for subsequent element population. This may be
// successfully called when an array is currently open or when nothing else has
// been added to the builder instance.
func (j *JSONBuilder) ArrayOpen() {
	if j.err == nil {
		switch {
		case j.topPtr == nil:
			if j.count == 0 {
				j.count++ // No more items in base after this
				j.push(cnArray)
				j.bld.WriteString("[")
			} else {
				j.errorf("multiple values are not permitted except in an object or array")
			}
		case j.topPtr.state == cnArray:
			if j.topPtr.count > 0 {
				j.bld.WriteString(",")
			}
			j.topPtr.count++
			j.push(cnArray)
			j.bld.WriteString("[")
		default:
			j.errorf("keyless array can be opened only in array or as only value in base")
		}
	}
}

// ArrayClose closes the currently open array.
func (j *JSONBuilder) ArrayClose() {
	if j.err == nil {
		j.pop(cnArray)
		if j.err == nil {
			j.bld.WriteString("]")
		}
	}
}

// Element adds the specified element to the builder instance. This may be
// successfully called when an array is currently open or when nothing else has
// been added to the builder instance.
func (j *JSONBuilder) Element(val interface{}) {
	var b []byte
	if j.err == nil {
		b, j.err = json.Marshal(val)
		if j.err == nil {
			switch {
			case j.topPtr == nil:
				if j.count == 0 {
					j.bld.Write(b)
					j.count = 1
				} else {
					j.errorf("multiple values are not permitted except in an object or array")
				}
			case j.topPtr.state == cnArray:
				if j.topPtr.count > 0 {
					j.bld.WriteString(",")
				}
				j.bld.Write(b)
				j.topPtr.count++
			default:
				j.errorf("keyless element can be opened only in array or as only value in base")
			}
		}
	}
}

// ObjectOpen opens an object for subsequent key/value population. This may be
// successfully called when an array is currently open or when nothing else has
// been added to the builder instance.
func (j *JSONBuilder) ObjectOpen() {
	if j.err == nil {
		switch {
		case j.topPtr == nil:
			if j.count == 0 {
				j.count++ // No more items in base after this
				j.push(cnObject)
				j.bld.WriteString("{")
			} else {
				j.errorf("multiple values are not permitted except in an object or array")
			}
		case j.topPtr.state == cnArray:
			if j.topPtr.count > 0 {
				j.bld.WriteString(",")
			}
			j.topPtr.count++
			j.push(cnObject)
			j.bld.WriteString("{")
		default:
			j.errorf("keyless object can be opened only in array or as only value in base")
		}
	}
}

// ObjectClose closes the currently open object.
func (j *JSONBuilder) ObjectClose() {
	if j.err == nil {
		j.pop(cnObject)
		if j.err == nil {
			j.bld.WriteString("}")
		}
	}
}

// KeyElement assigns a key/value pair in the open object.
func (j *JSONBuilder) KeyElement(key string, val interface{}) {
	var k, v []byte
	if j.err == nil {
		if j.topPtr != nil && j.topPtr.state == cnObject {
			k, j.err = json.Marshal(key)
			if j.err == nil {
				v, j.err = json.Marshal(val)
				if j.err == nil {
					if j.topPtr.count > 0 {
						j.bld.WriteString(",")
					}
					fmt.Fprintf(&j.bld, "%s:%s", k, v)
					j.topPtr.count++
				}
			}
		} else {
			j.errorf("a keyed element can only be assigned in an open object")
		}
	}
}

func (j *JSONBuilder) keyContainerOpen(key string, state int, openStr string) {
	var k []byte
	if j.err == nil {
		if j.topPtr != nil && j.topPtr.state == cnObject {
			k, j.err = json.Marshal(key)
			if j.err == nil {
				if j.topPtr.count > 0 {
					j.bld.WriteString(",")
				}
				j.topPtr.count++
				j.push(state)
				fmt.Fprintf(&j.bld, "%s:%s", k, openStr)
			}
		} else {
			j.errorf("a keyed element can only be assigned in an open object")
		}
	}
}

// KeyObjectOpen opens an object for subsequent key/value population and
// associates it with the specified key. This may be successfully called when
// an object is currently open.
func (j *JSONBuilder) KeyObjectOpen(key string) {
	j.keyContainerOpen(key, cnObject, "{")
}

// KeyArrayOpen opens an object for subsequent element population and
// associates it with the specified key. This may be successfully called when
// an object is currently open.
func (j *JSONBuilder) KeyArrayOpen(key string) {
	j.keyContainerOpen(key, cnArray, "[")
}

// Reset preapres the JSONBuilder receiver for reuse.
func (j *JSONBuilder) Reset() {
	*j = JSONBuilder{} // reset for reuse
}

// String returns the completed JSON object in compact string form. All open
// containers (arrays and objects) are closed in proper order. The receiver
// instance is reset and is available for reuse after this call returns.
func (j *JSONBuilder) String() (jsonStr string) {
	if j.err == nil {
		// Close all open containers
		for j.err == nil && j.topPtr != nil {
			if j.topPtr.state == cnArray {
				j.ArrayClose()
			} else {
				j.ObjectClose()
			}
		}
		jsonStr = j.bld.String()
		j.Reset()
	}
	return
}

// Error returns the current internal error value. If no error has occurred nil
// is returned.
func (j *JSONBuilder) Error() (err error) {
	return j.err
}
