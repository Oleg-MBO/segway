package client

// EspData represent data which
type EspData struct {
	// time from start esp
	Milis uint64

	AngleX, AngleY, AngleZ float64

	// axceleration by axes
	AccX, AccY, AccZ float64

	// angle from axceleration
	AAngleX, AAngleY, AAngleZ float64

	// gyro angle data
	GyroX, GyroY, GyroZ float64
}

// HasAngles check if all angles != 0
func (ed *EspData) HasAngles() bool {
	return ed.AngleX != 0 || ed.AngleY != 0 || ed.AngleZ != 0
}

// HasAcc check if any of axceleration data != 0
func (ed *EspData) HasAcc() bool {
	return ed.AccX != 0 || ed.AccY != 0 || ed.AccZ != 0
}

// HasAAngles check if any of angle from axceleration != 0
func (ed *EspData) HasAAngles() bool {
	return ed.AAngleX != 0 || ed.AAngleY != 0 || ed.AAngleZ != 0
}

// HasGyros check if any of gyros data angles != 0
func (ed *EspData) HasGyros() bool {
	return ed.GyroX != 0 || ed.GyroY != 0 || ed.GyroZ != 0
}
