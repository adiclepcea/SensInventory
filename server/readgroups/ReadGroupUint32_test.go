package readgroups_test

import (
	"testing"

	"github.com/adiclepcea/SensInventory/server/common"
	"github.com/adiclepcea/SensInventory/server/readgroups"
)

func TestCalculateUint32ShouldOk(t *testing.T) {
	rgf32, err := readgroups.ReadGroupUint32{}.NewReadGroup(1, 10)
	if err != nil {
		t.Error("No error expeccted when creating a ReadGroupFloat32")
		t.FailNow()
	}

	sensor := common.Sensor{Address: 1}
	sensor.Registers = []common.Register{
		common.Register{Location: 8, Type: common.Holding},
		common.Register{Location: 9, Type: common.Holding},
		common.Register{Location: 10, Type: common.Holding},
		common.Register{Location: 11, Type: common.Holding},
		common.Register{Location: 12, Type: common.Holding},
	}

	reading := common.Reading{Sensor: 1, Type: common.Holding,
		StartLocation: 8, Count: 5, ReadValues: []uint16{0, 0, 0xEF11, 0xFFDD, 0}}

	rez, err := rgf32.Calculate(reading)
	if rez != uint32(0xEF11FFDD) {
		t.Fatalf("The expected result should have been 0xEF11FFDD got %x", rez)
	}
}

func TestCalculateUint32ShouldLocationToHigh(t *testing.T) {
	rgf32, err := readgroups.ReadGroupUint32{}.NewReadGroup(1, 10)
	if err != nil {
		t.Error("No error expeccted when creating a ReadGroupFloat32")
		t.FailNow()
	}

	sensor := common.Sensor{Address: 1}
	sensor.Registers = []common.Register{
		common.Register{Location: 18, Type: common.Holding},
		common.Register{Location: 19, Type: common.Holding},
		common.Register{Location: 20, Type: common.Holding},
		common.Register{Location: 21, Type: common.Holding},
		common.Register{Location: 22, Type: common.Holding},
	}

	reading := common.Reading{Sensor: 1, Type: common.Holding,
		StartLocation: 18, Count: 5, ReadValues: []uint16{0, 0, 0xEF11, 0xFFDD, 0}}

	rez, err := rgf32.Calculate(reading)
	if err == nil {
		t.Fatal("Expected error because location to high but got nil,", rez)
	}
	t.Log(err.Error())

}

func TestCalculateUint32ShouldReadingToShort(t *testing.T) {
	rgf32, err := readgroups.ReadGroupUint32{}.NewReadGroup(1, 10)
	if err != nil {
		t.Error("No error expeccted when creating a ReadGroupFloat32")
		t.FailNow()
	}

	sensor := common.Sensor{Address: 1}
	sensor.Registers = []common.Register{
		common.Register{Location: 6, Type: common.Holding},
		common.Register{Location: 7, Type: common.Holding},
		common.Register{Location: 8, Type: common.Holding},
		common.Register{Location: 9, Type: common.Holding},
		common.Register{Location: 10, Type: common.Holding},
	}

	reading := common.Reading{Sensor: 1, Type: common.Holding,
		StartLocation: 6, Count: 5, ReadValues: []uint16{0, 0, 0xEF11, 0xFFDD, 0}}

	rez, err := rgf32.Calculate(reading)
	if err == nil {
		t.Fatal("Expected error because reading too short,", rez)
	}
	t.Log(err.Error())

}
