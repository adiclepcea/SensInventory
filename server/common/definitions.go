package common

//A sensor has an address and several registers
//each register has a type and a loccation.
//A reading is done by specifying a sensor address,
//the type of reading (which implies the type of registers)
//and the number of regsters to read
//You can group the values from several registers into a ReadGroup
//This ReadGroup will transform the values from this sensors into a
//resulting value
import "time"

//Constants for register types and
//calculated result types
const (
	//register types
	Coil          = "coil"
	Input         = "input"
	Holding       = "holding"
	InputDiscrete = "inputDiscrete"
	//calculated result types
	Float32 = "float32"
	Long    = "long"
	Byte    = "byte"
)

//Sensor represents a sensor with several configured registers
type Sensor struct {
	Address     uint8       `json:"address"` //485 address
	Description string      `json:"description,omitempty"`
	Registers   []Register  `json:"registers"`
	ReadGroups  []ReadGroup `json:"readGroup"`
}

//Register represents a register (coil, holding, input, input discrete)
type Register struct {
	Name            string `json:"name,omitempty"`
	RegisterAddress uint16 `json:"registerAddress"`
	RegisterType    string `json:"registerType"`
}

//Reading represents a reading from a Sensor
//and the representation of its values
type Reading struct {
	Sensor           uint8         `json:"sensor"`
	StartLocation    uint16        `json:"startLocation"`
	Count            uint16        `json:"count"`
	ReadValues       []uint16      `json:"readValues"`
	Time             time.Time     `json:"time"`
	CalculatedValues []interface{} `json:"calculatedValues"`
}

//ReadGroup uses the values of a group of registers
//to calculate a resultant value
type ReadGroup struct {
	SensorAddress int    `json:"sensorAddress"`
	Addresses     []int  `json:"addresses"`
	ResultType    string `json:"resultType"`
	ReadGroupCalculation
}

//ReadGroupCalculation defines the methods needed to
//obtain the value defined by a Grouping
type ReadGroupCalculation interface {
	Calculate(Reading) (interface{}, error)
}
