package freeipa

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccFreeIPASudoRuleOption(t *testing.T) {
	testSudoRuleOption := map[string]string{
		"name":   "sudo-rule-test",
		"option": "!authenticate",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeIPASudoRuleOptionResource_basic(testSudoRuleOption),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_sudo_rule_option.sudoopt", "option", testSudoRuleOption["option"]),
				),
			},
		},
	})
}

func testAccFreeIPASudoRuleOptionResource_basic(dataset map[string]string) string {
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

	resource "freeipa_sudo_rule" "test_rule" {
		name       = "%s"
	}

	 resource freeipa_sudo_rule_option "sudoopt" {
		name = freeipa_sudo_rule.test_rule.name
		option = "%s"
	 }
	`, provider_host, provider_user, provider_pass, dataset["name"], dataset["option"])
}
