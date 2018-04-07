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
	"math"

	"github.com/jung-kurt/etc/go/util"
)

// This example demonstrates float formatting.
// page-breaking.
func Example_Format64ToStrSig() {
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
