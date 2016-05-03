package common

import "time"

type ConfiguredValue struct {
	Name            string `json:"name"`
	RegisterAddress int    `json:"registerAddress"`
	RegisterLength  int8   `json:"registerLength"` //1 or 2
}

type Reading struct {
	ReadSensor Sensor        `json:"readSensor"`
	ReadValues []interface{} `json:"readValues"`
	Time       time.Time     `json:"time"`
}

type Sensor struct {
	Address          int               `json:"address"` //485 address
	Description      string            `json:"description"`
	ConfiguredValues []ConfiguredValue `json:"configuredValues"`
}
