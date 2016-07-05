package persistenceprovider

import (
	"testing"
	"time"

	"github.com/adiclepcea/SensInventory/server/common"
	"github.com/adiclepcea/SensInventory/server/readgroups"
)

func TestMockPersistenceProvider(t *testing.T) {
	mp, _ := MockPersistenceProvider{}.NewPersistenceProvider()

	reading1 := common.Reading{}
	now := time.Now()
	reading1.Type = common.Holding
	reading1.Count = 3
	reading1.Sensor = 10
	reading1.StartLocation = 0
	reading1.ReadValues = []uint16{100, 200, 300}
	reading1.Time = now.Format(common.TimeFormat)
	reading1.InitCalculatedValues()
	tNow, _ := time.Parse(common.TimeFormat, reading1.Time)
	rg, _ := readgroups.ReadGroupFloat32{}.NewReadGroup(10, 0)
	rg.ResultType = common.Float32
	rg.Calculate(&reading1)

	mp.SaveSensorReading(reading1)

	reading, _ := mp.GetSensorReading(10, tNow)

	if reading == nil {
		t.Fatal("Expected reading, got", reading)
	}

	if reading.Sensor != reading1.Sensor {
		t.Fatal("Expected sensor address", reading1.Sensor, "got", reading.Sensor)
	}

	countReadings, _ := mp.GetSensorReadingCountInPeriod(reading.Sensor, tNow.Add(time.Second*1), tNow.Add(time.Second*2))
	if countReadings != 0 {
		t.Fatal("Expected no reading in period, got ", countReadings, "readings")
	}

	readings, _ := mp.GetSensorReadingsInPeriod(reading.Sensor, tNow.Add(time.Second*-1), tNow.Add(time.Second*2))
	if len(readings) != 1 {
		t.Fatal("Expected one reading in period, got ", len(readings), "readings")
	}

	mp.DeleteSensorReading(reading1.Sensor, tNow)
	countReadings, _ = mp.GetSensorReadingCountInPeriod(reading.Sensor, tNow.Add(time.Second*-1), tNow.Add(time.Second*2))

	if countReadings != 0 {
		t.Fatal("No reading expected after deletion, got", countReadings)
	}

	mp.SaveSensorReading(reading1)

	mp.DeleteSensorReadingsInPeriod(reading.Sensor, tNow.Add(time.Second*1), tNow.Add(time.Second*2))

	countReadings, _ = mp.GetSensorReadingCountInPeriod(reading.Sensor, tNow.Add(time.Second*-1), tNow.Add(time.Second*2))

	if countReadings != 1 {
		t.Fatal("Expected 1 reading, got ", countReadings)
	}

	mp.DeleteAllReadingsInPeriod(tNow.Add(time.Second*-1), tNow.Add(time.Second*2))

	countReadings, _ = mp.GetSensorReadingCountInPeriod(reading.Sensor, tNow.Add(time.Second*-1), tNow.Add(time.Second*2))

	if countReadings != 0 {
		t.Fatal("No reading expected after deleting all, got", countReadings)
	}

}
