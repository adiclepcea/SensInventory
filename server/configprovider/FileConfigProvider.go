package configprovider

import (
	"errors"
	"fmt"
	"log"

	"github.com/adiclepcea/SensInventory/server/common"
)

const defaultFileName = "config.json"

//FileConfigProvider contains the configuration for the server
type FileConfigProvider struct {
	Sensors        map[uint8]common.Sensor `json:"Sensors"`
	ReadGroups     []common.ReadGroup      `json:"ReadGroups"`
	MinAddress     uint8
	MaxAddress     uint8
	FileConfigName string
	ConfigProvider
}

//NewConfigProvider creates a new ConfigProvider
func (FileConfigProvider) NewConfigProvider(minAddress uint8, maxAddress uint8, params ...string) *FileConfigProvider {
	c := FileConfigProvider{FileConfigName: defaultFileName}
	if len(params) == 1 {
		c.FileConfigName = params[0]
	}
	c.MinAddress = minAddress
	c.MaxAddress = maxAddress
	c.Sensors = make(map[uint8]common.Sensor)
	return &c
}

//IsSensorAddressTaken checks to see if there is already a slave with
//the passed address defined
func (configProvider *FileConfigProvider) IsSensorAddressTaken(address uint8) bool {
	if _, ok := configProvider.Sensors[address]; ok {
		return true
	}

	return false
}

//IsSensorValid checks to see if the sensot passed in is valid
func (configProvider *FileConfigProvider) IsSensorValid(sensor common.Sensor) error {
	if sensor.Address < configProvider.MinAddress || sensor.Address > configProvider.MaxAddress {
		err := fmt.Errorf("The sensor adresses must be between %d and %d", configProvider.MinAddress, configProvider.MaxAddress)
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
func (configProvider *FileConfigProvider) AddSensor(sensor common.Sensor) error {
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
func (configProvider *FileConfigProvider) RemoveSensorByAddress(address uint8) error {
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
func (configProvider *FileConfigProvider) RemoveSensor(sensor common.Sensor) error {

	return configProvider.RemoveSensorByAddress(sensor.Address)

}

//GetSensorByAddress returns the sensor with the given address
func (configProvider *FileConfigProvider) GetSensorByAddress(address uint8) (*common.Sensor, error) {
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
func (configProvider *FileConfigProvider) ChangeSensorAddress(addressBefore uint8, addressAfter uint8) error {
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
func (configProvider *FileConfigProvider) ChangeSensor(address uint8, after common.Sensor) error {
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
func (configProvider *FileConfigProvider) GetSensors() map[uint8]common.Sensor {
	return configProvider.Sensors
}
