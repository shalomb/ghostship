package format

import (
	"regexp"
    "fmt"
)

func Parse(s string) (string, string) {
	r, _ := regexp.Compile(`^\[(.*?)\]\s*\((.*?)\)(.*)$`)
	comps := r.FindAllStringSubmatch(s, -1)
	var format, style string
	for _, v := range comps {
        format = fmt.Sprintf("%s%s", v[1], v[3])
		style = v[2]
	}
	return format, style
}
