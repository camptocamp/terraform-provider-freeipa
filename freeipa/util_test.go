package freeipa

import (
	"testing"
)

var idTests = []struct {
	id       string
	name     string
	zoneName string
	_type    string
}{
	{"foo", "foo", "", ""},
	{"foo/example.tld", "foo", "example.tld", ""},
	{"foo/example.tld/A", "foo", "example.tld", "A"},
}

func TestSplitId(t *testing.T) {
	for _, tt := range idTests {
		name, zoneName, _type := splitID(tt.id)
		if name != tt.name || zoneName != tt.zoneName || _type != tt._type {
			t.Errorf("splitId(%s) => [%s, %s, %s]) want [%s, %s, %s]", tt.id, name, zoneName, _type, tt.name, tt.zoneName, tt._type)
		}
	}
}
