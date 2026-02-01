package jobs

import (
	"os"
	"strconv"
	"strings"
	"testing"

	pg "github.com/Calmantara/lis-backend/internal/adaptors/storage/mysql"
	"github.com/Calmantara/lis-backend/internal/core/configurations"
	"github.com/stretchr/testify/assert"
)

func TestMigration(t *testing.T) {
	assert.NotPanics(t, func() {
		Migrate(nil, []string{
			"up",
			"../../../tools/migrations",
		})
		// get latest migration file
		files, _ := os.ReadDir("../../../tools/migrations")
		latestVersion := strings.Split(files[len(files)-1].Name(), "_")[0]

		masterDB := pg.NewClient(configurations.Load().DatabaseMaster)
		slaveDB := pg.NewClient(configurations.Load().DatabaseSlave)
		gormConnection := pg.NewConnection(masterDB, slaveDB)

		version := 0
		err := gormConnection.
			Master().
			Select("version").
			Table("schema_migrations").
			Order("version DESC").
			Limit(1).
			Scan(&version).Error
		assert.NoError(t, err)

		// assert version
		assert.EqualValues(t, latestVersion, strconv.Itoa(version))
	})
}

func TestMigrationPanic(t *testing.T) {
	assert.Panics(t, func() {
		Migrate(nil, []string{
			"up",
			"tools/migrations_invalid_path",
		})
	})
}
