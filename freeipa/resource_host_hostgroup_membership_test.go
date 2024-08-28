package freeipa

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccFreeIPADNSHostGroupMembership(t *testing.T) {
	testDatasetHostgroup := map[string]string{
		"name": "test-hostgroup",
	}
	testDatasetHostgroup2 := map[string]string{
		"name": "test-hostgroup-2",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeIPADNSHostGroupMembershipResource_hostgroup(testDatasetHostgroup, testDatasetHostgroup2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_host_hostgroup_membership.groupmembership", "name", testDatasetHostgroup["name"]),
					resource.TestCheckResourceAttr("freeipa_host_hostgroup_membership.groupmembership", "hostgroup", testDatasetHostgroup2["name"]),
				),
			},
		},
	})
}

func testAccFreeIPADNSHostGroupMembershipResource_hostgroup(dataset_hostgroup map[string]string, dataset_hostgroup2 map[string]string) string {
	return fmt.Sprintf(`
	resource "freeipa_hostgroup" "hostgroup" {
		name       = "%s"
	}
	resource "freeipa_hostgroup" "subgroup" {
		name       = "%s"
	}

	resource freeipa_host_hostgroup_membership "groupmembership" {
	   name = freeipa_hostgroup.hostgroup.id
	   hostgroup = freeipa_hostgroup.subgroup.id
	}
	`, dataset_hostgroup["name"], dataset_hostgroup2["name"])
}
