package freeipa

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccFreeIPASudoRuleDenyCommand(t *testing.T) {
	testSudoRuleDenyCommand := map[string]string{
		"name":       "sudo-rule-test",
		"denycmd1":   "/bin/bash",
		"denycmd2":   "/bin/fish",
		"denycmdgrp": "terminals",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeIPASudoRuleDenyCommandResource_basic(testSudoRuleDenyCommand),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_sudo_rule_denycmd_membership.test_rule_denycmd", "name", testSudoRuleDenyCommand["name"]),
					resource.TestCheckResourceAttr("freeipa_sudo_rule_denycmd_membership.test_rule_denycmd", "sudocmd", testSudoRuleDenyCommand["denycmd1"]),
				),
			},
			{
				Config: testAccFreeIPASudoRuleDenyCommandResource_full(testSudoRuleDenyCommand),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_sudo_rule_denycmd_membership.test_rule_denycmdgrp", "name", testSudoRuleDenyCommand["name"]),
					resource.TestCheckResourceAttr("freeipa_sudo_rule_denycmd_membership.test_rule_denycmdgrp", "sudocmd_group", testSudoRuleDenyCommand["denycmdgrp"]),
				),
			},
		},
	})
}

func testAccFreeIPASudoRuleDenyCommandResource_basic(dataset map[string]string) string {
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
	resource "freeipa_sudo_rule_denycmd_membership" "test_rule_denycmd" {
		name       = freeipa_sudo_rule.test_rule.name
		sudocmd    = freeipa_sudo_cmd.cmd1.name
	}
	`, provider_host, provider_user, provider_pass, dataset["denycmd1"], dataset["name"])
}

func testAccFreeIPASudoRuleDenyCommandResource_full(dataset map[string]string) string {
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
	resource "freeipa_sudo_rule_denycmd_membership" "test_rule_denycmdgrp" {
		name       = freeipa_sudo_rule.test_rule.name
		sudocmd_group = freeipa_sudo_cmdgroup.cmdgroup.name
	}
	`, provider_host, provider_user, provider_pass, dataset["denycmd1"], dataset["denycmd2"], dataset["denycmdgrp"], dataset["name"])
}
