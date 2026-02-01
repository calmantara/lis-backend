package configurations

import (
	"strings"
	"sync"
	"time"

	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"
)

var (
	Config *Configuration
)

type Environment string

const (
	DEVELOPMENT Environment = "development"
	TEST        Environment = "test"
	PRODUCTION  Environment = "production"
)

type Configuration struct {
	DatabaseMaster Database
	DatabaseSlave  Database
	Application    Application
	JWT            JWT
	LisPlatform    LisPlatform

	mx sync.Mutex
}

func (c *Configuration) Lock() {
	c.mx.Lock()
}

func (c *Configuration) Unlock() {
	c.mx.Unlock()
}

func Load() *Configuration {
	vp := viper.New()
	vp.AutomaticEnv()
	vp.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	var config Configuration

	config.Application.load(vp)
	config.JWT.load(vp)
	config.LisPlatform.load(vp)

	config.DatabaseMaster.load(
		WithDatabaseViper(vp),
		WithDatabasePrefix("MASTER"),
		WithDatabaseFallbackName(
			config.Application.ApplicationName(),
		),
	)
	config.DatabaseSlave.load(
		WithDatabaseViper(vp),
		WithDatabasePrefix("SLAVE"),
		WithDatabaseFallbackName(
			config.Application.ApplicationName(),
		),
	)

	Config = &config

	return &config
}

func keyBind(input any, vp *viper.Viper) {
	mapKey := map[string]any{}
	mapstructure.Decode(input, &mapKey)
	for key := range mapKey {
		vp.BindEnv(key)
	}
}

type Application struct {
	Name        string        `mapstructure:"APP_NAME"`
	Environment Environment   `mapstructure:"APP_ENVIRONMENT"`
	Port        uint64        `mapstructure:"APP_PORT"`
	Graceful    time.Duration `mapstructure:"APP_GRACEFUL"`
	EnableHTTP2 bool          `mapstructure:"APP_ENABLE_HTTP2"`
	Secret      string        `mapstructure:"APP_SECRET"`
	Key         string        `mapstructure:"APP_KEY"`
}

func (a *Application) load(vp *viper.Viper) {
	keyBind(a, vp)
	vp.Unmarshal(&a)

	if len(a.Secret) != 32 {
		panic("APP_SECRET must be at least 32 characters long")
	}
}

func (a *Application) ApplicationName() string {
	return strings.ToLower(a.Name + "_" + string(a.Environment))
}

type JWT struct {
	Issuer            string `mapstructure:"JWT_ISSUER"`
	Secret            string `mapstructure:"JWT_SECRET"`
	SessionExpiration int    `mapstructure:"JWT_SESSION_EXPIRATION"`
	UserExpiration    int    `mapstructure:"JWT_USER_EXPIRATION"`
}

func (j *JWT) load(vp *viper.Viper) {
	keyBind(j, vp)
	vp.Unmarshal(&j)
}

type LisPlatform struct {
	Url string `mapstructure:"LIS_PLATFORM_URL"`
}

func (j *LisPlatform) load(vp *viper.Viper) {
	keyBind(j, vp)
	vp.Unmarshal(&j)
}
