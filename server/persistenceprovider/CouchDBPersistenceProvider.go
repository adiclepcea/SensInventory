package persistenceprovider

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/adiclepcea/SensInventory/server/common"
	"github.com/patrickjuchli/couch"
)

const (
	defaultCouchDBDatabase        = "sensinventory"
	viewOneQueryPrefix            = "_design/readings/_view/sensorTime?key="
	viewSensorInPeriodQueryPrefix = "_design/readings/_view/sensorTime?startkey="
	viewAllInPeriodQueryPrefix    = "_design/readings/_view/byTime?startkey="
)

//CouchDBPersistenceProvider defines the structure for a
//CouchDB persistence layer for this project
type CouchDBPersistenceProvider struct {
	CouchServer      string
	CouchDatabase    string
	CouchCredentials *couch.Credentials
	Database         couch.Database
	PersistenceProvider
}

//CouchDBReading is used to operate with Readings in CouchDB
type CouchDBReading struct {
	couch.Doc
	Reading common.Reading
}

type couchDBRow struct {
	ID      string         `json:"id"`
	Key     interface{}    `json:"key"`
	Value   interface{}    `json:"value"`
	Reading CouchDBReading `json:"doc,omitempty"`
}

type couchDBResult struct {
	TotalRows int          `json:"total_rows"`
	Offset    int          `json:"offset"`
	Rows      []couchDBRow `json:"rows"`
}

//NewPersistenceProvider creates a new persistence
//provider that will save data in CouchDB
func (CouchDBPersistenceProvider) NewPersistenceProvider(params ...string) (PersistenceProvider, error) {
	if len(params) == 0 {
		return nil, fmt.Errorf("No parameters given for connection")
	}

	var couchProvider CouchDBPersistenceProvider
	lenParams := len(params)
	if lenParams == 2 || lenParams == 4 {
		couchProvider = CouchDBPersistenceProvider{CouchDatabase: params[lenParams-1]}
	} else {
		couchProvider = CouchDBPersistenceProvider{CouchDatabase: defaultCouchDBDatabase}
	}

	couchProvider.CouchServer = params[0]
	if lenParams >= 3 {
		couchProvider.CouchCredentials = couch.NewCredentials(params[1], params[2])
	}

	log.Printf("Using server %s", couchProvider.CouchServer)

	if _, err := couchProvider.CreateDB(); err != nil {
		return nil, err
	}

	return &couchProvider, nil
}

func (couchProvider *CouchDBPersistenceProvider) createViews(views map[string]string) error {
	viewCode := make(map[string]interface{})
	for viewName, viewFunction := range views {
		mapFunction := struct {
			Map string `json:"map"`
		}{viewFunction}
		viewCode[viewName] = mapFunction
	}

	view := make(map[string]interface{})
	view["_id"] = "_design/readings"
	view["language"] = "javascript"
	view["views"] = viewCode

	var response interface{}

	_, err := couch.Do(couchProvider.CouchServer+"/"+couchProvider.CouchDatabase,
		"POST", couchProvider.CouchCredentials, view, &response)

	return err
}

//CreateDB creates the database if it does not exist
func (couchProvider *CouchDBPersistenceProvider) CreateDB() (*couch.Database, error) {
	server := couch.NewServer(couchProvider.CouchServer, couchProvider.CouchCredentials)

	db := server.Database(couchProvider.CouchDatabase)

	if !db.Exists() {
		log.Printf("Database %s does not exist yet. Creating it", couchProvider.CouchDatabase)
		if err := db.Create(); err != nil {
			return nil, err
		}
		mapViews := make(map[string]string)
		mapViews["sensorTime"] = "function(doc){if(doc.Reading) emit([doc.Reading.sensor,doc.Reading.time],null);}"
		mapViews["byTime"] = "function(doc){if(doc.Reading) emit(doc.Reading.time,null);}"
		if err := couchProvider.createViews(mapViews); err != nil {
			couchProvider.DeleteDB()
			return nil, err
		}
	}
	if !db.Exists() {
		return nil, fmt.Errorf("Database %s could not be created", couchProvider.CouchDatabase)
	}

	return db, nil
}

//DeleteDB deletes the database from CouchDB
func (couchProvider *CouchDBPersistenceProvider) DeleteDB() error {
	server := couch.NewServer(couchProvider.CouchServer, couchProvider.CouchCredentials)
	db := server.Database(couchProvider.CouchDatabase)

	if !db.Exists() {
		return nil
	}

	return db.DropDatabase()
}

//SaveSensorReading saves the reading in the database
func (couchProvider *CouchDBPersistenceProvider) SaveSensorReading(reading common.Reading) error {
	db, err := couchProvider.CreateDB()
	if err != nil {
		return err
	}
	couchReading := &CouchDBReading{Reading: reading}

	return db.Insert(couchReading)

}

func (couchProvider CouchDBPersistenceProvider) getBaseQueryString() string {
	server := couchProvider.CouchServer
	if !strings.HasSuffix(server, "/") {
		server = server + "/"
	}
	database := couchProvider.CouchDatabase
	if !strings.HasSuffix(database, "/") {
		database = database + "/"
	}

	return server + database
}

func getViewQueryStringOneOnly(sensor string, timeIn time.Time) string {
	strTime := timeIn.Format(common.TimeFormat)
	return fmt.Sprintf("%s[%s,\"%s\"]&include_docs=true", viewOneQueryPrefix, sensor, strTime)
}

func getViewQueryStringSensorInPeriod(sensor string, startTime time.Time, endTime time.Time) string {
	strStartTime := startTime.Format(common.TimeFormat)
	strEndTime := endTime.Format(common.TimeFormat)
	return fmt.Sprintf("%s[%s,\"%s\"]&endkey=[%s,\"%s\"]&ascending=true&include_docs=true",
		viewSensorInPeriodQueryPrefix, sensor, strStartTime, sensor, strEndTime)
}

func getViewQueryStringInPeriod(startTime time.Time, endTime time.Time) string {
	strStartTime := startTime.Format(common.TimeFormat)
	strEndTime := endTime.Format(common.TimeFormat)
	return fmt.Sprintf("%s\"%s\"&endkey=\"%s\"&ascending=true&include_docs=true",
		viewAllInPeriodQueryPrefix, strStartTime, strEndTime)
}

//GetCouchDBReadings returns the CouchDBReadings (incuding _id and _rev)
//from the query
func (couchProvider *CouchDBPersistenceProvider) GetCouchDBReadings(query string) (*[]CouchDBReading, error) {
	var resp couchDBResult
	_, err := couch.Do(couchProvider.getBaseQueryString()+query,
		"GET", couchProvider.CouchCredentials, nil, &resp)
	log.Printf(couchProvider.getBaseQueryString() + query)
	if err != nil {
		log.Printf("Error  asking for results %s returned %s",
			couchProvider.getBaseQueryString()+query, err.Error())
		return nil, err
	}
	var rows []CouchDBReading
	for _, row := range resp.Rows {
		rows = append(rows, row.Reading)
		log.Println(row.Reading)
	}
	return &rows, nil
}

//GetSensorReading returns the reading for the sensor with address
//"sensorAddress" at the time "time"
func (couchProvider *CouchDBPersistenceProvider) GetSensorReading(sensorAddress uint8, time time.Time) (*common.Reading, error) {

	query := getViewQueryStringOneOnly(strconv.Itoa(int(sensorAddress)), time)

	resp, err := couchProvider.GetCouchDBReadings(query)

	if err != nil {
		return nil, err
	}

	if len(*resp) == 0 {
		return nil, nil
	}

	return &(*resp)[0].Reading, nil
	// curl 'http://localhost:5984/sensinventory/_design/readings/_view/sensorTime?key=\[10,"2016-06-13%2023:28:48"\]&include_docs=true'
}

//GetSensorReadingsInPeriod returns the readings for a sensor in a given period
func (couchProvider *CouchDBPersistenceProvider) GetSensorReadingsInPeriod(
	sensorAddress uint8, startTime time.Time,
	endTime time.Time) ([]common.Reading, error) {

	query := getViewQueryStringSensorInPeriod(strconv.Itoa(int(sensorAddress)), startTime, endTime)

	resp, err := couchProvider.GetCouchDBReadings(query)
	if err != nil {
		return nil, err
	}
	if len(*resp) == 0 {
		return nil, nil
	}

	var rows []common.Reading
	for _, row := range *resp {
		rows = append(rows, row.Reading)
		log.Println(row.Reading)
	}
	return rows, nil

}

//GetSensorReadingCountInPeriod returns the number of readings
//for the given sensorAddress in the given period
func (couchProvider *CouchDBPersistenceProvider) GetSensorReadingCountInPeriod(
	sensorAddress uint8, startTime time.Time,
	endTime time.Time) (uint, error) {

	query := getViewQueryStringSensorInPeriod(strconv.Itoa(int(sensorAddress)), startTime, endTime)

	resp, err := couchProvider.GetCouchDBReadings(query)

	if err != nil {
		return 0, nil
	}

	return uint(len(*resp)), nil
}

//GetAllReadingsInPeriod returns all the readings in the period,
//without filtering by the sensorAddress
func (couchProvider *CouchDBPersistenceProvider) GetAllReadingsInPeriod(
	startTime time.Time, endTime time.Time) (*[]common.Reading, error) {

	query := getViewQueryStringInPeriod(startTime, endTime)

	resp, err := couchProvider.GetCouchDBReadings(query)

	if err != nil {
		return nil, err
	}

	if len(*resp) == 0 {
		return nil, nil
	}

	var rows []common.Reading
	for _, row := range *resp {
		rows = append(rows, row.Reading)
		log.Println(row.Reading)
	}
	return &rows, nil
}

//GetAllReadingsCountInPeriod returns the number of readings in the given period
func (couchProvider *CouchDBPersistenceProvider) GetAllReadingsCountInPeriod(
	startTime time.Time, endTime time.Time) (uint, error) {

	query := getViewQueryStringInPeriod(startTime, endTime)

	resp, err := couchProvider.GetCouchDBReadings(query)

	if err != nil {
		return 0, nil
	}

	return uint(len(*resp)), nil

}

//DeleteSensorReading deletes the specified sensor Reading
func (couchProvider *CouchDBPersistenceProvider) DeleteSensorReading(
	sensorAddress uint8, time time.Time) error {
	query := getViewQueryStringOneOnly(strconv.Itoa(int(sensorAddress)), time)

	resp, err := couchProvider.GetCouchDBReadings(query)

	if err != nil {
		return err
	}

	server := couch.NewServer(couchProvider.CouchServer, couchProvider.CouchCredentials)

	db := server.Database(couchProvider.CouchDatabase)
	if !db.Exists() {
		return fmt.Errorf("Database missing")
	}

	if len(*resp) > 0 {
		return db.Delete((*resp)[0].ID, (*resp)[0].Rev)
	}

	return fmt.Errorf("This reading was not found in the database")
}

//DeleteSensorReadingsInPeriod deletes all the readings from the specified
//sensor in the specified period of time
func (couchProvider *CouchDBPersistenceProvider) DeleteSensorReadingsInPeriod(
	sensorAddress uint8, startTime time.Time, endTime time.Time) error {

	query := getViewQueryStringSensorInPeriod(strconv.Itoa(int(sensorAddress)), startTime, endTime)

	resp, err := couchProvider.GetCouchDBReadings(query)

	if err != nil {
		return err
	}

	server := couch.NewServer(couchProvider.CouchServer, couchProvider.CouchCredentials)

	db := server.Database(couchProvider.CouchDatabase)

	if !db.Exists() {
		return fmt.Errorf("Database missing")
	}

	if len(*resp) > 0 {
		for _, r := range *resp {
			if err = db.Delete(r.ID, r.Rev); err != nil {
				return err
			}
		}
	} else {
		return fmt.Errorf(
			"No reading was not found in the database for this sensor in this period")
	}

	return nil

}

//DeleteAllReadingsInPeriod deletes all the readings in the specified period
//no mather the sensor
func (couchProvider *CouchDBPersistenceProvider) DeleteAllReadingsInPeriod(
	startTime time.Time, endTime time.Time) error {

	query := getViewQueryStringInPeriod(startTime, endTime)

	resp, err := couchProvider.GetCouchDBReadings(query)

	if err != nil {
		return err
	}

	server := couch.NewServer(couchProvider.CouchServer, couchProvider.CouchCredentials)

	db := server.Database(couchProvider.CouchDatabase)

	if !db.Exists() {
		return fmt.Errorf("Database missing")
	}

	if len(*resp) > 0 {
		for _, r := range *resp {
			if err = db.Delete(r.ID, r.Rev); err != nil {
				return err
			}
		}
	} else {
		return fmt.Errorf(
			"No reading was not found in the database for this sensor in this period")
	}

	return nil
}
