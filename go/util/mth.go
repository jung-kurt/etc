package util

import (
	"math"

	"gonum.org/v1/gonum/optimize"
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

// DownhillSimplex finds the lowest value reported by fnc. The number of
// dimensions is specified by the number of elements in init, the initial
// location. The same number of elements will be passed to fnc, the callback
// function, when probing. The final result, if err is nil, will contain this
// number of elements as well. Two parameters can be adjusted to avoid
// converging on suboptimal local minima: len specifies the simplex size, and
// expansion (some value greater than 1) specifies the multiplier used when
// expanding the simplex.
func DownhillSimplex(fnc func(x []float64) float64, init []float64, len, expansion float64) (res []float64, err error) {
	var prb optimize.Problem
	var r *optimize.Result

	prb.Func = fnc
	r, err = optimize.Local(prb, init, nil, &optimize.NelderMead{
		SimplexSize: len,
		Expansion:   expansion, // 1.25,
	})
	if err == nil {
		res = r.X
	}
	return
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

// Center uses the downhill simplex method to calculate the center of a circle
// of known radius to a set of observed points specified by pairs.
func Center(pairs []PairType, radius float64) (x, y float64, err error) {
	var res []float64
	var count = float64(len(pairs))

	centerFnc := func(x []float64) (val float64) {
		var alpha, beta float64
		alpha = x[0]
		beta = x[1]
		for _, pr := range pairs {
			xa := pr.X - alpha
			yb := pr.Y - beta
			e := math.Sqrt(xa*xa+yb*yb) - radius
			val += e * e
		}
		return
	}

	if count > 0 {
		lf, rt, tp, bt := BoundingBox(pairs)
		res, err = DownhillSimplex(centerFnc, []float64{(lf + rt) / 2, bt + bt - tp}, (rt-lf)/32, 1.2)
		if err == nil {
			x = res[0]
			y = res[1]
		}
	} else {
		err = errf("insufficient number of points to locate center")
	}
	return
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

// Perpendicular returns an equation that is perpendicular to eq and intersects
// it at x.
func (eq LinearEquationType) Perpendicular(x float64) (p LinearEquationType) {
	if eq.Slope != 0 {
		p.Slope = -1 / eq.Slope
		p.Intercept = (eq.Slope*x + eq.Intercept) - p.Slope*x
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
		if weight < 0 {
			weight = 0
		} else if weight > 1 {
			weight = 1
		}
		oldWeight := avg.weight
		avg.weight += weight
		// 	val = weight*val + (1-weight)*avg.value
		// 	avg.weight += 1
		avg.value = (avg.value*oldWeight + val*weight) / avg.weight
	}
}

// Value returns the current average of submitted values
func (avg AverageType) Value() float64 {
	return avg.value
}
