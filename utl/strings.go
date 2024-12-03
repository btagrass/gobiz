package utl

import (
	"strings"
)

func Contains(s string, substrs ...string) bool {
	for _, b := range substrs {
		if strings.Contains(s, b) {
			return true
		}
	}
	return false
}

func HasPrefix(s string, prefixs ...string) bool {
	for _, p := range prefixs {
		if strings.HasPrefix(s, p) || strings.HasPrefix(s, strings.ToUpper(p)) {
			return true
		}
	}
	return false
}

func HasSuffix(s string, suffixs ...string) bool {
	for _, f := range suffixs {
		if strings.HasSuffix(s, f) || strings.HasSuffix(s, strings.ToUpper(f)) {
			return true
		}
	}
	return false
}

func InsertAfter(s string, findinserts ...string) string {
	var oldnews []string
	for i := 0; i < len(findinserts); i++ {
		if i%2 == 1 {
			oldnews = append(oldnews, findinserts[i-1]+findinserts[i])
		} else {
			oldnews = append(oldnews, findinserts[i])
		}
	}
	return Replace(s, oldnews...)
}

func InsertBefore(s string, findinserts ...string) string {
	var oldnews []string
	for i := 0; i < len(findinserts); i++ {
		if i%2 == 1 {
			oldnews = append(oldnews, findinserts[i]+findinserts[i-1])
		} else {
			oldnews = append(oldnews, findinserts[i])
		}
	}
	return Replace(s, oldnews...)
}

func Replace(s string, oldnews ...string) string {
	return strings.NewReplacer(oldnews...).Replace(s)
}

func Split(s string, seps ...rune) []string {
	return strings.FieldsFunc(s, func(r rune) bool {
		for _, sp := range seps {
			if r == sp {
				return true
			}
		}
		return false
	})
}

func Trim(s string, cutsets ...string) string {
	var oldnews []string
	for _, c := range cutsets {
		oldnews = append(oldnews, c, "")
	}
	return Replace(s, oldnews...)
}
