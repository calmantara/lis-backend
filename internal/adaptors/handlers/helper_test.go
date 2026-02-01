package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/Calmantara/lis-backend/internal/core/configurations"
	"github.com/Calmantara/lis-backend/internal/core/models"
	"github.com/Calmantara/lis-backend/internal/core/services"
	"github.com/Calmantara/lis-backend/internal/utils"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/dig"

	ms "github.com/Calmantara/lis-backend/internal/adaptors/storage/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

var (
	gormConnection ms.MySql
)

func mockData[T any](t *testing.T, data T) (result T) {
	t.Helper()

	err := gormConnection.
		Master().
		Create(&data).
		Error
	assert.NoError(t, err)

	t.Cleanup(func() {
		gormConnection.
			Master().
			Unscoped().
			Delete(&data)
	})

	return data
}

func injector(t *testing.T, digger *dig.Container) *echo.Echo {
	t.Helper()

	configurations.Load()
	utils.NewZap(
		utils.WithAppName(configurations.Config.Application.Name),
		utils.WithEnvironment(configurations.Config.Application.Environment),
	)

	engine := mockEcho(t)
	digger.Provide(func() models.Router {
		return models.Router{
			External: engine.V1External,
			Internal: engine.V1Internal,
		}
	}, dig.Name("routerV1"))

	ms.NewInjector(digger)
	services.NewInjector(digger)
	NewInjector(digger)

	err := digger.Invoke(Invoke)
	assert.NoError(t, err)

	if gormConnection == nil {
		slave := ms.NewClient(configurations.Config.DatabaseSlave)
		master := ms.NewClient(configurations.Config.DatabaseMaster)
		gormConnection = ms.NewConnection(master, slave)
	}

	return engine.Engine
}

func setupMigration(t *testing.T) {
	t.Helper()
	// setup migration
	env := configurations.Load()
	slave := ms.NewClient(env.DatabaseSlave)
	master := ms.NewClient(env.DatabaseMaster)
	// dependencies injection
	gormConnection := ms.NewConnection(master, slave)
	db, _ := gormConnection.Master().DB()
	driver, _ := mysql.WithInstance(db, &mysql.Config{})
	dir, _ := os.Getwd()
	m, _ := migrate.NewWithDatabaseInstance(
		"file://"+dir+"/../../../../tools/migrations",
		"postgres", driver)

	err := m.Up()
	fmt.Println("error in migration", err)

	query := `
	CREATE TABLE IF NOT EXISTS device_sensors_001
	PARTITION OF device_sensors
		FOR VALUES FROM ('1970-01-01')
					TO ('3000-01-01')
	`
	gormConnection.
		Master().
		Exec(query)
}

type echoEngine struct {
	Engine     *echo.Echo
	V1External *echo.Group
	V1Internal *echo.Group
}

func mockEcho(t *testing.T) echoEngine {
	t.Helper()
	// echo engine
	engine := echo.New()
	// v1
	externalV1 := engine.Group("/api/v1")
	internalV1 := engine.Group("/api/internal/v1")
	return echoEngine{
		Engine:     engine,
		V1External: externalV1,
		V1Internal: internalV1,
	}
}

func setupRequest(method, path string, body any, headers map[string]any) (rec *httptest.ResponseRecorder, req *http.Request) {
	// Setup
	if body != nil {
		str := body.(string)
		req = httptest.NewRequest(
			method,
			path,
			strings.NewReader(str),
		)
	} else {
		req = httptest.NewRequest(method, path, nil)
	}
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("User-Agent", "mock")

	basicHeader := map[string]any{
		"X-Request-ID":     "mock-request-id",
		"X-Application-ID": "mock-application-id",
	}
	for key, value := range basicHeader {
		req.Header.Set(key, fmt.Sprintf("%v", value))
	}

	for key, value := range headers {
		req.Header.Set(key, fmt.Sprintf("%v", value))
	}
	rec = httptest.NewRecorder()

	return
}
