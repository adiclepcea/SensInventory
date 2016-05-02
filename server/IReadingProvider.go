package server

import "github.com/adiclepcea/SensInventory/server/common"

type IReadingProvider interface {
	NewReadingProvider(*ConfigProvider) IReadingProvider
	GetReading(int) (*common.Reading, *error)
}
