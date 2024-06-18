package freeipa

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccFreeIPADNSHBACUserMembership(t *testing.T) {
	var testHbacUser map[string]string
	testHbacUser = map[string]string{
		"name":       "hbac_test",
		"user":       "test-user",
		"first_name": "Bill",
		"last_name":  "Bachnov",
	}
	var testHbacGroup map[string]string
	testHbacGroup = map[string]string{
		"name":  "hbac_test",
		"group": "test-group",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeIPADNSHBACUserMembershipResource_user(testHbacUser),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_hbac_policy_user_membership.hbac_user", "name", testHbacUser["name"]),
					resource.TestCheckResourceAttr("freeipa_hbac_policy_user_membership.hbac_user", "user", testHbacUser["user"]),
				),
			},
			{
				Config: testAccFreeIPADNSHBACUserMembershipResource_group(testHbacGroup),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_hbac_policy_user_membership.hbac_group", "name", testHbacGroup["name"]),
					resource.TestCheckResourceAttr("freeipa_hbac_policy_user_membership.hbac_group", "group", testHbacGroup["group"]),
				),
			},
		},
	})
}

func testAccFreeIPADNSHBACUserMembershipResource_user(dataset map[string]string) string {
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
	  
	resource "freeipa_hbac_policy" "hbac_policy" {
		name       = "%s"
	}

	resource "freeipa_user" "user-0" {
		name       = "%s"
		first_name = "%s"
		last_name  = "%s"
	}

	resource "freeipa_hbac_policy_user_membership" "hbac_user" {
		name = freeipa_hbac_policy.hbac_policy.name
		user = freeipa_user.user-0.name
	}
	`, provider_host, provider_user, provider_pass, dataset["name"], dataset["user"], dataset["first_name"], dataset["last_name"])
}

func testAccFreeIPADNSHBACUserMembershipResource_group(dataset map[string]string) string {
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
	  
	resource "freeipa_hbac_policy" "hbac_policy" {
		name       = "%s"
	}

	resource "freeipa_group" "group-0" {
		name       = "%s"
	}

	resource "freeipa_hbac_policy_user_membership" "hbac_group" {
		name = freeipa_hbac_policy.hbac_policy.name
		group = freeipa_group.group-0.name
	}
	`, provider_host, provider_user, provider_pass, dataset["name"], dataset["group"])
}
