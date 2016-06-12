package configprovider_test

import (
	"testing"

	"reflect"

	"github.com/adiclepcea/SensInventory/server/common"
	"github.com/adiclepcea/SensInventory/server/configprovider"
)

func TestMockAddSensorWithInvalidAddressShoudlFail(t *testing.T) {
	conf, _ := configprovider.MockConfigProvider{}.NewConfigProvider()
	conf.SetAddressLimits(1, 32)

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

func TestMockAddSensorShouldFail(t *testing.T) {
	conf, _ := configprovider.MockConfigProvider{}.NewConfigProvider()
	conf.SetAddressLimits(1, 32)

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

func TestMockAddSensorShouldOk(t *testing.T) {
	conf, _ := configprovider.MockConfigProvider{}.NewConfigProvider()
	conf.SetAddressLimits(1, 32)

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

func TestMockGetSensorByAddressShouldFail(t *testing.T) {
	conf, _ := configprovider.MockConfigProvider{}.NewConfigProvider()
	conf.SetAddressLimits(1, 32)

	_, err := conf.GetSensorByAddress(1)

	if err == nil {
		t.Error("Expected", "not nil", "got", err)
		t.Fail()
	}

}

func TestMockGetSensorByAddressShouldOk(t *testing.T) {
	conf, _ := configprovider.MockConfigProvider{}.NewConfigProvider()
	conf.SetAddressLimits(1, 32)

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

func TestMockRemoveSensorByAddressShouldFail(t *testing.T) {
	conf, _ := configprovider.MockConfigProvider{}.NewConfigProvider()
	conf.SetAddressLimits(1, 32)

	err := conf.RemoveSensorByAddress(1)
	if err == nil {
		t.Error("Expected", "not nil", "got", err)
		t.Fail()
	}
}

func TestMockRemoveSensorByAddressShouldOk(t *testing.T) {
	conf, _ := configprovider.MockConfigProvider{}.NewConfigProvider()
	conf.SetAddressLimits(1, 32)

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

func TestMockRemoveSensorShouldFail(t *testing.T) {
	conf, _ := configprovider.MockConfigProvider{}.NewConfigProvider()
	conf.SetAddressLimits(1, 32)

	sensor := common.Sensor{Address: 1, Description: "test"}
	sensor.Registers = []common.Register{common.Register{
		Name: "test ReadValue", Location: 100, Type: common.Holding}}

	err := conf.RemoveSensor(sensor)
	if err == nil {
		t.Error("Expected", "not nil", "got", err)
		t.Fail()
	}
}

func TestMockRemoveSensorShouldOk(t *testing.T) {
	conf, _ := configprovider.MockConfigProvider{}.NewConfigProvider()
	conf.SetAddressLimits(1, 32)

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

func TestMockChangeSensorAddressShouldOk(t *testing.T) {
	conf, _ := configprovider.MockConfigProvider{}.NewConfigProvider()
	conf.SetAddressLimits(1, 32)
	sensor1 := common.Sensor{Address: 1, Description: "test"}
	sensor1.Registers = []common.Register{common.Register{
		Name: "test ReadValue", Location: 100, Type: common.Input}}
	err := conf.AddSensor(sensor1)
	if err != nil {
		t.Fatalf("No error expected when adding a sensor, got %s", err.Error())
	}
	err = conf.ChangeSensorAddress(1, 3)

	if err != nil {
		t.Error("When changing to an address unallocated expected", "no error", "got", err)
		t.Fail()
	}
}

func TestMockChangeSensorAddressShouldFail(t *testing.T) {
	conf, _ := configprovider.MockConfigProvider{}.NewConfigProvider()
	conf.SetAddressLimits(1, 32)
	sensor1 := common.Sensor{Address: 1, Description: "test"}
	sensor2 := common.Sensor{Address: 2, Description: "test"}
	sensor1.Registers = []common.Register{common.Register{
		Name: "test1 ReadValue", Location: 101, Type: common.Input}}
	sensor2.Registers = []common.Register{common.Register{
		Name: "test2 ReadValue", Location: 100, Type: common.Input}}

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
	if err.Error() != "There is allready a sensor registered with address 2" {
		t.Errorf("There is an unexpected error message: %s", err.Error())
		t.Fail()
	}

}

func TestMockChangeSensorShouldFail(t *testing.T) {
	conf, _ := configprovider.MockConfigProvider{}.NewConfigProvider()
	conf.SetAddressLimits(1, 32)

	sensor := common.Sensor{Address: 1, Description: "Test"}

	err := conf.ChangeSensor(1, sensor)

	if err == nil {
		t.Error("Expected", "not nil", "got", err)
		t.Fail()
	}
}

func TestMockChangeSensorShouldOk(t *testing.T) {
	conf, _ := configprovider.MockConfigProvider{}.NewConfigProvider()
	conf.SetAddressLimits(1, 32)

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

func TestMockGetSensors(t *testing.T) {
	conf, _ := configprovider.MockConfigProvider{}.NewConfigProvider()
	conf.SetAddressLimits(1, 32)

	sensor1 := common.Sensor{Address: 1, Description: "Test"}
	sensor1.Registers = []common.Register{common.Register{
		Name: "test ReadValue", Location: 100, Type: common.Input}}
	sensor2 := common.Sensor{Address: 2, Description: "Test"}
	sensor2.Registers = []common.Register{common.Register{
		Name: "test ReadValue", Location: 100, Type: common.Coil}}

	err := conf.AddSensor(sensor1)
	if err != nil {
		t.Fatalf("No error extepected while adding sensor, got %s", err.Error())
	}
	err = conf.AddSensor(sensor2)
	if err != nil {
		t.Fatalf("No error extepected while adding sensor, got %s", err.Error())
	}

	sensors := conf.GetSensors()
	if len(sensors) != 2 {
		t.Fatalf("Should have 2 sensors, got: %d", len(sensors))
	}

	val, ok := sensors["1"]
	if !ok {
		t.Fatal("Should have sensor with address 1, got:nil")
	}
	if val.Address != 1 {
		t.Fatalf("Sensor shoudl have address 1, got:%d", val.Address)
	}
	val, ok = sensors["2"]
	if !ok {
		t.Fatal("Should have sensor with address 2, got:nil")
	}
	if val.Address != 2 {
		t.Fatalf("Sensor shoudl have address 2, got:%d", val.Address)
	}
}
