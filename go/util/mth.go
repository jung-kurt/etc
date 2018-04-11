package util

import (
	"fmt"
	"math"
)

// PairType defines a two-dimensional coordianate.
type PairType struct {
	X, Y float64
}

// Cluster breaks apart pairs into zero or more slices that each contain at
// least minPts pairs and have gaps between X values no greater than gapX.
// Elements in pairs must be ordered from low X value to high.
func Cluster(pairs []PairType, minPts int, gapX float64) [][]PairType {
	var ln int
	var clList [][]PairType
	var list []PairType

	list = make([]PairType, 0, 32)

	place := func() {
		// list has contiguous points; discard if fewer than minPts
		if len(list) >= minPts {
			clList = append(clList, list)
		}
		list = make([]PairType, 0, 32)
	}

	for _, pr := range pairs {
		ln = len(list)
		if ln == 0 {
			list = append(list, pr)
		} else if pr.X < list[ln-1].X+gapX {
			list = append(list, pr)
		} else {
			place()
			list = append(list, pr)
		}
	}
	place()
	return clList
}

// BoundingBox returns the smallest and greatest values of X and Y in the
// specified slice of coordinates.
func BoundingBox(pairs []PairType) (lf, rt, tp, bt float64) {
	var xr, yr RangeType

	for j, pr := range pairs {
		xr.Set(pr.X, j == 0)
		yr.Set(pr.Y, j == 0)
	}
	return xr.Min, xr.Max, yr.Max, yr.Min
}

// RangeType holds the minimum and maximum values of a range.
type RangeType struct {
	Min, Max float64
}

// Set adjusts the fields Min and Max such that Min holds the smallest value
// encountered and Max the largest. If init is true, val is assigned to both
// Min and Max. If init is false, val is assigned to Min only if it is smaller,
// and val is assigned to Max only if it is greater.
func (r *RangeType) Set(val float64, init bool) {
	if init {
		r.Min = val
		r.Max = val
	} else {
		if val < r.Min {
			r.Min = val
		}
		if val > r.Max {
			r.Max = val
		}
	}
}

// LinearEquationType describes a line with its slope and intercept
type LinearEquationType struct {
	Slope, Intercept float64
}

// String implements the fmt Stringer interface
func (eq LinearEquationType) String() string {
	var op string
	switch math.Signbit(eq.Intercept) {
	case true:
		op = "-"
		eq.Intercept = -eq.Intercept
	default:
		op = "+"
	}
	return fmt.Sprintf("y = %f x %s %f", eq.Slope, op, eq.Intercept)
}

// Perpendicular returns an equation that is perpendicular to eq and intersects
// it at x.
func (eq LinearEquationType) Perpendicular(x float64) (p LinearEquationType) {
	if eq.Slope != 0 {
		p.Slope = -1 / eq.Slope
		p.Intercept = (eq.Slope*x + eq.Intercept) - p.Slope*x
	}
	return
}

// PerpendicularPoint returns an equation that is perpendicular to eq and
// includes the point specified by x and y.
func (eq LinearEquationType) PerpendicularPoint(x, y float64) (p LinearEquationType) {
	if eq.Slope != 0 {
		p.Slope = -1 / eq.Slope
		p.Intercept = y - p.Slope*x
	}
	return
}

// DistanceToPoint returns the shortest distance from the specified point to
// eq.
func (eq LinearEquationType) DistanceToPoint(x, y float64) (d float64) {
	if eq.Slope != 0 {
		p := eq.PerpendicularPoint(x, y)
		xi := (eq.Intercept - p.Intercept) / (p.Slope - eq.Slope)
		yi := eq.Slope*xi + eq.Intercept
		dx := xi - x
		dy := yi - y
		d = math.Sqrt(dx*dx + dy*dy)
	} else {
		d = math.Abs(y - eq.Intercept)
	}
	return
}

// LinearY returns the value of the linear function defined by intercept and
// slope at the specified x value.
func LinearY(slope, intercept, x float64) (y float64) {
	y = slope*x + intercept
	return
}

// Linear returns the y-intercept and slope of the straight line joining the
// two specified points. For scaling purposes, associate the arguments as
// follows: x1: observed low value, y1: desired low value, x2: observed high
// value, y2: desired high value.
func Linear(x1, y1, x2, y2 float64) (eq LinearEquationType) {
	if x2 != x1 {
		eq.Slope = (y2 - y1) / (x2 - x1)
		eq.Intercept = y2 - x2*eq.Slope
	}
	return
}

// LinearPointSlope returns the y-intercept of the straight line joining the
// specified arbitrary point (not necessarily an intercept) and the line's
// slope.
func LinearPointSlope(x, y, slope float64) (intercept float64) {
	intercept = y - slope*x
	return
}

// AverageType manages the calculation of a running average
type AverageType struct {
	weight float64
	value  float64
}

// Add adds a value to a running average. weight is quietly constrained to the
// range [0, 1].
func (avg *AverageType) Add(val, weight float64) {
	if weight > 0 {
		if weight > 1 {
			weight = 1
		}
		oldWeight := avg.weight
		avg.weight += weight
		avg.value = (avg.value*oldWeight + val*weight) / avg.weight
	}
}

// Value returns the current average of submitted values
func (avg AverageType) Value() float64 {
	return avg.value
}

// RootMeanSquareLinear returns the RMS value for the specified regression
// variables.
func RootMeanSquareLinear(xList, yList []float64, intercept, slope float64) (rms float64) {
	count := len(xList)
	if count > 0 && count == len(yList) {
		for j := 0; j < count; j++ {
			yEst := xList[j]*slope + intercept
			diff := yEst - yList[j]
			rms += diff * diff
		}
		rms = math.Sqrt(rms / float64(count))
	}
	return
}

// RootMeanSquare returns the RMS value for the specified slice of values. From
// https://en.wikipedia.org/wiki/Root_mean_square: In Estimation theory, the
// root mean square error of an estimator is a measure of the imperfection of
// the fit of the estimator to the data.
func RootMeanSquare(list []float64) (rms float64) {
	count := len(list)
	if count > 0 {
		for _, val := range list {
			rms += val * val
		}
		rms = math.Sqrt(rms / float64(count))
	}
	return
}

// GeometricMean returns the nth root of the product of the n values in list.
func GeometricMean(list []float64) (mean float64) {
	count := len(list)
	if count > 0 {
		mean = 1
		for _, val := range list {
			mean *= val
		}
		if count > 1 {
			mean = math.Pow(mean, float64(1)/float64(count))
		}
	}
	return
}

// ArithmeticMean returns the sum of the values in list divided by the number
// of values.
func ArithmeticMean(list []float64) (mean float64) {
	count := len(list)
	if count > 0 {
		for _, val := range list {
			mean += val
		}
		if count > 1 {
			mean /= float64(count)
		}
	}
	return
}
