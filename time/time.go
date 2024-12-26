package time

import (
	"bytes"
	"time"

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
func (r *TimeRenderer) Render() (string, error) {
	// return pwd, filepath.Base(pwd)
	var buf bytes.Buffer

	f, err := strftime.New(`%H%M%S`)
	if err := f.Format(&buf, time.Now()); err != nil {
		log.Println(err.Error())
	}
	return buf.String(), err
}
