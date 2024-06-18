package freeipa

import (
	"fmt"
	"os"
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
	`, provider_host, provider_user, provider_pass, dataset["sudocmd"], dataset["sudocmd2"], dataset["name"])
}

func testAccFreeIPASudocmdgroupMembershipResource_full(dataset map[string]string) string {
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
	`, provider_host, provider_user, provider_pass, dataset["sudocmd"], dataset["sudocmd2"], dataset["name"])
}
