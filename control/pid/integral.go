package pid

import (
	"math"
)

// IRegul represent integral regulator
type IRegul struct {
	i        float64
	integral float64
	limits   float64
}

// NewIRegul return integral regulator
func NewIRegul(ikoef, limits float64) *IRegul {
	return &IRegul{i: ikoef, limits: limits}
}

// Update using for update value
func (i *IRegul) Update(dt float64, inputs float64) {
	i.integral += i.i * inputs * dt
	if i.limits != 0 && math.Abs(i.integral) > i.limits {
		i.integral = i.limits * sign(i.integral)
	}
}

// Output using for get value
func (i *IRegul) Output() float64 {
	return i.integral
}
