package persistenceprovider

import (
	"time"

	"github.com/adiclepcea/SensInventory/server/common"
)

//PersistenceProvider is the prototype to be
//implemented when creating a persistence provider
type PersistenceProvider interface {
	NewPersistenceProvider(params ...string) (PersistenceProvider, error)
	SaveSensorReading(common.Reading) error
	GetSensorReading(uint8, time.Time) (*common.Reading, error)
	GetSensorReadingsInPeriod(uint8, time.Time, time.Time) ([]common.Reading, error)
	GetSensorReadingCountInPeriod(uint8, time.Time, time.Time) (uint, error)
	GetAllReadingsInPeriod(time.Time, time.Time) (*[]common.Reading, error)
	GetAllReadingsCountInPeriod(time.Time, time.Time) (uint, error)
	DeleteSensorReading(uint8, time.Time) error
	DeleteSensorReadingsInPeriod(uint8, time.Time, time.Time) error
	DeleteAllReadingsInPeriod(time.Time, time.Time) error
	SaveItem(string, interface{}) error
	ReadItem(string) (interface{}, error)
}
