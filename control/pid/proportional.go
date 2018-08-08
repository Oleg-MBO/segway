package pid

// PRegul represent proportional repulator
type PRegul struct {
	k, output float64
}

// NewPRegul return proportional repulator
func NewPRegul(k float64) *PRegul {
	return &PRegul{k: k}
}

// Update using for update value
func (p *PRegul) Update(dt float64, inputs float64) {
	p.output = inputs * p.k
}

// Output using for get value
func (p *PRegul) Output() float64 {
	return p.output
}
