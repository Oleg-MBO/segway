package client

type espConf int

const (

	// SendGyroAngle is used to enable sending angle from gyroscope
	SendGyroAngle espConf = iota

	// SendAccAngle is used to enable sending angle form accelerometer
	SendAccAngle

	// SendAcc is used to enable sending axeleration form accelerometer
	SendAcc
)
