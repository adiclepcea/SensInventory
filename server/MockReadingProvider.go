package server

import (
	"math/rand"
	"time"

	"github.com/adiclepcea/SensInventory/server/common"
	"github.com/adiclepcea/SensInventory/server/configprovider"
)

//MockReadingProvider is a mock provider for reading sensors. Used in tests
type MockReadingProvider struct {
	Conf configprovider.MockConfigProvider
	ReadingProvider
}

//NewReadingProvider returns a new reading provider having the configuration
//provided by "cp"
func (mockReadingProvider MockReadingProvider) NewReadingProvider(cp *configprovider.MockConfigProvider) *MockReadingProvider {
	mockReadingProvider.Conf = *cp

	return &mockReadingProvider
}

func (mockReadingProvider *MockReadingProvider) getRandValuesForSensor(sensor common.Sensor) []uint16 {
	rez := make([]uint16, len(sensor.Registers))
	for i, confValue := range sensor.Registers {
		rez[i] = mockReadingProvider.getRandValueForConfiguredValue(confValue)
	}

	return rez
}

func (mockReadingProvider *MockReadingProvider) getRandValueForConfiguredValue(configuredValue common.Register) uint16 {

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	return (uint16)(r1.Intn(255))
}

//GetReading returns a mock random read from the sensor having address "address"
func (mockReadingProvider *MockReadingProvider) GetReading(address uint8) (*common.Reading, error) {

	sensor, err := mockReadingProvider.Conf.GetSensorByAddress(address)

	if err != nil {
		return nil, err
	}

	reading := common.Reading{Sensor: sensor.Address, Time: time.Now().Format(common.TimeFormat),
		ReadValues: mockReadingProvider.getRandValuesForSensor(*sensor)}
	return &reading, nil
}
