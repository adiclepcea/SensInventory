package configprovider

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/adiclepcea/SensInventory/server/common"
)

//MockConfigProvider contains the configuration for the server
type MockConfigProvider struct {
	Sensors    map[string]common.Sensor
	MinAddress uint8
	MaxAddress uint8
	ConfigProvider
}

//NewConfigProvider creates a new ConfigProvider
func (MockConfigProvider) NewConfigProvider(params ...string) (ConfigProvider, error) {
	c := MockConfigProvider{}

	c.Sensors = make(map[string]common.Sensor)
	return &c, nil
}

//SetAddressLimits adds the minimum and maximum limits for the sensor addreses
func (configProvider *MockConfigProvider) SetAddressLimits(minAddress uint8, maxAddress uint8) error {
	configProvider.MinAddress = minAddress
	configProvider.MaxAddress = maxAddress
	return nil
}

//IsSensorAddressTaken checks to see if there is already a slave with
//the passed address defined
func (configProvider *MockConfigProvider) IsSensorAddressTaken(address uint8) (bool, error) {
	if _, ok := configProvider.Sensors[strconv.Itoa(int(address))]; ok {
		return true, nil
	}

	return false, nil
}

//IsSensorValid checks to see if the sensot passed in is valid
func (configProvider *MockConfigProvider) IsSensorValid(sensor common.Sensor) error {
	if sensor.Address < configProvider.MinAddress || sensor.Address > configProvider.MaxAddress {
		err := fmt.Errorf("The sensor adresses must be between %d and %d", configProvider.MinAddress, configProvider.MaxAddress)
		log.Println(err.Error())
		return err
	}

	if len(sensor.Registers) == 0 {
		err := errors.New("The sensor must have at least one configured register")
		log.Println(err.Error())
		return err
	}

	return nil
}

//AddSensor adds a new sensor that the server should interrogate
func (configProvider *MockConfigProvider) AddSensor(sensor common.Sensor) error {
	if err := configProvider.IsSensorValid(sensor); err != nil {
		log.Println((err).Error())
		return err
	}
	taken, err := configProvider.IsSensorAddressTaken(sensor.Address)
	if err != nil {
		return err
	}
	if taken {
		err := fmt.Errorf("AddSensor. A sensor with address %d has already been registered", sensor.Address)
		log.Println(err.Error())
		return err
	}

	configProvider.Sensors[strconv.Itoa(int(sensor.Address))] = sensor

	return nil
}

//RemoveSensorByAddress removes the sensor having the specified address
//from the collection of sensors that the server interrogates
func (configProvider *MockConfigProvider) RemoveSensorByAddress(address uint8) error {
	taken, err := configProvider.IsSensorAddressTaken(address)
	if err != nil {
		return err
	}
	if !taken {
		err := fmt.Errorf("No sensor with address %d is registered", address)
		log.Println(err.Error())
		return err
	}

	delete(configProvider.Sensors, strconv.Itoa(int(address)))

	return nil
}

//RemoveSensor removes the specified sensor from the collection of sensors that
//the server interrogates
func (configProvider *MockConfigProvider) RemoveSensor(sensor common.Sensor) error {

	return configProvider.RemoveSensorByAddress(sensor.Address)

}

//GetSensorByAddress returns the sensor with the given address
func (configProvider *MockConfigProvider) GetSensorByAddress(address uint8) (*common.Sensor, error) {
	var sensor common.Sensor
	var ok bool

	if sensor, ok = configProvider.Sensors[strconv.Itoa(int(address))]; !ok {
		err := fmt.Errorf("Getting sensor. No sensor with address %d is registered", address)
		log.Println(err.Error())
		return nil, err
	}

	return &sensor, nil
}

//ChangeSensorAddress changes the address of the sensor that currently has
//address "addressBefore" with the "addressAfter"
func (configProvider *MockConfigProvider) ChangeSensorAddress(addressBefore uint8, addressAfter uint8) error {
	sensorBefore, err := configProvider.GetSensorByAddress(addressBefore)
	if err != nil {
		return err
	}
	taken, err := configProvider.IsSensorAddressTaken(addressAfter)
	if err != nil {
		return err
	}
	if taken {
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
func (configProvider *MockConfigProvider) ChangeSensor(address uint8, after common.Sensor) error {
	var sensorBefore *common.Sensor
	var err error
	if sensorBefore, err = configProvider.GetSensorByAddress(address); err != nil {
		return err
	}

	if err = configProvider.IsSensorValid(after); err != nil {
		return err
	}
	//TODO - check for validity of ReadGroups
	sensorBefore.Description = after.Description
	sensorBefore.Registers = after.Registers
	sensorBefore.ReadGroups = after.ReadGroups
	configProvider.Sensors[strconv.Itoa(int(sensorBefore.Address))] = *sensorBefore

	return nil

}

//GetSensors returns a map of the sensor addresses mapped to the sensors themselves
func (configProvider *MockConfigProvider) GetSensors() map[string]common.Sensor {
	return configProvider.Sensors
}
