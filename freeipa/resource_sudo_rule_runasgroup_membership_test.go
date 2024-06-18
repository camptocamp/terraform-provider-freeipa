package freeipa

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccFreeIPASudoRuleRunAsGroup(t *testing.T) {
	testSudoRuleRunAsGroup := map[string]string{
		"name":      "sudo-rule-test",
		"user":      "test-user",
		"firstname": "Test",
		"lastname":  "User",
		"usergroup": "test-group",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeIPASudoRuleRunAsGroupResource_basic(testSudoRuleRunAsGroup),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_sudo_rule_runasgroup_membership.runasgroupmember", "runasgroup", testSudoRuleRunAsGroup["usergroup"]),
				),
			},
		},
	})
}

func testAccFreeIPASudoRuleRunAsGroupResource_basic(dataset map[string]string) string {
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

	  resource "freeipa_user" "user" {
		name       = "%s"
		first_name = "%s"
		last_name  = "%s"
	}

	resource "freeipa_group" "group" {
		name       = "%s"
	}
	resource freeipa_user_group_membership "groupemembership" {
	   name = resource.freeipa_group.group.id
	   user = resource.freeipa_user.user.id
	}

	resource "freeipa_sudo_rule" "test_rule" {
		name       = "%s"
	}

	 resource freeipa_sudo_rule_runasgroup_membership "runasgroupmember" {
		name = freeipa_sudo_rule.test_rule.name
		runasgroup = freeipa_group.group.name
	 }
	`, provider_host, provider_user, provider_pass, dataset["user"], dataset["firstname"], dataset["lastname"], dataset["usergroup"], dataset["name"])
}
