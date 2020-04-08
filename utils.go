package di

import "strings"

func parseTag(tag string) (name string, optional bool) {
	options := strings.Split(tag, ",")
	if len(options) == 0 {
		return "", false
	}
	if len(options) == 1 && options[0] == "optional" {
		return "", true
	}
	if len(options) == 1 {
		return options[0], false
	}
	if len(options) == 2 && options[1] == "optional" {
		return options[0], true
	}
	panic("incorrect di tag")
}
