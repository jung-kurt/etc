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

package regression_test

import (
	"fmt"
	"testing"

	"github.com/jung-kurt/etc/go/regression"
	"github.com/jung-kurt/etc/go/util"
)

// Error should be returned when Center is called with zero-length list.
func Test01(t *testing.T) {
	_, _, err := regression.Center([]util.PairType{}, 10)
	if err == nil {
		t.Fatalf("expecting error with zero-length point list")
	}
}

// linearRegressionExample demonstrates fitting a straight line to observation
// points
func ExampleLinearFit() {
	var (
		yList = []float64{12.2, 13.6, 15.9, 18.3}
		xList = []float64{1.0, 2.0, 3.0, 4.0}
		le    regression.LinearFitType
	)

	le = regression.LinearFit(xList, yList)
	fmt.Printf("%s\n", le)
	le.Eq.Intercept = -le.Eq.Intercept
	fmt.Printf("%s\n", le)
	// Output:
	// y(x) = 2.06 * x + 9.85 (r squared 0.987, RMS 0.266)
	// y(x) = 2.06 * x - 9.85 (r squared 0.987, RMS 0.266)
}

func ExampleCenter() {
	var err error
	var x, y float64
	var pairs = []util.PairType{
		{X: 376.25, Y: 519.21},
		{X: 387.5, Y: 503.0749999999998},
		{X: 398.75, Y: 487.1149999999998},
		{X: 410, Y: 471.1199999999999},
		{X: 421.25, Y: 459.57000000000016},
		{X: 432.5, Y: 451.55499999999984},
		{X: 443.75, Y: 444.9749999999999},
		{X: 455, Y: 436.6100000000001},
		{X: 466.25, Y: 425.27},
		{X: 477.5, Y: 419.17999999999984},
		{X: 488.75, Y: 416.30999999999995},
		{X: 500, Y: 412.91499999999996},
		{X: 511.25, Y: 409.9050000000002},
		{X: 522.5, Y: 407.98},
		{X: 533.75, Y: 407.03499999999985},
		{X: 545, Y: 405.91499999999996},
		{X: 556.25, Y: 405.32000000000016},
		{X: 567.5, Y: 404.9000000000001},
		{X: 578.75, Y: 406.5799999999999},
		{X: 590, Y: 409.625},
		{X: 601.25, Y: 409.3449999999998},
		{X: 612.5, Y: 413.05499999999984},
		{X: 623.75, Y: 418.30499999999984},
		{X: 635, Y: 423.0300000000002},
		{X: 646.25, Y: 427.05499999999984},
		{X: 657.5, Y: 431.9200000000001},
		{X: 668.75, Y: 443.4000000000001},
		{X: 680, Y: 447.80999999999995},
		{X: 691.25, Y: 452.42999999999984},
		{X: 702.5, Y: 463.4200000000001},
		{X: 713.75, Y: 484.5250000000001},
		{X: 725, Y: 498.98},
		{X: 736.25, Y: 516.3049999999998},
	}

	f3 := func(val float64) string {
		return util.Float64ToStrSig(val, ".", ",", 3, 3)
	}

	x, y, err = regression.Center(pairs, 200)
	if err == nil {
		fmt.Printf("downhill simplex: center [%s, %s]", f3(x), f3(y))
	} else {
		fmt.Printf("downhill simplex error: %s", err)
	}
	// Output:
	// downhill simplex: center [554, 263]
}
