package freeipa

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccFreeIPASudoRuleHost(t *testing.T) {
	testSudoRuleHost := map[string]string{
		"name":      "sudo-rule-test",
		"host":      "host.testacc.ipatest.lan",
		"hostip":    "192.168.10.1",
		"hostgroup": "test-hosts",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeIPASudoRuleHostResource_basic(testSudoRuleHost),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_sudo_rule_host_membership.hostmember", "name", testSudoRuleHost["name"]),
					resource.TestCheckResourceAttr("freeipa_sudo_rule_host_membership.hostmember", "host", testSudoRuleHost["host"]),
				),
			},
		},
	})
}

func testAccFreeIPASudoRuleHostResource_basic(dataset map[string]string) string {
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
	  
	  resource "freeipa_dns_zone" "testacc_ipatest_lan" {
		zone_name          = "testacc.ipatest.lan"
	}
	
	  
	resource "freeipa_host" "host" {
		name       = "%s"
		ip_address = "%s"
		depends_on = [
			freeipa_dns_zone.testacc_ipatest_lan
		]
	}

	resource "freeipa_hostgroup" "hostgroup" {
		name       = "%s"
	}
	
	resource freeipa_host_hostgroup_membership "groupmembership" {
	   name = freeipa_hostgroup.hostgroup.id
	   host = freeipa_host.host.id
	}

	resource "freeipa_sudo_rule" "test_rule" {
		name       = "%s"
	}

	resource freeipa_sudo_rule_host_membership "hostmember" {
		name = freeipa_sudo_rule.test_rule.name
		host = freeipa_host.host.name
	 }

	 resource freeipa_sudo_rule_host_membership "hostgroupmember" {
		name = freeipa_sudo_rule.test_rule.name
		hostgroup = freeipa_hostgroup.hostgroup.name
	 }
	`, provider_host, provider_user, provider_pass, dataset["host"], dataset["hostip"], dataset["hostgroup"], dataset["name"])
}
