package format

import (
	"regexp"
    "fmt"
)

func Parse(s string) (string, string) {
	r, _ := regexp.Compile(`^(\[.*?\]\s*)\((.*?)\)(.*)$`)
	comps := r.FindAllStringSubmatch(s, -1)
	var format, style string
	for _, v := range comps {
		style = v[1]
        format = fmt.Sprintf("%s%s", v[0], v[2])
	}
	return format, style
}
