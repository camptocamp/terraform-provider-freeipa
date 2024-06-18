package freeipa

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccFreeIPASudoRule(t *testing.T) {
	testSudoRule := map[string]string{
		"name":               "sudo-rule-test",
		"description":        "Test sudo rule",
		"enabled":            "true",
		"usercategory":       "all",
		"hostcategory":       "all",
		"commandcategory":    "all",
		"runasusercategory":  "all",
		"runasgroupcategory": "all",
		"order":              "2",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeIPASudoRuleResource_basic(testSudoRule),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_sudo_rule.test_rule", "name", testSudoRule["name"]),
				),
			},
			{
				Config: testAccFreeIPASudoRuleResource_full(testSudoRule),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_sudo_rule.test_rule", "name", testSudoRule["name"]),
					resource.TestCheckResourceAttr("freeipa_sudo_rule.test_rule", "description", testSudoRule["description"]),
				),
			},
		},
	})
}

func testAccFreeIPASudoRuleResource_basic(dataset map[string]string) string {
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
	`, provider_host, provider_user, provider_pass, dataset["name"])
}

func testAccFreeIPASudoRuleResource_full(dataset map[string]string) string {
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
		name        = "%s"
		description  = "%s"
		enabled = %s
		usercategory = "%s"
		hostcategory = "%s"
		commandcategory = "%s"
		runasusercategory = "%s"
		runasgroupcategory = "%s"
		order = %s
	}
	`, provider_host, provider_user, provider_pass, dataset["name"], dataset["description"], dataset["enabled"], dataset["usercategory"], dataset["hostcategory"],
		dataset["commandcategory"], dataset["runasusercategory"], dataset["runasgroupcategory"], dataset["order"])
}
