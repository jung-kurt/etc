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
	"math/rand"
	"net"
	"os"
	"strings"
	"testing"

	"github.com/jung-kurt/etc/go/util"
)

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
