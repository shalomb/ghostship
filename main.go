package main

import (
	"fmt"

	xdg "github.com/adrg/xdg"

	config "github.com/shalomb/ghostship/config"
	directory "github.com/shalomb/ghostship/directory"
	gitstatus "github.com/shalomb/ghostship/gitstatus"
	log "github.com/sirupsen/logrus"
)

func main() {
	cfgFile, err := xdg.ConfigFile("ghostship/config.toml")
	if err != nil {
		log.Fatalf("Unable to source config file: %v", err)
	}

	requiredModules, configFields := config.Parse(cfgFile)
	log.Printf("config:\n%+v\n%+v\n", requiredModules, configFields)

	for _, v := range requiredModules {
		log.Printf("render required mod: %v\n", v)
	}

	log.Printf("gitstatus: %+v\n", gitstatus.Status())
	a, b := directory.Status()
	log.Printf("directory: %+v %+v\n", a, b)
	fmt.Printf("\n")
}
