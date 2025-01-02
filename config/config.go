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

// DirectoryConfig ...
type DirectoryConfig struct {
	BaseComponentConfig
}

// GitStatusConfig ...
type GitStatusConfig struct {
	BaseComponentConfig
}

// TimeConfig ...
type TimeConfig struct {
	BaseComponentConfig
}

// TomlConfig ...
type TomlConfig struct {
	Format          string
	CharacterConfig `toml:"character"`
	DirectoryConfig `toml:"directory"`
	GitStatusConfig `toml:"gitstatus"`
	TimeConfig      `toml:"time"`
}

// DefaultConfig ...
func DefaultConfig() TomlConfig {
	return TomlConfig{
		Format: "foo",
		CharacterConfig: CharacterConfig{
			BaseComponentConfig{
				Format: "$symbol ",
				Style:  "blue",
			},
			"[❯](bold green) ", // SuccessSymbol
			"[❯](bold red) ",   // ErrorSymbol
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
func Parse(cfgFile string) (TomlConfig, []string, map[string]string) {
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
