package freeipa

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccFreeIPADNSHostgroup(t *testing.T) {
	testHostgroup := map[string]string{
		"name":        "testhostgroup",
		"description": "Host group test",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeIPADNSHostgroupResource_basic(testHostgroup),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_hostgroup.hostgroup", "name", testHostgroup["name"]),
				),
			},
			{
				Config: testAccFreeIPADNSHostgroupResource_full(testHostgroup),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_hostgroup.hostgroup", "name", testHostgroup["name"]),
					resource.TestCheckResourceAttr("freeipa_hostgroup.hostgroup", "description", testHostgroup["description"]),
				),
			},
		},
	})
}

func testAccFreeIPADNSHostgroupResource_basic(dataset map[string]string) string {
	return fmt.Sprintf(`
	resource "freeipa_hostgroup" "hostgroup" {
		name       = "%s"
	}
	`, dataset["name"])
}

func testAccFreeIPADNSHostgroupResource_full(dataset map[string]string) string {
	return fmt.Sprintf(`
	resource "freeipa_hostgroup" "hostgroup" {
		name        = "%s"
		description  = "%s"
	}
	`, dataset["name"], dataset["description"])
}
