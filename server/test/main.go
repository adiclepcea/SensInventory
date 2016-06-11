package main

import (
	"fmt"
	"strconv"

	"github.com/adiclepcea/SensInventory/server"
	"github.com/adiclepcea/SensInventory/server/common"
	"github.com/adiclepcea/SensInventory/server/configprovider"
	"github.com/gin-gonic/gin"
)

func main() {
	conf := configprovider.MockConfigProvider{}.NewConfigProvider()
	conf.SetAddressLimits(1, 32)
	sensor1 := common.Sensor{Address: 1, Description: "", Registers: []common.Register{common.Register{Name: "test ReadValue",
		Location: 100, Type: common.Holding}}}
	sensor2 := common.Sensor{Address: 2, Description: "", Registers: []common.Register{common.Register{Name: "Sensor 2 Value 1",
		Location: 100, Type: common.Holding}, common.Register{Name: "Sensor 2 value 2",
		Location: 102, Type: common.Holding}}}

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
		if _, err := conf.GetSensorByAddress(uint8(nAddress)); err != nil {
			c.JSON(404, gin.H{"error": "Sensor address not registered"})
			return
		}

		readValues, e := mockServer.GetReading(uint8(nAddress))
		if e != nil {
			c.JSON(500, gin.H{"error": (e).Error()})
		} else {
			c.JSON(200, readValues)
		}
		fmt.Println(readValues)
	})

	r.POST("/sensor", func(c *gin.Context) {
		var sensor common.Sensor
		if c.BindJSON(&sensor) == nil {

			err := conf.AddSensor(sensor)
			if err != nil {
				c.JSON(400, gin.H{"error": (err).Error()})
				return
			}
			c.JSON(201, nil)
		} else {
			c.JSON(400, gin.H{"error": "Invalid JSON"})
		}
	})

	r.PUT("/sensor/:address", func(c *gin.Context) {
		address := c.Params.ByName("address")
		nAddress, err := strconv.Atoi(address)
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid address"})
			return
		}

		var sensor common.Sensor

		if c.BindJSON(&sensor) == nil {
			if _, err := conf.GetSensorByAddress(sensor.Address); err != nil {
				c.JSON(404, gin.H{"error": "Sensor address not registered"})
				return
			}
			if sensor.Address != uint8(nAddress) {
				if err := conf.ChangeSensorAddress(sensor.Address, uint8(nAddress)); err != nil {
					c.JSON(500, gin.H{"error": (err).Error()})
					return
				}
			}
			err := conf.ChangeSensor(uint8(nAddress), sensor)
			if err != nil {
				c.JSON(500, gin.H{"error": (err).Error()})
				return
			}
			c.JSON(200, nil)
		} else {
			c.JSON(400, gin.H{"error": "Invalid JSON"})
		}

	})

	r.DELETE("/sensor/:address", func(c *gin.Context) {
		address := c.Params.ByName("address")
		nAddress, err := strconv.Atoi(address)
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid address"})
			return
		}
		if _, err := conf.GetSensorByAddress(uint8(nAddress)); err != nil {
			c.JSON(404, gin.H{"error": "Sensor address not registered"})
			return
		}

		if err := conf.RemoveSensorByAddress(uint8(nAddress)); err != nil {
			c.JSON(500, gin.H{"error": (err).Error()})
			return
		}
		c.JSON(200, nil)
	})

	r.Run(":8081")
}
