package common

//A sensor has an address and several registers
//each register has a type and a location.
//A reading is done by specifying a sensor address,
//the type of reading (which implies the type of registers)
//and the number of registers to read
//You can group the values from several registers into a ReadGroup
//This ReadGroup will transform the values from this sensors into a
//resulting value

//Constants for register types and
//calculated result types
const (
	TimeFormat = "2006-01-02T15:04:05"
	//register types
	Coil          = "coil"
	Input         = "input"
	Holding       = "holding"
	InputDiscrete = "inputDiscrete"
	//calculated result types
	Float32 = "float32"
	Uint32  = "uint32"
	Int32   = "int32"
)

//Sensor represents a sensor with several configured registers
type Sensor struct {
	Address     uint8       `json:"address"` //485 address
	Description string      `json:"description,omitempty"`
	Registers   []Register  `json:"registers"`
	ReadGroups  []ReadGroup `json:"readGroups"`
}

//Register represents a register (coil, holding, input, input discrete)
type Register struct {
	Name     string `json:"name,omitempty"`
	Location uint16 `json:"location"`
	Type     string `json:"type"`
}

//Reading represents a reading from a Sensor
//and the representation of its values
type Reading struct {
	Sensor           uint8                  `json:"sensor"`
	Type             string                 `json:"type"`
	StartLocation    uint16                 `json:"startLocation"`
	Count            uint16                 `json:"count"`
	ReadValues       []uint16               `json:"readValues"`
	Time             string                 `json:"time"`
	CalculatedValues map[string]interface{} `json:"calculatedValues"`
}

//InitCalculatedValues initiates the map that will hold the calculated values
//for a reading
func (reading *Reading) InitCalculatedValues() {
	reading.CalculatedValues = make(map[string]interface{})
}

//ReadGroup uses the values of a group of registers
//to calculate a resultant value
type ReadGroup struct {
	SensorAddress   uint8  `json:"sensorAddress"`
	StartLocation   uint16 `json:"startLocation"`
	ResultType      string `json:"resultType"`
	ReadGroupWorker `json:"-"`
}

//ReadGroupWorker defines the methods needed to
//initialize a ReadGroup and obtain the value defined by it
type ReadGroupWorker interface {
	Calculate(*Reading) (interface{}, error)
	NewReadGroup(uint8, uint16) (*ReadGroup, error)
}
