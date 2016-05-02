package main

import (
	"fmt"
	"strconv"

	"github.com/adiclepcea/SensInventory/server"
	"github.com/adiclepcea/SensInventory/server/common"
	"github.com/gin-gonic/gin"
)

func main() {
	conf := server.ConfigProvider{}.NewConfigProvider()
	sensor1 := common.Sensor{Address: 1, Description: "", ConfiguredValues: []common.ConfiguredValue{common.ConfiguredValue{Name: "test ReadValue",
		RegisterAddress: 100, RegisterLength: 2}}}
	sensor2 := common.Sensor{Address: 2, Description: "", ConfiguredValues: []common.ConfiguredValue{common.ConfiguredValue{Name: "Sensor 2 Value 1",
		RegisterAddress: 100, RegisterLength: 2}, common.ConfiguredValue{Name: "Sensor 2 value 2",
		RegisterAddress: 102, RegisterLength: 1}}}

	conf.AddSensor(sensor1)
	conf.AddSensor(sensor2)

	mockServer := server.MockReadingProvider{}.NewReadingProvider(conf)

	r := gin.Default()

	r.GET("/sensor/:address", func(c *gin.Context) {
		address := c.Params.ByName("address")
		nAddress, err := strconv.Atoi(address)
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid address"})
			return
		}
		if _, err := conf.GetSensorByAddress(nAddress); err != nil {
			c.JSON(404, gin.H{"error": "Sensor address not registered"})
			return
		}

		readValues, e := mockServer.GetReading(nAddress)
		if e != nil {
			c.JSON(500, gin.H{"error": (*e).Error()})
		} else {
			c.JSON(200, readValues)
		}
		fmt.Println(readValues)
	})

	r.Run(":8081")
	/*
		readValues, err := mockServer.GetReading(1)

		if err != nil {
			fmt.Println((*err).Error())
			return
		}

		fmt.Println(readValues)

		readValues, err = mockServer.GetReading(2)

		if err != nil {
			fmt.Println((*err).Error())
			return
		}
		fmt.Println(readValues)
	*/
}
