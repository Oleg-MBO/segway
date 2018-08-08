package pid

import (
	"math"
)

// DRegul represent derivative regulator
type DRegul struct {
	d         float64
	prevValue float64
	deriv     float64
	limits    float64
}

// NewDRegul return derivative regulator
func NewDRegul(dkoef, limits float64) *DRegul {
	return &DRegul{d: dkoef, limits: limits}
}

// Update using for update value
func (d *DRegul) Update(dt float64, input float64) {

	d.deriv = d.d * (input - d.prevValue) / dt
	d.prevValue = input
	if d.limits != 0 && math.Abs(d.deriv) > d.limits {
		d.deriv = d.limits * sign(d.deriv)
	}
}

// Output using for get value
func (d *DRegul) Output() float64 {
	return d.deriv
}
