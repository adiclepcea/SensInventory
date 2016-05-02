package server

import (
	"math/rand"
	"time"

	"github.com/adiclepcea/SensInventory/server/common"
)

type MockReadingProvider struct {
	Conf *ConfigProvider
}

func (this MockReadingProvider) NewReadingProvider(cp *ConfigProvider) IReadingProvider {
	this.Conf = cp

	return &this
}

func (this *MockReadingProvider) getRandValuesForSensor(sensor common.Sensor) []interface{} {
	rez := make([]interface{}, len(sensor.ConfiguredValues))
	for i, confValue := range sensor.ConfiguredValues {
		rez[i] = this.getRandValueForConfiguredValue(confValue)
	}

	return rez
}

func (this *MockReadingProvider) getRandValueForConfiguredValue(configuredValue common.ConfiguredValue) interface{} {

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	if configuredValue.RegisterLength == 1 {
		return r1.Intn(255)
	} else if configuredValue.RegisterLength == 2 {
		return r1.Float32()
	}

	return 0 //for now we only provide mock values for integers - for single registry reads and for floats - for double registry reads
}

func (this *MockReadingProvider) GetReading(address int) (*common.Reading, *error) {

	sensor, err := this.Conf.GetSensorByAddress(address)

	if err != nil {
		return nil, err
	}

	reading := common.Reading{ReadSensor: *sensor, Time: time.Now(), ReadValues: this.getRandValuesForSensor(*sensor)}
	return &reading, nil
}
