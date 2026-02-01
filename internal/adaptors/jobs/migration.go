package jobs

import (
	"os"
	"time"

	"github.com/Calmantara/lis-backend/internal/core/configurations"
	"github.com/Calmantara/lis-backend/internal/utils"
	"github.com/golang-migrate/migrate/v4"
	"github.com/spf13/cobra"

	pg "github.com/Calmantara/lis-backend/internal/adaptors/storage/mysql"

	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Migrate(cmd *cobra.Command, args []string) {
	command := args[0]
	path := args[1]

	// Dependency injection
	configuration := configurations.Load()
	// load logger
	log := utils.NewZap(
		utils.WithAppName(configuration.Application.ApplicationName()),
		utils.WithEnvironment(configuration.Application.Environment),
	)
	log.Infow("Running lis migration job")
	// load database client
	masterDB := pg.NewClient(configuration.DatabaseMaster)
	slaveDB := pg.NewClient(configuration.DatabaseSlave)
	gormConnection := pg.NewConnection(masterDB, slaveDB)

	db, _ := gormConnection.Master().DB()
	driver, _ := mysql.WithInstance(db, &mysql.Config{
		DatabaseName:     configuration.DatabaseMaster.Name,
		StatementTimeout: 30 * time.Second,
	})

	dir := ""
	if d, _ := os.Getwd(); d != "/" && d != "//" {
		dir = d + "/"
	}
	log.Infow("Running migration job", map[string]any{
		"command": command,
		"path":    "file://" + dir + path,
	})
	m, err := migrate.NewWithDatabaseInstance(
		"file://"+dir+path,
		configuration.DatabaseMaster.Name,
		driver,
	)
	if err != nil {
		log.Errorw("failed to find migration file", map[string]any{
			"error": err.Error(),
		})

		panic(err)
	}

	job := m.Down
	if command == "up" {
		job = m.Up
	}
	err = job()

	log.Infow("Successfully executed migration", map[string]any{
		"error": err,
	})
}
