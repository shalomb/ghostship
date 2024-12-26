package directory

import (
	"fmt"
	"os"
	"path/filepath"
)

var (
    directory string
)

func Status() (string, string) {
    pwd, err := os.Getwd()
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    return pwd, filepath.Base(pwd)
}
