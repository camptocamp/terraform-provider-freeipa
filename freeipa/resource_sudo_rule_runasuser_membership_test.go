package freeipa

import (
	"fmt"
	"os"
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

	resource "freeipa_sudo_rule" "test_rule" {
		name       = "%s"
	}

	resource freeipa_sudo_rule_runasuser_membership "runasusermember" {
		name = freeipa_sudo_rule.test_rule.name
		runasuser = freeipa_user.user.name
	 }

	`, provider_host, provider_user, provider_pass, dataset["user"], dataset["firstname"], dataset["lastname"], dataset["name"])
}
