package util

import (
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Address string

func (a Address) String() string {
	return string(a)
}

func (a Address) MarhsalJSON() (buf []byte, err error) {
	buf = []byte("\"" + string(a) + "\"")
	return
}

func (a *Address) UnmarshalJSON(buf []byte) (err error) {
	str := string(buf[1 : len(buf)-1])
	_, _, err = ParseAddrPort(str)
	if err == nil {
		*a = Address(str)
	}
	return
}

type Duration time.Duration

func (d Duration) Dur() time.Duration {
	return time.Duration(d)
}

func (d Duration) String() string {
	return d.Dur().String()
}

func (d Duration) MarshalJSON() (buf []byte, err error) {
	buf = []byte("\"" + d.String() + "\"")
	return
}

func (d *Duration) UnmarshalJSON(buf []byte) (err error) {
	var td time.Duration
	td, err = time.ParseDuration(string(buf[1 : len(buf)-1]))
	if err == nil {
		*d = Duration(td)
	}
	return
}

type Distance float64

// In returns the receiver value in inches
func (d Distance) In() float64 {
	return float64(d)
}

func (d Distance) String() string {
	return fmt.Sprintf("%fin", float64(d))
}

func (d Distance) MarshalJSON() (buf []byte, err error) {
	buf = []byte("\"" + d.String() + "\"")
	return
}

var reDistance = regexp.MustCompile("^(\\d*\\.?\\d*)(m|mm|cm|in|inch|ft|foot|'|\")$")

func (d *Distance) UnmarshalJSON(buf []byte) (err error) {
	var match []string
	var str string
	var val float64

	str = string(buf)[1 : len(buf)-1]
	match = reDistance.FindStringSubmatch(str)
	if match != nil {
		val, err = strconv.ParseFloat(match[1], 64)
		if err == nil {
			switch match[2] {
			case "in", "inch", `"`:
				// val already in inches
			case "m":
				val *= 39.3701
			case "mm":
				val *= 0.0393701
			case "cm":
				val *= 0.393701
			case "ft", "foot", "'":
				val *= 12
			default:
				err = errf("unrecognized distance unit \"%s\"", match[2])
			}
		}
	} else {
		err = errf("unrecognized distance \"%s\"", str)
	}
	if err == nil {
		*d = Distance(val)
	}
	return
}

// ParseAddrPort takes a string like "20.30.40.50:60" and returns the IP
// address and port if successful, otherwise an error. The colon and port must
// be specified.
func ParseAddrPort(str string) (ip net.IP, port uint16, err error) {
	// See func net.SplitHostPort(hostport string) (host, port string, err error)
	pair := strings.Split(str, ":")
	if len(pair) == 2 {
		ip = net.ParseIP(pair[0])
		if ip != nil {
			var v uint64
			v, err = strconv.ParseUint(pair[1], 10, 16)
			if err == nil {
				port = uint16(v)
			} else {
				err = errf("\"%s\" is invalid port specifier", pair[1])
			}
		} else {
			err = errf("\"%s\" not a valid IP address", pair[0])
		}
	} else {
		err = errf("\"%s\" is missing port specifier", str)
	}
	return
}
func errf(format string, args ...interface{}) error {
	return fmt.Errorf(format, args...)
}
