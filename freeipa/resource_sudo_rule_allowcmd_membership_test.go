package freeipa

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccFreeIPASudoRuleAllowCommand(t *testing.T) {
	testSudoRuleAllowCommand := map[string]string{
		"name":        "sudo-rule-test",
		"allowcmd1":   "/bin/bash",
		"allowcmd2":   "/bin/fish",
		"allowcmdgrp": "terminals",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeIPASudoRuleAllowCommandResource_basic(testSudoRuleAllowCommand),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_sudo_rule_allowcmd_membership.test_rule_allowcmd", "name", testSudoRuleAllowCommand["name"]),
					resource.TestCheckResourceAttr("freeipa_sudo_rule_allowcmd_membership.test_rule_allowcmd", "sudocmd", testSudoRuleAllowCommand["allowcmd1"]),
				),
			},
			{
				Config: testAccFreeIPASudoRuleAllowCommandResource_full(testSudoRuleAllowCommand),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_sudo_rule_allowcmd_membership.test_rule_allowcmdgrp", "name", testSudoRuleAllowCommand["name"]),
					resource.TestCheckResourceAttr("freeipa_sudo_rule_allowcmd_membership.test_rule_allowcmdgrp", "sudocmd_group", testSudoRuleAllowCommand["allowcmdgrp"]),
				),
			},
		},
	})
}

func testAccFreeIPASudoRuleAllowCommandResource_basic(dataset map[string]string) string {
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
	
	
	  
	resource "freeipa_sudo_cmd" "cmd1" {
		name        = "%s"
	}
	resource "freeipa_sudo_rule" "test_rule" {
		name       = "%s"
	}
	resource "freeipa_sudo_rule_allowcmd_membership" "test_rule_allowcmd" {
		name       = freeipa_sudo_rule.test_rule.name
		sudocmd    = freeipa_sudo_cmd.cmd1.name
	}
	`, provider_host, provider_user, provider_pass, dataset["allowcmd1"], dataset["name"])
}

func testAccFreeIPASudoRuleAllowCommandResource_full(dataset map[string]string) string {
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
	  
		    
	resource "freeipa_sudo_cmd" "cmd1" {
		name        = "%s"
	}
	resource "freeipa_sudo_cmd" "cmd2" {
		name        = "%s"
	}
	resource "freeipa_sudo_cmdgroup" "cmdgroup" {
		name       = "%s"
	}
	resource "freeipa_sudo_cmdgroup_membership" "cmdgroup_member1" {
		name       = freeipa_sudo_cmdgroup.cmdgroup.name
		sudocmd    = freeipa_sudo_cmd.cmd1.name
	}
	resource "freeipa_sudo_cmdgroup_membership" "cmdgroup_member2" {
		name       = freeipa_sudo_cmdgroup.cmdgroup.name
		sudocmd    = freeipa_sudo_cmd.cmd2.name
	}
	resource "freeipa_sudo_rule" "test_rule" {
		name       = "%s"
	}
	resource "freeipa_sudo_rule_allowcmd_membership" "test_rule_allowcmdgrp" {
		name       = freeipa_sudo_rule.test_rule.name
		sudocmd_group = freeipa_sudo_cmdgroup.cmdgroup.name
	}
	`, provider_host, provider_user, provider_pass, dataset["allowcmd1"], dataset["allowcmd2"], dataset["allowcmdgrp"], dataset["name"])
}
