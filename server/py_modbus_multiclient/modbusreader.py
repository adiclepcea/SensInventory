import Queue
import threading
from time import sleep

NO_FULL_PACKAGE_AVAILABLE = 1
NO_SLAVE_MATCH = 2
INCORRECT_CRC = 3
PACKAGE_OK = 4

#ModBusReader reads the q Queue and
#passes the requests to the slaves it finds in the requests
#if they are betwees the passed "slaves"
class ModBusReader:
    def __init__(self, slaves, q):
        self.slaves = slaves
        for s in slaves:
            print "listening on behalf of slave:", s.address
        self.q = q


    def crc16(self, reqIn):
        crc = 0xFFFF
        l = len(reqIn)
        i = 0
        while i < l:
            j = 0
            #crc = crc ^ ord(reqIn[i])
            crc = crc ^ reqIn[i]
            while j < 8:
                if (crc & 0x1):
                    mask = 0xA001
                else:
                    mask = 0x00
                crc = ((crc >> 1) & 0x7FFF) ^ mask
                j += 1
            i += 1
        if crc < 0:
            crc -= 256
        #crc% 256* ; crc / 256 -will provide the pachet order
        return crc

    #this will only work if the first byte in the Queue is the first byte of a
    #packet (i.e. the queue does not start with the middle of the packet)
    #the requests 15 (Force Multiple Coils) and
    #16(Preset Multiple Registers) are not yet supported
    #the function eliminates every 0 from the beginning of the queue, but nothing more
    def checkForFullPackage(self):
        while self.q.qsize()>=8 and self.q.queue[0]==0:
            self.q.get()
        if self.q.qsize()>=8:
            self.package = []
            for i in range(0,8):
                self.package.append(self.q.get())
            #################################################################
            ####    force multiple coils and preset multiple registers   ####
            ####        are not tested yet, as I have no use of them     ####
            #################################################################
            if self.package[1] == 15: #force multiple coils
                number_of_coils_to_write = self.package[4]*256+self.package[5]
                plus_read = number_of_coils_to_write/8
                if number_of_coils_to_write%8>0:
                    plus_read+=1
                    #package[6] should be now equal to plus_read
                for i in range(0,plus_read+1):
                    package.append(self.q.get()) #this will block waiting for data
            if self.package[1] == 16: #preset multiple registers
                plus_read =  self.package[4]*256+self.package[5]
                #package[6] should be now equal to plus_read
                for i in range(0,plus_read+1):
                    package.append(self.q.get()) #this will block waiting for data
            ##################################################################
            ####                    end of untested code                  ####
            ##################################################################
            for slave in self.slaves:
                if slave.address == self.package[0]:
                    self.requested_slave = slave
                    crc = self.crc16(self.package[:6])
                    if crc%256==self.package[6] and crc/256==self.package[7]:
                        self.response = self.createResponse(slave)
                        return PACKAGE_OK
                    else:
                        return INCORRECT_CRC
            print "i do not know slave:", self.package[0]
            return NO_SLAVE_MATCH
        else:
            return NO_FULL_PACKAGE_AVAILABLE

    def createResponse(self, sl):
        resp = sl.respondToRequest(self.package)
        if resp[1] == 5 or resp[1] == 6:
            return resp
        crc = self.crc16(resp)
        resp.append(crc%256)
        resp.append(crc/256)
        return resp

    def run(self):
        print "Started reading values"
