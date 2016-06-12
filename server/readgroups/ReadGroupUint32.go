package readgroups

//this follows converts the values of two uint16 to a uint32
import (
	"fmt"
	"log"

	"github.com/adiclepcea/SensInventory/server/common"
)

//ReadGroupUint32 calculates a uint32 value from 2 holding or input registers
type ReadGroupUint32 struct {
	common.ReadGroup
}

//NewReadGroup initializes a ReadGroupUint32
func (ReadGroupUint32) NewReadGroup(sensorAddress uint8,
	startLocation uint16) (*ReadGroupUint32, error) {
	rgf32 := ReadGroupUint32{common.ReadGroup{}}
	rgf32.SensorAddress = sensorAddress
	rgf32.StartLocation = startLocation
	rgf32.ResultType = common.Uint32
	return &rgf32, nil
}

//Calculate performs the transformation between registries values and Uint32
func (rgf32 ReadGroupUint32) Calculate(reading common.Reading) (interface{}, error) {
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
	xa := uint32(reading.ReadValues[poz1])
	xb := uint32(reading.ReadValues[poz1+1])
	x = (xa << 16) + xb

	return x, nil
}
