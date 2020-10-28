package freeipa

import (
	"strings"
)

func splitID(id string) (a, b string) {
	if strings.Contains(id, "/") {
		return id[0:strings.Index(id, "/")], id[strings.Index(id, "/")+1:]
	}
	return "", id
}
