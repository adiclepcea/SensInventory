package sensinventory

import "github.com/adiclepcea/SensInventory/Server/common"

type IReadingProvider interface {
	NewReadingProvider(*ConfigProvider) IReadingProvider
	GetReading(common.Sensor) *common.Reading
}
