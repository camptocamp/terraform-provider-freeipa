package freeipa

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccFreeIPADNSUserGroupMembership(t *testing.T) {
	testDatasetUser := map[string]string{
		"login":     "testuser",
		"firstname": "Test",
		"lastname":  "User",
	}
	testDatasetGroup := map[string]string{
		"name": "testgroup",
	}
	testDatasetGroup2 := map[string]string{
		"name": "testgroup-2",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeIPADNSUserGroupMembershipResource_user(testDatasetUser, testDatasetGroup),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_user_group_membership.groupemembership", "name", testDatasetGroup["name"]),
					resource.TestCheckResourceAttr("freeipa_user_group_membership.groupemembership", "user", testDatasetUser["login"]),
				),
			},
			{
				Config: testAccFreeIPADNSUserGroupMembershipResource_group(testDatasetGroup, testDatasetGroup2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_user_group_membership.groupemembership", "name", testDatasetGroup["name"]),
					resource.TestCheckResourceAttr("freeipa_user_group_membership.groupemembership", "group", testDatasetGroup2["name"]),
				),
			},
		},
	})
}

func testAccFreeIPADNSUserGroupMembershipResource_user(dataset_user map[string]string, dataset_group map[string]string) string {
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
	`, provider_host, provider_user, provider_pass, dataset_user["login"], dataset_user["firstname"], dataset_user["lastname"], dataset_group["name"])
}

func testAccFreeIPADNSUserGroupMembershipResource_group(dataset_group map[string]string, dataset_group2 map[string]string) string {
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
	
	resource "freeipa_group" "group" {
		name       = "%s"
	}
	resource "freeipa_group" "subgroup" {
		name       = "%s"
	}

	resource freeipa_user_group_membership "groupemembership" {
	   name = resource.freeipa_group.group.id
	   group = resource.freeipa_group.subgroup.id
	}
	`, provider_host, provider_user, provider_pass, dataset_group["name"], dataset_group2["name"])
}
