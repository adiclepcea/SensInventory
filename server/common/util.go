package common

import "time"

//Coil represents the coil type of register
var Coil = 0

//Input represents the input type of register
var Input = 1

//Holding represents the holding type of register
var Holding = 2

//InputDiscrete represents the input discrete type of register
var InputDiscrete = 3

//ConfiguredValue represents a register (coil, holding, input, input discrete)
type ConfiguredValue struct {
	Name            string `json:"name"`
	RegisterAddress int    `json:"registerAddress"`
	RegisterType    int    `json:"registerType"`
}

//Reading represents a reading from a Sensor
type Reading struct {
	ReadSensor Sensor        `json:"readSensor"`
	ReadValues []interface{} `json:"readValues"`
	Time       time.Time     `json:"time"`
}

//Sensor represents a sensor with several configured registers
type Sensor struct {
	Address          int               `json:"address"` //485 address
	Description      string            `json:"description"`
	ConfiguredValues []ConfiguredValue `json:"configuredValues"`
}
