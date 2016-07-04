package persistenceprovider

import (
	"fmt"
	"log"
	"time"

	"github.com/adiclepcea/SensInventory/server/common"
)

type timedReading struct {
	Reading  common.Reading
	ReadTime time.Time
}

//MockPersistenceProvider is a fake persistence provider
type MockPersistenceProvider struct {
	timedReadings []timedReading
	PersistenceProvider
}

//NewPersistenceProvider - initiate a new MockReadingProvider
func (MockPersistenceProvider) NewPersistenceProvider(params ...string) (PersistenceProvider, error) {
	return &MockPersistenceProvider{}, nil
}

//SaveSensorReading - mocks saving a sensor
func (mpp *MockPersistenceProvider) SaveSensorReading(reading common.Reading) error {
	tr := timedReading{Reading: reading, ReadTime: time.Now()}
	log.Printf("Saving %d\n", reading.Sensor)
	mpp.timedReadings = append(mpp.timedReadings, tr)
	return nil
}

//GetSensorReading returns the reading for sensor with sensorAddress at the exact time t
func (mpp *MockPersistenceProvider) GetSensorReading(sensorAddress uint8, t time.Time) (*common.Reading, error) {
	for _, tr := range mpp.timedReadings {
		if tr.ReadTime.Equal(t) && tr.Reading.Sensor == sensorAddress {
			return &tr.Reading, nil
		}
	}
	return nil, nil
}

//GetSensorReadingsInPeriod returns all the readings for the sensor with address
//sensorAddress in the period between start and end
func (mpp *MockPersistenceProvider) GetSensorReadingsInPeriod(sensorAddress uint8, start time.Time, end time.Time) ([]common.Reading, error) {
	readings := []common.Reading{}
	for _, tr := range mpp.timedReadings {
		if tr.ReadTime.Before(end) && tr.ReadTime.After(start) && tr.Reading.Sensor == sensorAddress {
			readings = append(readings, tr.Reading)
		}
	}
	return readings, nil
}

//GetSensorReadingCountInPeriod returnns the number of readings for the sensor
//with address sensorAddress between start and end time
func (mpp *MockPersistenceProvider) GetSensorReadingCountInPeriod(sensorAddress uint8, start time.Time, end time.Time) (uint, error) {
	count := 0
	for _, tr := range mpp.timedReadings {
		if tr.ReadTime.Before(end) && tr.ReadTime.After(start) && tr.Reading.Sensor == sensorAddress {
			count++
		}
	}
	return uint(count), nil
}

//GetAllReadingsInPeriod returns all the readings in the period between start and end
func (mpp *MockPersistenceProvider) GetAllReadingsInPeriod(start time.Time, end time.Time) (*[]common.Reading, error) {
	readings := []common.Reading{}
	for _, tr := range mpp.timedReadings {
		if tr.ReadTime.Before(end) && tr.ReadTime.After(start) {
			readings = append(readings, tr.Reading)
		}
	}
	return &readings, nil
}

//GetAllReadingsCountInPeriod returns the count all the readings in the period between start and end
func (mpp *MockPersistenceProvider) GetAllReadingsCountInPeriod(start time.Time, end time.Time) (uint, error) {
	count := 0
	for _, tr := range mpp.timedReadings {
		if tr.ReadTime.Before(end) && tr.ReadTime.After(start) {
			count++
		}
	}
	return uint(count), nil
}

//DeleteSensorReading deletes the reading from the sensowith address sensorAddress
//made at the t time
func (mpp *MockPersistenceProvider) DeleteSensorReading(sensorAddress uint8, t time.Time) error {
	for i, tr := range mpp.timedReadings {
		if tr.ReadTime.Equal(t) && tr.Reading.Sensor == sensorAddress {
			mpp.timedReadings = append(mpp.timedReadings[:i], mpp.timedReadings[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Not found")
}

//DeleteSensorReadingsInPeriod deletes all the readings for the sensor with address
//sensorAddress between start and end times
func (mpp *MockPersistenceProvider) DeleteSensorReadingsInPeriod(sensorAddress uint8, start time.Time, end time.Time) error {
	rez := mpp.timedReadings[:0]
	for _, tr := range mpp.timedReadings {
		if tr.ReadTime.Before(end) && tr.ReadTime.After(start) && tr.Reading.Sensor == sensorAddress {
			rez = append(rez, tr)
		}
	}
	mpp.timedReadings = rez
	return nil
}

//DeleteAllReadingsInPeriod deletes all reading between start and end times
func (mpp *MockPersistenceProvider) DeleteAllReadingsInPeriod(start time.Time, end time.Time) error {
	rez := mpp.timedReadings[:0]
	for _, tr := range mpp.timedReadings {
		if tr.ReadTime.Before(end) && tr.ReadTime.After(start) {
			rez = append(rez, tr)
		}
	}
	mpp.timedReadings = rez
	return nil
}
