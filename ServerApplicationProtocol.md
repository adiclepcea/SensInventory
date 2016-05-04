#The protocol used to communicate between the application and the server

###Concept

The server should be accesible from no mather what application as long as the application uses the described protocol.

###The base protocol

The protocol is REST based. That being said, it means the information can be accessed by anyone through normal HTTP requests.

###Details

The data format used to communicate is json.

We have the following available operations on the server side (the server has the address *server* and the sensor has the address *10*:

#### Read the last value from a sensor:
* GET request to http://server/sensor/10 (curl -i http://localhost:8082/sensor/10)
* json response example: 
 
```
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Wed, 04 May 2016 06:54:29 GMT
Content-Length: 199
{
	"readSensor": {
		"address": 10,
		"description": "Sensor 10 - Weight sensor",
		"configuredValues": [{
			"name": "test ReadValue",
			"registerAddress": 100,
			"registerLength": 2
		}]
	},
	"readValues": [0.8980626],
	"time": "2016-05-04T06:54:29.158158472Z"
}
```

#### Add a new sensor for reading:
* POST request to http://server (in the example below we add sensor 100 - the valid values are 1 - 32 and 100 -, the value 100 could be used for any new sensor. 
```
curl -H "Content-Type: application/json" -X POST -i http://localhost:8082/sensor -d "{\"address\":100,\"description\":\"New Sensor\",\"configuredValues\":[{\"name\":\"test ReadValue\",\"registerAddress\":80,\"registerLength\":2},{\"name\":\"byte read value\",\"registerAddress\":82,\"registerLength\":1}]}"
```
* Response when OK:
```
HTTP/1.1 201 Created
Content-Type: application/json; charset=utf-8
Date: Wed, 04 May 2016 07:07:41 GMT
Content-Length: 5

null
```

* Notice that the server might also respons with an error status when the sensor cannot be added (i.e. there might be allready a sensor with this address or the json might be malformed etc.)

####Change a sensor
* PUT request to http://server (in the example below we change the sensor added above. Notice that the address 100 will be moved to address 3 and also the name of the first read value will be changed to "First read value"
```
curl -H "Content-Type: application/json" -X PUT -i http://localhost:8082/sensor/3 -d "{\"address\":100,\"description\":\"New Sensor\",\"configuredValues\":[{\"name\":\"First read value\",\"registerAddress\":80,\"registerLength\":2},{\"name\":\"byte read value\",\"registerAddress\":82,\"registerLength\":1}]}"
```

* Response from server when OK:
```
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Wed, 04 May 2016 07:13:33 GMT
Content-Length: 5

null
```
* Notice that the server might also respond with error status when there are problems while changing the sensor

#### Remove a sensor from reading:

* DELETE request to http://server:
```
curl -X DELETE -i http://localhost:8082/sensor/100
```

* Response when OK:
```        
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Wed, 04 May 2016 07:16:56 GMT
Content-Length: 5

null	
```

### Future

* We could also provide a possibility to ask for several sensor values. Either the last ones read or the values read in a time interval.
* We could add the possibility to set a trigger. Thus, when a certain limit will be reached, the server could either send a mail, a SMS or access a web address with certain parameters
* Ideas?

