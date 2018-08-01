package main

import (
	"time"
)

// AperiodicLink represent aperiodic link
// https://all4study.ru/modelirovanie/aperiodicheskoe-zveno-v-forme-raznostnogo-uravneniya.html
type AperiodicLink struct {
	isFirst      bool
	value        float64
	tConst       float64
	prevCalcTime time.Time
}

// NewAperiodicLink return *AperiodicLink structure
func NewAperiodicLink(T float64) *AperiodicLink {
	return &AperiodicLink{
		isFirst: true,
		tConst:  T,
	}
}

// Update is used for update current value
func (al *AperiodicLink) Update(v float64) float64 {
	if al.isFirst {
		al.value = v
		al.prevCalcTime = time.Now()
		al.isFirst = false
		return v
	}
	dt := time.Now().Sub(al.prevCalcTime).Seconds()
	al.value = 1/al.tConst*dt*(v-al.value) + al.value
	al.prevCalcTime = time.Now()
	return al.value
}

// GetValue is used for get current value
func (al *AperiodicLink) GetValue() float64 {
	return al.value
}

// SetTConst is used for set time contant T
func (al *AperiodicLink) SetTConst(T float64) {
	al.tConst = T
}
