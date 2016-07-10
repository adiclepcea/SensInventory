package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/adiclepcea/SensInventory/server/common"
	"github.com/adiclepcea/SensInventory/server/configprovider"
	"github.com/adiclepcea/SensInventory/server/persistenceprovider"
	"github.com/adiclepcea/SensInventory/server/readingprovider"
	"github.com/julienschmidt/httprouter"
)

//ErrorMessage is used to transmit an error message
type ErrorMessage struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

var configProvider configprovider.ConfigProvider
var persistenceProvider persistenceprovider.PersistenceProvider
var readingProvider readingprovider.ReadingProvider
var scheduleProvider *readingprovider.ScheduleProvider

func initialize() {
	var err error
	//choose the desired config provider
	configProvider, err = configprovider.FileConfigProvider{}.NewConfigProvider()
	if err != nil {
		log.Fatalf("Error initializing config provider: %s\n", err.Error())
	}
	configProvider.SetAddressLimits(0, 30)

	//choose the desired persistence provider
	persistenceProvider, err = persistenceprovider.CouchDBPersistenceProvider{}.NewPersistenceProvider("http://127.0.0.1:5984")
	if err != nil {
		log.Fatalf("Error initializing persistence provider: %s\n", err.Error())
	}

	//For now the ModBUSReadingProvider is the only one
	//This was the purpose anyway
	readingProvider = readingprovider.ModBUSReadingProvider{}.NewReadingProvider(&configProvider)

	scheduleProvider = readingprovider.ScheduleProvider{}.NewScheduleProvider(readingProvider, &persistenceProvider)
	scheduleProvider.Start()

}

func errorToJSONByteArray(errorString string, err error) []byte {
	errMsg := ErrorMessage{Error: errorString, Message: err.Error()}
	msg, err := json.Marshal(errMsg)
	if err != nil {
		log.Printf("Could not marshal message: %v", errMsg)
		return nil
	}
	log.Printf("%s", msg)
	return msg
}

func returnSuccess(w http.ResponseWriter) {
	success := struct {
		Result string `json:"result"`
	}{Result: "OK"}
	encoder := json.NewEncoder(w)
	encoder.Encode(success)
}

func saveSchedule(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	err := scheduleProvider.Save()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errorToJSONByteArray("could not save schedule", err))
		return
	}
	returnSuccess(w)
}

func loadSchedule(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	err := scheduleProvider.Load()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errorToJSONByteArray("could not save schedule", err))
		return
	}
	returnSuccess(w)
}

func getTimers(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	it := scheduleProvider.Timers

	encoder := json.NewEncoder(w)
	if it == nil {
		encoder.Encode(map[string]interface{}{})
		return
	}
	encoder.Encode(it)
}

func deleteTimer(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var err error
	var itAddress int
	w.Header().Add("Content-Type", "application/json")
	itString := p.ByName("timer")
	if itAddress, err = strconv.Atoi(itString); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errorToJSONByteArray("could not convert to valid sensor address", err))
		return
	}
	log.Printf("Deleting sensor %d\n", itAddress)
	err = scheduleProvider.RemoveTimer(itAddress)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write(errorToJSONByteArray("could not delete timer", err))
		return
	}
	returnSuccess(w)
}

func addTimer(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Header().Add("Content-Type", "application/json")
	it, err := getTimerFromBody(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errorToJSONByteArray("no valid timer received", err))
		return
	}
	log.Printf("Adding timer for sensor %d\n", it.SensorAddress)

	err = scheduleProvider.AddTimer(*it)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errorToJSONByteArray("could not add timer", err))
		return
	}
	w.WriteHeader(http.StatusCreated)
	returnSuccess(w)
}

func getTimerFromBody(r *http.Request) (*readingprovider.IntervalTimer, error) {
	decoder := json.NewDecoder(r.Body)
	var it readingprovider.IntervalTimer
	err := decoder.Decode(&it)
	if err != nil {
		return nil, err
	}

	return &it, nil

}

func getSensors(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	sensors := configProvider.GetSensors()
	rez := make(map[string]string)
	for addr, sensor := range sensors {
		rez[addr] = sensor.Description
	}
	encoder := json.NewEncoder(w)
	encoder.Encode(rez)
}

func getSensor(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var sensorAddress int
	var err error
	var sensor *common.Sensor

	w.Header().Add("Content-Type", "application/json")
	sensorString := p.ByName("sensor")

	if sensorAddress, err = strconv.Atoi(sensorString); err != nil {
		w.Write(errorToJSONByteArray("could not convert to valid sensor address", err))
		return
	}
	sensor, err = configProvider.GetSensorByAddress(uint8(sensorAddress))

	if err != nil {
		w.Write(errorToJSONByteArray("could not get sensor", err))
		return
	}
	encoder := json.NewEncoder(w)
	encoder.Encode(sensor)
}

func getSensorFromBody(r *http.Request) (*common.Sensor, error) {
	decoder := json.NewDecoder(r.Body)
	var sensor common.Sensor
	err := decoder.Decode(&sensor)
	if err != nil {
		return nil, err
	}

	for _, reg := range sensor.Registers {
		if !isTypeOk(reg.Type) {
			return nil, fmt.Errorf("register type %s unknown", reg.Type)
		}
	}

	if sensor.ReadGroups != nil {
		for _, rg := range sensor.ReadGroups {
			if rg.ResultType != common.Float32 &&
				rg.ResultType != common.Int32 &&
				rg.ResultType != common.Uint32 {
				return nil, fmt.Errorf("Type %s unknown", rg.ResultType)
			}
		}
	}

	return &sensor, nil

}

func addSensor(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Header().Add("Content-Type", "application/json")
	sensor, err := getSensorFromBody(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errorToJSONByteArray("no valid sensor received", err))
		return
	}
	log.Printf("Adding sensor %d\n", sensor.Address)
	err = configProvider.AddSensor(*sensor)
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		w.Write(errorToJSONByteArray("could not add sensor", err))
		return
	}
	w.WriteHeader(http.StatusCreated)
	returnSuccess(w)
}

func deleteSensor(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var err error
	var sensorAddress int
	w.Header().Add("Content-Type", "application/json")
	sensorString := p.ByName("sensor")
	if sensorAddress, err = strconv.Atoi(sensorString); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errorToJSONByteArray("could not convert to valid sensor address", err))
		return
	}
	log.Printf("Deleting sensor %d\n", sensorAddress)
	err = configProvider.RemoveSensorByAddress(uint8(sensorAddress))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write(errorToJSONByteArray("could not delete sensor", err))
		return
	}
	returnSuccess(w)
}
func changeSensor(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var err error
	var sensorAddress int
	var sensor *common.Sensor
	w.Header().Add("Content-Type", "application/json")
	sensorString := p.ByName("sensor")
	if sensorAddress, err = strconv.Atoi(sensorString); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errorToJSONByteArray("could not convert to valid sensor address", err))
		return
	}
	log.Printf("Changing sensor %d\n", sensorAddress)
	sensor, err = getSensorFromBody(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errorToJSONByteArray("could use request body as sensor", err))
		return
	}

	err = configProvider.ChangeSensor(uint8(sensorAddress), *sensor)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errorToJSONByteArray("could not change sensor", err))
		return
	}

	returnSuccess(w)

}

func isTypeOk(typeString string) bool {
	if typeString != common.Coil &&
		typeString != common.Holding && typeString != common.Input &&
		typeString != common.InputDiscrete {
		return false
	}
	return true
}

func readSensor(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	//TODO -- see if it is not better to put everything into
	//a scheduler that will read automatically and store it in
	//the persistenceprovider
	var sensorAddress int
	var startLocation int
	var length int
	var err error
	sensorString := p.ByName("sensor")
	typeString := p.ByName("type")
	startString := p.ByName("start")
	lengthString := p.ByName("length")

	if !isTypeOk(typeString) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errorToJSONByteArray("could not read sensor", fmt.Errorf("register type %s unknown", typeString)))
		return
	}
	if sensorAddress, err = strconv.Atoi(sensorString); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errorToJSONByteArray("could not convert to valid sensor address", err))
		return
	}
	if startLocation, err = strconv.Atoi(startString); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errorToJSONByteArray("could not convert to valid start location", err))
		return
	}
	if length, err = strconv.Atoi(lengthString); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errorToJSONByteArray("could not convert to valid length", err))
		return
	}
	log.Printf("Reading sensor %d, type %s, start %d, length %d ",
		sensorAddress, typeString, startLocation, length)
	err = scheduleProvider.Read(uint8(sensorAddress), typeString, uint16(startLocation), uint16(length), true, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errorToJSONByteArray("could not read sensor", err))
		return
	}
	encoder := json.NewEncoder(w)
	encoder.Encode(map[string]string{"Status": "OK"})
}

func main() {
	initialize()
	mux := httprouter.New()
	mux.ServeFiles("/static/*filepath", http.Dir("static"))
	mux.GET("/sensors/:sensor", getSensor)
	mux.POST("/sensors", addSensor)
	mux.DELETE("/sensors/:sensor", deleteSensor)
	mux.PUT("/sensors/:sensor", changeSensor)
	mux.GET("/sensors", getSensors)
	mux.POST("/schedule/timers", addTimer)
	mux.DELETE("/schedule/timers/:timer", deleteTimer)
	mux.GET("/schedule/timers", getTimers)
	mux.PUT("/schedule/save", saveSchedule)
	mux.PUT("/schedule/load", loadSchedule)

	mux.GET("/read/:sensor/:type/:start/:length", readSensor)

	server := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: mux,
	}
	server.ListenAndServe()
}
