package freeipa

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccFreeIPAAutomemberaddCondition(t *testing.T) {
	testDatasetHostgroup := map[string]string{
		"name":        "group-testautomembercond",
		"description": "Host group test",
	}
	testDatasetAutomemberadd := map[string]string{
		"name":        "group-testautomembercond",
		"description": "automember",
		"type":        "hostgroup",
	}
	testAutomemberaddCondition := map[string]string{
		"name":           "group-testautomembercond",
		"description":    "automembercond",
		"type":           "hostgroup",
		"key":            "fqdn",
		"inclusiveregex": "\\.foo\\.bar\\.net$",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeIPAAutomemberaddConditionResource(testDatasetHostgroup, testDatasetAutomemberadd, testAutomemberaddCondition),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_automemberadd_condition.automembercondition", "name", testAutomemberaddCondition["name"]),
					resource.TestCheckResourceAttr("freeipa_automemberadd_condition.automembercondition", "description", testAutomemberaddCondition["description"]),
					resource.TestCheckResourceAttr("freeipa_automemberadd_condition.automembercondition", "type", testAutomemberaddCondition["type"]),
					resource.TestCheckResourceAttr("freeipa_automemberadd_condition.automembercondition", "key", testAutomemberaddCondition["key"]),
					resource.TestCheckResourceAttr("freeipa_automemberadd_condition.automembercondition", "inclusiveregex.0", testAutomemberaddCondition["inclusiveregex"]),
				),
			},
		},
	})
}

func testAccFreeIPAAutomemberaddConditionResource(dataset_group map[string]string, dataset_automemberadd map[string]string, dataset_automemberaddcondition map[string]string) string {
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
	resource "freeipa_automemberadd" "automember" {
		name       = resource.freeipa_hostgroup.hostgroup.name
		type       = "%s"
	}
	resource "freeipa_automemberadd_condition" "automembercondition" {
		name           = freeipa_automemberadd.automember.name
		description    = "%s"
		type           = "%s"
		key            = "%s"
	  inclusiveregex = [%#v]
	}
	`, provider_host, provider_user, provider_pass, dataset_group["name"], dataset_automemberadd["type"], dataset_automemberaddcondition["description"], dataset_automemberaddcondition["type"], dataset_automemberaddcondition["key"], dataset_automemberaddcondition["inclusiveregex"])
}
