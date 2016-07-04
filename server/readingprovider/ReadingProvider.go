package readingprovider

import (
	"github.com/adiclepcea/SensInventory/server/common"
	"github.com/adiclepcea/SensInventory/server/configprovider"
)

//ReadingProvider provides an interface(blueprint) for a provider
//that reads sensors
type ReadingProvider interface {
	NewReadingProvider(*configprovider.ConfigProvider) ReadingProvider
	GetReading(uint8, string, uint16, uint16) (*common.Reading, error)
}
