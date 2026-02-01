package configurations

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type (
	DatabaseOption func(d *Database)

	Database struct {
		*viper.Viper
		Prefix                  string
		FallbackApplicationName string        `mapstructure:"DATABASE_FALLBACK_APPLICATION_NAME"`
		Driver                  string        `mapstructure:"DATABASE_DRIVER"`
		Username                string        `mapstructure:"DATABASE_USERNAME"`
		Password                string        `mapstructure:"DATABASE_PASSWORD"`
		Host                    string        `mapstructure:"DATABASE_HOST"`
		Name                    string        `mapstructure:"DATABASE_NAME"`
		SSL                     string        `mapstructure:"DATABASE_SSL"`
		Port                    uint64        `mapstructure:"DATABASE_PORT"`
		MaxConnection           int           `mapstructure:"DATABASE_MAX_CONNECTION"`
		MaxIDleConnection       int           `mapstructure:"DATABASE_MAX_IDLE_CONNECTION"`
		MaxIDleTime             time.Duration `mapstructure:"DATABASE_MAX_IDLE_TIME"`
		MaxLifeTime             time.Duration `mapstructure:"DATABASE_MAX_LIFE_TIME"`
	}
)

func WithDatabasePrefix(s string) DatabaseOption {
	return func(d *Database) {
		d.Prefix = s
	}
}

func WithDatabaseFallbackName(s string) DatabaseOption {
	return func(d *Database) {
		d.FallbackApplicationName = s
	}
}

func WithDatabaseViper(vp *viper.Viper) DatabaseOption {
	return func(d *Database) {
		d.Viper = vp
	}
}

func (d *Database) ConnectionString() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%v)/%s",
		d.Username,
		d.Password,
		d.Host,
		d.Port,
		d.Name,
	)
}

func (d *Database) load(options ...DatabaseOption) {
	for _, fn := range options {
		fn(d)
	}
	// make new viper
	vp := viper.New()
	if d.Viper != nil {
		vp = d.Viper
	}
	// make key binding with suffix
	if d.Prefix != "" {
		vp.SetEnvPrefix(strings.ToUpper(d.Prefix))
	}
	vp.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	vp.AutomaticEnv()
	// binding to struct
	keyBind(d, vp)
	vp.Unmarshal(d)
}
