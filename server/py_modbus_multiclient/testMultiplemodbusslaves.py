import unittest
import multiplemodbusslaves

class TestSlave(unittest.TestCase):
    def init_slave_for_test(self,value):
        dReg = {}
        dReg["length"] = 1
        dReg["location"] = 10
        dReg["max"] = 122
        dReg["min"] = 0
        dReg["type"] = multiplemodbusslaves.TYPE_HOLDING
        dReg["value"] = value

        dSlave = {}
        dSlave["address"] = 10
        dSlave["description"] = "test_slave"
        dSlave["registries"] = []
        dSlave["registries"].append(dReg)
        slave = multiplemodbusslaves.Slave(dSlave)
        return slave

    def test_slave_build_from_dict(self):
        slave = self.init_slave_for_test(multiplemodbusslaves.VALUE_FIXED)
        self.assertEqual(slave.address,10)
        self.assertEqual(len(slave.registries),1)
        self.assertEqual(slave.registries[0].length,1)
        self.assertEqual(slave.registries[0].location,10)
        self.assertEqual(slave.registries[0].max,122)
        self.assertEqual(slave.registries[0].min,0)
        self.assertEqual(slave.registries[0].type,multiplemodbusslaves.TYPE_HOLDING)
        self.assertEqual(slave.registries[0].value,multiplemodbusslaves.VALUE_FIXED)

    def test_slave_askValues(self):
        slave = self.init_slave_for_test(multiplemodbusslaves.VALUE_FIXED)
        #length is different - should fail
        with self.assertRaises(multiplemodbusslaves.NoSuchRegistryError):
            slave.askValues(multiplemodbusslaves.TYPE_HOLDING,10,2)
        #location is different - should fail
        with self.assertRaises(multiplemodbusslaves.NoSuchRegistryError):
            slave.askValues(multiplemodbusslaves.TYPE_HOLDING,101,1)
        #type is different - should fail
        with self.assertRaises(multiplemodbusslaves.NoSuchRegistryError):
            slave.askValues(multiplemodbusslaves.TYPE_COIL,10,1)
        #everything is ok. should succed
        self.assertEqual(slave.askValues(multiplemodbusslaves.TYPE_HOLDING,10,1),[122])

        #we add another register with a different type at the next location
        dReg = {}
        dReg["length"] = 1
        dReg["location"] = 11
        dReg["max"] = 1
        dReg["min"] = 1
        dReg["type"] = multiplemodbusslaves.TYPE_COIL
        dReg["value"] = multiplemodbusslaves.VALUE_FIXED
        reg = multiplemodbusslaves.Registry(dReg)
        slave.registries.append(reg)
        #now we should get an error as the registries are not both the same type as the asked one
        with self.assertRaises(multiplemodbusslaves.NoSuchRegistryError):
            slave.askValues(multiplemodbusslaves.TYPE_HOLDING,10,2)

        #should be ok for the seccond registry only
        self.assertEqual(slave.askValues(multiplemodbusslaves.TYPE_COIL,11,1),[1])

        #now we fix that
        dReg["type"] = multiplemodbusslaves.TYPE_HOLDING
        reg = multiplemodbusslaves.Registry(dReg)
        slave.registries[1]=reg
        #everything should be ok
        self.assertEqual(slave.askValues(multiplemodbusslaves.TYPE_HOLDING,10,2),[122,1])

    def test_slave_random_values(self):
        for i in range(1,100):
            slave = self.init_slave_for_test(multiplemodbusslaves.VALUE_FIXED)
            slave.registries[0].value = multiplemodbusslaves.VALUE_RANDOM
            slave.registries[0].min = 1
            slave.registries[0].max = 2
            self.assertIn(slave.askValues(multiplemodbusslaves.TYPE_HOLDING,10,1),[[1],[2]])


if __name__ == '__main__':
    unittest.main()
