package freeipa

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"freeipa": testAccProvider,
	}
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("FREEIPA_HOST"); v == "" {
		t.Fatal("FREEIPA_HOST must be set for acceptance tests")
	}
	if v := os.Getenv("FREEIPA_USERNAME"); v == "" {
		t.Fatal("FREEIPA_USERNAME must be set for acceptance tests")
	}
	if v := os.Getenv("FREEIPA_PASSWORD"); v == "" {
		t.Fatal("FREEIPA_PASSWORD must be set for acceptance tests")
	}
}
