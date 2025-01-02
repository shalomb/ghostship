package main

import (
	"fmt"
	// "reflect"

	xdg "github.com/adrg/xdg"

	character "github.com/shalomb/ghostship/character"
	config "github.com/shalomb/ghostship/config"
	directory "github.com/shalomb/ghostship/directory"
	gitstatus "github.com/shalomb/ghostship/gitstatus"
	linebreak "github.com/shalomb/ghostship/linebreak"
	renderer "github.com/shalomb/ghostship/renderer"
	time "github.com/shalomb/ghostship/time"

	log "github.com/sirupsen/logrus"
)

func main() {
	cfgFile, err := xdg.ConfigFile("ghostship/config.toml")
	if err != nil {
		log.Fatalf("Unable to source config file: %v", err)
	}

	conf, requiredModules, configFields := config.Parse(cfgFile)
	log.Debugf("config:\n%+v\n%+v\n", requiredModules, configFields)

	handler := renderer.New()
	renderers := make(map[string]renderer.Renderer)

	renderers["character"] = character.Renderer()
	renderers["directory"] = directory.Renderer()
	renderers["gitstatus"] = gitstatus.Renderer()
	renderers["linebreak"] = linebreak.Renderer()
	renderers["time"] = time.Renderer()

	for _, v := range requiredModules {
		val, ok := renderers[v]
		if !ok {
			// log.Warnf("renderers has no key: %+v", v)
			continue
		}
		handler.SetRenderer(val)

		// if val.Name() == "time" {
		// 	log.Printf("time: %+v", configFields["time"])
		// }

		rendered, err := handler.Render(conf, configFields[val.Name()])
		if err == nil {
			fmt.Printf("%s", rendered)
		} else {
			log.Warnf("Failure in Renderer: %v", rendered)
		}
	}
}
