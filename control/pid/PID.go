package pid

// PID represent PID regulator
type PID struct {
	P *PRegul
	D *DRegul
	I *IRegul
}

// NewPID create new PID regulator
func NewPID(pRegul *PRegul, dRegul *DRegul, iRegul *IRegul) *PID {
	if pRegul == nil {
		pRegul = NewPRegul(0)
	}
	if dRegul == nil {
		dRegul = NewDRegul(0, 0)
	}
	if iRegul == nil {
		iRegul = NewIRegul(0, 0)
	}

	return &PID{P: pRegul, D: dRegul, I: iRegul}
}

type ConfNewPIDSimple struct {
	P, D, DLim, I, Ilim float64
}

// NewPIDSimple create clasic new PID regulator from config
func NewPIDSimple(conf ConfNewPIDSimple) *PID {

	p := NewPRegul(conf.P)
	d := NewDRegul(conf.D, conf.DLim)
	i := NewIRegul(conf.I, conf.Ilim)
	NewPID(p, d, i)

	return NewPID(p, d, i)
}

// Update using for update value
func (pid *PID) Update(dt float64, input float64) {
	pid.P.Update(dt, input)
	pid.D.Update(dt, input)
	pid.I.Update(dt, input)
}

// Output using for get value
func (pid *PID) Output() float64 {
	return pid.P.Output() + pid.D.Output() + pid.I.Output()
}
