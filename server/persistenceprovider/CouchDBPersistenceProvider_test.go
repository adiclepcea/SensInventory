// +build integration

package persistenceprovider_test

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/adiclepcea/SensInventory/server/common"
	pp "github.com/adiclepcea/SensInventory/server/persistenceprovider"
	"github.com/adiclepcea/SensInventory/server/readgroups"
)

var (
	testServer    string = "http://127.0.0.1:5984"
	testServerBad string = "http://127.0.0.1"
	username      *string
	password      *string
	reading1      common.Reading
	reading2      common.Reading
	reading3      common.Reading
)

func initTest() {
	envServer := os.Getenv("COUCH_TEST_SERVER")
	envServerBad := os.Getenv("COUCH_TEST_SERVER_BAD")
	envUsername := os.Getenv("COUCH_USER")
	envPassword := os.Getenv("COUCH_PASSWORD")

	if envServer != "" {
		testServer = envServer
	}
	if envServerBad != "" {
		testServerBad = envServerBad
	}

	if envUsername != "" && envPassword != "" {
		username = &envUsername
		password = &envPassword
	}

	log.Printf("Using server=%s, bad server=%s\n", testServer, testServerBad)
	if username != nil {
		log.Printf("Using username=%s, password=%s", username, password)
	}

	reading1.Type = common.Holding
	reading1.Count = 3
	reading1.Sensor = 10
	reading1.StartLocation = 0
	reading1.ReadValues = []uint16{100, 200, 300}
	reading1.Time = time.Now().Format(common.TimeFormat)
	reading1.InitCalculatedValues()

	rg, _ := readgroups.ReadGroupFloat32{}.NewReadGroup(10, 0)
	rg.ResultType = common.Float32
	rg.Calculate(&reading1)

	reading2.Type = common.Coil
	reading2.Count = 8
	reading2.Sensor = 12
	reading2.StartLocation = 0
	reading2.ReadValues = []uint16{0x255}
	reading2.Time = time.Now().Add(time.Duration(15) * time.Second).Format(common.TimeFormat)

	reading3.Type = common.Coil
	reading3.Count = 8
	reading3.Sensor = 13
	reading3.StartLocation = 0
	reading3.ReadValues = []uint16{0x255}
	reading3.Time = time.Now().Add(time.Duration(65) * time.Second).Format(common.TimeFormat)

}

func ConnectToCouch() (*pp.CouchDBPersistenceProvider, error) {
	initTest()
	var cdbp *pp.CouchDBPersistenceProvider
	var err error
	if username != nil {
		cdbp, err = pp.CouchDBPersistenceProvider{}.NewPersistenceProvider(testServer, *username, *password)
	} else {
		cdbp, err = pp.CouchDBPersistenceProvider{}.NewPersistenceProvider(testServer)
	}
	return cdbp, err
}

func TestNewPersistenceProviderShouldOK(t *testing.T) {
	initTest()
	var cdbp *pp.CouchDBPersistenceProvider
	var err error

	if username != nil {
		t.Log("Trying to create a CouchDBProvider with user and pass arguments")
		cdbp, err = pp.CouchDBPersistenceProvider{}.NewPersistenceProvider(testServer, *username, *password)
		if err != nil {
			t.Fatalf("No error is expected here, got %s", err.Error())
		}
		if cdbp.CouchCredentials == nil {
			t.Fatalf("Credentials should be set here, got nil")
		}

		cdbp.DeleteDB()

		t.Log("Trying to create a CouchDBProvider with user and pass arguments and dbname")
		cdbp, err = pp.CouchDBPersistenceProvider{}.NewPersistenceProvider(testServer, *username, *password, "fakedb")
		if err != nil {
			t.Fatalf("No error is expected here, got %s", err.Error())
		}
		defer cdbp.DeleteDB()
	} else {
		t.Log("Trying to create a CouchDBProvider with right arguments")
		cdbp, err = pp.CouchDBPersistenceProvider{}.NewPersistenceProvider(testServer)
		if err != nil {
			t.Fatalf("No error is expected here, got %s", err.Error())
		}
		if cdbp.CouchCredentials != nil {
			t.Fatalf("No credentials should be set here, got not nil")
		}
		cdbp.DeleteDB()
		t.Log("Trying to create a CouchDBProvider with dbname")
		cdbp, err = pp.CouchDBPersistenceProvider{}.NewPersistenceProvider(testServer, "fakedb")
		if err != nil {
			t.Fatalf("No error is expected here, got %s", err.Error())
		}
		defer cdbp.DeleteDB()
	}
	if cdbp.CouchCredentials != nil {
		t.Fatalf("No credentials should be set here, got not nil")
	}

	if cdbp.CouchDatabase != "fakedb" {
		t.Fatalf("Wrong database found. Expected fakedb got %s", cdbp.CouchDatabase)
	}

}

func TestNewPersistenceProviderShouldFail(t *testing.T) {
	initTest()
	t.Log("Trying to create a CouchDBProvider without arguments")
	cdbp, err := pp.CouchDBPersistenceProvider{}.NewPersistenceProvider()
	if err == nil {
		defer cdbp.DeleteDB()
		t.Fatalf("An error should have occured when calling without params")
	}

	t.Log("Trying to create a CouchDBProvider with wrong arguments")

	_, err = pp.CouchDBPersistenceProvider{}.NewPersistenceProvider(testServerBad)
	if err == nil {
		t.Fatalf("Error is expected here, got nil")
	}

}

func TestSaveSensorReadingShouldOK(t *testing.T) {
	cdbp, err := ConnectToCouch()
	if err != nil {
		t.Fatal("No error expected when connecting to couchdb, got ", err.Error())
	}
	defer cdbp.DeleteDB()
	err = cdbp.SaveSensorReading(reading1)
	if err != nil {
		t.Fatal("No error expected when saving a reading, got ", err.Error())
	}
}

func TestGetSensorReading(t *testing.T) {
	cdbp, err := ConnectToCouch()
	if err != nil {
		t.Fatal("No error Expected when connecting to couchdb, got ", err.Error())
	}
	defer cdbp.DeleteDB()
	timeIn, err := time.Parse(common.TimeFormat, reading1.Time)
	if err != nil {
		t.Fatal("No error expected when getting time from string")
	}

	err = cdbp.SaveSensorReading(reading1)
	if err != nil {
		t.Fatal("No error expected when saving a reading, got ", err.Error())
	}

	reading, err := cdbp.GetSensorReading(reading1.Sensor, timeIn)

	if err != nil {
		t.Fatalf("No error expected when retrieving a record. Got %s", err.Error())
	}
	if reading == nil {
		t.Fatal("No nul result expected when retrieving a record. Got nil")
	}
	if reading.Sensor != reading1.Sensor || reading.Time != reading1.Time ||
		reading.StartLocation != reading1.StartLocation {
		t.Fatalf("Expected %d, %s, %d, %v, got %d,%s,%d,%v",
			reading1.Sensor, reading1.Time, reading1.StartLocation, reading1.CalculatedValues,
			reading.Sensor, reading.Time, reading.StartLocation, reading.CalculatedValues)
	}

}

func TestGetSensorReadingsInPeriod(t *testing.T) {
	cdbp, err := ConnectToCouch()
	if err != nil {
		t.Fatal("No error Expected when connecting to couchdb, got ", err.Error())
	}
	defer cdbp.DeleteDB()
	startTime, err := time.Parse(common.TimeFormat, reading1.Time)
	intermediaryTime, err := time.Parse(common.TimeFormat, reading2.Time)
	endTime, err := time.Parse(common.TimeFormat, reading3.Time)
	if err != nil {
		t.Fatal("No error expected when getting time from string")
	}
	reading2.Sensor = 10
	reading3.Sensor = 10
	err = cdbp.SaveSensorReading(reading1)
	cdbp.SaveSensorReading(reading2)
	cdbp.SaveSensorReading(reading3)
	reading2.Sensor = 12
	reading3.Sensor = 13
	if err != nil {
		t.Fatal("No error expected when saving a reading, got ", err.Error())
	}

	readings, err := cdbp.GetSensorReadingsInPeriod(reading1.Sensor, startTime, endTime)
	readingsIntermediary, err := cdbp.GetSensorReadingsInPeriod(reading1.Sensor, startTime, intermediaryTime)

	if err != nil {
		t.Fatalf("No error expected when retrieving records in a period. Got %s", err.Error())
	}
	if readings == nil {
		t.Fatal("No nul result expected when retrieving records in a period. Got nil")
	}
	if len(*readings) != 3 {
		t.Fatalf("Expected 3 readings, got %s", len(*readings))
	}

	if readingsIntermediary == nil {
		t.Fatal("Intermediary: No nul result expected when retrieving records in a period. Got nil")
	}
	if len(*readingsIntermediary) != 2 {
		t.Fatalf("Expected 2 readings, got %d", len(*readingsIntermediary))
	}

}

func TestGetCountSensorReadingInPeriod(t *testing.T) {
	cdbp, err := ConnectToCouch()
	if err != nil {
		t.Fatal("No error Expected when connecting to couchdb, got ", err.Error())
	}

	defer cdbp.DeleteDB()

	startTime, err := time.Parse(common.TimeFormat, reading1.Time)
	intermediaryTime, err := time.Parse(common.TimeFormat, reading2.Time)
	endTime, err := time.Parse(common.TimeFormat, reading3.Time)
	if err != nil {
		t.Fatal("No error expected when getting time from string")
	}
	reading2.Sensor = 10
	reading3.Sensor = 10
	err = cdbp.SaveSensorReading(reading1)
	cdbp.SaveSensorReading(reading2)
	cdbp.SaveSensorReading(reading3)
	reading2.Sensor = 12
	reading3.Sensor = 13
	if err != nil {
		t.Fatal("No error expected when saving a reading, got ", err.Error())
	}

	noOfReadings, err := cdbp.GetSensorReadingCountInPeriod(reading1.Sensor, startTime, endTime)
	noOfReadingsIntermediary, err := cdbp.GetSensorReadingCountInPeriod(reading1.Sensor, startTime, intermediaryTime)

	if err != nil {
		t.Fatalf("No error expected when retrieving records in a period. Got %s", err.Error())
	}
	if noOfReadings != 3 {
		t.Fatalf("3 results expected got %d", noOfReadings)
	}

	if noOfReadingsIntermediary != 2 {
		t.Fatalf("Expected 2 readings, got %d", noOfReadingsIntermediary)
	}

}

func TestGetReadingsInPeriod(t *testing.T) {
	cdbp, err := ConnectToCouch()
	if err != nil {
		t.Fatal("No error Expected when connecting to couchdb, got ", err.Error())
	}
	defer cdbp.DeleteDB()
	startTime, err := time.Parse(common.TimeFormat, reading1.Time)
	intermediaryTime, err := time.Parse(common.TimeFormat, reading2.Time)
	endTime, err := time.Parse(common.TimeFormat, reading3.Time)
	if err != nil {
		t.Fatal("No error expected when getting time from string")
	}
	err = cdbp.SaveSensorReading(reading1)
	cdbp.SaveSensorReading(reading2)
	cdbp.SaveSensorReading(reading3)
	if err != nil {
		t.Fatal("No error expected when saving a reading, got ", err.Error())
	}

	readings, err := cdbp.GetAllReadingsInPeriod(startTime, endTime)
	readingsIntermediary, err := cdbp.GetAllReadingsInPeriod(startTime, intermediaryTime)

	if err != nil {
		t.Fatalf("No error expected when retrieving records in a period. Got %s", err.Error())
	}
	if readings == nil {
		t.Fatal("No nul result expected when retrieving records in a period. Got nil")
	}
	if len(*readings) != 3 {
		t.Fatalf("Expected 3 readings, got %s", len(*readings))
	}

	if readingsIntermediary == nil {
		t.Fatal("Intermediary: No nul result expected when retrieving records in a period. Got nil")
	}
	if len(*readingsIntermediary) != 2 {
		t.Fatalf("Expected 2 readings, got %d", len(*readingsIntermediary))
	}

}

func TestGetReadingsCountInPeriod(t *testing.T) {
	cdbp, err := ConnectToCouch()
	if err != nil {
		t.Fatal("No error Expected when connecting to couchdb, got ", err.Error())
	}
	defer cdbp.DeleteDB()
	startTime, err := time.Parse(common.TimeFormat, reading1.Time)
	intermediaryTime, err := time.Parse(common.TimeFormat, reading2.Time)
	endTime, err := time.Parse(common.TimeFormat, reading3.Time)
	if err != nil {
		t.Fatal("No error expected when getting time from string")
	}
	err = cdbp.SaveSensorReading(reading1)
	cdbp.SaveSensorReading(reading2)
	cdbp.SaveSensorReading(reading3)
	if err != nil {
		t.Fatal("No error expected when saving a reading, got ", err.Error())
	}

	noOfReadings, err := cdbp.GetAllReadingsCountInPeriod(startTime, endTime)
	noOfReadingsIntermediary, err := cdbp.GetAllReadingsCountInPeriod(startTime, intermediaryTime)

	if err != nil {
		t.Fatalf("No error expected when retrieving records in a period. Got %s", err.Error())
	}
	if noOfReadings != 3 {
		t.Fatalf("3 results expected got %d", noOfReadings)
	}

	if noOfReadingsIntermediary != 2 {
		t.Fatalf("Expected 2 readings, got %d", noOfReadingsIntermediary)
	}

}

func TestDeleteSensorReading(t *testing.T) {
	cdbp, err := ConnectToCouch()
	if err != nil {
		t.Fatal("No error Expected when connecting to couchdb, got ", err.Error())
	}
	startTime, err := time.Parse(common.TimeFormat, reading1.Time)
	intermediaryTime, err := time.Parse(common.TimeFormat, reading2.Time)
	if err != nil {
		t.Fatal("No error expected when getting time from string")
	}
	err = cdbp.SaveSensorReading(reading1)
	cdbp.SaveSensorReading(reading2)
	cdbp.SaveSensorReading(reading3)

	if err != nil {
		cdbp.DeleteDB()
		t.Fatal("No error expected when saving a reading, got ", err.Error())
	}

	err = cdbp.DeleteSensorReading(reading1.Sensor, startTime)

	if err != nil {
		cdbp.DeleteDB()
		t.Fatalf("No error expected while deleting an existing reading. Got", err.Error())
	}

	err = cdbp.DeleteSensorReading(reading1.Sensor, startTime)
	if err == nil {
		cdbp.DeleteDB()
		t.Fatalf("Error expected while deleting an inexistent reading. Got", err)
	}

	cdbp.CouchDatabase = "inexistent"
	err = cdbp.DeleteSensorReading(reading2.Sensor, intermediaryTime)
	if err == nil {
		t.Fatalf("Error expected while deleting a reading from an inexistent db. Got", err)
	}
	cdbp.DeleteDB()
}

func TestDeleteSensorReadingsInPeriod(t *testing.T) {
	cdbp, err := ConnectToCouch()
	if err != nil {
		t.Fatal("No error Expected when connecting to couchdb, got ", err.Error())
	}

	startTime, err := time.Parse(common.TimeFormat, reading1.Time)
	endTime, err := time.Parse(common.TimeFormat, reading3.Time)

	if err != nil {
		t.Fatal("No error expected when getting time from string")
	}

	reading2.Sensor = 10
	reading3.Sensor = 10
	err = cdbp.SaveSensorReading(reading1)
	cdbp.SaveSensorReading(reading2)
	cdbp.SaveSensorReading(reading3)
	reading2.Sensor = 12
	reading3.Sensor = 13

	if err != nil {
		cdbp.DeleteDB()
		t.Fatal("No error expected when saving a reading, got ", err.Error())
	}

	err = cdbp.DeleteSensorReadingsInPeriod(reading1.Sensor, startTime, endTime)

	if err != nil {
		cdbp.DeleteDB()
		t.Fatalf("No error expected while deleting an existing reading. Got", err.Error())
	}

	err = cdbp.DeleteSensorReadingsInPeriod(reading1.Sensor, startTime, endTime)
	if err == nil {
		cdbp.DeleteDB()
		t.Fatalf("Error expected while deleting an inexistent reading. Got", err)
	}

	cdbp.CouchDatabase = "inexistent"
	err = cdbp.DeleteSensorReadingsInPeriod(reading3.Sensor, startTime, endTime)
	if err == nil {
		t.Fatalf("Error expected while deleting a reading from an inexistent db. Got", err)
	}
}

func TestDeleterReadingsInPeriod(t *testing.T) {
	cdbp, err := ConnectToCouch()
	if err != nil {
		t.Fatal("No error Expected when connecting to couchdb, got ", err.Error())
	}

	startTime, err := time.Parse(common.TimeFormat, reading1.Time)
	endTime, err := time.Parse(common.TimeFormat, reading3.Time)

	if err != nil {
		t.Fatal("No error expected when getting time from string")
	}

	err = cdbp.SaveSensorReading(reading1)
	cdbp.SaveSensorReading(reading2)
	cdbp.SaveSensorReading(reading3)

	if err != nil {
		cdbp.DeleteDB()
		t.Fatal("No error expected when saving a reading, got ", err.Error())
	}

	err = cdbp.DeleteAllReadingsInPeriod(startTime, endTime)

	if err != nil {
		cdbp.DeleteDB()
		t.Fatalf("No error expected while deleting an existing reading. Got", err.Error())
	}

	err = cdbp.DeleteAllReadingsInPeriod(startTime, endTime)
	if err == nil {
		cdbp.DeleteDB()
		t.Fatalf("Error expected while deleting an inexistent reading. Got", err)
	}

	cdbp.CouchDatabase = "inexistent"
	err = cdbp.DeleteAllReadingsInPeriod(startTime, endTime)
	if err == nil {
		t.Fatalf("Error expected while deleting a reading from an inexistent db. Got", err)
	}
}
