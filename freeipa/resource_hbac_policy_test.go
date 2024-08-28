package freeipa

import (
	"fmt"
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
	return fmt.Sprintf(`
	resource "freeipa_hbac_policy" "hbac_policy" {
		name       = "%s"
	}
	`, dataset["name"])
}

func testAccFreeIPADNSHBACResource_full(dataset map[string]string) string {
	return fmt.Sprintf(`
	resource "freeipa_hbac_policy" "hbac_policy" {
		name        = "%s"
		description  = "%s"
		enabled = %s
		usercategory = "%s"
		hostcategory = "%s"
	}
	`, dataset["name"], dataset["description"], dataset["enabled"], dataset["usercategory"], dataset["hostcategory"])
}
