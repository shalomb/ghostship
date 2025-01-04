package time

import (
	"bytes"
	"fmt"
	"time"

	colors "github.com/shalomb/ghostship/colors"
	config "github.com/shalomb/ghostship/config"

	strftime "github.com/lestrrat-go/strftime"
	log "github.com/sirupsen/logrus"
)

const (
	name = "time"
)

// TimeRenderer ...
type TimeRenderer struct{}

// Renderer ...
// TODO: Rename this to New()
func Renderer() *TimeRenderer {
	return &TimeRenderer{}
}

// Name ...
func (r *TimeRenderer) Name() string {
	return name
}

// Render ...
func (r *TimeRenderer) Render(c config.AppConfig, e config.EnvironmentConfig) (string, error) {
	// return pwd, filepath.Base(pwd)
	var buf bytes.Buffer

	cfg := c.TimeConfig
	log.Debugf("TimeConfig: %+v, s: %+v, v: %+v", c, e, cfg.Format)

	tf, err := strftime.New(cfg.Format)
	if err := tf.Format(&buf, time.Now()); err != nil {
		log.Println(err.Error())
	}

	style := cfg.Style
	ret := fmt.Sprintf(
		"%s%s%s",
		colors.ByName(style),
		buf.String(),
		colors.Reset,
	)

	return ret, err
}
