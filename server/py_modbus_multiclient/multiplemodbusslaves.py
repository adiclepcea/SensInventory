import yaml
import sys
import serial
import serial.rs485
import modbusreader
import Queue
import random
from time import sleep
import os

VALUE_RANDOM = "random_generated"
VALUE_FIXED = "fixed"
VALUE_READ = "read"

TYPE_HOLDING = "holding"
TYPE_INPUT = "input"
TYPE_COIL = "coil"
TYPE_INPUT_DISCRETE = "input_discrete"

dtypes = {}
dtypes[1]=TYPE_COIL
dtypes[2]=TYPE_INPUT
dtypes[3]=TYPE_HOLDING
dtypes[4]=TYPE_INPUT_DISCRETE

class Connection:
    def __init__(self,speed, port):
        self.speed = speed
        self.port = port
class Registry:
    def __init__(self,dict):
        self.length = dict["length"]
        self.location = dict["location"]
        self.max = dict["max"]
        self.min = dict["min"]
        self.type = dict["type"]
        self.value = dict["value"]

class NoSuchRegistryError(Exception):
    def __init__(self, address,location,reg_type):
        self.address = address
        self.location = location
        self.type = reg_type
    def __str__(self):
        return repr("There is no matching registry for slave {} at location {} with type {}".format(self.address,self.location,self.type))

class Slave:
    def __init__(self,dict):
        self.address = dict["address"]
        self.description = dict["description"]
        self.registries = []
        for r in dict["registries"]:
            self.registries.append(Registry(r))

    def generateRandomValue(self, min, max):
        return random.randint(min,max)

    def writeRegistry(self, package):
        #a successful write single registry should return the initial package.
        #for now we do nothing with this, so we consider it a success
        return package


    def respondToRequest(self,package):
        if package[0]!=self.address:
            print "Invalid address for me. I'm {} and request address is {}".format(self.address,package[0])
            return None
        request_type = 0
        try:
            request_type = package[1]
            if request_type == 0x06 or request_type==0x05:
                return writeRegistry(package)
            location_to_read = package[2] * 256 + package[3]
            length_to_read = package[4] * 256 + package[5]
            values = self.askValues(dtypes[request_type],location_to_read,length_to_read)
            package = []
            package.append(self.address)
            package.append(request_type)
            print package
            bytes_to_follow = 0
            if request_type == 1 or request_type==2: #for coils we need to calculate the value to send
                bytes_to_follow = len(values)/8
                if len(values)%8>0:
                    bytes_to_follow+=1
                package.append(bytes_to_follow)
                for i in xrange(0,len(values),8):
                    r = 0
                    for j in range(0,min(8,len(values)-i)):
                        r = (r<<1)+(1 if values[i+j] else 0)
                    package.append(r)
            elif request_type == 3 or request_type == 4: #for registers and inputs we put 2 bytes for register
                bytes_to_follow = len(values)*2
                package.append(bytes_to_follow)
                for v in values:
                    package.append(v/256)
                    package.append(v%256)

            return package
        except NoSuchRegistryError as ex:
            print ex
            package = []
            package.append(self.address)
            package.append(request_type | 0b10000000)
            package.append(0x02) #error code for Illegal Data Address
            return package


    #this only works if the registries are passed in ascending order in config
    def askValues(self,reg_type,location,length):
        loc = location
        values = []
        #our maximum location was surpassed - raise an error
        if location+length>=65535 or location+length>self.registries[len(self.registries)-1].location+1:
            raise NoSuchRegistryError(self.address,loc,reg_type)
        #is there a registry here
        if location<self.registries[0].location:
            raise NoSuchRegistryError(self.address,loc,reg_type)
        #search for first location
        start = 0
        for reg in self.registries:
            if reg.location == location:
                break
            else:
                start = start+1
        #identify the values
        for i in range(start,len(self.registries)):
            reg = self.registries[i]
            if reg.location == loc:
                if reg.type == reg_type:
                    #we add the corresponding value
                    #todo - move this in corresponding methods as below todos describe
                    if reg.value==VALUE_FIXED:
                        values.append(reg.max)
                    elif reg.value==VALUE_RANDOM:
                        values.append(self.generateRandomValue(reg.min, reg.max))
                    elif reg.value==VALUE_READ:
                        values.append(11)   #todo - change to read the value from somewhere
                else:
                    #if the register is not the valid type, we raise an error
                    raise NoSuchRegistryError(self.address,loc,reg_type)
            elif len(values)>0:
                raise NoSuchRegistryError(self.address,loc,reg_type)

            loc = loc+1
            if loc == location+length and len(values)>0 :
                return values
        raise NoSuchRegistryError(self.address,loc,reg_type)


class Configuration:
    def __init__(self, configFileName):
        yamlConfig = yaml.load(file(configFileName))
        self.load(yamlConfig["connection"],yamlConfig["slaves"])

    def load(self,connection,slaves):
        self.connection = Connection(connection[0]["speed"],connection[0]["port"])
        self.slaves = []
        for s in slaves:
            self.slaves.append(Slave(s))


if __name__ == '__main__':
    config = Configuration("config.yaml")
    sys.stdout.write('ModBus bus on port {}, speed={}.\nActing as {} slaves.\n'.format(config.connection.port,config.connection.speed,len(config.slaves)))
    sys.stdout.write("Press Ctrl+C to stop\n")

    try:
        ser = serial.rs485.RS485(config.connection.port,baudrate=config.connection.speed,timeout=10)
        ser.rs485_mode=serial.rs485.RS485Settings(delay_before_tx=0.2)
    except:
        logging.exception("error openning serial port")

    q = Queue.Queue()
    mr = modbusreader.ModBusReader(config.slaves,q)
    while 1:
        if ser.inWaiting()>0:
            s = ser.read(ser.inWaiting())
            for c in s:
                q.put(ord(c))
            ser.flushInput()
            ser.flushOutput()
            resp = mr.checkForFullPackage()
            if resp==modbusreader.PACKAGE_OK:
                resp = ''.join(chr(x) for x in mr.response)
                #print " ".join(map(lambda x:x.encode('hex'),resp))
                ser.write(resp)
                ser.flushOutput()

        sleep(0.01)


    ser.close()
