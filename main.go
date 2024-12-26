package main

import (
	"fmt"

	xdg "github.com/adrg/xdg"

	config "github.com/shalomb/ghostship/config"
	directory "github.com/shalomb/ghostship/directory"
	gitstatus "github.com/shalomb/ghostship/gitstatus"
	renderer "github.com/shalomb/ghostship/renderer"
	time "github.com/shalomb/ghostship/time"

	log "github.com/sirupsen/logrus"
)

func main() {
	cfgFile, err := xdg.ConfigFile("ghostship/config.toml")
	if err != nil {
		log.Fatalf("Unable to source config file: %v", err)
	}

	requiredModules, configFields := config.Parse(cfgFile)
	log.Printf("config:\n%+v\n%+v\n", requiredModules, configFields)

	handler := renderer.New()
	renderers := make(map[string]renderer.Renderer)

	renderers["gitstatus"] = gitstatus.Renderer()
	renderers["directory"] = directory.Renderer()
	renderers["time"] = time.Renderer()

	for _, v := range requiredModules {
		val, ok := renderers[v]
		if !ok {
			// log.Warnf("renderers has no key: %+v", v)
			continue
		}
		handler.SetRenderer(val)

		rendered, err := handler.Render()
		if err == nil {
			fmt.Printf("%s", rendered)
		} else {
			log.Warnf("Failure in Renderer: %v", rendered)
		}
	}
}
