#Modbus client

##Description

The purpose of this solution is to implement configurable clients for modbus.

The configuration is stored in [yaml](http://yaml.org) format.

The library used for reading the configuration is [libyaml](http://pyyaml.org/wiki/LibYAML). You can find the MIT Licence in this repository.

The ideea is that you could use this software to simulate an entire bus if you want. It is supposed to connect to a serial port and respond to one or several addresses. For now it is still very much a work in progress.

