package control

// Outputer interface represent element whitch has output
type Outputer interface {
	Output() float64
}

// Elementer interface represent one element somthening like PID repulator
type Elementer interface {
	Update(dt float64, inputs ...float64) float64
	Outputer
}

// type Piper interface {

// }
