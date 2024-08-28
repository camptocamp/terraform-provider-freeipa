package freeipa

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccFreeIPADNSHBACHostMembership(t *testing.T) {
	var testHbacHostgroup map[string]string
	testHbacHostgroup = map[string]string{
		"name":      "hbac_test",
		"hostgroup": "test-hostgroup",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeIPADNSHBACHostMembership_hostgroup(testHbacHostgroup),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_hbac_policy_host_membership.hbac_hostgroup", "name", testHbacHostgroup["name"]),
					resource.TestCheckResourceAttr("freeipa_hbac_policy_host_membership.hbac_hostgroup", "hostgroup", testHbacHostgroup["hostgroup"]),
				),
			},
		},
	})
}

func testAccFreeIPADNSHBACHostMembership_hostgroup(dataset map[string]string) string {
	return fmt.Sprintf(`
	resource "freeipa_hbac_policy" "hbac_policy" {
		name       = "%s"
	}

	  
	resource "freeipa_hostgroup" "hostgroup" {
		name       = "%s"
	}

	resource "freeipa_hbac_policy_host_membership" "hbac_hostgroup" {
		name = freeipa_hbac_policy.hbac_policy.name
		hostgroup = freeipa_hostgroup.hostgroup.name
		depends_on = [
			freeipa_hbac_policy.hbac_policy,
			freeipa_hostgroup.hostgroup
		]
	}
	`, dataset["name"], dataset["hostgroup"])
}
