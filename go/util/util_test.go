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
	"math"
	"net"
	"os"
	"testing"

	"github.com/jung-kurt/etc/go/util"
)

// This example demonstrates JSON handling
func ExampleJSONPut() {
	const fileStr = "example.json"
	type cfgType struct {
		Addr             util.Address
		DstA, DstB, DstC util.Distance
	}
	var str = `{"Addr": "10.20.30.40:50","DstA": "3.5in","DstB": "2.54cm","DstC": "72cm"}`
	var err error
	var cfg cfgType

	show := func(lfStr, rtStr string, args ...interface{}) {
		fmt.Println(util.StrDotPairFormat(24, lfStr, rtStr, args...))
	}

	showDst := func(lfStr string, val util.Distance) {
		show(lfStr, "%s", util.Float64ToStrSig(float64(val), ".", ",", 3, 3))
	}

	err = ioutil.WriteFile(fileStr, []byte(str), 0644)
	if err == nil {
		err = util.JSONGetFile(fileStr, &cfg)
		if err == nil {
			show("Addr", "%s", cfg.Addr)
			showDst("DstA", cfg.DstA)
			showDst("DstB", cfg.DstB)
			showDst("DstC", cfg.DstC)
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
	// Addr......10.20.30.40:50
	// DstA................3.50
	// DstB................1.00
	// DstC................28.3
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
	fmt.Printf("Geometric %.6f, Arithmetic %.6f",
		util.GeometricMean(list), util.ArithmeticMean(list))
	// Output:
	// Geometric 3.053326, Arithmetic 3.066667
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
