package configurations

import (
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestWithDatabasePrefix(t *testing.T) {
	db := &Database{}
	prefix := "test_prefix"

	WithDatabasePrefix(prefix)(db)

	assert.Equal(t, prefix, db.Prefix)
}

func TestWithDatabaseFallbackName(t *testing.T) {
	db := &Database{}
	fallbackName := "test_fallback"

	WithDatabaseFallbackName(fallbackName)(db)

	assert.Equal(t, fallbackName, db.FallbackApplicationName)
}

func TestWithDatabaseViper(t *testing.T) {
	db := &Database{}
	vp := viper.New()

	WithDatabaseViper(vp)(db)

	assert.Equal(t, vp, db.Viper)
}

func TestDatabase_ConnectionString(t *testing.T) {
	db := &Database{
		Host:                    "localhost",
		Port:                    5432,
		Username:                "user",
		Password:                "password",
		Name:                    "testdb",
		SSL:                     "disable",
		FallbackApplicationName: "test_app",
	}

	expected := "host=localhost port=5432 user=user password=password dbname=testdb sslmode=disable fallback_application_name=test_app"
	assert.Equal(t, expected, db.ConnectionString())
}

func TestDatabase_Load_WithCustomViper(t *testing.T) {
	vp := viper.New()
	vp.Set("DATABASE_DRIVER", "postgres")
	vp.Set("DATABASE_HOST", "localhost")
	vp.Set("DATABASE_PORT", 5432)

	db := &Database{}
	db.load(WithDatabaseViper(vp))

	assert.Equal(t, "postgres", db.Driver)
	assert.Equal(t, "localhost", db.Host)
	assert.Equal(t, uint64(5432), db.Port)
}

func TestDatabase_Load_WithPrefix(t *testing.T) {
	vp := viper.New()
	t.Setenv("TEST_DATABASE_DRIVER", "mysql")

	db := &Database{}
	db.load(
		WithDatabaseViper(vp),
		WithDatabasePrefix("TEST"),
	)

	assert.Equal(t, "mysql", db.Driver)
}

func TestDatabase_Load_Unmarshal(t *testing.T) {
	vp := viper.New()
	vp.Set("DATABASE_MAX_CONNECTION", 10)
	vp.Set("DATABASE_MAX_IDLE_CONNECTION", 5)
	vp.Set("DATABASE_MAX_IDLE_TIME", "30s")
	vp.Set("DATABASE_MAX_LIFE_TIME", "1h")

	db := &Database{}
	db.load(WithDatabaseViper(vp))

	assert.Equal(t, 10, db.MaxConnection)
	assert.Equal(t, 5, db.MaxIDleConnection)
	assert.Equal(t, 30*time.Second, db.MaxIDleTime)
	assert.Equal(t, 1*time.Hour, db.MaxLifeTime)
}
