package tool

// AperiodicLink represent aperiodic link
// https://all4study.ru/modelirovanie/aperiodicheskoe-zveno-v-forme-raznostnogo-uravneniya.html
type AperiodicLink struct {
	isFirst bool
	value   float64
	tConst  float64
	// prevCalcTime time.Time
}

// NewAperiodicLink return *AperiodicLink structure
func NewAperiodicLink(T float64) *AperiodicLink {
	return &AperiodicLink{
		isFirst: true,
		tConst:  T,
	}
}

// Update is used for update current value
// https://ru.wikipedia.org/wiki/Апериодическое_звено
func (al *AperiodicLink) Update(dt, v float64) {
	if al.isFirst {
		al.isFirst = false
		al.value = v
		// al.prevCalcTime = time.Now()
		// al.isFirst = false
		return
	}
	// dt := time.Now().Sub(al.prevCalcTime).Seconds()
	al.value = al.tConst*(v-al.value)*dt + al.value
	// al.prevCalcTime = time.Now()
}

// Output is used for get current value
func (al *AperiodicLink) Output() float64 {
	return al.value
}

// SetTConst is used for set time contant T
func (al *AperiodicLink) SetTConst(T float64) {
	al.tConst = T
}
