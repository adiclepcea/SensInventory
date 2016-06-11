package configprovider

import "github.com/adiclepcea/SensInventory/server/common"

//ConfigProvider is a prototype for a configuration manager
type ConfigProvider interface {
	NewConfigProvider(params ...string) (*ConfigProvider, error)
	SetAddressLimits(minAddress uint8, maxAddress uint8) error
	IsSensorAddressTaken(address uint8) (bool, error)
	IsSensorValid(sensor common.Sensor) error
	AddSensor(sensor common.Sensor) error
	RemoveSensorByAddress(address uint8) error
	RemoveSensor(sensor common.Sensor) error
	GetSensorByAddress(address uint8) (*common.Sensor, error)
	ChangeSensorAddress(addressBefore uint8, addressAfter uint8) error
	ChangeSensor(address uint8, after common.Sensor) error
	GetSensors() map[uint8]common.Sensor
}
