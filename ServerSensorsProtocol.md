#Description of the communication protocol between server and sensors.

###Method of communication

The protocol used for communication is ModBus RTU. This means that the communication will take place through the serial port. 

<<<<<<< HEAD
The wire protocol used will be 485.

###Details of the read data

=======
The wire protocol used will be 485 over 2 wires.

###Details of the read data

####Ideas

1. We will have a maximum of 32 sensors to read
2. Each sensor can have an address in range 1 to 32 when functioning. 
3. Each new sensor should have an initial address of 100. This way it cannot be read until the server changes the address of this sensor and thus knows about it.
4. Each sensor must be capable of changing its address when asked by the server.
5. Each sensor should be capable of responding to the server with a read value and a time of reading in the same package.
6. Each sensor can have (but is optional) the possibility to have a limit set. This limit can be used (optionally) to emit an alarm.
7. In the eventuality the sensor has a limit to set. This will have a default value that can be changed by the server.
8. Each sensor can have visual signling means (LCD, Leds etc.). These can signal when a value is read, when a communication with the server takes place etc.

####Consequences

1. The sensors should have at least one Holding Register (for the address)
2. The sensors should have at least one Input Register (for the read value - weight in our case).
3. The sensors should have a RTC. It is possible to also have another Input Resgister that will hold the time at which the value was read, but the sensor can also send the time at which it responds to server in the response.

####Packages

**Codes used**

The codes used in the packages will be:

* 03 - Read Holding Register - used to read the address of one sensor (optional)
* 04 - Read Input Register - used for reading values from the sensor
* 06 - Write Single Register - used for writing an address to the sensor and possibly to set the time
* 07 - Read exception status - used to get possible errors from the sensor
* 08 - Disgnostic (optional)

**Details**
*03 - Read Holding Register*
The sensor will store the address at register 10

*04 - Read Input Register*
The sensor will store the read values in 2 registers. The first register will have address 100 and will be the integer part of the value and the seccond register will have address 101 and will hold the decimal value (two decimals).
For example for a value of 23.56 kgs, register 100 will hold the value 23 (0x17) and register 101 will hold the value 56 (0x38).

The sensor will store the time in the following format:
* Register 102 - year
* Register 103 - month
* Register 104 - day
* Register 105 - hour
* Register 106 - minute
* Register 107 - second

When reading input register the server will normally ask for 8 registers (100 - 107) because the time a value was read is also needed.

So the package from the server would contain for address ADDR:
ADDR 04 64 00 08 CRC

The response would contain for a value of 23.56 read in 2016-04-03 12:10:50 :

ADDR 04 10 00 17 00 38 07 E0 00 04 00 03 00 0C 00 0A 00 32 CRC 

(TODO - verify please)

>>>>>>> 17640ce24b1c80fbb64e351d9a6eee860313ec85

