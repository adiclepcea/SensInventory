package server

import "github.com/adiclepcea/SensInventory/server/common"

//ReadingProvider provides an interface(blueprint) for a provider
//that reads sensors
type ReadingProvider interface {
	NewReadingProvider(*ConfigProvider) ReadingProvider
	GetReading(int) (*common.Reading, *error)
}
