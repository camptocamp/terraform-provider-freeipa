package freeipa

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccFreeIPADNSUser(t *testing.T) {
	testDataset := map[string]string{
		"login":                    "testuser",
		"firstname":                "Test",
		"lastname":                 "User",
		"account_disabled":         "false",
		"car_license":              "A-111-B",
		"city":                     "El Mundo",
		"display_name":             "Test User",
		"email_address":            "testuser@example.com",
		"employee_number":          "000001",
		"employee_type":            "Developer",
		"full_name":                "Test User",
		"gecos":                    "Test User",
		"gid_number":               "10001",
		"home_directory":           "/home/testuser",
		"initials":                 "TU",
		"job_title":                "Developer",
		"krb_principal_name":       "tuser@IPATEST.LAN",
		"login_shell":              "/bin/bash",
		"manager":                  "Dev Manager",
		"mobile_numbers":           "0123456789",
		"organisation_unit":        "Devs",
		"postal_code":              "12345",
		"preferred_language":       "English",
		"province":                 "England",
		"random_password":          "false",
		"ssh_public_key":           "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDDmMkNHn3R+DzSamQDSW60a0iVlAvzbuC3auu8lNoi3u6lvMemsZqPTuvfY4Xlf7uzm+dya3fTRdPKn8sYgPwQ4saUpCSlegN44PjJMhonR1a7FbpHLWj8CRRfzdUSznQhzFcFff0wMBYAklXlyjvdFM8ahl7zHO08HR6469XOVwO1Tb3OGPrXB2lzStK5PKfk5DO/IKl4vHSKhVNVnsZe52rHiZrxOqdGyCijtvwmW2YfIAGc1k4Seqn/Nn7NxKIFBH3hxaUDqgpZneXzuw9GI/F0M8phnHxXNFVZvIWZVcanEeXtH9Z+vVx1ujNcB2QhiPfLMqkNl9db7uykSGKFM4jD0UjGj5kJ8TOC39Safk7XzpQTnrqvIi158zBHVSgugth+QsE1I9/PL2wlzx1qWV2991JKIOc8m52Iwq02tyO8JaSssFTk9szkLTAHedPnZeBbdnlRYcHqX+NPaUh3hqRTZBIR79Ruk6WAliFkED1L0SgwDfGFlevn1Kde9ok=",
		"street_address":           "1600, Pensylvania av.",
		"telephone_numbers":        "1234567890",
		"uid_number":               "10001",
		"userpassword":             "P@ssword",
		"krb_principal_expiration": "2049-12-31T23:59:59Z",
		"krb_password_expiration":  "2049-12-31T23:59:59Z",
		"userclass":                "user-account",
	}
	testDataset2 := map[string]string{
		"login":     "testuser2",
		"firstname": "Test",
		"lastname":  "User",
	}
	testDataset3 := map[string]string{
		"login":     "testuser2",
		"firstname": "Chuck",
		"lastname":  "Norris",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeIPADNSUserResource_basic(testDataset),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_user.user", "name", testDataset["login"]),
					resource.TestCheckResourceAttr("freeipa_user.user", "first_name", testDataset["firstname"]),
					resource.TestCheckResourceAttr("freeipa_user.user", "last_name", testDataset["lastname"]),
				),
			},
			{
				Config: testAccFreeIPADNSUserResource_full(testDataset),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_user.user", "name", testDataset["login"]),
					resource.TestCheckResourceAttr("freeipa_user.user", "first_name", testDataset["firstname"]),
					resource.TestCheckResourceAttr("freeipa_user.user", "last_name", testDataset["lastname"]),
				),
			},
			{
				Config: testAccFreeIPADNSUserResource_basic(testDataset2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_user.user", "name", testDataset2["login"]),
					resource.TestCheckResourceAttr("freeipa_user.user", "first_name", testDataset2["firstname"]),
					resource.TestCheckResourceAttr("freeipa_user.user", "last_name", testDataset2["lastname"]),
				),
			},
			{
				Config: testAccFreeIPADNSUserResource_basic(testDataset3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_user.user", "name", testDataset2["login"]),
					resource.TestCheckResourceAttr("freeipa_user.user", "first_name", testDataset3["firstname"]),
					resource.TestCheckResourceAttr("freeipa_user.user", "last_name", testDataset3["lastname"]),
				),
			},
		},
	})
}

func testAccFreeIPADNSUserResource_basic(dataset map[string]string) string {
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
	`, provider_host, provider_user, provider_pass, dataset["login"], dataset["firstname"], dataset["lastname"])
}

func testAccFreeIPADNSUserResource_full(dataset map[string]string) string {
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
		name        = "%s"
		first_name  = "%s"
		last_name   = "%s"
		account_disabled = %s
		car_license = ["%s"]
		city = "%s"
		display_name = "%s"
		email_address = ["%s"]
		employee_number = "%s"
		employee_type = "%s"
		full_name = "%s"
		gecos = "%s"
		gid_number = %s
		home_directory = "%s"
		initials = "%s"
		job_title = "%s"
		krb_principal_name = ["%s"]
		login_shell = "%s"
		manager = "%s"
		mobile_numbers = ["%s"]
		organisation_unit = "%s"
		postal_code = "%s"
		preferred_language = "%s"
		province = "%s"
		random_password = %s
		ssh_public_key = ["%s"]
		street_address = "%s"
		telephone_numbers = ["%s"]
		uid_number = %s
		userpassword = "%s"
		krb_principal_expiration = "%s"
		krb_password_expiration = "%s"
		userclass = ["%s"]
	}
	`, provider_host, provider_user, provider_pass, dataset["login"], dataset["firstname"], dataset["lastname"], dataset["account_disabled"],
		dataset["car_license"], dataset["city"], dataset["display_name"], dataset["email_address"], dataset["employee_number"], dataset["employee_type"],
		dataset["full_name"], dataset["gecos"], dataset["gid_number"], dataset["home_directory"], dataset["initials"], dataset["job_title"],
		dataset["krb_principal_name"], dataset["login_shell"], dataset["manager"], dataset["mobile_numbers"], dataset["organisation_unit"], dataset["postal_code"],
		dataset["preferred_language"], dataset["province"], dataset["random_password"], dataset["ssh_public_key"], dataset["street_address"], dataset["telephone_numbers"],
		dataset["uid_number"], dataset["userpassword"], dataset["krb_principal_expiration"], dataset["krb_password_expiration"], dataset["userclass"])
}
