package main

// GyroDataStruct is used for save filtered data from gyroscop
type GyroDataStruct struct {
	rotX *AperiodicLink
	rotY *AperiodicLink
	rotZ *AperiodicLink
}

// NewGyroData return *GyroData to save data from gyroscop
func NewGyroData(aperiodicT float64) *GyroDataStruct {
	return &GyroDataStruct{
		rotX: NewAperiodicLink(aperiodicT),
		rotY: NewAperiodicLink(aperiodicT),
		rotZ: NewAperiodicLink(aperiodicT),
	}
}

// Update is used for update data from gyroscop
// gets angle rotation by axes
func (g *GyroDataStruct) Update(rotX, rotY, rotZ float64) (rX, rY, rZ float64) {
	rX = g.rotX.Update(rotX)
	rY = g.rotY.Update(rotY)
	rZ = g.rotZ.Update(rotZ)
	return
}

// GetValues is used for get saved values in structures
func (g *GyroDataStruct) GetValues() (rotX, rotY, rotZ float64) {
	rotX = g.rotX.GetValue()
	rotY = g.rotY.GetValue()
	rotZ = g.rotZ.GetValue()
	return
}

// SetTConst is used for change T constant in aperiodics links
func (g *GyroDataStruct) SetTConst(T float64) {
	g.rotX.SetTConst(T)
	g.rotY.SetTConst(T)
	g.rotZ.SetTConst(T)
}
