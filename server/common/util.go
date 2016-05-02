package common

import "time"

type ReadValue struct {
	Name            string
	RegisterAddress int
	RegisterLength  int8 //1 or 2
}

type Reading struct {
	ReadSensor Sensor
	Time       time.Time
}

type Sensor struct {
	Address     int //485 address
	Description string
	ReadValues  []ReadValue
}
