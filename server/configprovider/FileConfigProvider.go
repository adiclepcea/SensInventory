package configprovider

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/adiclepcea/SensInventory/server/common"
)

const defaultFileName = "./config.json"

//FileConfigProvider contains the configuration for the server
type FileConfigProvider struct {
	Sensors        map[string]common.Sensor `json:"Sensors"`
	MinAddress     uint8                    `json:"minAddress"`
	MaxAddress     uint8                    `json:"maxAddress"`
	FileConfigName string                   `json:"-"`
	ConfigProvider `json:"-"`
}

//NewConfigProvider creates a new ConfigProvider
func (FileConfigProvider) NewConfigProvider(params ...string) (*FileConfigProvider, error) {
	c := FileConfigProvider{FileConfigName: defaultFileName}
	if len(params) == 1 {
		c.FileConfigName = params[0]
	}

	c.Sensors = make(map[string]common.Sensor)

	_, err := c.LoadConfig()

	if err != nil {
		log.Printf(err.Error())
		return nil, err
	}

	return &c, nil
}

//LoadConfig loads the configuration from the file
func (configProvider *FileConfigProvider) LoadConfig() (bool, error) {
	if _, err := os.Stat(configProvider.FileConfigName); err != nil {
		log.Println("Config file not found. Creating a new one")
		return false, nil
	}

	configFile, err := os.Open(configProvider.FileConfigName)
	if err != nil {
		log.Printf("Error while opening the config file %s: %s",
			configProvider.FileConfigName, err.Error())
		return false, err
	}
	defer configFile.Close()
	fi, err := configFile.Stat()
	if err != nil {
		return false, err
	}
	if fi.Size() == 0 {
		return false, nil
	}
	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(&configProvider); err != nil {
		return false, err
	}

	return true, nil
}

//Save saves the configuration into the config file
func (configProvider *FileConfigProvider) Save() error {
	configFile, err := os.Create(configProvider.FileConfigName)
	if err != nil {
		return err
	}
	defer configFile.Close()
	jsonEncoder := json.NewEncoder(configFile)

	if err = jsonEncoder.Encode(configProvider); err != nil {
		return err
	}

	return nil
}

//SetAddressLimits adds the minimum and maximum limits for the sensor addreses
func (configProvider *FileConfigProvider) SetAddressLimits(minAddress uint8, maxAddress uint8) error {
	configProvider.MinAddress = minAddress
	configProvider.MaxAddress = maxAddress
	return configProvider.Save()
}

//IsSensorAddressTaken checks to see if there is already a slave with
//the passed address defined
func (configProvider *FileConfigProvider) IsSensorAddressTaken(address uint8) (bool, error) {
	if _, ok := configProvider.Sensors[strconv.Itoa(int(address))]; ok {
		return true, nil
	}

	return false, nil
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
	return configProvider.Save()
}

//RemoveSensorByAddress removes the sensor having the specified address
//from the collection of sensors that the server interrogates
func (configProvider *FileConfigProvider) RemoveSensorByAddress(address uint8) error {
	taken, err := configProvider.IsSensorAddressTaken(address)
	if err != nil {
		return err
	}
	if !taken {
		err := fmt.Errorf("No sensor with %d address is registered", address)
		log.Println(err.Error())
		return err
	}

	delete(configProvider.Sensors, strconv.Itoa(int(address)))

	return configProvider.Save()

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

	if sensor, ok = configProvider.Sensors[strconv.Itoa(int(address))]; !ok {
		err := fmt.Errorf("No sensor with address %d is registered", address)
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
	configProvider.Sensors[strconv.Itoa(int(sensorBefore.Address))] = *sensorBefore

	return configProvider.Save()

}

//GetSensors returns a map of the sensor addresses mapped to the sensors themselves
func (configProvider *FileConfigProvider) GetSensors() map[string]common.Sensor {
	return configProvider.Sensors
}
