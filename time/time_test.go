package time

import (
	"regexp"
	"testing"
    "time"
    "bytes"
    "fmt"

	// log "github.com/sirupsen/logrus"
	strftime "github.com/lestrrat-go/strftime"
	config "github.com/shalomb/ghostship/config"
	log "github.com/sirupsen/logrus"
	assert "github.com/stretchr/testify/assert"
)

func TestTimeFormat(t *testing.T) {
	renderer := Renderer()
	actual, err := renderer.Render(config.DefaultConfig(), config.EnvironmentConfig{})
	assert.Equal(t, err, nil, "Error not empty")

	var buf bytes.Buffer
	tf, _ := strftime.New("%H%M%S")
	if err := tf.Format(&buf, time.Now()); err != nil {
		log.Println(err.Error())
	}
	re, _ := regexp.Compile(`\d{6}`)
	expected := re.FindString(actual)

	log.Printf("TestTimeFormat out: %+v match:%v", actual, expected)
	assert.Equal(t,
        fmt.Sprintf("\x1b[38;5;28m%s\x1b[0m", buf.String()),
        actual,
        "Rendered time string must match",
    )
	assert.Equal(t,
        buf.String(),
        expected,
        "Rendered date must match",
    )
}
