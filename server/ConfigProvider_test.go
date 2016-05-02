package sensinventory

import (
	"testing"

	"reflect"

	"github.com/adiclepcea/SensInventory/Server/common"
)

func TestAddSensorShouldFail(t *testing.T) {
	conf := ConfigProvider{}.NewConfigProvider()

	sensor1 := common.Sensor{}
	sensor1.Address = 1
	sensor1.Description = "Mock"

	sensor2 := common.Sensor{}
	sensor2.Address = 1
	sensor2.Description = "Should fail"

	err := conf.AddSensor(sensor1)
	if err != nil {
		t.Log("No fail should happen here")
		t.Error("Expected", nil, "got", *err)
	}

	err = conf.AddSensor(sensor2)

	if err == nil {
		t.Error("Expected", "not nil", "got", err)
	}

}

func TestAddSensorShouldOk(t *testing.T) {
	conf := ConfigProvider{}.NewConfigProvider()

	sensor1 := common.Sensor{Address: 2, Description: "sensor 1",
		ReadValues: []common.ReadValue{common.ReadValue{Name: "test ReadValue",
			RegisterAddress: 100, RegisterLength: 2}}}

	err := conf.AddSensor(sensor1)

	if err != nil {
		t.Error("Expected", nil, "got", err)
	}

	sensorBack, err := conf.GetSensorByAddress(sensor1.Address)
	if err != nil {
		t.Error("Expected", nil, "from GetSensorByAddress with param", sensor1.Address, "got", *err)
	}

	if !reflect.DeepEqual(sensor1, *sensorBack) {
		t.Error("Expected", sensor1.Address, ",", sensor1.Description, ",", sensor1.ReadValues,
			"got", sensorBack.Address, ",", sensorBack.Description, ",", sensorBack.ReadValues)
	}
}

func TestGetSensorByAddressShouldFail(t *testing.T) {
	conf := ConfigProvider{}.NewConfigProvider()

	_, err := conf.GetSensorByAddress(1)

	if err == nil {
		t.Error("Expected", "not nil", "got", err)
	}

}

func TestGetSensorByAddressShouldOk(t *testing.T) {
	conf := ConfigProvider{}.NewConfigProvider()
	sensor := common.Sensor{Address: 1, Description: "test"}

	conf.AddSensor(sensor)

	sensorBack, err := conf.GetSensorByAddress(sensor.Address)
	if err != nil {
		t.Error("Expected ", nil, "got", *err)
	}
	if !reflect.DeepEqual(sensor, *sensorBack) {
		t.Error("Expected", sensor, "got", sensorBack)
	}
}

func TestRemoveSensorByAddressShouldFail(t *testing.T) {
	conf := ConfigProvider{}.NewConfigProvider()

	err := conf.RemoveSensorByAddress(1)
	if err == nil {
		t.Error("Expected", "not nil", "got", err)
	}
}

func TestRemoveSensorByAddressShouldOk(t *testing.T) {
	conf := ConfigProvider{}.NewConfigProvider()
	sensor := common.Sensor{Address: 1, Description: "test"}
	conf.AddSensor(sensor)
	err := conf.RemoveSensorByAddress(sensor.Address)
	if err != nil {
		t.Error("Expected", "nil", "got", *err)
	}
}

func TestRemoveSensorShouldFail(t *testing.T) {
	conf := ConfigProvider{}.NewConfigProvider()
	sensor := common.Sensor{Address: 1, Description: "test"}
	err := conf.RemoveSensor(sensor)
	if err == nil {
		t.Error("Expected", "not nil", "got", err)
	}
}

func TestRemoveSensorShouldOk(t *testing.T) {
	conf := ConfigProvider{}.NewConfigProvider()
	sensor := common.Sensor{Address: 1, Description: "test"}
	conf.AddSensor(sensor)
	err := conf.RemoveSensor(sensor)
	if err != nil {
		t.Error("Expected", "nil", "got", *err)
	}
}

func TestChangeSensorAddressShouldFail(t *testing.T) {
	conf := ConfigProvider{}.NewConfigProvider()
	sensor1 := common.Sensor{Address: 1, Description: "test"}
	sensor2 := common.Sensor{Address: 2, Description: "test"}

	err := conf.ChangeSensorAddress(1, 2)

	if err == nil {
		t.Error("When no sensor added expected", "nil", "got", *err)
	}

	conf.AddSensor(sensor1)
	conf.AddSensor(sensor2)
	err = conf.ChangeSensorAddress(1, 2)

	if err == nil {
		t.Error("When address already exists expected", "not nil", "got", err)
	}

}

func TestChangeSensorShouldFail(t *testing.T) {
	conf := ConfigProvider{}.NewConfigProvider()

	sensor := common.Sensor{Address: 1, Description: "Test"}

	err := conf.ChangeSensor(1, sensor)

	if err == nil {
		t.Error("Expected", "not nil", "got", err)
	}
}

func TestChangeSensorShouldOk(t *testing.T) {
	conf := ConfigProvider{}.NewConfigProvider()

	sensor := common.Sensor{Address: 1, Description: "Test"}
	conf.AddSensor(sensor)

	sensor2 := common.Sensor{Address: 2, Description: "Test 2",
		ReadValues: []common.ReadValue{common.ReadValue{Name: "test ReadValue",
			RegisterAddress: 100, RegisterLength: 2}}}

	err := conf.ChangeSensor(sensor.Address, sensor2)

	if err != nil {
		t.Error("Expected", "nil", "got", *err)
	}

	sensor1, _ := conf.GetSensorByAddress(sensor.Address)

	if sensor1.Description != sensor2.Description || !reflect.DeepEqual(sensor1.ReadValues, sensor2.ReadValues) {
		t.Error("Expected", sensor2, "got", *sensor1)
	}
}