/*
 * Copyright (c) 2012-2018 Kurt Jung (Gmail: kurt.w.jung)
 *
 * Permission to use, copy, modify, and distribute this software for any
 * purpose with or without fee is hereby granted, provided that the above
 * copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
 * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
 * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
 * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 */

package util_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/jung-kurt/etc/go/util"
)

// Test Duration methods
func TestDuration(t *testing.T) {
	var d util.Duration
	var err error

	err = d.Set("1.5h")
	if err != nil {
		t.Fatalf("error setting Duration value with flag.Value interface")
	}
}

// Test various responses to network address parse.
func TestJSON(t *testing.T) {
	var strList = []string{`{"Dst": "3.5x"}`, `{"Dst": "foo"}`, `{"Wait": "12f"}`}
	var err error
	var cfg struct {
		Dst  util.Distance
		Wait util.Duration
	}

	for _, str := range strList {
		err = util.JSONGet(strings.NewReader(str), &cfg)
		if err == nil {
			t.Fatalf("JSON \"%s\" should not have parsed successfully (%s)", str, cfg.Dst)
		}
	}
}

// Test bad source for file copy
func TestShell(t *testing.T) {
	_, err := util.ShellBuf([]byte(""), "/bin/sh", "./xyz")
	if err == nil {
		t.Fatalf("ShellBuf(\"./xyz\") should have generated error")
	}
}

// Test bad source for file copy
func TestFileCopy(t *testing.T) {
	err := util.FileCopy(".", ".foo")
	if err == nil {
		t.Fatalf("FileCopy(\".\", \".foo\") should have generated error")
	}
}

// This example tests the Shell routine
func ExampleShell() {
	var err error
	if runtime.GOOS == "linux" {
		var w []byte
		w, err = util.ShellBuf([]byte(""), "ls", "-1")
		if err == nil {
			if len(w) > 0 {
				// OK
			} else {
				err = fmt.Errorf("empty output from ShellBuf()")
			}
		}
	}
	if err == nil {
		fmt.Printf("OK\n")
	} else {
		fmt.Printf("%s\n", err)
	}
	// Output:
	// OK
}

// This example tests the file copy routine
func ExampleFileCopy() {
	var err error
	var a, b []byte
	a = []byte("12345678")
	err = ioutil.WriteFile("_a", a, 0600)
	if err == nil {
		err = util.FileCopy("_a", "_b")
		if err == nil {
			b, err = ioutil.ReadFile("_b")
			if err == nil {
				if bytes.Equal(a, b) {
					// OK
				} else {
					err = fmt.Errorf("destination file unequal to source file")
				}
			}
			os.Remove("_b")
		}
		os.Remove("_a")
	}
	if err == nil {
		fmt.Printf("file copy successful\n")
	} else {
		fmt.Printf("%s\n", err)
	}
	// Output:
	// file copy successful
}

// This example demonstrates JSON handling
func ExampleJSONPut() {
	const fileStr = "example.json"
	type cfgType struct {
		Addr util.Address
		DstA, DstB, DstC,
		DstD, DstE, DstF util.Distance
		Wait util.Duration
	}
	var str = `{"Addr": "10.20.30.40:50","DstA": "3.5in","DstB": "2.54cm",
	"DstC": "72cm", "DstD": "1.23m", "DstE": "0.5ft", "DstF": "48mm",
	"Wait": "12s"}`
	var err error
	var cfg cfgType

	show := func(lfStr, rtStr string, args ...interface{}) {
		fmt.Println(util.StrDotPairFormat(24, lfStr, rtStr, args...))
	}

	showDst := func(lfStr string, val util.Distance) {
		show(lfStr, "%s", util.Float64ToStrSig(val.In(), ".", ",", 3, 3))
	}

	err = ioutil.WriteFile(fileStr, []byte(str), 0644)
	if err == nil {
		err = util.JSONGetFile(fileStr, &cfg)
		if err == nil {
			fmt.Println(util.StrHeaderFormat(24, "", " %s ", "JSON Fields"))
			show("Addr", "%s", cfg.Addr)
			showDst("DstA", cfg.DstA)
			showDst("DstB", cfg.DstB)
			showDst("DstC", cfg.DstC)
			showDst("DstD", cfg.DstD)
			showDst("DstE", cfg.DstE)
			showDst("DstF", cfg.DstF)
			show("Wait", "%.2f", cfg.Wait.Dur().Seconds())
			err = util.JSONPutFile(fileStr, cfg)
			if err == nil {
				os.Remove(fileStr)
			}
		}
	}
	if err != nil {
		fmt.Printf("Error %s\n", err)
	}
	// Output:
	// ----- JSON Fields ------
	// Addr......10.20.30.40:50
	// DstA................3.50
	// DstB................1.00
	// DstC................28.3
	// DstD................48.4
	// DstE................6.00
	// DstF................1.89
	// Wait...............12.00
}

// This example demonstrates float formatting.
// page-breaking.
func ExampleFloat64ToStrSig() {
	var val float64
	var str string
	var buf bytes.Buffer

	val = -math.Pi * 1000000
	for j := 0; j < 15; j++ {
		buf.Reset()
		fmt.Fprintf(&buf, "[%20f]", val)
		for prec := 1; prec < 6; prec++ {
			str = util.Float64ToStrSig(val, ".", ",", prec, 3)
			fmt.Fprintf(&buf, " [%s]", str)
		}
		fmt.Println(buf.String())
		val /= 10
	}
	// Output:
	// [     -3141592.653590] [-3,000,000] [-3,100,000] [-3,140,000] [-3,142,000] [-3,141,600]
	// [      -314159.265359] [-300,000] [-310,000] [-314,000] [-314,200] [-314,160]
	// [       -31415.926536] [-30,000] [-31,000] [-31,400] [-31,420] [-31,416]
	// [        -3141.592654] [-3,000] [-3,100] [-3,140] [-3,142] [-3,141.6]
	// [         -314.159265] [-300] [-310] [-314] [-314.2] [-314.16]
	// [          -31.415927] [-30] [-31] [-31.4] [-31.42] [-31.416]
	// [           -3.141593] [-3] [-3.1] [-3.14] [-3.142] [-3.1416]
	// [           -0.314159] [-0.3] [-0.31] [-0.314] [-0.3142] [-0.31416]
	// [           -0.031416] [-0.03] [-0.031] [-0.0314] [-0.03142] [-0.031416]
	// [           -0.003142] [-0.003] [-0.0031] [-0.00314] [-0.003142] [-0.0031416]
	// [           -0.000314] [-0.0003] [-0.00031] [-0.000314] [-0.0003142] [-0.00031416]
	// [           -0.000031] [-0.00003] [-0.000031] [-0.0000314] [-0.00003142] [-0.000031416]
	// [           -0.000003] [-0.000003] [-0.0000031] [-0.00000314] [-0.000003142] [-0.0000031416]
	// [           -0.000000] [-0.0000003] [-0.00000031] [-0.000000314] [-0.0000003142] [-0.00000031416]
	// [           -0.000000] [-0.00000003] [-0.000000031] [-0.0000000314] [-0.00000003142] [-0.000000031416]
}

// This example demonstrates the GeometricMean function
func ExampleGeometricMean() {
	list := []float64{3.1, 2.8, 3.0, 2.9, 2.9, 3.7}
	fmt.Printf("Geometric %.6f, Arithmetic %.6f, RMS %.6f",
		util.GeometricMean(list), util.ArithmeticMean(list), util.RootMeanSquare(list))
	// Output:
	// Geometric 3.053326, Arithmetic 3.066667, RMS 3.081125
}

// This example demonstrates the DistanceToPoint method
func ExampleLinearEquationType_DistanceToPoint() {
	le := util.Linear(1, 1, 3, 2)
	for j := -1; j < 5; j++ {
		fmt.Printf("%.3f\n", le.DistanceToPoint(float64(j), 2))
	}
	// Output:
	// 1.789
	// 1.342
	// 0.894
	// 0.447
	// 0.000
	// 0.447
}

// Test various responses to network address parse.
func TestParseAddrPort(t *testing.T) {
	var ip net.IP
	var port uint16
	var err error

	ip, port, err = util.ParseAddrPort("1.2.3.4:5")
	if err != nil || port != 5 || "1.2.3.4" != ip.String() {
		t.Fatalf("unexpected error")
	}
	_, _, err = util.ParseAddrPort("1.2.3.4:x")
	if err == nil {
		t.Fatalf("unexpected success")
	}
	_, _, err = util.ParseAddrPort("1.x.3.4:5")
	if err == nil {
		t.Fatalf("unexpected success")
	}
	_, _, err = util.ParseAddrPort("1.2.3.4")
	if err == nil {
		t.Fatalf("unexpected success")
	}
}

func eprintf(err error, format string, args ...interface{}) {
	if err == nil {
		fmt.Printf(format, args...)
	} else {
		fmt.Printf("%s\n", err)
	}
}

// This example demonstrates various string functions.
func ExampleStrDelimit() {
	var valInt32 int32
	var valUint32 uint32
	var valInt64 int64
	var str string
	var err error

	valInt64 = int64(math.Round(math.Pi * 1000000))
	valInt32 = int32(math.Round(math.Pi * 1000000))

	fmt.Println(util.Float64ToStr(math.Pi*1000000, 3))
	fmt.Println(util.Int32ToStr(valInt32))
	fmt.Println(util.Int32ToStr(-valInt32))
	str = util.Int32ToStr(-valInt32)
	valInt32, err = util.ToInt32(str)
	eprintf(err, "%d\n", valInt32)
	str = util.Int32ToStr(-valInt32)
	valUint32, err = util.ToUint32(str)
	eprintf(err, "%d\n", valUint32)
	fmt.Println(util.StrCurrency100(valInt64))
	fmt.Println(util.StrCurrency100(-valInt64))
	fmt.Println(util.StrCurrency100(31))
	fmt.Println(util.StrDots("left", -20))
	fmt.Println(util.StrDots("right", 20))
	fmt.Println(util.StrDots("two", 2))
	fmt.Println(util.StrDots("two", 3))
	fmt.Println(util.StrDotPairFormat(20, "Prologue", "%d", 12))
	fmt.Println(util.StrDotPairFormat(4, "Prologue", "%d", 12))
	fmt.Println(util.StrDotPair(20, "Epilogue", "467"))
	fmt.Println(util.StrDotPair(2, "Epilogue", "467"))
	fmt.Println(util.StrIf(math.Pi > 3, "G", "L"))
	fmt.Println(util.StrIf(math.Pi < 3, "L", "G"))
	// Output:
	// 3,141,592.654
	// 3,141,593
	// -3,141,593
	// -3141593
	// 3141593
	// $31,415.93
	// -$31,415.93
	// $0.31
	// left ...............
	// .............. right
	// tw
	// two
	// Prologue..........12
	// Prologue..12
	// Epilogue.........467
	// Epilogue..467
	// G
	// G
}

// This example demonstrates the generic sorting function
func ExampleSort() {
	var list = []string{"red", "green", "blue"}

	util.Sort(len(list), func(a, b int) bool {
		return list[a] < list[b]
	}, func(a, b int) {
		list[a], list[b] = list[b], list[a]
	})
	for _, str := range list {
		fmt.Println(str)
	}
	// Output:
	// blue
	// green
	// red
}

// This example demonstrates data point clustering
func ExampleCluster() {
	var list []util.PairType
	var clList [][]util.PairType
	var x float64

	rnd := rand.New(rand.NewSource(42))
	for j := 0; j < 100; j++ {
		x += 6 * rnd.Float64()
		list = append(list, util.PairType{X: x, Y: x})
	}
	clList = util.Cluster(list, 6, 4)
	for _, list = range clList {
		fmt.Printf("[")
		for _, val := range list {
			fmt.Printf("%8.3f", val.X)
		}
		lf, rt, _, _ := util.BoundingBox(list)
		fmt.Printf("] (%.3f - %.3f)\n", lf, rt)
	}
	// Output:
	// [   2.238   2.634   6.259   7.512   7.775  10.074] (2.238 - 10.074)
	// [  27.848  29.156  31.326  32.053  36.037  38.834] (27.848 - 38.834)
	// [  58.879  59.595  61.636  63.355  64.718  68.633  68.885] (58.879 - 68.885)
	// [ 166.203 169.399 170.803 173.522 175.232 175.467 177.020] (166.203 - 177.020)
	// [ 223.779 227.606 230.851 233.541 234.132 237.015] (223.779 - 237.015)
}

// Demonstrate range of values
func ExampleRangeType() {
	var rng util.RangeType

	rng.Set(7, true)
	rng.Set(3, false)
	rng.Set(8, false)
	rng.Set(8, false)
	rng.Set(0, false)
	rng.Set(2, false)
	fmt.Printf("Range: %.3f - %.3f\n", rng.Min, rng.Max)
	// Output:
	// Range: 0.000 - 8.000
}

// Demonstrate weighted averages
func ExampleAverageType() {
	var ave util.AverageType

	rnd := rand.New(rand.NewSource(42))
	for j := 0; j < 12; j++ {
		ave.Add(3+2*rnd.Float64(), 0.2+rnd.Float64())
	}
	fmt.Printf("Weighted average: %.3f\n", ave.Value())
	// Output:
	// Weighted average: 4.094
}

// Demonstrate linear regression root mean square
func ExampleRootMeanSquareLinear() {
	var xList, yList []float64

	rnd := rand.New(rand.NewSource(42))
	eq := util.LinearEquationType{Slope: 1.32, Intercept: 54.6}
	for j := 0; j < 12; j++ {
		x := float64(j)
		xList = append(xList, x)
		y := eq.Slope*x + eq.Intercept + (0.2*rnd.Float64() - 0.1)
		yList = append(yList, y)
	}
	fmt.Printf("RMS: %.3f\n", util.RootMeanSquareLinear(xList, yList, eq.Intercept, eq.Slope))
	fmt.Printf("Intercept: %.3f\n", util.LinearPointSlope(100, 100, eq.Slope))
	fmt.Printf("LinearY(12): %.3f\n", util.LinearY(eq.Slope, eq.Intercept, 12))
	fmt.Printf("Distance to (100, 100): %.3f\n", eq.DistanceToPoint(100, 100))
	fmt.Printf("Perpendicular: %s\n", eq.Perpendicular(100))
	eq.Intercept = -12
	fmt.Printf("Shifted equation: %s\n", eq)
	ceq := util.LinearEquationType{Slope: 0, Intercept: 5}
	fmt.Printf("Distance to constant line: %.3f\n", ceq.DistanceToPoint(6, 6))
	// Output:
	// RMS: 0.052
	// Intercept: -32.000
	// LinearY(12): 70.440
	// Distance to (100, 100): 52.294
	// Perpendicular: y = -0.757576 x + 262.357576
	// Shifted equation: y = 1.320000 x - 12.000000
	// Distance to constant line: 1.000
}

// Demonstrate archiving a directory
func ExampleArchive() {
	var err error
	const dirStr = "test"
	const zipStr = "test.zip"

	os.RemoveAll(dirStr)
	os.Remove(zipStr)
	err = os.Mkdir(dirStr, 0755)
	if err == nil {
		for j := 0; j < 3 && err == nil; j++ {
			fileStr := filepath.Join(dirStr, fmt.Sprintf("file%02d.txt", j))
			err = ioutil.WriteFile(fileStr, []byte(fileStr), 0644)
		}
		if err == nil {
			err = util.ArchiveFile(dirStr, "test.zip", func(fileStr string) {
				fmt.Printf("archiving %s...\n", fileStr)
			})
			if err == nil {
				fmt.Printf("successfully created %s\n", zipStr)
				os.RemoveAll(dirStr)
				err = util.UnarchiveFile(zipStr, dirStr, func(fileStr string) {
					fmt.Printf("extracting %s...\n", fileStr)
				})
				if err == nil {
					fmt.Printf("successfully unarchived %s\n", zipStr)
				}
				os.RemoveAll(dirStr)
				os.Remove(zipStr)
			}
		}
	}
	if err != nil {
		fmt.Printf("archive error: %s\n", err)
	}
	// Output:
	// archiving test/file00.txt...
	// archiving test/file01.txt...
	// archiving test/file02.txt...
	// successfully created test.zip
	// extracting test/file00.txt...
	// extracting test/file01.txt...
	// extracting test/file02.txt...
	// successfully unarchived test.zip
}

// Demonstrate the binary read, write, and log routines
func ExampleBinaryLog() {
	type recType struct {
		a uint8
		b uint16
		c uint32
		d uint64
		e int8
		f int16
		g int32
		h int64
	}
	var err error

	show := func(hdrStr string, r recType) {
		log.Printf("%s", hdrStr)
		log.Printf("a: %d, b: %d, c: %d, d: %d, e: %d, f: %d, g: %d, h: %d",
			r.a, r.b, r.c, r.d, r.e, r.f, r.g, r.h)
	}

	load := func() (sl []byte) {
		var r recType
		var buf bytes.Buffer
		r = recType{a: 12, b: 142, c: 48456, d: 6581, e: -2, f: -765, g: 2776, h: -54232}
		show("BinaryWrite", r)
		util.BinaryWrite(&buf, r.a, r.b, r.c, r.d, r.e, r.f, r.g, r.h)
		sl = buf.Bytes()
		return
	}

	log.SetFlags(0)
	log.SetOutput(os.Stdout)
	sl := load()
	log.Printf("packed bytes")
	util.BinaryLog(sl)
	var r recType
	err = util.BinaryRead(bytes.NewReader(sl), &r.a, &r.b, &r.c, &r.d, &r.e, &r.f, &r.g, &r.h)
	if err == nil {
		show("BinaryRead", r)
	} else {
		log.Printf("error reading packed bytes: %s", err)
	}
	// Output:
	// BinaryWrite
	// a: 12, b: 142, c: 48456, d: 6581, e: -2, f: -765, g: 2776, h: -54232
	// packed bytes
	// 00000000  0c 8e 00 48 bd 00 00 b5  19 00 00 00 00 00 00 fe  |...H............|
	// 00000000  03 fd d8 0a 00 00 28 2c  ff ff ff ff ff ff        |......(,......|
	// BinaryRead
	// a: 12, b: 142, c: 48456, d: 6581, e: -2, f: -765, g: 2776, h: -54232
}

// Demonstrate the use of the logging mechanism as an io.Writer
func ExampleLogWriter() {
	lg := log.New(os.Stdout, "", 0)
	lw := util.LogWriter(lg)
	lg.Printf("simple log line")
	// lw implements the io.Writer interface
	fmt.Fprintf(lw, "Line one\nLine two\n")
	fmt.Fprintf(lw, "Line three\n")
	fmt.Fprintf(lw, "\n")
	fmt.Fprintf(lw, "\n\nLast line\n")
	lw.Close()
	log.SetFlags(0)
	log.SetOutput(os.Stdout)
	lw = util.LogWriter(nil)
	fmt.Fprintf(lw, "written with log.Print() from the standard library\n")
	lw.Close()
	lg.Printf("another simple log line")
	// Output:
	// simple log line
	// Line one
	// Line two
	// Line three
	// Last line
	// written with log.Print() from the standard library
	// another simple log line
}

func ExampleCaptureOutput() {
	var err error
	var getOut, getErr func() *strings.Builder
	var errBuf, outBuf *strings.Builder

	getOut, err = util.CaptureOutput(&os.Stdout)
	if err == nil {
		getErr, err = util.CaptureOutput(&os.Stderr)
		if err == nil {
			for j := 0; j < 5; j++ {
				fmt.Printf("line %d\n", j)
				fmt.Fprintf(os.Stderr, "error %d\n", j)
			}
			errBuf = getErr()
			outBuf = getOut()
			fmt.Printf("--- Stderr (%d) ---\n%s", errBuf.Len(), errBuf.String())
			fmt.Printf("--- Stdout (%d)---\n%s", outBuf.Len(), outBuf.String())
		}
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
	}
	// Output:
	// --- Stderr (40) ---
	// error 0
	// error 1
	// error 2
	// error 3
	// error 4
	// --- Stdout (35)---
	// line 0
	// line 1
	// line 2
	// line 3
	// line 4
}

func ExampleCaptureStdOutAndErr() {
	get, err := util.CaptureStdOutAndErr()
	if err == nil {
		for j := 0; j < 3; j++ {
			fmt.Printf("line %d\n", j)
			fmt.Fprintf(os.Stderr, "error %d\n", j)
		}
		outStr, errStr := get() // This terminates buffering and returns accumulated content
		fmt.Printf("--- Stderr (%d) ---\n%s", len(errStr), errStr)
		fmt.Printf("--- Stdout (%d)---\n%s", len(outStr), outStr)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
	}
	// Output:
	// --- Stderr (24) ---
	// error 0
	// error 1
	// error 2
	// --- Stdout (21)---
	// line 0
	// line 1
	// line 2
}

// Test various JSONBuilder errors
func TestJSONBuilder(t *testing.T) {
	var bld util.JSONBuilder

	bld.Element(1)
	bld.Element(2)
	if bld.Error() == nil {
		t.Fatalf("multiple elements are root are not allowed in JSON")
	}

	bld.Reset()
	bld.ArrayOpen()
	bld.ObjectClose()
	if bld.Error() == nil {
		t.Fatalf("cannot close incorrect container type")
	}

	bld.Reset()
	bld.ObjectClose()
	if bld.Error() == nil {
		t.Fatalf("cannot close unopened container")
	}

	bld.Reset()
	bld.Element(1)
	bld.ArrayOpen()
	if bld.Error() == nil {
		t.Fatalf("cannot open container in full root")
	}

	bld.Reset()
	bld.ObjectOpen()
	bld.ArrayOpen()
	if bld.Error() == nil {
		t.Fatalf("cannot open keyless array in open object")
	}

	bld.Reset()
	bld.ArrayOpen()
	bld.Element(1)
	bld.ArrayOpen()
	if bld.Error() != nil {
		t.Fatalf("arrays within arrays are permitted")
	}

	bld.Reset()
	bld.ObjectOpen()
	bld.Element(1)
	if bld.Error() == nil {
		t.Fatalf("keyless element allowed only in empty root or open array")
	}

	bld.Reset()
	bld.ObjectOpen()
	bld.ObjectOpen()
	if bld.Error() == nil {
		t.Fatalf("keyless object cannot be opened in object")
	}

	bld.Reset()
	bld.Element(1)
	bld.ObjectOpen()
	if bld.Error() == nil {
		t.Fatalf("keyless object cannot be opened in non-empty root")
	}

	bld.Reset()
	bld.ArrayOpen()
	bld.Element(1)
	bld.ObjectOpen()
	if bld.Error() != nil {
		t.Fatalf("objects within arrays are permitted")
	}

	bld.Reset()
	bld.KeyElement("Num", 1)
	if bld.Error() == nil {
		t.Fatalf("keyed elements only permitted in objects")
	}

	bld.Reset()
	bld.ObjectOpen()
	bld.KeyElement("A", 1)
	bld.KeyObjectOpen("B")
	if bld.Error() != nil {
		t.Fatalf("keyed objects within objects are permitted")
	}

	bld.Reset()
	bld.ArrayOpen()
	bld.KeyObjectOpen("A")
	if bld.Error() == nil {
		t.Fatalf("keyed objects only permitted in objects")
	}

	bld.Reset()
	bld.ArrayOpen()
	bld.Element("A")
	_ = bld.String()
	if bld.Error() != nil {
		t.Fatalf("arrays should be closed automatically when String() called")
	}

}

func ExampleJSONBuilder() {
	var bld util.JSONBuilder

	report := func() {
		err := bld.Error()
		if err == nil {
			fmt.Printf("%s\n", bld.String())
		} else {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			bld.Reset()
		}
	}

	bld.ObjectOpen()
	bld.KeyElement("ID", 42)
	bld.KeyElement("Name", "Prairie")
	report() // Close object automatically

	bld.ObjectOpen()                 // {
	bld.KeyObjectOpen("Alpha")       // {{
	bld.KeyArrayOpen("List")         // {{[
	bld.Element("Shiawassee")        // {{[
	bld.Element(true)                // {{[
	bld.Element(3.14)                // {{[
	bld.ArrayClose()                 // {{
	bld.KeyElement("Name", "Brilla") // {{
	bld.ObjectClose()                // {
	bld.KeyElement("Beta", 1234)     // {
	bld.ObjectClose()                //
	report()

	// JSON-encodable structures can be elements
	bld.Element(struct {
		ID   int
		Name string
	}{22, "Tess"})
	report()

	// Output:
	// {"ID":42,"Name":"Prairie"}
	// {"Alpha":{"List":["Shiawassee",true,3.14],"Name":"Brilla"},"Beta":1234}
	// {"ID":22,"Name":"Tess"}
}

// func simple(val interface{}) {
// 	var bld JSONBuilder
//
// 	bld.Element(val)
// 	fmt.Printf("JSON: %s\n", bld.String())
// }
//
// func arraygen(bld *JSONBuilder, stk []int) {
// 	if len(stk) > 0 {
// 		bld.ArrayOpen()
// 		for j := 0; j < stk[0]; j++ {
// 			if len(stk) == 1 {
// 				switch j {
// 				case 0:
// 					bld.Element("ok")
// 				case 1:
// 					bld.Element(rec.b)
// 					rec.b = !rec.b
// 				default:
// 					bld.Element(j)
// 				}
// 			} else {
// 				arraygen(bld, stk[1:])
// 			}
// 		}
// 		bld.ArrayClose()
// 	}
// }
//
// func array(counts ...int) {
// 	var bld JSONBuilder
//
// 	arraygen(&bld, counts)
// 	report(&bld)
// }
//
//
// func customA() {
// 	var bld JSONBuilder
//
// }
//
//
// func customB() {
// 	var bld JSONBuilder
//
// }
//
// func main() {
// 	simple(42)
// 	simple("Shiawassee")
// 	simple(true)
// 	simple(rec)
// 	array(3)
// 	array(2, 3)
// 	array(2, 3, 2)
// 	customA()
// 	customB()
