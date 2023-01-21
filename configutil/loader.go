package configutil

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Loader loads configs.
type Loader struct {
	File      string       // required
	EnvPrefix string       // optional
	Viper     *viper.Viper // optional
}

// Load loads a config from a file and/or the env.
func (l Loader) Load(ptr interface{}) error {
	v := l.Viper
	if v == nil {
		v = NewViper()
	}
	if l.EnvPrefix != "" {
		v.SetEnvPrefix(l.EnvPrefix)
	}
	v.SetConfigFile(l.File)
	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("could not read-in config: %w", err)
	}
	if err := v.Unmarshal(ptr); err != nil {
		return fmt.Errorf("could not load config: %w", err)
	}
	return nil
}

// NewViper creates a new viper instance with sane defaults.
func NewViper() *viper.Viper {
	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	return v
}
