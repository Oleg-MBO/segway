package pid

// PID represent PID regulator
type PID struct {
	pRegul *PRegul
	dRegul *DRegul
	iRegul *IRegul
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

	return &PID{pRegul: pRegul, dRegul: dRegul, iRegul: iRegul}
}

// Update using for update value
func (pid *PID) Update(dt float64, input float64) {
	pid.pRegul.Update(dt, input)
	pid.dRegul.Update(dt, input)
	pid.iRegul.Update(dt, input)
}

// Output using for get value
func (pid *PID) Output() float64 {
	return pid.pRegul.Output() + pid.dRegul.Output() + pid.iRegul.Output()
}
