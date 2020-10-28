package freeipa

import (
	"testing"
)

var idTests = []struct {
	id         string
	zoneName   string
	resourceID string
}{
	{"foo", "foo", ""},
	{"foo/example.tld", "example.tld", "foo"},
}

func TestSplitId(t *testing.T) {
	for _, tt := range idTests {
		resourceID, zoneName := splitID(tt.id)
		if zoneName != tt.zoneName || resourceID != tt.resourceID {
			t.Errorf("splitId(%s) => [%s, %s]) want [%s, %s]", tt.id, resourceID, zoneName, tt.resourceID, tt.zoneName)
		}
	}
}
