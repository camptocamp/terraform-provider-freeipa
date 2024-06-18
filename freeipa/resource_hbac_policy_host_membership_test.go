package freeipa

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccFreeIPADNSHBACHostMembership(t *testing.T) {
	var testHbacHost map[string]string
	testHbacHost = map[string]string{
		"name":     "hbac_test",
		"hostname": "test-host.testacc.ipatest.lan",
		"ip":       "192.168.1.52",
	}
	var testHbacHostgroup map[string]string
	testHbacHostgroup = map[string]string{
		"name":      "hbac_test",
		"hostgroup": "test-hostgroup",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeIPADNSHBACHostMembership_host(testHbacHost),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_hbac_policy_host_membership.hbac_host", "name", testHbacHost["name"]),
					resource.TestCheckResourceAttr("freeipa_hbac_policy_host_membership.hbac_host", "host", testHbacHost["hostname"]),
				),
			},
			{
				Config: testAccFreeIPADNSHBACHostMembership_hostgroup(testHbacHostgroup),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_hbac_policy_host_membership.hbac_hostgroup", "name", testHbacHostgroup["name"]),
					resource.TestCheckResourceAttr("freeipa_hbac_policy_host_membership.hbac_hostgroup", "hostgroup", testHbacHostgroup["hostgroup"]),
				),
			},
		},
	})
}

func testAccFreeIPADNSHBACHostMembership_host(dataset map[string]string) string {
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

	resource "freeipa_dns_zone" "testacc_ipatest_lan" {
		zone_name          = "testacc.ipatest.lan"
		skip_overlap_check = true
		disable_zone       = false
		dynamic_updates    = true
	  
	}

	resource "freeipa_dns_zone" "reverse_192_168_1_0" {
		zone_name          = "192.168.1.0/24"
		is_reverse_zone = true
		skip_overlap_check = true
		disable_zone       = false
		dynamic_updates    = true
	  }
	  
	  
	resource "freeipa_host" "host-0" {
		name       = "%s"
		ip_address = "%s"
		depends_on = [
			freeipa_dns_zone.testacc_ipatest_lan,
			freeipa_dns_zone.reverse_192_168_1_0
		]
	}

	resource "freeipa_hbac_policy_host_membership" "hbac_host" {
		name = freeipa_hbac_policy.hbac_policy.name
		host = freeipa_host.host-0.name
		depends_on = [
			freeipa_host.host-0,
			freeipa_hbac_policy.hbac_policy
		]
	}
	`, provider_host, provider_user, provider_pass, dataset["name"], dataset["hostname"], dataset["ip"])
}

func testAccFreeIPADNSHBACHostMembership_hostgroup(dataset map[string]string) string {
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

	  
	resource "freeipa_hostgroup" "hostgroup" {
		name       = "%s"
	}

	resource "freeipa_hbac_policy_host_membership" "hbac_hostgroup" {
		name = freeipa_hbac_policy.hbac_policy.name
		hostgroup = freeipa_hostgroup.hostgroup.name
		depends_on = [
			freeipa_hbac_policy.hbac_policy,
			freeipa_hostgroup.hostgroup
		]
	}
	`, provider_host, provider_user, provider_pass, dataset["name"], dataset["hostgroup"])
}
