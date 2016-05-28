package server

import "github.com/adiclepcea/SensInventory/server/common"

//IReadingProvider provides an interface(blueprint) for a provider
//that reads sensors
type IReadingProvider interface {
	NewReadingProvider(*ConfigProvider) IReadingProvider
	GetReading(int) (*common.Reading, *error)
}
