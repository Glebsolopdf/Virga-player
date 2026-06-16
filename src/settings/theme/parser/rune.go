package parser

import "strings"

func ParseRune(value string) (rune, bool) {
	v := strings.TrimSpace(value)
	v = strings.Trim(v, `"'`)
	v = strings.ReplaceAll(v, `\\`, `\`)
	v = strings.ReplaceAll(v, `\t`, "\t")
	v = strings.ReplaceAll(v, `\n`, "\n")
	if v == "" {
		return 0, false
	}
	r := []rune(v)
	if len(r) == 0 {
		return 0, false
	}
	return r[0], true
}
