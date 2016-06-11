package configprovider_test

import (
	"testing"

	"reflect"

	"github.com/adiclepcea/SensInventory/server/common"
	"github.com/adiclepcea/SensInventory/server/configprovider"
)

func TestAddSensorWithInvalidAddressShoudlFail(t *testing.T) {
	conf := configprovider.MockConfigProvider{}.NewConfigProvider(1, 32)

	sensor1 := common.Sensor{}
	sensor1.Address = 33
	sensor1.Description = "Mock"
	sensor1.Registers = []common.Register{common.Register{
		Name: "test ReadValue", Location: 100, Type: common.Holding}}

	err := conf.AddSensor(sensor1)

	if err == nil {
		t.Error("There should be an error when adding a sensor with an invalid address")
		t.Fail()
	}
}

func TestAddSensorShouldFail(t *testing.T) {
	conf := configprovider.MockConfigProvider{}.NewConfigProvider(1, 32)

	sensor1 := common.Sensor{}
	sensor1.Address = 1
	sensor1.Description = "Mock"
	sensor1.Registers = []common.Register{common.Register{
		Name: "test ReadValue", Location: 100, Type: common.Holding}}

	sensor2 := common.Sensor{}
	sensor2.Address = 1
	sensor2.Description = "Should fail"
	sensor2.Registers = []common.Register{common.Register{
		Name: "test ReadValue", Location: 100, Type: common.Input}}

	err := conf.AddSensor(sensor1)
	if err != nil {
		t.Log("No fail should happen here")
		t.Error("Expected", nil, "got", err)
		t.FailNow()
	} else {
		t.Log("TestAddSensorShouldFail - OK: No error")
	}

	err = conf.AddSensor(sensor2)

	if err == nil {
		t.Error("Expected", "not nil", "got", err)
		t.FailNow()
	} else {
		t.Log("TestAddSensorShouldFail - OK: ", err.Error())
	}

}

func TestAddSensorShouldOk(t *testing.T) {
	conf := configprovider.MockConfigProvider{}.NewConfigProvider(1, 32)

	sensor1 := common.Sensor{Address: 2, Description: "sensor 1",
		Registers: []common.Register{common.Register{
			Name: "test ReadValue", Location: 100, Type: common.Input}}}

	err := conf.AddSensor(sensor1)

	if err != nil {
		t.Error("Expected", nil, "got", err)
		t.FailNow()
	}

	sensorBack, err := conf.GetSensorByAddress(sensor1.Address)
	if err != nil {
		t.Error("Expected", nil, "from GetSensorByAddress with param",
			sensor1.Address, "got", err)
		t.FailNow()
	}

	if !reflect.DeepEqual(sensor1, *sensorBack) {
		t.Error("Expected", sensor1.Address, ",", sensor1.Description, ",",
			sensor1.Registers, "got", sensorBack.Address, ",",
			sensorBack.Description, ",", sensorBack.Registers)
		t.Fail()
	}
}

func TestGetSensorByAddressShouldFail(t *testing.T) {
	conf := configprovider.MockConfigProvider{}.NewConfigProvider(1, 32)

	_, err := conf.GetSensorByAddress(1)

	if err == nil {
		t.Error("Expected", "not nil", "got", err)
		t.Fail()
	}

}

func TestGetSensorByAddressShouldOk(t *testing.T) {
	conf := configprovider.MockConfigProvider{}.NewConfigProvider(1, 32)
	sensor := common.Sensor{Address: 1, Description: "test"}
	sensor.Registers = []common.Register{common.Register{
		Name: "test ReadValue", Location: 100, Type: common.Holding}}

	err := conf.AddSensor(sensor)
	if err != nil {
		t.Error("Expected", nil, "got", err)
		t.FailNow()
	}

	sensorBack, err := conf.GetSensorByAddress(sensor.Address)
	if err != nil {
		t.Error("Expected ", nil, "got", err)
		t.FailNow()
	}
	if !reflect.DeepEqual(sensor, *sensorBack) {
		t.Error("Expected", sensor, "got", sensorBack)
		t.Fail()
	}
}

func TestRemoveSensorByAddressShouldFail(t *testing.T) {
	conf := configprovider.MockConfigProvider{}.NewConfigProvider(1, 32)

	err := conf.RemoveSensorByAddress(1)
	if err == nil {
		t.Error("Expected", "not nil", "got", err)
		t.Fail()
	}
}

func TestRemoveSensorByAddressShouldOk(t *testing.T) {
	conf := configprovider.MockConfigProvider{}.NewConfigProvider(1, 32)
	sensor := common.Sensor{Address: 1, Description: "test"}
	sensor.Registers = []common.Register{common.Register{
		Name: "test ReadValue", Location: 100, Type: common.Input}}

	conf.AddSensor(sensor)
	err := conf.RemoveSensorByAddress(sensor.Address)
	if err != nil {
		t.Error("Expected", "nil", "got", err)
		t.Fail()
	}
}

func TestRemoveSensorShouldFail(t *testing.T) {
	conf := configprovider.MockConfigProvider{}.NewConfigProvider(1, 32)
	sensor := common.Sensor{Address: 1, Description: "test"}
	sensor.Registers = []common.Register{common.Register{
		Name: "test ReadValue", Location: 100, Type: common.Holding}}

	err := conf.RemoveSensor(sensor)
	if err == nil {
		t.Error("Expected", "not nil", "got", err)
		t.Fail()
	}
}

func TestRemoveSensorShouldOk(t *testing.T) {
	conf := configprovider.MockConfigProvider{}.NewConfigProvider(1, 32)
	sensor := common.Sensor{Address: 1, Description: "test"}
	sensor.Registers = []common.Register{common.Register{
		Name: "test ReadValue", Location: 100, Type: common.Input}}

	err := conf.AddSensor(sensor)
	if err != nil {
		t.Errorf("Expected no error when adding a new sensor, but got %s", err.Error())
		t.FailNow()
	}
	err = conf.RemoveSensor(sensor)
	if err != nil {
		t.Error("Expected", "nil", "got", err)
		t.Fail()
	}
}

func TestChangeSensorAddressShouldFail(t *testing.T) {
	conf := configprovider.MockConfigProvider{}.NewConfigProvider(1, 32)
	sensor1 := common.Sensor{Address: 1, Description: "test"}
	sensor2 := common.Sensor{Address: 2, Description: "test"}

	err := conf.ChangeSensorAddress(1, 2)

	if err == nil {
		t.Error("When no sensor added expected", "nil", "got", err)
		t.FailNow()
	}

	conf.AddSensor(sensor1)
	conf.AddSensor(sensor2)
	err = conf.ChangeSensorAddress(1, 2)

	if err == nil {
		t.Error("When address already exists expected", "not nil", "got", err)
		t.Fail()
	}

}

func TestChangeSensorShouldFail(t *testing.T) {
	conf := configprovider.MockConfigProvider{}.NewConfigProvider(1, 32)

	sensor := common.Sensor{Address: 1, Description: "Test"}

	err := conf.ChangeSensor(1, sensor)

	if err == nil {
		t.Error("Expected", "not nil", "got", err)
		t.Fail()
	}
}

func TestChangeSensorShouldOk(t *testing.T) {
	conf := configprovider.MockConfigProvider{}.NewConfigProvider(1, 32)

	sensor := common.Sensor{Address: 1, Description: "Test"}
	sensor.Registers = []common.Register{common.Register{
		Name: "test ReadValue", Location: 100, Type: common.Input}}

	err := conf.AddSensor(sensor)

	if err != nil {
		t.Error("Expected no error when adding a new sensor", "got", err)
		t.FailNow()
	}

	sensor2 := common.Sensor{Address: 2, Description: "Test 2",
		Registers: []common.Register{common.Register{
			Name: "test ReadValue", Location: 100, Type: common.Holding}}}

	err = conf.ChangeSensor(sensor.Address, sensor2)

	if err != nil {
		t.Error("Expected", "nil", "got", err)
		t.FailNow()
	}

	sensor1, _ := conf.GetSensorByAddress(sensor.Address)

	if sensor1.Description != sensor2.Description || !reflect.DeepEqual(sensor1.Registers, sensor2.Registers) {
		t.Error("Expected", sensor2, "got", *sensor1)
	}
}
