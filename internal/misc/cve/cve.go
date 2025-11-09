package cve

import (
	"regexp"
)

func IsValidCVEFormat(s string) bool {
	pattern := regexp.MustCompile(`^CVE-\d{4}-\d{4,}$`)
	return pattern.MatchString(s)
}
