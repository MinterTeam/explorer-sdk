package helpers

import (
	"regexp"
)

func RemovePrefix(s string) string {
	re, err := regexp.Compile(`(?i)(mx|mt|mp|mc)`)
	if err != nil {
		return ""
	}
	return re.ReplaceAllString(s, "")
}
