package freeipa

// import (
// 	"fmt"
// 	"testing"

// 	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
// )

// func TestAccFreeIPASudoRuleUser(t *testing.T) {
// 	testSudoRuleUser := map[string]string{
// 		"name":      "sudo-rule-test",
// 		"user":      "test-user",
// 		"firstname": "Test",
// 		"lastname":  "User",
// 		"usergroup": "test-group",
// 	}

// 	resource.Test(t, resource.TestCase{
// 		PreCheck:  func() { testAccPreCheck(t) },
// 		Providers: testAccProviders,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccFreeIPASudoRuleUserResource_basic(testSudoRuleUser),
// 				Check: resource.ComposeTestCheckFunc(
// 					resource.TestCheckResourceAttr("freeipa_sudo_rule_user_membership.usermember", "user", testSudoRuleUser["user"]),
// 					resource.TestCheckResourceAttr("freeipa_sudo_rule_user_membership.groupmember", "group", testSudoRuleUser["usergroup"]),
// 				),
// 			},
// 		},
// 	})
// }

// func testAccFreeIPASudoRuleUserResource_basic(dataset map[string]string) string {
// 	return fmt.Sprintf(`
// 	  resource "freeipa_user" "user" {
// 		name       = "%s"
// 		first_name = "%s"
// 		last_name  = "%s"
// 	}

// 	resource "freeipa_group" "group" {
// 		name       = "%s"
// 	}
// 	resource freeipa_user_group_membership "groupemembership" {
// 	   name = resource.freeipa_group.group.id
// 	   user = resource.freeipa_user.user.id
// 	}

// 	resource "freeipa_sudo_rule" "test_rule" {
// 		name       = "%s"
// 	}

// 	resource freeipa_sudo_rule_user_membership "usermember" {
// 		name = freeipa_sudo_rule.test_rule.name
// 		user = freeipa_user.user.id
// 	 }

// 	 resource freeipa_sudo_rule_user_membership "groupmember" {
// 		name = freeipa_sudo_rule.test_rule.name
// 		group = freeipa_group.group.id
// 	 }
// 	`, dataset["user"], dataset["firstname"], dataset["lastname"], dataset["usergroup"], dataset["name"])
// }
