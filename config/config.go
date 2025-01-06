// Package config ...
package config

import (
	"os"
	"reflect"
	"regexp"
	"strings"
	// "fmt"
	toml "github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
)

// EnvironmentConfig ...
type EnvironmentConfig map[string]any

// BaseComponentConfig ...
type BaseComponentConfig struct {
	Disabled bool   `toml:"disabled"`
	Format   string `toml:"format"`
	Style    string `toml:"style"`
}

// CharacterConfig ...
type CharacterConfig struct {
	BaseComponentConfig
	SuccessSymbol string `toml:"success_symbol"`
	ErrorSymbol   string `toml:"error_symbol"`
}

// CommandNumberConfig ...
type CommandNumberConfig struct {
	BaseComponentConfig
}

// DurationConfig ...
type DurationConfig struct {
	BaseComponentConfig
	MinTime uint16 `toml:"min_time"`
}

// DirectoryConfig ...
type DirectoryConfig struct {
	BaseComponentConfig
}

// GitStatusConfig ...
type GitStatusConfig struct {
	BaseComponentConfig
	DirtyStyle  string `toml:"dirty_style"`
	DriftStyle  string `toml:"drift_style"`
	NormalStyle string `toml:"normal_style"`
	StagedStyle string `toml:"staged_style"`
	SymbolStyle string `toml:"symbol_style"`
}

// StatusConfig ...
type StatusConfig struct {
	BaseComponentConfig
}

// TimeConfig ...
type TimeConfig struct {
	BaseComponentConfig
}

// AppConfig ...
type AppConfig struct {
	Format              string
	CharacterConfig     `toml:"character"`
	CommandNumberConfig `toml:"commandno"`
	DirectoryConfig     `toml:"directory"`
	DurationConfig      `toml:"duration"`
	GitStatusConfig     `toml:"gitstatus"`
	StatusConfig        `toml:"status"`
	TimeConfig          `toml:"time"`
}

// DefaultConfig ...
func DefaultConfig() AppConfig {
	return AppConfig{
		Format: "foo",
		CharacterConfig: CharacterConfig{
			BaseComponentConfig{
				Format: "$symbol ",
				Style:  "blue",
			},
			"[❯](bold green) ", // SuccessSymbol
			"[❯](bold red) ",   // ErrorSymbol
		},
		CommandNumberConfig: CommandNumberConfig{
			BaseComponentConfig{
				Format: "!",
				Style:  "gray bold",
			},
		},
		DurationConfig: DurationConfig{
			BaseComponentConfig{
				Format: "took %ds",
				Style:  "gray bold",
			},
			2, // MinTime
		},
		StatusConfig: StatusConfig{
			BaseComponentConfig{
				Style: "green",
			},
		},
		TimeConfig: TimeConfig{
			BaseComponentConfig{
				Format: "%H%M%S",
				Style:  "green",
			},
		},
	}
}

// Parse ...
func Parse(cfgFile string) (AppConfig, []string, map[string]string) {
	conf := DefaultConfig()

	tomlData, err := os.ReadFile(cfgFile)
	if err != nil {
		log.Fatalf("Unable to read config file: %v", err)
	}

	if _, err := toml.Decode(string(tomlData), &conf); err != nil {
		log.Fatalf("Unable to parse config file: %+v", err)
	}

	var requiredFields []string
	if conf.Format != "" {
		r, _ := regexp.Compile(`\$([^\$]+)`)
		comps := r.FindAllStringSubmatch((strings.TrimRight(conf.Format, "\n")), -2)
		for _, v := range comps {
			requiredFields = append(requiredFields, v[1])
		}
	}
	log.Debugf("requiredFields: %+v", requiredFields)

	confFields := make(map[string]string)
	val := reflect.ValueOf(conf)
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		// fieldByName := reflect.ValueOf(&conf).Elem().FieldByName(field.Name).Interface()
		kind := val.Type().Field(i).Type.Kind().String()
		tag := val.Type().Field(i).Tag.Get("toml")
		if (kind) == "struct" {
			confFields[tag] = field.Name
		}
	}
	log.Debugf("confFields: %+v\n", confFields)

	return conf, requiredFields, confFields
}
