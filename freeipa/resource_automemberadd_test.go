package freeipa

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccFreeIPAAutomemberadd_Hostgroup(t *testing.T) {
	testDatasetHostgroup := map[string]string{
		"name":        "group-testautomember",
		"description": "Host group test",
	}
	testAutomemberadd := map[string]string{
		"name":        "group-testautomember",
		"description": "automember",
		"type":        "hostgroup",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeIPAAutomemberaddResource_basic(testDatasetHostgroup, testAutomemberadd),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_automemberadd.automemberadd", "name", testAutomemberadd["name"]),
					resource.TestCheckResourceAttr("freeipa_automemberadd.automemberadd", "type", "hostgroup"),
				),
			},
			{
				Config: testAccFreeIPAAutomemberaddResource_full(testDatasetHostgroup, testAutomemberadd),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_automemberadd.automemberadd", "name", testAutomemberadd["name"]),
					resource.TestCheckResourceAttr("freeipa_automemberadd.automemberadd", "description", testAutomemberadd["description"]),
					resource.TestCheckResourceAttr("freeipa_automemberadd.automemberadd", "type", "hostgroup"),
				),
			},
		},
	})
}

func testAccFreeIPAAutomemberaddResource_basic(dataset_group map[string]string, dataset_automemberadd map[string]string) string {
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
	  
	resource "freeipa_hostgroup" "hostgroup" {
		name       = "%s"
	}
	resource "freeipa_automemberadd" "automemberadd" {
		name       = resource.freeipa_hostgroup.hostgroup.name
		type       = "%s"
	}
	`, provider_host, provider_user, provider_pass, dataset_group["name"], dataset_automemberadd["type"])
}

func testAccFreeIPAAutomemberaddResource_full(dataset_group map[string]string, dataset_automemberadd map[string]string) string {
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
	  
	resource "freeipa_automemberadd" "automemberadd" {
		name        = "%s"
		description       = "%s"
		type       = "%s"
	}
	`, provider_host, provider_user, provider_pass, dataset_group["name"], dataset_automemberadd["description"], dataset_automemberadd["type"])
}
