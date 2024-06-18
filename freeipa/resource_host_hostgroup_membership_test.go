package freeipa

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccFreeIPADNSHostGroupMembership(t *testing.T) {
	testDatasetHost := map[string]string{
		"name":       "testhost.testacc.ipatest.lan",
		"ip_address": "192.168.1.10",
	}
	testDatasetHostgroup := map[string]string{
		"name": "test-hostgroup",
	}
	testDatasetHostgroup2 := map[string]string{
		"name": "test-hostgroup-2",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeIPADNSHostGroupMembershipResource_host(testDatasetHost, testDatasetHostgroup),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_host_hostgroup_membership.groupmembership", "name", testDatasetHostgroup["name"]),
					resource.TestCheckResourceAttr("freeipa_host_hostgroup_membership.groupmembership", "host", testDatasetHost["name"]),
				),
			},
			{
				Config: testAccFreeIPADNSHostGroupMembershipResource_hostgroup(testDatasetHostgroup, testDatasetHostgroup2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_host_hostgroup_membership.groupmembership", "name", testDatasetHostgroup["name"]),
					resource.TestCheckResourceAttr("freeipa_host_hostgroup_membership.groupmembership", "hostgroup", testDatasetHostgroup2["name"]),
				),
			},
		},
	})
}

func testAccFreeIPADNSHostGroupMembershipResource_host(dataset_host map[string]string, dataset_hostgroup map[string]string) string {
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
	`, provider_host, provider_user, provider_pass, dataset_host["name"], dataset_host["ip_address"], dataset_hostgroup["name"])
}

func testAccFreeIPADNSHostGroupMembershipResource_hostgroup(dataset_hostgroup map[string]string, dataset_hostgroup2 map[string]string) string {
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
	resource "freeipa_hostgroup" "subgroup" {
		name       = "%s"
	}

	resource freeipa_host_hostgroup_membership "groupmembership" {
	   name = freeipa_hostgroup.hostgroup.id
	   hostgroup = freeipa_hostgroup.subgroup.id
	}
	`, provider_host, provider_user, provider_pass, dataset_hostgroup["name"], dataset_hostgroup2["name"])
}
