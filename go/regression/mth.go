package regression

import (
	"errors"
	"fmt"
	"math"

	"github.com/jung-kurt/etc/go/util"
	"gonum.org/v1/gonum/optimize"
	"gonum.org/v1/gonum/stat"
)

// LinearFitType groups together the slope, intercept, coefficient of
// determination, and root-mean-square average deviation of a regression line.
type LinearFitType struct {
	Eq       util.LinearEquationType
	RSquared float64
	RMS      float64
}

func f3(val float64) string {
	return util.Float64ToStrSig(val, ".", ",", 3, 3)
}

// Strings implements the fmt Stringer interface.
func (fit LinearFitType) String() string {
	var b float64
	var op string

	b = fit.Eq.Intercept
	if b < 0 {
		b = -b
		op = "-"
	} else {
		op = "+"
	}
	return fmt.Sprintf("y(x) = %s * x %s %s (r squared %s, RMS %s)",
		f3(fit.Eq.Slope), op, f3(b), f3(fit.RSquared), f3(fit.RMS))
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

// Center uses the downhill simplex method to calculate the center of a circle
// of known radius to a set of observed points specified by pairs.
func Center(pairs []util.PairType, radius float64) (x, y float64, err error) {
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
		lf, rt, tp, bt := util.BoundingBox(pairs)
		res, err = DownhillSimplex(centerFnc, []float64{(lf + rt) / 2, bt + bt - tp}, (rt-lf)/32, 1.2)
		if err == nil {
			x = res[0]
			y = res[1]
		}
	} else {
		err = errors.New("insufficient number of points to locate center")
	}
	return
}

// LinearFit return the slope, intercept, r-squared values, and RMS value for
// the least squares regression fit of the points specified by xList and yList.
func LinearFit(xList, yList []float64) (le LinearFitType) {
	le.Eq.Intercept, le.Eq.Slope = stat.LinearRegression(xList, yList, nil, false)
	le.RSquared = stat.RSquared(xList, yList, nil, le.Eq.Intercept, le.Eq.Slope)
	le.RMS = util.RootMeanSquareLinear(xList, yList, le.Eq.Intercept, le.Eq.Slope)
	// logf("RSquared: Gonum %.3f, local %.3f", le.RSquared, rSquared(xList, yList, le.Eq.Intercept, le.Eq.Slope))
	return
}
