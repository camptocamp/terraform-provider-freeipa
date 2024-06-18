package freeipa

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccFreeIPADNSHBAC(t *testing.T) {
	var testHbac map[string]string
	testHbac = map[string]string{
		"name":         "hbac_test",
		"description":  "Automatic test HBAC policy",
		"enabled":      "true",
		"usercategory": "all",
		"hostcategory": "all",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeIPADNSHBACResource_basic(testHbac),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_hbac_policy.hbac_policy", "name", testHbac["name"]),
				),
			},
			{
				Config: testAccFreeIPADNSHBACResource_full(testHbac),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_hbac_policy.hbac_policy", "name", testHbac["name"]),
					resource.TestCheckResourceAttr("freeipa_hbac_policy.hbac_policy", "description", testHbac["description"]),
					resource.TestCheckResourceAttr("freeipa_hbac_policy.hbac_policy", "usercategory", testHbac["usercategory"]),
					resource.TestCheckResourceAttr("freeipa_hbac_policy.hbac_policy", "hostcategory", testHbac["hostcategory"]),
				),
			},
		},
	})
}

func testAccFreeIPADNSHBACResource_basic(dataset map[string]string) string {
	provider_host := os.Getenv("FREEIPA_HOST")
	provider_user := os.Getenv("FREEIPA_USERNAME")
	provider_pass := os.Getenv("FREEIPA_PASSWORD")
	return fmt.Sprintf(`
	provider "freeipa" {
		host     = "%s"
		username = "%s"
		password = "%s"
		insecure = true
	  }
	  
	resource "freeipa_hbac_policy" "hbac_policy" {
		name       = "%s"
	}
	`, provider_host, provider_user, provider_pass, dataset["name"])
}

func testAccFreeIPADNSHBACResource_full(dataset map[string]string) string {
	provider_host := os.Getenv("FREEIPA_HOST")
	provider_user := os.Getenv("FREEIPA_USERNAME")
	provider_pass := os.Getenv("FREEIPA_PASSWORD")
	return fmt.Sprintf(`
	provider "freeipa" {
		host     = "%s"
		username = "%s"
		password = "%s"
		insecure = true
	  }
	  
	resource "freeipa_hbac_policy" "hbac_policy" {
		name        = "%s"
		description  = "%s"
		enabled = %s
		usercategory = "%s"
		hostcategory = "%s"
	}
	`, provider_host, provider_user, provider_pass, dataset["name"], dataset["description"], dataset["enabled"], dataset["usercategory"], dataset["hostcategory"])
}
