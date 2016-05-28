package server

import (
	"errors"
	"fmt"
	"log"

	"github.com/adiclepcea/SensInventory/server/common"
)

const minAddress int = 1
const maxAddress int = 32
const initAddress int = 100

//ConfigProvider contains the configuration for the server
type ConfigProvider struct {
	Sensors map[int]common.Sensor
}

//NewConfigProvider creates a new ConfigProvider
func (ConfigProvider) NewConfigProvider() *ConfigProvider {
	c := ConfigProvider{}
	c.Sensors = make(map[int]common.Sensor)
	return &c
}

//IsSensorAddressTaken checks to see if there is already a slave with
//the passed address defined
func (configProvider *ConfigProvider) IsSensorAddressTaken(address int) bool {
	if _, ok := configProvider.Sensors[address]; ok {
		return true
	}

	return false
}

//IsSensorValid checks to see if the sensot passed in is valid
func (configProvider *ConfigProvider) IsSensorValid(sensor common.Sensor) error {
	if (sensor.Address < minAddress || sensor.Address > maxAddress) && sensor.Address != initAddress {
		err := fmt.Errorf("The sensor adresses must be between %d and %d or exactly %d", minAddress, maxAddress, initAddress)
		log.Println(err.Error())
		return err
	}

	if len(sensor.ConfiguredValues) == 0 {
		err := errors.New("The sensor must have at least one configured address")
		log.Println(err.Error())
		return err
	}

	return nil
}

//AddSensor adds a new sensor that the server should interrogate
func (configProvider *ConfigProvider) AddSensor(sensor common.Sensor) error {
	if err := configProvider.IsSensorValid(sensor); err != nil {
		log.Println((err).Error())
		return err
	}
	if configProvider.IsSensorAddressTaken(sensor.Address) {
		err := fmt.Errorf("AddSensor. A sensor with address %d has already been registered", sensor.Address)
		log.Println(err.Error())
		return err
	}

	configProvider.Sensors[sensor.Address] = sensor

	return nil
}

//RemoveSensorByAddress removes the sensor having the specified address
//from the collection of sensors that the server interrogates
func (configProvider *ConfigProvider) RemoveSensorByAddress(address int) error {
	if !configProvider.IsSensorAddressTaken(address) {
		err := fmt.Errorf("No sensor with %d address is registered", address)
		log.Println(err.Error())
		return err
	}

	delete(configProvider.Sensors, address)

	return nil
}

//RemoveSensor removes the specified sensor from the collection of sensors that
//the server interrogates
func (configProvider *ConfigProvider) RemoveSensor(sensor common.Sensor) error {

	return configProvider.RemoveSensorByAddress(sensor.Address)

}

//GetSensorByAddress returns the sensor with the given address
func (configProvider *ConfigProvider) GetSensorByAddress(address int) (*common.Sensor, error) {
	var sensor common.Sensor
	var ok bool

	if sensor, ok = configProvider.Sensors[address]; !ok {
		err := fmt.Errorf("No sensor with %d address is registered", address)
		log.Println(err.Error())
		return nil, err
	}

	return &sensor, nil
}

//ChangeSensorAddress changes the address of the sensor that currently has
//address "addressBefore" with the "addressAfter"
func (configProvider *ConfigProvider) ChangeSensorAddress(addressBefore int, addressAfter int) error {
	sensorBefore, err := configProvider.GetSensorByAddress(addressBefore)
	if err != nil {
		return err
	}

	if configProvider.IsSensorAddressTaken(addressAfter) {
		err := fmt.Errorf("There is allready a sensor registered with address %d", addressAfter)
		log.Println(err.Error())
		return err
	}

	if err := configProvider.RemoveSensorByAddress(sensorBefore.Address); err != nil {
		return err
	}

	sensorBefore.Address = addressAfter

	return configProvider.AddSensor(*sensorBefore)

}

//ChangeSensor changes the sensor having address "address" to be similar with
//the sensor "after"
func (configProvider *ConfigProvider) ChangeSensor(address int, after common.Sensor) error {
	var sensorBefore *common.Sensor
	var err error
	if sensorBefore, err = configProvider.GetSensorByAddress(address); err != nil {
		return err
	}

	if err = configProvider.IsSensorValid(after); err != nil {
		return err
	}

	sensorBefore.Description = after.Description
	sensorBefore.ConfiguredValues = after.ConfiguredValues
	configProvider.Sensors[sensorBefore.Address] = *sensorBefore

	return nil

}

//GetSensors returns a map of the sensor addresses mapped to the sensors themselves
func (configProvider *ConfigProvider) GetSensors() map[int]common.Sensor {
	return configProvider.Sensors
}
