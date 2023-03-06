package configutil

import (
	"fmt"
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/spf13/viper"
)

// Loader loads configs.
type Loader struct {
	File      string       // optional
	EnvPrefix string       // optional
	Koanf     *koanf.Koanf // optional, default
	Viper     *viper.Viper // optional
}

// Load loads a config from a file and/or the env.
func (l Loader) Load(ptr interface{}) error {
	k := l.Koanf
	if k == nil && l.Viper == nil {
		k = NewKoanf() // Use koand by default
	}

	if k != nil {
		if l.File != "" {
			if err := k.Load(file.Provider(l.File), yaml.Parser()); err != nil {
				return fmt.Errorf("could not load config file %s: %w", l.File, err)
			}
		}

		prefix := strings.TrimSuffix(l.EnvPrefix, "_") + "_" // add trailing underscore if not present
		err := k.Load(env.Provider(prefix, ".", func(s string) string {
			return strings.Replace(strings.ToLower(strings.TrimPrefix(s, prefix)), "_", ".", -1)
		}), nil)
		if err != nil {
			return fmt.Errorf("could not load env config: %w", err)
		}

		if err = k.Unmarshal("", ptr); err != nil {
			return fmt.Errorf("could not load config: %w", err)
		}
		return nil
	}

	// Fallback to viper
	v := l.Viper
	if l.Viper == nil {
		v = NewViper()
	}
	if l.EnvPrefix != "" {
		v.SetEnvPrefix(l.EnvPrefix)
	}
	v.SetConfigFile(l.File) // viper requires to set a file: https://github.com/spf13/viper/issues/584
	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("could not read-in config: %w", err)
	}
	if err := v.Unmarshal(ptr); err != nil {
		return fmt.Errorf("could not load config: %w", err)
	}

	return nil
}

// NewKoanf creates a new koanf instance.
func NewKoanf() *koanf.Koanf {
	return koanf.New(".")
}

// NewViper creates a new viper instance with sane defaults.
func NewViper() *viper.Viper {
	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	return v
}
