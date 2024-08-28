package freeipa

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccFreeIPASudocmdgroup(t *testing.T) {
	testSudocmdgroup := map[string]string{
		"name":         "services",
		"description":  "Service management related sudo commands",
		"description2": "Service management related sudo commands but not the same",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeIPASudocmdgroupResource_basic(testSudocmdgroup),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_sudo_cmdgroup.cmdgroup", "name", testSudocmdgroup["name"]),
				),
			},
			{
				Config: testAccFreeIPASudocmdgroupResource_full(testSudocmdgroup),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_sudo_cmdgroup.cmdgroup", "name", testSudocmdgroup["name"]),
					resource.TestCheckResourceAttr("freeipa_sudo_cmdgroup.cmdgroup", "description", testSudocmdgroup["description"]),
				),
			},
			{
				Config: testAccFreeIPASudocmdgroupResource_update(testSudocmdgroup),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_sudo_cmdgroup.cmdgroup", "name", testSudocmdgroup["name"]),
					resource.TestCheckResourceAttr("freeipa_sudo_cmdgroup.cmdgroup", "description", testSudocmdgroup["description2"]),
				),
			},
		},
	})
}

func testAccFreeIPASudocmdgroupResource_basic(dataset map[string]string) string {
	return fmt.Sprintf(`
	resource "freeipa_sudo_cmdgroup" "cmdgroup" {
		name       = "%s"
	}
	`, dataset["name"])
}

func testAccFreeIPASudocmdgroupResource_full(dataset map[string]string) string {
	return fmt.Sprintf(`
	resource "freeipa_sudo_cmdgroup" "cmdgroup" {
		name        = "%s"
		description  = "%s"
	}
	`, dataset["name"], dataset["description"])
}

func testAccFreeIPASudocmdgroupResource_update(dataset map[string]string) string {
	return fmt.Sprintf(`
	resource "freeipa_sudo_cmdgroup" "cmdgroup" {
		name        = "%s"
		description  = "%s"
	}
	`, dataset["name"], dataset["description2"])
}
