package readingprovider

import (
	"fmt"
	"math"
	"time"

	"github.com/adiclepcea/SensInventory/server/common"
	"github.com/adiclepcea/SensInventory/server/configprovider"
	"github.com/goburrow/modbus"
)

//SerialConfig is the struct holding the configuration for the serial port
type SerialConfig struct {
	Port     string
	BaudRate int
	DataBits int
	Parity   string
	StopBits int
	Timeout  time.Duration
}

//ModBUSReadingProvider is the type used for reading a modbus bus
type ModBUSReadingProvider struct {
	ConfigProvider *configprovider.ConfigProvider
	serialConfig   *SerialConfig
	handler        *modbus.RTUClientHandler
	ReadingProvider
}

//NewReadingProvider is the function that builds a new ModBUSReadingProvider
func (modbusProvider ModBUSReadingProvider) NewReadingProvider(configProvider *configprovider.ConfigProvider) ReadingProvider {
	modbus := ModBUSReadingProvider{ConfigProvider: configProvider}
	modbus.initialize()
	return &modbus
}

func getUInt16FromBytes(input []byte) []uint16 {
	var rez []uint16
	var temp1, temp2 uint16
	for i := 0; i < len(input); i += 2 {
		temp1 = uint16(input[i])
		if len(input) >= i+2 {
			temp2 = uint16(input[i+1])
		} else {
			temp2 = 0
		}
		temp1 = (temp1 << 8) | temp2
		rez = append(rez, temp1)
	}
	return rez
}

func getBitFromBytes(input []byte, length int) []uint16 {
	var rez []uint16
	for _, byt := range input {
		fmt.Println("Using", byt)
		for i := 0; i < 8; i++ {
			if byte(math.Pow(2, float64(i)))&byt > 0 {
				rez = append(rez, uint16(1))
			} else {
				rez = append(rez, uint16(0))
			}
			if len(rez) >= length {
				break
			}
		}
	}
	return rez
}

//GetReading is the function that will read a sensor and return the reading
func (modbusProvider *ModBUSReadingProvider) GetReading(sensor uint8, registerType string, startLocation uint16, length uint16) (*common.Reading, error) {
	modbusProvider.handler.SlaveId = sensor
	var results []byte
	var results16 []uint16
	reading := common.Reading{}
	err := modbusProvider.handler.Connect()
	if err != nil {
		return nil, err
	}
	defer modbusProvider.handler.Close()
	client := modbus.NewClient(modbusProvider.handler)
	switch registerType {
	case common.Coil:
		results, err = client.ReadCoils(startLocation, length)
		if err != nil {
			return nil, err
		}
		results16 = getBitFromBytes(results, int(length))
		break
	case common.Holding:
		results, err = client.ReadHoldingRegisters(startLocation, length)
		if err != nil {
			return nil, err
		}
		results16 = getUInt16FromBytes(results)
		break
	case common.Input:
		results, err = client.ReadInputRegisters(startLocation, length)
		if err != nil {
			return nil, err
		}
		results16 = getUInt16FromBytes(results)
		break
	case common.InputDiscrete:
		results, err = client.ReadDiscreteInputs(startLocation, length)
		if err != nil {
			return nil, err
		}
		results16 = getBitFromBytes(results, int(length))
		break
	default:
		return nil, fmt.Errorf("Register type %s not supported", registerType)
	}

	reading.Sensor = sensor
	reading.StartLocation = startLocation
	reading.Count = length
	reading.Time = time.Now().Format(common.TimeFormat)
	reading.ReadValues = results16
	return &reading, nil
}

func (modbusProvider *ModBUSReadingProvider) initialize() {
	//TODO this should not be reached
	//find a method to configure this
	//perhaps in the configProvider?
	if modbusProvider.serialConfig == nil {
		modbusProvider.serialConfig = &SerialConfig{Port: "/dev/ttyUSB1",
			BaudRate: 115200, DataBits: 8, Parity: "N", StopBits: 1,
			Timeout: 5 * time.Second}
	}
	if modbusProvider.handler == nil {
		modbusProvider.handler = modbus.NewRTUClientHandler(modbusProvider.serialConfig.Port)
		modbusProvider.handler.StopBits = modbusProvider.serialConfig.StopBits
		modbusProvider.handler.BaudRate = modbusProvider.serialConfig.BaudRate
		modbusProvider.handler.DataBits = modbusProvider.serialConfig.DataBits
		modbusProvider.handler.Parity = modbusProvider.serialConfig.Parity
		modbusProvider.handler.Timeout = modbusProvider.serialConfig.Timeout
	}
}
