#Server SensInventory

##Server for SensInventory system

###Description

It reads a modbus serial bus. Saves the data ito a database and it serves the data
through  simple REST Api.

The entire configuration can be made throgh the REST Api.

You have the following possibilities when you work with sensors
(the examples use curl and use a server based on loalhost, port 8080):

* See the configuration of a sensor (here we read the configuration of the sensor with address 10):
```
curl http://localhost:8080/sensors/10
```

* Add the configuration for a new sensor to read:
```
curl -H "Content-Type: application/json" -X POST http://localhost:8080/sensors -d '{"address":10,"description":"Test","registers":[{"name":"test ReadValue","location":100,"type":"input"}],"readGroups":null}'
```

* Delete an existing sensor (this will delete a sensor with the address 10):
```
curl -X DELETE http://localhost:8080/sensors/10
```

* Change the configuration for a certain sensor:
```
curl -H "Content-Type: application/json" -X PUT http://localhost:8080/sensors/10 -d '{"address":10,"description":"Test for fun","registers":[{"name":" good known value","location":90,"type":"input"}],"readGroups":null}'
```
this will change the sensor with address 10 we added before to have a different description,
a different name for a read value and a different location (from 100 to 90) of the read register.

* Get a list of all the read sensors:
```
curl http://localhost:8080/sensors
```

As you modify a sensor, it will be read from the modbus bus.
The sensors are read based on a schedule. You can add or delete schedules as you wish.

To work with the scheduler, you use timers:

* Add a timer:
```
curl -H "Content-Type: application/json" -X POST http://localhost:8080/timers -d '{"sensorAddress":1,"readType":"holding","startLocation":11,"readLength":1,"interval":4000000000,"repeat":true,"store":true}'
```
This will add a timer that will read the sensor with address 1, from location 11, 1 register every 4 secconds.
The result will be stored.

* Get a list with all the registered timers:
```
curl http://localhost:8080/timers
```

* Delete a timer:
```
curl -X DELETE http://localhost:8080/timers/1
```
This will delete the timer with id 1. The ids can be retrieved by using the GET request above

* You can also save and load a schedule
```
curl -X PUT http://localhost:8080/schedule/save
```
and
```
curl -X PUT http://localhost:8080/schedule/load
```

The ideea is that you can add several timers and they will run for as long as the application runs
When you are happy with the way the application reads the data, you can save the schedule
containing all the timers you have configured.

That way, if you stop the application (or some external factor like a power shortage intervenes),
you can "load" the saved scheduler and it works with the saved timers.

Notice that the sensors are saved automatically, you do not need to save or load them.

An example application is found in main.go

### Observations

* You can also use the scheduler to read a certain sensor without the need of having a timer setup.
* There are interfaces allowing easy implementations of other kind of persistence or configurations.
  For now the persistence is done in CouchDB and the configurations (what values the sensors have defined) are simple json files.
