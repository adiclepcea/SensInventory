package server

import (
	"math/rand"
	"time"

	"github.com/adiclepcea/SensInventory/server/common"
)

//MockReadingProvider is a mock provider for reading sensors. Used in tests
type MockReadingProvider struct {
	IReadingProvider
	Conf *ConfigProvider
}

//NewReadingProvider returns a new reading provider having the configuration
//provided by "cp"
func (mockReadingProvider MockReadingProvider) NewReadingProvider(cp *ConfigProvider) *MockReadingProvider {
	mockReadingProvider.Conf = cp

	return &mockReadingProvider
}

func (mockReadingProvider *MockReadingProvider) getRandValuesForSensor(sensor common.Sensor) []interface{} {
	rez := make([]interface{}, len(sensor.ConfiguredValues))
	for i, confValue := range sensor.ConfiguredValues {
		rez[i] = mockReadingProvider.getRandValueForConfiguredValue(confValue)
	}

	return rez
}

func (mockReadingProvider *MockReadingProvider) getRandValueForConfiguredValue(configuredValue common.ConfiguredValue) interface{} {

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	if configuredValue.RegisterLength == 1 {
		return r1.Intn(255)
	} else if configuredValue.RegisterLength == 2 {
		return r1.Float32()
	}

	return 0 //for now we only provide mock values for integers - for single registry reads and for floats - for double registry reads
}

//GetReading returns a mock random read from the sensor having address "address"
func (mockReadingProvider *MockReadingProvider) GetReading(address int) (*common.Reading, error) {

	sensor, err := mockReadingProvider.Conf.GetSensorByAddress(address)

	if err != nil {
		return nil, err
	}

	reading := common.Reading{ReadSensor: *sensor, Time: time.Now(), ReadValues: mockReadingProvider.getRandValuesForSensor(*sensor)}
	return &reading, nil
}
