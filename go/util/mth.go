package util

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
