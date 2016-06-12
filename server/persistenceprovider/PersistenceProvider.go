package persistenceprovider

import (
	"time"

	"github.com/adiclepcea/SensInventory/server/common"
)

//PersistenceProvider is the prototype to be
//implemented when creting a persistence provider
type PersistenceProvider interface {
	SaveSensorReading(common.Reading)
	GetSensorReading(uint8, time.Time)
	GetSensorReadingsInPeriod(uint8, time.Time, time.Time)
	GetAllReadingsInPeriod(time.Time, time.Time)
	DeleteSensorReading(uint8, time.Time)
	DeleteSensorReadingsInPeriod(uint8, time.Time, time.Time)
	DeleteAllReadingsInPeriod(time.Time, time.Time)
}
