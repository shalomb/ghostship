package character

import (
	"fmt"

	colors "github.com/shalomb/ghostship/colors"
	config "github.com/shalomb/ghostship/config"
	format "github.com/shalomb/ghostship/format"
	log "github.com/sirupsen/logrus"
)

const (
	name = ""
)

type CharacterRenderer struct{}

// Renderer ...
func Renderer() *CharacterRenderer {
	return &CharacterRenderer{}
}

// Name ...
func (r *CharacterRenderer) Name() string {
	return name
}

// func LookupFieldByTag(cfg any, target string) string {
// 	val := reflect.ValueOf(cfg)
// 	for i := 0; i < val.NumField(); i++ {
// 		field := val.Type().Field(i)
// 		// fieldByName := reflect.ValueOf(&conf).Elem().FieldByName(field.Name).Interface()
// 		kind := val.Type().Field(i).Type.Kind().String()
// 		tag := val.Type().Field(i).Tag.Get("toml")
//         if t == tag {
//         }
//         log.Printf("f:%v, k:%v, t:%+v", field, kind, tag)
// 	}
//     return ""
// }

// Render ...
func (r *CharacterRenderer) Render(c config.AppConfig, e config.EnvironmentConfig) (string, error) {
	cfg := c.CharacterConfig
	log.Debugf("CharacterConfig: %+v, s: %+v, v: %+v", c, e, cfg.Format)

	var symbol, stl string
	if e["status"] == uint16(0) && e["pipestatus"] == uint16(0) {
		symbol, stl = format.Parse(cfg.SuccessSymbol)
	} else {
		symbol, stl = format.Parse(cfg.ErrorSymbol)
	}

	style := colors.ByExpression(stl)

	ret := fmt.Sprintf(
		"\\[%s\\]%s\\[%s\\] ",
		style,
		symbol,
		colors.Reset,
	)

	return ret, nil
}
