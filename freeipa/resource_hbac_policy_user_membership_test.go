package freeipa

import (
	"fmt"
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
		},
	})
}

func testAccFreeIPADNSHBACUserMembershipResource_user(dataset map[string]string) string {
	return fmt.Sprintf(`
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
	`, dataset["name"], dataset["user"], dataset["first_name"], dataset["last_name"])
}

// func testAccFreeIPADNSHBACUserMembershipResource_group(dataset map[string]string) string {
// 	return fmt.Sprintf(`
// 	resource "freeipa_hbac_policy" "hbac_policy" {
// 		name       = "%s"
// 	}

// 	resource "freeipa_group" "group-0" {
// 		name       = "%s"
// 	}

// 	resource "freeipa_hbac_policy_user_membership" "hbac_group" {
// 		name = freeipa_hbac_policy.hbac_policy.name
// 		group = freeipa_group.group-0.name
// 	}
// 	`, dataset["name"], dataset["group"])
// }
