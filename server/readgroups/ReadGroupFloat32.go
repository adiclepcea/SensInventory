package readgroups

//this follows the conversion rules according to IEEE 754
//to convert between 4 bytes (32 bits) to a float32

import (
	"fmt"
	"log"
	"math"

	"github.com/adiclepcea/SensInventory/server/common"
)

//ReadGroupFloat32 calculates a float 32 value from 2 holding or input registers
type ReadGroupFloat32 struct {
	common.ReadGroup
}

//NewReadGroup initializes a ReadGroupFloat32
func (ReadGroupFloat32) NewReadGroup(sensorAddress uint8,
	startLocation uint16) (*ReadGroupFloat32, error) {
	rgf32 := ReadGroupFloat32{common.ReadGroup{}}
	rgf32.SensorAddress = sensorAddress
	rgf32.StartLocation = startLocation
	rgf32.ResultType = common.Float32
	return &rgf32, nil
}

//Calculate performs the transformation between registries values and Float32
func (rgf32 ReadGroupFloat32) Calculate(reading common.Reading) (interface{}, error) {
	log.Printf("SensorAddress %d \n StartLocation %d\n", rgf32.SensorAddress, rgf32.StartLocation)
	if reading.StartLocation > rgf32.StartLocation {
		log.Printf("Reading start location (%d > %d)",
			reading.StartLocation,
			rgf32.StartLocation)
		return nil, fmt.Errorf(
			"Could not calculate the value, Reading startLocation (%d > %d)",
			reading.StartLocation, rgf32.StartLocation)
	}

	if reading.StartLocation+reading.Count-1 < rgf32.StartLocation+1 {
		log.Printf("Reading end location (%d < %d)",
			reading.StartLocation+reading.Count-1,
			rgf32.StartLocation+1)
		return nil, fmt.Errorf(
			"Could not calculate the value, Reading startLocation + count (%d < %d)",
			reading.StartLocation+reading.Count-1,
			rgf32.StartLocation+1)
	}

	if reading.Type != common.Holding && reading.Type != common.Input {
		log.Printf("Reading type %s should be Holding or Input", reading.Type)
		return nil, fmt.Errorf("Reading type %s should be Holding or Input", reading.Type)
	}
	poz1 := rgf32.StartLocation - reading.StartLocation
	var x uint32
	byte1 := uint32(reading.ReadValues[poz1+1] & 0xFF)
	byte2 := uint32((reading.ReadValues[poz1+1] & 0xFF00) / 0xFF)
	byte3 := uint32(reading.ReadValues[poz1] & 0xFF)
	byte4 := uint32((reading.ReadValues[poz1] & 0xFF00) / 0xFF)
	x = (byte1 << 24) + (byte2 << 16) + (byte3 << 8) + byte4
	log.Printf("Value x=%x", x)

	return math.Float32frombits(x), nil
}
