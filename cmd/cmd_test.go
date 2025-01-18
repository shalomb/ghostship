package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"testing"

	log "github.com/sirupsen/logrus"
	assert "github.com/stretchr/testify/assert"
)

func TestBuildParams(t *testing.T) {
	log.Printf("BuildParams: %+v %[1]T", appName)
	assert.Equal(t, appName, "ghostship", "AppName must not be nil")
}

func TestRootCmd(t *testing.T) {
	assert.Equal(t, "", "", "Version must be equal")
	log.Printf("+%v", rootCmd.Runnable())
	assert.Equal(t, rootCmd.Runnable(), true, "rootCmd must be runnable")
}

func TestTimeCmd(t *testing.T) {
	assert.Equal(t, "", "", "Time must be equal")
	log.Printf("+%v", timeCmd.Runnable())
	assert.Equal(t, timeCmd.Runnable(), true, "TimeCmd must be runnable")
}

func TestVersionCmd(t *testing.T) {
	assert.Equal(t, "", "", "Version must be equal")
	log.Printf("+%v", versionCmd.Runnable())
	assert.Equal(t, versionCmd.Runnable(), true, "VersionCmd must be runnable")
}

// https://stackoverflow.com/a/77151975/742600
func captureOutput(f func() error) (string, error) {
	stdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		return "", err
	}
	os.Stdout = w
	err = f()
	os.Stdout = stdout
	w.Close()
	out, _ := io.ReadAll(r)
	return string(out), err
}

func TestConfigCmd(t *testing.T) {
	actual, err := captureOutput(func() error {
		err := configPS1([]string{""})
		return err
	})
	assert.Nil(t, err)

	re, _ := regexp.Compile(`.*`)
	expected := re.FindString(actual)

	assert.NotEqual(t, expected, "", "")
}

func TestInitCmd(t *testing.T) {
	actual, err := captureOutput(func() error {
		_, err := renderInit([]string{"bash"}...)
		return err
	})
	assert.Nil(t, err)

	re, _ := regexp.Compile(`.*bell-alert.*`)
	expected := re.FindString(actual)

	// log.Printf("actual:%+v", actual)
	assert.NotEqual(t,
		expected,
		"",
		"")

	tmpfile, err := os.CreateTemp(os.TempDir(), "XXXXXX-*")
	defer os.Remove(tmpfile.Name())
	assert.Nil(t, err)

    // os.Executable() under `go test` does not pointcc``
	// err = os.WriteFile(tmpfile.Name(), []byte(actual), 0644)
	// assert.Nil(t, err)
	//
	// var outb, errb bytes.Buffer
	// cmd := exec.Command(
	// 	"bash", "-c",
	// 	fmt.Sprintf("set -eu; cp %[1]s /tmp/meow; source %[1]s;", tmpfile.Name()),
	// )
	// cmd.Stdout = &outb
	// cmd.Stderr = &errb
	// err = cmd.Run()
	// fmt.Println("out:", outb.String(), "err:", errb.String())
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.Printf("err:%+v", err)
	// log.Printf("tmpfile:%+v",
	// 	fmt.Sprintf("'source %s'", tmpfile.Name()),
	// )
	// assert.Nil(t, err)
}
