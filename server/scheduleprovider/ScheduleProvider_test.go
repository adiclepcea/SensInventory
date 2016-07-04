package scheduleprovider_test

import (
	"testing"
	"time"

	"github.com/adiclepcea/SensInventory/server/common"
	"github.com/adiclepcea/SensInventory/server/configprovider"
	"github.com/adiclepcea/SensInventory/server/persistenceprovider"
	"github.com/adiclepcea/SensInventory/server/readingprovider"
	"github.com/adiclepcea/SensInventory/server/scheduleprovider"
)

func TestScheduleProviderNotNewShouldFail(t *testing.T) {
	schprovider := scheduleprovider.ScheduleProvider{}
	it := scheduleprovider.IntervalTimer{}
	err := schprovider.AddTimer(it)
	if err == nil {
		t.Fatal("Expected error while adding an intervalTimer on an not initialized scheduleprovider got nil")
	}
}

func TestScheduleProviderShouldFail(t *testing.T) {
	cp, _ := configprovider.MockConfigProvider{}.NewConfigProvider()
	rp := readingprovider.MockReadingProvider{}.NewReadingProvider(&cp)
	schprovider := scheduleprovider.ScheduleProvider{}.NewScheduleProvider(rp, nil)
	it := scheduleprovider.IntervalTimer{}
	it.Persist = true
	err := schprovider.AddTimer(it)
	if err == nil {
		t.Fatal("Expected error while adding an intervalTimer", "with persist",
			"while shecdule provider has no persist")
	}
}

func TestScheduleProviderShouldOk(t *testing.T) {
	cp, _ := configprovider.MockConfigProvider{}.NewConfigProvider()
	cp.SetAddressLimits(0, 50)
	sensor1 := common.Sensor{}
	sensor1.Address = 33
	sensor1.Description = "Mock"
	sensor1.Registers = []common.Register{common.Register{
		Name: "test ReadValue", Location: 100, Type: common.Holding}}

	sensor2 := common.Sensor{}
	sensor2.Address = 22
	sensor2.Description = "Mock"
	sensor2.Registers = []common.Register{common.Register{
		Name: "test ReadValue", Location: 10, Type: common.Holding}}
	cp.AddSensor(sensor1)
	cp.AddSensor(sensor2)
	rp := readingprovider.MockReadingProvider{}.NewReadingProvider(&cp)
	pp, _ := persistenceprovider.MockPersistenceProvider{}.NewPersistenceProvider()
	schprovider := scheduleprovider.ScheduleProvider{}.NewScheduleProvider(rp, &pp)
	it := scheduleprovider.IntervalTimer{}
	it.Persist = true
	it.ReadType = common.Holding
	it.SensorAddress = 33
	it.StartLocation = 100
	it.ReadLength = 1
	interv := (time.Second * 6)
	it.Interval = &interv

	it2 := scheduleprovider.IntervalTimer{}
	it2.Persist = true
	it2.ReadType = common.Holding
	it2.SensorAddress = 22
	it2.StartLocation = 10
	it2.ReadLength = 1
	interv2 := (time.Second * 2)
	it2.Interval = &interv2
	err := schprovider.AddTimer(it)
	err = schprovider.AddTimer(it2)
	if err != nil {
		t.Fatal("No error expected when adding an interval timer with",
			"persistence", "to a schedule provider with PersistenceProvider not nil")
	}

	schprovider.Start()
	time.Sleep(20 * time.Second)
	schprovider.Stop()
}
