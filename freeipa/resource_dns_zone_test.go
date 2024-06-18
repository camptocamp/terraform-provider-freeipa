package freeipa

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccFreeIPADNSDNSZone(t *testing.T) {
	testDnsZone := map[string]string{
		"zone_name":                   "testacc.ipatest.lan",
		"admin_email_address":         "hostmaster@ipatest.lan",
		"allow_inline_dnssec_signing": "false",
		"allow_prt_sync":              "true",
		"allow_query":                 "192.168.1.0/24",
		"allow_transfer":              "1.1.1.1",
		"authoritative_nameserver":    "192.168.1.53",
		"bind_update_policy":          "grant IPATEST.LAN krb5-self * A; grant IPATEST.LAN krb5-self * AAAA; grant IPATEST.LAN krb5-self * SSHFP;",
		"default_ttl":                 "3600",
		"disable_zone":                "false",
		"dynamic_updates":             "true",
		"is_reverse_zone":             "false",
		"skip_nameserver_check":       "true",
		"skip_overlap_check":          "true",
		"soa_expire":                  "86400",
		"soa_minimum":                 "3600",
		"soa_refresh":                 "3600",
		"soa_retry":                   "2",
		"ttl":                         "3600",
		"zone_forwarders":             "1.1.1.1",
	}
	testDnsZoneReverse := map[string]string{
		"zone_name":                   "192.168.1.0/24",
		"admin_email_address":         "hostmaster@ipatest.lan",
		"allow_inline_dnssec_signing": "false",
		"allow_prt_sync":              "true",
		"allow_query":                 "192.168.1.0/24",
		"allow_transfer":              "1.1.1.1",
		"authoritative_nameserver":    "192.168.1.53",
		"bind_update_policy":          "grant IPATEST.LAN krb5-self * A; grant IPATEST.LAN krb5-self * AAAA; grant IPATEST.LAN krb5-self * SSHFP;",
		"default_ttl":                 "3600",
		"disable_zone":                "false",
		"dynamic_updates":             "true",
		"is_reverse_zone":             "true",
		"skip_nameserver_check":       "true",
		"skip_overlap_check":          "true",
		"soa_expire":                  "86400",
		"soa_minimum":                 "3600",
		"soa_refresh":                 "3600",
		"soa_retry":                   "2",
		"ttl":                         "3600",
		"zone_forwarders":             "1.1.1.1",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeIPADNSDNSZoneResource_basic(testDnsZone),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_dns_zone.dnszone", "zone_name", testDnsZone["zone_name"]),
				),
			},
			{
				Config: testAccFreeIPADNSDNSZoneResource_full(testDnsZone),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_dns_zone.dnszone", "zone_name", testDnsZone["zone_name"]),
					resource.TestCheckResourceAttr("freeipa_dns_zone.dnszone", "admin_email_address", testDnsZone["admin_email_address"]),
					resource.TestCheckResourceAttr("freeipa_dns_zone.dnszone", "bind_update_policy", testDnsZone["bind_update_policy"]),
				),
			},
			{
				Config: testAccFreeIPADNSDNSZoneReverseResource_basic(testDnsZoneReverse),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_dns_zone.dnszonereverse", "zone_name", testDnsZoneReverse["zone_name"]),
				),
			},
			{
				Config: testAccFreeIPADNSDNSZoneReverseResource_full(testDnsZoneReverse),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_dns_zone.dnszonereverse", "zone_name", testDnsZoneReverse["zone_name"]),
					resource.TestCheckResourceAttr("freeipa_dns_zone.dnszonereverse", "admin_email_address", testDnsZoneReverse["admin_email_address"]),
					resource.TestCheckResourceAttr("freeipa_dns_zone.dnszonereverse", "bind_update_policy", testDnsZoneReverse["bind_update_policy"]),
				),
			},
		},
	})
}

func testAccFreeIPADNSDNSZoneResource_basic(dataset map[string]string) string {
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
	  
	resource "freeipa_dns_zone" "dnszone" {
		zone_name = "%s"
	}
	`, provider_host, provider_user, provider_pass, dataset["zone_name"])
}

func testAccFreeIPADNSDNSZoneResource_full(dataset map[string]string) string {
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
	  
	resource "freeipa_dns_zone" "dnszone" {
		zone_name                   = "%s"
		admin_email_address         = "%s"
		allow_inline_dnssec_signing = %s
		allow_prt_sync              = %s
		allow_query                 = "%s"
		allow_transfer              = "%s"
		authoritative_nameserver    = "%s."
		bind_update_policy          = "%s"
		default_ttl                 = %s
		disable_zone                = %s
		dynamic_updates             = %s
		is_reverse_zone             = %s
		skip_nameserver_check       = %s
		skip_overlap_check          = %s
		soa_expire                  = %s
		soa_minimum                 = %s
		soa_refresh                 = %s
		soa_retry                   = %s
		ttl                         = %s
		zone_forwarders             = ["%s"]
	}
	`, provider_host, provider_user, provider_pass, dataset["zone_name"], dataset["admin_email_address"], dataset["allow_inline_dnssec_signing"], dataset["allow_prt_sync"],
		dataset["allow_query"], dataset["allow_transfer"], provider_host, dataset["bind_update_policy"], dataset["default_ttl"], dataset["disable_zone"],
		dataset["dynamic_updates"], dataset["is_reverse_zone"], dataset["skip_nameserver_check"], dataset["skip_overlap_check"], dataset["soa_expire"],
		dataset["soa_minimum"], dataset["soa_refresh"], dataset["soa_retry"], dataset["ttl"], dataset["zone_forwarders"])
}

func testAccFreeIPADNSDNSZoneReverseResource_basic(dataset map[string]string) string {
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
	  
	resource "freeipa_dns_zone" "dnszonereverse" {
		zone_name       = "%s"
		is_reverse_zone = %s
	}
	`, provider_host, provider_user, provider_pass, dataset["zone_name"], dataset["is_reverse_zone"])
}

func testAccFreeIPADNSDNSZoneReverseResource_full(dataset map[string]string) string {
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
	  
	resource "freeipa_dns_zone" "dnszonereverse" {
		zone_name                   = "%s"
		admin_email_address         = "%s"
		allow_inline_dnssec_signing = %s
		allow_prt_sync              = %s
		allow_query                 = "%s"
		allow_transfer              = "%s"
		authoritative_nameserver    = "%s."
		bind_update_policy          = "%s"
		default_ttl                 = %s
		disable_zone                = %s
		dynamic_updates             = %s
		is_reverse_zone             = %s
		skip_nameserver_check       = %s
		skip_overlap_check          = %s
		soa_expire                  = %s
		soa_minimum                 = %s
		soa_refresh                 = %s
		soa_retry                   = %s
		ttl                         = %s
		zone_forwarders             = ["%s"]
	}
	`, provider_host, provider_user, provider_pass, dataset["zone_name"], dataset["admin_email_address"], dataset["allow_inline_dnssec_signing"], dataset["allow_prt_sync"],
		dataset["allow_query"], dataset["allow_transfer"], provider_host, dataset["bind_update_policy"], dataset["default_ttl"], dataset["disable_zone"],
		dataset["dynamic_updates"], dataset["is_reverse_zone"], dataset["skip_nameserver_check"], dataset["skip_overlap_check"], dataset["soa_expire"],
		dataset["soa_minimum"], dataset["soa_refresh"], dataset["soa_retry"], dataset["ttl"], dataset["zone_forwarders"])
}
