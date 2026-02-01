package configurations

import (
	"os"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfiguration(t *testing.T) {
	t.Setenv("APP_NAME", "TestApp")
	t.Setenv("APP_ENVIRONMENT", "test")
	t.Setenv("APP_PORT", "8080")
	t.Setenv("APP_GRACEFUL", "5s")
	t.Setenv("APP_ENABLE_HTTP2", "true")
	t.Setenv("JWT_ISSUER", "TestIssuer")
	t.Setenv("JWT_SECRET", "TestSecret")
	t.Setenv("JWT_SESSION_EXPIRATION", "3600")
	t.Setenv("JWT_USER_EXPIRATION", "7200")

	config := Load()

	assert.Equal(t, "TestApp", config.Application.Name)
	assert.Equal(t, TEST, config.Application.Environment)
	assert.Equal(t, uint64(8080), config.Application.Port)
	assert.Equal(t, 5*time.Second, config.Application.Graceful)
	assert.True(t, config.Application.EnableHTTP2)

	assert.Equal(t, "TestIssuer", config.JWT.Issuer)
	assert.Equal(t, "TestSecret", config.JWT.Secret)
	assert.Equal(t, 3600, config.JWT.SessionExpiration)
	assert.Equal(t, 7200, config.JWT.UserExpiration)
}

func TestLockAndUnlock(t *testing.T) {
	config := &Configuration{}

	// Test Lock
	config.Lock()
	assert.NotNil(t, &config.mx, "Mutex should be locked")

	// Test Unlock
	config.Unlock()
	assert.NotNil(t, &config.mx, "Mutex should be unlocked")
}

func TestLoadConfigurationWithDefaults(t *testing.T) {
	// Clear environment variables to test defaults
	os.Clearenv()

	assert.Panics(t, func() {
		Load()
	})
}

func TestApplicationLoad(t *testing.T) {
	vp := viper.New()
	vp.AutomaticEnv()

	secret := ""
	for len(secret) < 32 {
		secret += "1"
	}

	t.Setenv("APP_NAME", "TestApp")
	t.Setenv("APP_ENVIRONMENT", "test")
	t.Setenv("APP_PORT", "8080")
	t.Setenv("APP_GRACEFUL", "5s")
	t.Setenv("APP_ENABLE_HTTP2", "true")
	t.Setenv("APP_SECRET", secret)

	app := Application{}
	app.load(vp)

	assert.Equal(t, "TestApp", app.Name)
	assert.Equal(t, TEST, app.Environment)
	assert.Equal(t, uint64(8080), app.Port)
	assert.Equal(t, 5*time.Second, app.Graceful)
	assert.True(t, app.EnableHTTP2)
}

func TestApplicationLoad_EmptySecret(t *testing.T) {
	vp := viper.New()
	vp.AutomaticEnv()
	t.Setenv("APP_NAME", "TestApp")
	t.Setenv("APP_ENVIRONMENT", "test")
	t.Setenv("APP_PORT", "8080")
	t.Setenv("APP_GRACEFUL", "5s")
	t.Setenv("APP_ENABLE_HTTP2", "true")
	t.Setenv("APP_SECRET", "")

	app := Application{}
	assert.Panics(t, func() {
		app.load(vp)
	})

	assert.Equal(t, "TestApp", app.Name)
	assert.Equal(t, TEST, app.Environment)
	assert.Equal(t, uint64(8080), app.Port)
	assert.Equal(t, 5*time.Second, app.Graceful)
	assert.True(t, app.EnableHTTP2)
}

func TestJWTLoad(t *testing.T) {
	vp := viper.New()
	vp.AutomaticEnv()
	t.Setenv("JWT_ISSUER", "TestIssuer")
	t.Setenv("JWT_SECRET", "TestSecret")
	t.Setenv("JWT_SESSION_EXPIRATION", "3600")
	t.Setenv("JWT_USER_EXPIRATION", "7200")

	jwt := JWT{}
	jwt.load(vp)

	assert.Equal(t, "TestIssuer", jwt.Issuer)
	assert.Equal(t, "TestSecret", jwt.Secret)
	assert.Equal(t, 3600, jwt.SessionExpiration)
	assert.Equal(t, 7200, jwt.UserExpiration)
}
