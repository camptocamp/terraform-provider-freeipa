package freeipa

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccFreeIPASudocmd(t *testing.T) {
	testSudocmd := map[string]string{
		"name":         "/bin/bash",
		"description":  "The bash terminal",
		"description2": "The other bash terminal",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeIPASudocmdResource_basic(testSudocmd),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_sudo_cmd.cmd", "name", testSudocmd["name"]),
				),
			},
			{
				Config: testAccFreeIPASudocmdResource_full(testSudocmd),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_sudo_cmd.cmd", "name", testSudocmd["name"]),
					resource.TestCheckResourceAttr("freeipa_sudo_cmd.cmd", "description", testSudocmd["description"]),
				),
			},
			{
				Config: testAccFreeIPASudocmdResource_update(testSudocmd),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_sudo_cmd.cmd", "name", testSudocmd["name"]),
					resource.TestCheckResourceAttr("freeipa_sudo_cmd.cmd", "description", testSudocmd["description2"]),
				),
			},
		},
	})
}

func testAccFreeIPASudocmdResource_basic(dataset map[string]string) string {
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
	  
	resource "freeipa_sudo_cmd" "cmd" {
		name       = "%s"
	}
	`, provider_host, provider_user, provider_pass, dataset["name"])
}

func testAccFreeIPASudocmdResource_full(dataset map[string]string) string {
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
	  
	resource "freeipa_sudo_cmd" "cmd" {
		name        = "%s"
		description  = "%s"
	}
	`, provider_host, provider_user, provider_pass, dataset["name"], dataset["description"])
}

func testAccFreeIPASudocmdResource_update(dataset map[string]string) string {
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
	  
	resource "freeipa_sudo_cmd" "cmd" {
		name        = "%s"
		description  = "%s"
	}
	`, provider_host, provider_user, provider_pass, dataset["name"], dataset["description2"])
}
