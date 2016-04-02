# SensInventory
##System for keeping imventory based on sensors

###What it is meant to be

This system should read some sensors, for now it is all about weight sensors, and use the data for calculating the stock level for certain products.

*Example 1:*
You are responsible with delivering toners to your company (ok - you are the IT). Normally you have a minimum stock for each kind of toner you are using (let's say that for toner 38A you must have at least 3 toners. That is because the suppliers takes 3 weeks to deliver new ones and you normally need 1 per week). 
Using this system you could let the system tell you when you reached the limit.

Perhaps you could eventually even make an automatic order to you supplier.

*Example 2:*
You have a bar. You serve beer from barrels. You know that the empty barrel weights 10 Kgs and you can let the system tell you when you only have about 10 litres of beer left in the barrel so that you can order some more.

*Example 3:*
You are a bee keeper. You could use the system to tell you when the hives are too heavy, meaning that they are full of honey. Yummy.


###How do we see it

**The system has an electronic part:**

* Scale (weight sensor)
* Arduino based electronics
 * Board reading the scale
 * Communication part to send the signal to a central server


**A server part:**

The server should read data from several arduinos (*see electronic part*), save it in a database for future use, serve it using a REST based protocol.


**An application part:**
The applications can take several forms. We will only provide a reference implementation.
The applications should read the data from the sensor and provide further functionnality like:

* Show stocks level
* Add new stock parts
* Send notification on low levels

