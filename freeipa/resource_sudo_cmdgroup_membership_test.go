package freeipa

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccFreeIPASudocmdgroupMembership(t *testing.T) {
	testSudocmdgroupMembership := map[string]string{
		"name":     "terminals",
		"sudocmd":  "/bin/bash",
		"sudocmd2": "/bin/fish",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeIPASudocmdgroupMembershipResource_basic(testSudocmdgroupMembership),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_sudo_cmdgroup_membership.cmdgroup_member", "name", testSudocmdgroupMembership["name"]),
				),
			},
			{
				Config: testAccFreeIPASudocmdgroupMembershipResource_full(testSudocmdgroupMembership),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_sudo_cmdgroup_membership.cmdgroup_member", "name", testSudocmdgroupMembership["name"]),
					resource.TestCheckResourceAttr("freeipa_sudo_cmdgroup_membership.cmdgroup_member2", "name", testSudocmdgroupMembership["name"]),
				),
			},
		},
	})
}

func testAccFreeIPASudocmdgroupMembershipResource_basic(dataset map[string]string) string {
	return fmt.Sprintf(`
	resource "freeipa_sudo_cmd" "cmd" {
		name       = "%s"
	}

	resource "freeipa_sudo_cmd" "cmd2" {
		name       = "%s"
	}

	resource "freeipa_sudo_cmdgroup" "cmdgroup" {
		name       = "%s"
	}

	resource "freeipa_sudo_cmdgroup_membership" "cmdgroup_member" {
		name       = freeipa_sudo_cmdgroup.cmdgroup.name
		sudocmd    = freeipa_sudo_cmd.cmd.name
	}
	`, dataset["sudocmd"], dataset["sudocmd2"], dataset["name"])
}

func testAccFreeIPASudocmdgroupMembershipResource_full(dataset map[string]string) string {
	return fmt.Sprintf(`
	  resource "freeipa_sudo_cmd" "cmd" {
		name       = "%s"
	}

	resource "freeipa_sudo_cmd" "cmd2" {
		name       = "%s"
	}

	resource "freeipa_sudo_cmdgroup" "cmdgroup" {
		name       = "%s"
	}

	resource "freeipa_sudo_cmdgroup_membership" "cmdgroup_member" {
		name       = freeipa_sudo_cmdgroup.cmdgroup.name
		sudocmd    = freeipa_sudo_cmd.cmd.name
	}

	resource "freeipa_sudo_cmdgroup_membership" "cmdgroup_member2" {
		name       = freeipa_sudo_cmdgroup.cmdgroup.id
		sudocmd    = freeipa_sudo_cmd.cmd2.id
	}
	`, dataset["sudocmd"], dataset["sudocmd2"], dataset["name"])
}
