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
