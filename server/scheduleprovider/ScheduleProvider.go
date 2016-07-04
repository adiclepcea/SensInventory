package scheduleprovider

import (
	"fmt"
	"log"
	"time"

	"github.com/adiclepcea/SensInventory/server/persistenceprovider"
	"github.com/adiclepcea/SensInventory/server/readingprovider"
)

//ScheduleProvider is the base structure needed for
//a scheduled read/write of sensors
type ScheduleProvider struct {
	readingProvider     readingprovider.ReadingProvider
	readingChannel      chan readingprovider.ReadingProvider
	persistenceProvider *persistenceprovider.PersistenceProvider
	Timers              []IntervalTimer `json:"timers"`
}

//IntervalTimer defines an interval and a read configuration for that interval
type IntervalTimer struct {
	SensorAddress       uint8          `json:"sensorAddress"`
	ReadType            string         `json:"readType"`
	StartLocation       uint16         `json:"startLocation"`
	ReadLength          uint16         `json:"readLength"`
	Interval            *time.Duration `json:"interval"`
	Repeat              bool           `json:"repeat"`
	FirstTime           *time.Time     `json:"time"`
	Persist             bool           `json:"store"`
	readingChannel      chan readingprovider.ReadingProvider
	persistenceProvider *persistenceprovider.PersistenceProvider
	timer               *time.Timer
	ticker              *time.Ticker
}

//Start for IntervalTimer
//is meant to be called by schedule provider
//and will start a timer that will perform a read
func (intervalTimer *IntervalTimer) Start() {
	if intervalTimer.FirstTime == nil {
		intervalTimer.startReading()
	} else {
		if intervalTimer.Interval == nil {
			return
		}
		var firstTime time.Time
		firstTime = *intervalTimer.FirstTime
		//TODO Find a better method to calculate when to run first
		//as this can take a lot of time for a FirstTime set well before
		//and a short interval
		for firstTime.Before(time.Now()) {
			firstTime = firstTime.Add(*intervalTimer.Interval)
		}
		startIn := firstTime.Sub(time.Now())
		intervalTimer.timer = time.AfterFunc(startIn, intervalTimer.startReading)
	}
}

//startReading is called to start reading periodically
func (intervalTimer *IntervalTimer) startReading() {
	if intervalTimer.Interval != nil {
		intervalTimer.ticker = time.NewTicker(*intervalTimer.Interval)
		go func() {
			for t := range intervalTimer.ticker.C {
				log.Printf("Reading sensor %d, start location=%d, length=%d, type=%s, %v",
					intervalTimer.SensorAddress, intervalTimer.StartLocation,
					intervalTimer.ReadLength, intervalTimer.ReadType, t)
				intervalTimer.read()
			}
		}()
	}

}

func (intervalTimer *IntervalTimer) read() {
	if intervalTimer.readingChannel == nil {
		log.Println("No reading provider defined.")
		return
	}
	readingProvider := <-intervalTimer.readingChannel
	reading, err := readingProvider.GetReading(intervalTimer.SensorAddress,
		intervalTimer.ReadType, intervalTimer.StartLocation,
		intervalTimer.ReadLength)
	intervalTimer.readingChannel <- readingProvider

	if err != nil {
		log.Printf("Error: %s, sensor %d, start %d, length %d, type %s\n",
			err.Error(), intervalTimer.SensorAddress, intervalTimer.StartLocation,
			intervalTimer.ReadLength, intervalTimer.ReadType)
		return
	}
	if intervalTimer.Persist {
		if reading != nil {
			err = (*intervalTimer.persistenceProvider).SaveSensorReading(*reading)
			if err != nil {
				log.Printf("Error persisting %s\n", err.Error())
			}
		}
	}
}

//Stop will stop the ticker so that no more reading will happen
func (intervalTimer *IntervalTimer) Stop() {
	log.Printf("Stopping %d, %d,%d, %s\n", intervalTimer.SensorAddress,
		intervalTimer.StartLocation, intervalTimer.ReadLength,
		intervalTimer.ReadType)
	if intervalTimer.ticker != nil {
		intervalTimer.ticker.Stop()
	}
	if intervalTimer.timer != nil {
		intervalTimer.timer.Stop()
	}
}

//NewScheduleProvider initializes a ScheduleProvider and creates a channel for
//reading
func (ScheduleProvider) NewScheduleProvider(rp readingprovider.ReadingProvider, pp *persistenceprovider.PersistenceProvider) *ScheduleProvider {
	schProvider := ScheduleProvider{readingProvider: rp, persistenceProvider: pp}
	schProvider.readingChannel = make(chan readingprovider.ReadingProvider, 1)
	return &schProvider
}

//AddTimer adds an interval timer to the schedule provider
func (schProvider *ScheduleProvider) AddTimer(intervalTimer IntervalTimer) error {
	if schProvider.readingProvider == nil {
		return fmt.Errorf("No reading provider defined!")
	}
	if intervalTimer.Persist && schProvider.persistenceProvider == nil {
		return fmt.Errorf("Error adding timer with persistence: No persistece provider defined!")
	}
	intervalTimer.readingChannel = schProvider.readingChannel
	intervalTimer.persistenceProvider = schProvider.persistenceProvider
	schProvider.Timers = append(schProvider.Timers, intervalTimer)
	return nil
}

//Start for ScheduleProvider
//will generate a go routine for each IntervalTimer
func (schProvider *ScheduleProvider) Start() {
	for _, interval := range schProvider.Timers {
		//start a go routine for each interval
		i := interval
		go i.Start()
	}
	schProvider.readingChannel <- schProvider.readingProvider
}

//Stop send the signal to stop to all IntervalTimers
func (schProvider *ScheduleProvider) Stop() {
	for _, interval := range schProvider.Timers {
		interval.Stop()
	}
}
