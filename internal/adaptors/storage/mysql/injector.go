package mysql

import (
	"github.com/Calmantara/lis-backend/internal/core/configurations"
	"go.uber.org/dig"
)

func NewInjector(digger *dig.Container) {
	master := NewClient(configurations.Config.DatabaseMaster)
	slave := NewClient(configurations.Config.DatabaseSlave)
	mysql := NewConnection(master, slave)

	digger.Provide(func() MySql {
		return mysql
	})

	digger.Provide(NewHealthRepo)
	digger.Provide(NewDeviceMessageCommand)
	digger.Provide(NewTransaction)
}
