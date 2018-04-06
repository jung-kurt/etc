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

// StrIf returns aStr if cond is true, otherwise bStr.
func StrIf(cond bool, aStr string, bStr string) string {
	if cond {
		return aStr
	}
	return bStr
}

// StrDotPair fills in the blank space between two strings with dots. The total
// length displayed is indicated by fullLen. The left and right strings are
// specified by lfStr and rtStr respectively. At least two dots are shown,
// even if this means the total length exceeds fullLen.
func StrDotPair(fullLen int, lfStr, rtStr string) string {
	dotLen := fullLen - len(lfStr) - len(rtStr)
	if dotLen < 2 {
		dotLen = 2
	}
	return lfStr + strings.Repeat(".", dotLen) + rtStr
}

// StrDotPairFormat fills in the blank space between two strings with dots. The
// total length displayed is indicated by fullLen. The left string is specified
// by lfStr. The right string is formatted. At least two dots are shown, even
// if this means the total length exceeds fullLen.
func StrDotPairFormat(fullLen int, lfStr, rtFmtStr string, args ...interface{}) string {
	rtStr := fmt.Sprintf(rtFmtStr, args...)
	dotLen := fullLen - len(lfStr) - len(rtStr)
	if dotLen < 2 {
		dotLen = 2
	}
	return lfStr + strings.Repeat(".", dotLen) + rtStr
}

// StrDots fills in blank spaces with dots. A negative value for fullLen
// indicates a left justified string; positive indicates right. For example,
// ("abc", -12) returns "abc ........", ("abc", 12) returns "........ abc". If
// the length of str is longer than the absolute value of fullLen, it is
// truncated to that value.
func StrDots(str string, fullLen int) string {
	left := fullLen < 0
	if left {
		fullLen = -fullLen
	}
	slen := len(str)
	if slen > fullLen {
		return str[:fullLen]
	} else if slen < fullLen-1 {
		dotStr := strings.Repeat(".", fullLen-slen-1)
		if left {
			return str + " " + dotStr
		}
		return dotStr + " " + str
	}
	return str
}

// StrDelimit converts 'ABCDEFG' to, for example, 'A,BCD,EFG'
func StrDelimit(str string, sepstr string, sepcount int) string {
	pos := len(str) - sepcount
	for pos > 0 {
		str = str[:pos] + sepstr + str[pos:]
		pos = pos - sepcount
	}
	return str
}

// StrCurrency100 converts -123456789 to -$1,234,567.89
func StrCurrency100(amt100 int64) (str string) {
	var sign string
	if amt100 < 0 {
		sign = "-"
		amt100 = -amt100
	} else {
		sign = ""
	}
	if amt100 < 100 {
		str = fmt.Sprintf("%s$0.%02d", sign, amt100)
	} else {
		str = strconv.FormatInt(amt100, 10)
		ln := len(str)
		str = fmt.Sprintf("%s$%s.%s", sign, StrDelimit(str[:ln-2], ",", 3), str[ln-2:])
	}
	return
}

// ToUint32 converts the specified string to a 32-bit unsigned integer
func ToUint32(str string) (v uint32, err error) {
	var v64 uint64
	v64, err = strconv.ParseUint(str, 10, 32)
	if err == nil {
		v = uint32(v64)
	}
	return
}

// ToInt32 converts the specified string to a 32-bit signed integer
func ToInt32(str string) (v int32, err error) {
	var v64 int64
	v64, err = strconv.ParseInt(str, 10, 32)
	if err == nil {
		v = int32(v64)
	}
	return
}

var reFloat = regexp.MustCompile("^(\\-?)(\\d*)(\\.?\\d*)$")

// Float64ToStr returns a string with commas, for example, 1234 -> "1,234"
func Float64ToStr(v float64, precision int) (str string) {
	var match []string
	str = strconv.FormatFloat(v, 'f', precision, 64)
	match = reFloat.FindStringSubmatch(str)
	// log.Printf("str [%s], reFloat [%s], match %v", str, reFloat, match)
	if match != nil {
		str = match[1] + StrDelimit(match[2], ",", 3) + match[3]
	}
	return
}

var reExpFloat = regexp.MustCompile("^(\\-?)(\\d+)(\\.?)(\\d*)E(\\-|\\+)(\\d\\d)$")

// Float64ToStrSig returns a string representation of val adjusted to the
// specified number of significant digits.
func Float64ToStrSig(val float64, dec, sep string, sigDig, grpLen int) (str string) {
	var pad, exp int
	var match []string
	var sign, tailStr string

	str = strconv.FormatFloat(val, 'E', sigDig-1, 64)
	match = reExpFloat.FindStringSubmatch(str)
	if match != nil {
		sign = match[1]
		exp, _ = strconv.Atoi(match[6]) // regexp guarantees success
		if match[5] == "-" {
			str = strings.Repeat("0", exp) + match[2] + match[4]
			str = sign + str[:1] + dec + str[1:]
		} else {
			pad = exp - len(match[4])
			if pad > 0 {
				match[4] += strings.Repeat("0", pad)
			}
			str = match[2] + match[4]
			exp++
			tailStr = ""
			if len(str) > exp {
				tailStr = dec + str[exp:]
				str = str[:exp]
			}
			str = sign + StrDelimit(str, sep, grpLen) + tailStr
		}
	}
	return
}

// IntToStr returns a string with commas, for example, 1234 -> "1,234"
func IntToStr(v int64) (str string) {
	str = strconv.FormatInt(v, 10)
	neg := str[0:1] == "-"
	if neg {
		str = str[1:]
	}
	str = StrDelimit(str, ",", 3)
	if neg {
		str = "-" + str
	}
	return
}

// Int32ToStr returns a string with commas, for example 1234 -> "1,234"
func Int32ToStr(v int32) string {
	return IntToStr(int64(v))
}
