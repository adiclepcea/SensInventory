#The protocol used to communicate between the application and the server

###Concept

The server should be accesible from no mather what application as long as the application uses the described protocol.

###The base protocol

The protocol is REST based. That being said, it means the information can be accessed by anyone through normal HTTP requests.

###Details

The data format used to communicate is json.

We have the following available operations on the server side (the server has the address *server* and the sensor has the address *10*:

* Read the last value from a sensor:
 * GET request to http://server/sensor/10
 * json response example: 

 ```
 	{
 		"readings": [{
			"address": "10",
			"time": "2016-04-03 12:00:00",
			"value1": "33.21"
		}]
	}
 ```

* Read the list of values from a sensor in a certain frame ie. for reading starting from 2016-04-03 12:00:10 up to 2016-04-04 11:30:00: 
 * GET request to http://server/sensor/address=10&start=20160403_120010&end=20160404_113000
 * json response example:

 ```
	{
		"readings": [{
			"address": "10",
			"time": "2016:04:03 12:00:00",
			"value1": "33.21"
		}, {
			"address": "10",
			"time": "2016:04:03 12:00:10",
			"value1": "33.21"
		}, 

		...
		
		{
			"address": "10",
			"time": "2016:04:04 11:30:00",
			"value1": "30.21"
		}]
	}
 ```

* Write the time to a sensor (the time written to the sensor will be the time that is set on ther server when the request was made) ie. for setting the time to 2016-04-03 12:00:00: 
 * PUT request to http://server
 * json contents of the message:

 ```
 	{
		"settings": [{
			"address": "10",
			"setting": "time"
		}]
	}
 ```

* Assign an address to a new sensor (add sensor with address 10 here):
 * POST request to http://server
 * json contents of the message:

 ```        
 	{
		"address": "10"
	}
	
 ```

* Remove a sensor from reading:
 * DELETE request to http://server
 * json contents of the message:

 ```        
 	{
		"address": "10"
	}
	
 ```

### Future

* We could also provide a possibility to ask for several sensor values. Either the last ones read or the values read in a time interval.
* We could add the possibility to set a trigger. Thus, when a certain limit will be reached, the server could either send a mail, a SMS or access a web address with certain parameters
* Ideas?

