package server

import (
	"errors"
	"fmt"
	"log"

	"github.com/adiclepcea/SensInventory/server/common"
)

const (
	minAddress  uint8 = 1
	maxAddress  uint8 = 32
	initAddress uint8 = 100
)

//DefaultConfigProvider contains the configuration for the server
type DefaultConfigProvider struct {
	Sensors    map[uint8]common.Sensor
	ReadGroups []common.ReadGroup
	ConfigProvider
}

//NewConfigProvider creates a new ConfigProvider
func (DefaultConfigProvider) NewConfigProvider() *DefaultConfigProvider {
	c := DefaultConfigProvider{}
	c.Sensors = make(map[uint8]common.Sensor)
	return &c
}

//IsSensorAddressTaken checks to see if there is already a slave with
//the passed address defined
func (configProvider *DefaultConfigProvider) IsSensorAddressTaken(address uint8) bool {
	if _, ok := configProvider.Sensors[address]; ok {
		return true
	}

	return false
}

//IsSensorValid checks to see if the sensot passed in is valid
func (configProvider *DefaultConfigProvider) IsSensorValid(sensor common.Sensor) error {
	if (sensor.Address < minAddress || sensor.Address > maxAddress) && sensor.Address != initAddress {
		err := fmt.Errorf("The sensor adresses must be between %d and %d or exactly %d", minAddress, maxAddress, initAddress)
		log.Println(err.Error())
		return err
	}

	if len(sensor.Registers) == 0 {
		err := errors.New("The sensor must have at least one configured address")
		log.Println(err.Error())
		return err
	}

	return nil
}

//AddSensor adds a new sensor that the server should interrogate
func (configProvider *DefaultConfigProvider) AddSensor(sensor common.Sensor) error {
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
func (configProvider *DefaultConfigProvider) RemoveSensorByAddress(address uint8) error {
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
func (configProvider *DefaultConfigProvider) RemoveSensor(sensor common.Sensor) error {

	return configProvider.RemoveSensorByAddress(sensor.Address)

}

//GetSensorByAddress returns the sensor with the given address
func (configProvider *DefaultConfigProvider) GetSensorByAddress(address uint8) (*common.Sensor, error) {
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
func (configProvider *DefaultConfigProvider) ChangeSensorAddress(addressBefore uint8, addressAfter uint8) error {
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
func (configProvider *DefaultConfigProvider) ChangeSensor(address uint8, after common.Sensor) error {
	var sensorBefore *common.Sensor
	var err error
	if sensorBefore, err = configProvider.GetSensorByAddress(address); err != nil {
		return err
	}

	if err = configProvider.IsSensorValid(after); err != nil {
		return err
	}

	sensorBefore.Description = after.Description
	sensorBefore.Registers = after.Registers
	configProvider.Sensors[sensorBefore.Address] = *sensorBefore

	return nil

}

//GetSensors returns a map of the sensor addresses mapped to the sensors themselves
func (configProvider *DefaultConfigProvider) GetSensors() map[uint8]common.Sensor {
	return configProvider.Sensors
}
