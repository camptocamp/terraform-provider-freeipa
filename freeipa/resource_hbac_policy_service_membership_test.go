package freeipa

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccFreeIPADNSHBACServiceMembership(t *testing.T) {
	var testHbacSvc map[string]string
	testHbacSvc = map[string]string{
		"name":    "hbac_test",
		"service": "sshd",
	}
	var testHbacSvcGroup map[string]string
	testHbacSvcGroup = map[string]string{
		"name":         "hbac_test",
		"servicegroup": "Sudo",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeIPADNSHBACServiceMembershipResource_service(testHbacSvc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_hbac_policy_service_membership.hbac_svc", "name", testHbacSvc["name"]),
					resource.TestCheckResourceAttr("freeipa_hbac_policy_service_membership.hbac_svc", "service", testHbacSvc["service"]),
				),
			},
			{
				Config: testAccFreeIPADNSHBACServiceMembershipResource_servicegroup(testHbacSvcGroup),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_hbac_policy_service_membership.hbac_svcgrp", "name", testHbacSvcGroup["name"]),
					resource.TestCheckResourceAttr("freeipa_hbac_policy_service_membership.hbac_svcgrp", "servicegroup", testHbacSvcGroup["servicegroup"]),
				),
			},
		},
	})
}

func testAccFreeIPADNSHBACServiceMembershipResource_service(dataset map[string]string) string {
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

	resource "freeipa_hbac_policy_service_membership" "hbac_svc" {
		name = freeipa_hbac_policy.hbac_policy.name
		service = "%s"
	}
	`, provider_host, provider_user, provider_pass, dataset["name"], dataset["service"])
}

func testAccFreeIPADNSHBACServiceMembershipResource_servicegroup(dataset map[string]string) string {
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

	resource "freeipa_hbac_policy_service_membership" "hbac_svcgrp" {
		name = freeipa_hbac_policy.hbac_policy.name
		servicegroup = "%s"
	}
	`, provider_host, provider_user, provider_pass, dataset["name"], dataset["servicegroup"])
}
