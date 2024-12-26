package config

import (
	// "fmt"
    "os"
	"reflect"
	"regexp"
	"strings"
    log "github.com/sirupsen/logrus"
	toml "github.com/BurntSushi/toml"
)

type BaseComponentConfig struct {
	Disabled bool   `toml:"disabled"`
	Format   string `toml:"format"`
	Style    string `toml:"style"`
}

type GitStatusConfig struct {
	BaseComponentConfig
}

type DirectoryConfig struct {
	BaseComponentConfig
}

type tomlConfig struct {
	Format          string
	GitStatusConfig `toml:"gitstatus"`
	DirectoryConfig `toml:"directory"`
}

func Parse(cfgFile string) ([]string, map[string]any) {
    tomlData , err := os.ReadFile(cfgFile)
        if err != nil {
        log.Fatalf("Unable to read config file: %v", err)
    }

	var conf tomlConfig
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

	confFields := make(map[string]any)
	val := reflect.ValueOf(conf)
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		kind := val.Type().Field(i).Type.Kind().String()
		tag := val.Type().Field(i).Tag.Get("toml")
		if (kind) == "struct" {
			confFields[tag] = field
		}
	}
	log.Debugf("confFields: %+v\n", confFields)

    return requiredFields, confFields
}
