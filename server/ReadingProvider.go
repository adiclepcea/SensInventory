package server

import (
	"github.com/adiclepcea/SensInventory/server/common"
	"github.com/adiclepcea/SensInventory/server/configprovider"
)

//ReadingProvider provides an interface(blueprint) for a provider
//that reads sensors
type ReadingProvider interface {
	NewReadingProvider(*configprovider.ConfigProvider) ReadingProvider
	GetReading(int) (*common.Reading, *error)
}
