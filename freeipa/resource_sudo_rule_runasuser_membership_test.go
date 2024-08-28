package freeipa

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccFreeIPASudoRuleRunAsUser(t *testing.T) {
	testSudoRuleRunAsUser := map[string]string{
		"name":      "sudo-rule-test",
		"user":      "test-user",
		"firstname": "Test",
		"lastname":  "User",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeIPASudoRuleRunAsUserResource_basic(testSudoRuleRunAsUser),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_sudo_rule_runasuser_membership.runasusermember", "runasuser", testSudoRuleRunAsUser["user"]),
				),
			},
		},
	})
}

func testAccFreeIPASudoRuleRunAsUserResource_basic(dataset map[string]string) string {
	return fmt.Sprintf(`
	  resource "freeipa_user" "user" {
		name       = "%s"
		first_name = "%s"
		last_name  = "%s"
	}

	resource "freeipa_sudo_rule" "test_rule" {
		name       = "%s"
	}

	resource freeipa_sudo_rule_runasuser_membership "runasusermember" {
		name = freeipa_sudo_rule.test_rule.name
		runasuser = freeipa_user.user.name
	 }

	`, dataset["user"], dataset["firstname"], dataset["lastname"], dataset["name"])
}
