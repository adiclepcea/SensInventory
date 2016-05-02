package server

import (
	"errors"
	"fmt"
	"log"

	"github.com/adiclepcea/SensInventory/server/common"
)

const MIN_ADDRESS int = 1
const MAX_ADDRESS int = 32

type ConfigProvider struct {
	Sensors map[int]common.Sensor
}

func (ConfigProvider) NewConfigProvider() *ConfigProvider {
	c := ConfigProvider{}
	c.Sensors = make(map[int]common.Sensor)
	return &c
}

func (this *ConfigProvider) AddSensor(sensor common.Sensor) *error {
	if sensor.Address < MIN_ADDRESS || sensor.Address > MAX_ADDRESS {
		err := errors.New(fmt.Sprintf("The sensor adresses must be between %d and %d ", MIN_ADDRESS, MAX_ADDRESS))
		log.Println(err.Error())
		return &err
	}
	if _, ok := this.Sensors[sensor.Address]; ok {
		err := errors.New(fmt.Sprintf("AddSensor. A sensor with address %d has already been registered", sensor.Address))
		log.Println(err.Error())
		return &err
	}

	this.Sensors[sensor.Address] = sensor

	return nil
}

func (this *ConfigProvider) RemoveSensorByAddress(address int) *error {
	if _, ok := this.Sensors[address]; !ok {
		err := errors.New(fmt.Sprintf("No sensor with %d address is registered", address))
		log.Println(err.Error())
		return &err
	}

	delete(this.Sensors, address)

	return nil
}

func (this *ConfigProvider) RemoveSensor(sensor common.Sensor) *error {

	return this.RemoveSensorByAddress(sensor.Address)

}

func (this *ConfigProvider) GetSensorByAddress(address int) (*common.Sensor, *error) {
	var sensor common.Sensor
	var ok bool

	if sensor, ok = this.Sensors[address]; !ok {
		err := errors.New(fmt.Sprintf("No sensor with %d address is registered", address))
		log.Println(err.Error())
		return nil, &err
	}

	return &sensor, nil
}

func (this *ConfigProvider) ChangeSensorAddress(addressBefore int, addressAfter int) *error {
	sensorBefore, err := this.GetSensorByAddress(addressBefore)
	if err != nil {
		return err
	}

	if _, err = this.GetSensorByAddress(addressAfter); err == nil {
		e := errors.New(fmt.Sprintf("There is allready a sensor registered with address %d", addressAfter))
		log.Println(e.Error())
		return &e
	}

	sensorBefore.Address = addressAfter

	return nil

}

func (this *ConfigProvider) ChangeSensor(address int, after common.Sensor) *error {
	var sensorBefore *common.Sensor
	var err *error
	if sensorBefore, err = this.GetSensorByAddress(address); err != nil {
		return err
	}

	sensorBefore.Description = after.Description
	sensorBefore.ConfiguredValues = after.ConfiguredValues
	this.Sensors[sensorBefore.Address] = *sensorBefore

	return nil

}

func (this *ConfigProvider) GetSensors() map[int]common.Sensor {
	return this.Sensors
}
