package freeipa

import (
	"context"
	"log"
	"strings"
	"time"

	ipa "github.com/camptocamp/go-freeipa/freeipa"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceFreeIPAUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFreeIPADNSUserCreate,
		ReadContext:   resourceFreeIPADNSUserRead,
		UpdateContext: resourceFreeIPADNSUserUpdate,
		DeleteContext: resourceFreeIPADNSUserDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"first_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "First name",
			},
			"last_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Last name",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "UID",
			},
			"full_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Full name",
			},
			"display_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Display name",
			},
			"initials": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Initials",
			},
			"home_directory": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Home directory",
			},
			"gecos": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "GECOS",
			},
			"login_shell": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Login shell",
			},
			"krb_principal_name": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Principal alias",
			},
			"krb_principal_expiration": {
				Description: "Kerberos principal expiration " +
					"[RFC3339](https://datatracker.ietf.org/doc/html/rfc3339#section-5.8) format " +
					"(see [RFC3339 time string](https://tools.ietf.org/html/rfc3339#section-5.8) e.g., " +
					"`YYYY-MM-DDTHH:MM:SSZ`)",
				Type:     schema.TypeString,
				Optional: true,
			},
			"krb_password_expiration": {
				Description: "User password expiration " +
					"[RFC3339](https://datatracker.ietf.org/doc/html/rfc3339#section-5.8) format " +
					"(see [RFC3339 time string](https://tools.ietf.org/html/rfc3339#section-5.8) e.g., " +
					"`YYYY-MM-DDTHH:MM:SSZ`)",
				Type:     schema.TypeString,
				Optional: true,
			},
			"userpassword": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Prompt to set the user password",
			},
			"email_address": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Email address",
			},
			"telephone_numbers": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Telephone Number",
			},
			"mobile_numbers": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Mobile Telephone Number",
			},
			"random_password": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Generate a random user password",
			},
			"uid_number": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "User ID Number (system will assign one if not provided)",
			},
			"gid_number": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Group ID Number",
			},
			"street_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Street address",
			},
			"city": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "City",
			},
			"province": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "State/Province",
			},
			"postal_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ZIP code",
			},
			"organisation_unit": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Org. Unit",
			},
			"job_title": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Job Title",
			},
			"manager": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Manager",
			},
			"employee_number": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Employee Number",
			},
			"employee_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Employee Type",
			},
			"preferred_language": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Preferred Language",
			},
			"account_disabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Account disabled",
			},
			"ssh_public_key": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "SSH public key",
			},
			"car_license": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Car License",
			},
			"userclass": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "User category (semantics placed on this attribute are for local interpretation)",
			},
		},
	}
}

func resourceFreeIPADNSUserCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Creating freeipa user")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	optArgs := ipa.UserAddOptionalArgs{}

	args := ipa.UserAddArgs{
		Givenname: d.Get("first_name").(string),
		Sn:        d.Get("last_name").(string),
	}

	if _v, ok := d.GetOkExists("name"); ok {
		v := _v.(string)
		optArgs.UID = &v
	}
	if _v, ok := d.GetOkExists("full_name"); ok {
		v := _v.(string)
		optArgs.Cn = &v
	}
	if _v, ok := d.GetOkExists("display_name"); ok {
		v := _v.(string)
		optArgs.Displayname = &v
	}
	if _v, ok := d.GetOkExists("initials"); ok {
		v := _v.(string)
		optArgs.Initials = &v
	}
	if _v, ok := d.GetOkExists("home_directory"); ok {
		v := _v.(string)
		optArgs.Homedirectory = &v
	}
	if _v, ok := d.GetOkExists("gecos"); ok {
		v := _v.(string)
		optArgs.Gecos = &v
	}
	if _v, ok := d.GetOkExists("login_shell"); ok {
		v := _v.(string)
		optArgs.Loginshell = &v
	}
	if _v, ok := d.GetOkExists("krb_principal_name"); ok {
		v := utilsGetArry(_v.([]interface{}))
		optArgs.Krbprincipalname = &v
	}
	if _v, ok := d.GetOkExists("userpassword"); ok {
		v := _v.(string)
		optArgs.Userpassword = &v
	}
	if _v, ok := d.GetOkExists("email_address"); ok {
		v := utilsGetArry(_v.([]interface{}))
		optArgs.Mail = &v
	}
	if _v, ok := d.GetOkExists("telephone_numbers"); ok {
		v := utilsGetArry(_v.([]interface{}))
		optArgs.Telephonenumber = &v
	}
	if _v, ok := d.GetOkExists("mobile_numbers"); ok {
		v := utilsGetArry(_v.([]interface{}))
		optArgs.Mobile = &v
	}
	if _v, ok := d.GetOkExists("random_password"); ok {
		v := _v.(bool)
		optArgs.Random = &v
	}
	if _v, ok := d.GetOkExists("uid_number"); ok {
		v := _v.(int)
		optArgs.Uidnumber = &v
	}
	if _v, ok := d.GetOkExists("gid_number"); ok {
		v := _v.(int)
		optArgs.Gidnumber = &v
	}
	if _v, ok := d.GetOkExists("street_address"); ok {
		v := _v.(string)
		optArgs.Street = &v
	}
	if _v, ok := d.GetOkExists("city"); ok {
		v := _v.(string)
		optArgs.L = &v
	}
	if _v, ok := d.GetOkExists("province"); ok {
		v := _v.(string)
		optArgs.St = &v
	}
	if _v, ok := d.GetOkExists("postal_code"); ok {
		v := _v.(string)
		optArgs.Postalcode = &v
	}
	if _v, ok := d.GetOkExists("organisation_unit"); ok {
		v := _v.(string)
		optArgs.Ou = &v
	}
	if _v, ok := d.GetOkExists("job_title"); ok {
		v := _v.(string)
		optArgs.Title = &v
	}
	if _v, ok := d.GetOkExists("manager"); ok {
		v := _v.(string)
		optArgs.Manager = &v
	}
	if _v, ok := d.GetOkExists("employee_number"); ok {
		v := _v.(string)
		optArgs.Employeenumber = &v
	}
	if _v, ok := d.GetOkExists("employee_type"); ok {
		v := _v.(string)
		optArgs.Employeetype = &v
	}
	if _v, ok := d.GetOkExists("preferred_language"); ok {
		v := _v.(string)
		optArgs.Preferredlanguage = &v
	}
	if _v, ok := d.GetOkExists("account_disabled"); ok {
		v := _v.(bool)
		optArgs.Nsaccountlock = &v
	}
	if _v, ok := d.GetOkExists("ssh_public_key"); ok {
		v := utilsGetArry(_v.([]interface{}))
		optArgs.Ipasshpubkey = &v
	}
	if _v, ok := d.GetOkExists("car_license"); ok {
		v := utilsGetArry(_v.([]interface{}))
		optArgs.Carlicense = &v
	}
	if _v, ok := d.GetOkExists("krb_principal_expiration"); ok {
		v := _v.(string)
		timestamp, err := time.Parse(time.RFC3339, v)
		if err != nil {
			return diag.Errorf("The krb_principal_expiration timestamp could not be parsed as RFC3339: %s", err)
		}
		optArgs.Krbprincipalexpiration = &timestamp
	}
	if _v, ok := d.GetOkExists("krb_password_expiration"); ok {
		v := _v.(string)
		timestamp, err := time.Parse(time.RFC3339, v)
		if err != nil {
			return diag.Errorf("The krb_password_expiration timestamp could not be parsed as RFC3339: %s", err)
		}
		optArgs.Krbpasswordexpiration = &timestamp
	}
	if _v, ok := d.GetOkExists("userclass"); ok {
		v := utilsGetArry(_v.([]interface{}))
		optArgs.Userclass = &v
	}

	_, err = client.UserAdd(&args, &optArgs)
	if err != nil {
		return diag.Errorf("Error creating freeipa user: %s", err)
	}

	d.SetId(d.Get("name").(string))

	return resourceFreeIPADNSUserRead(ctx, d, meta)
}

func resourceFreeIPADNSUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Read freeipa user")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	all := true
	optArgs := ipa.UserShowOptionalArgs{
		All: &all,
	}

	if _v, ok := d.GetOkExists("name"); ok {
		v := _v.(string)
		optArgs.UID = &v
	} else {
		var tempId = d.Id()
		optArgs.UID = &tempId
		d.Set("name", d.Id())
	}

	log.Printf("[DEBUG] Read freeipa user %s", d.Id())
	res, err := client.UserShow(&ipa.UserShowArgs{}, &optArgs)
	if err != nil {
		if strings.Contains(err.Error(), "NotFound") {
			d.SetId("")
			log.Printf("[DEBUG] User not found")
			return nil
		} else {
			return diag.Errorf("Error reading freeipa user: %s", err)
		}
	}

	log.Printf("[DEBUG] Read freeipa user %s", res.Result.UID)

	return nil
}

func resourceFreeIPADNSUserUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Update freeipa user")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}
	var hasChange = false
	optArgs := ipa.UserModOptionalArgs{}

	if _v, ok := d.GetOkExists("name"); ok {
		v := _v.(string)
		optArgs.UID = &v
	}

	if d.HasChange("full_name") {
		if _v, ok := d.GetOkExists("full_name"); ok {
			v := _v.(string)
			if v != "" {
				optArgs.Cn = &v
				hasChange = true
			}
		}
	}
	if d.HasChange("first_name") {
		if _v, ok := d.GetOkExists("first_name"); ok {
			v := _v.(string)
			optArgs.Givenname = &v
			hasChange = true
		}
	}
	if d.HasChange("last_name") {
		if _v, ok := d.GetOkExists("last_name"); ok {
			v := _v.(string)
			optArgs.Sn = &v
			hasChange = true
		}
	}
	if d.HasChange("display_name") {
		if _v, ok := d.GetOkExists("display_name"); ok {
			v := _v.(string)
			optArgs.Displayname = &v
			hasChange = true
		}
	}
	if d.HasChange("initials") {
		if _v, ok := d.GetOkExists("initials"); ok {
			v := _v.(string)
			optArgs.Initials = &v
			hasChange = true
		}
	}
	if d.HasChange("home_directory") {
		if _v, ok := d.GetOkExists("home_directory"); ok {
			v := _v.(string)
			if v != "" {
				optArgs.Homedirectory = &v
				hasChange = true
			}
		}
	}
	if d.HasChange("gecos") {
		if _v, ok := d.GetOkExists("gecos"); ok {
			v := _v.(string)
			optArgs.Gecos = &v
			hasChange = true
		}
	}
	if d.HasChange("login_shell") {
		if _v, ok := d.GetOkExists("login_shell"); ok {
			v := _v.(string)
			if v != "" {
				optArgs.Loginshell = &v
				hasChange = true
			}
		}
	}
	if d.HasChange("userpassword") {
		if _v, ok := d.GetOkExists("userpassword"); ok {
			v := _v.(string)
			optArgs.Userpassword = &v
			hasChange = true
		}
	}
	if d.HasChange("random_password") {
		if _v, ok := d.GetOkExists("random_password"); ok {
			v := _v.(bool)
			optArgs.Random = &v
			hasChange = true
		}
	}
	if d.HasChange("uid_number") {
		if _v, ok := d.GetOkExists("uid_number"); ok {
			v := _v.(int)
			if v != 0 {
				optArgs.Uidnumber = &v
				hasChange = true
			}
		}
	}
	if d.HasChange("gid_number") {
		if _v, ok := d.GetOkExists("gid_number"); ok {
			v := _v.(int)
			if v != 0 {
				optArgs.Gidnumber = &v
				hasChange = true
			}
		}
	}
	if d.HasChange("street_address") {
		if _v, ok := d.GetOkExists("street_address"); ok {
			v := _v.(string)
			optArgs.Street = &v
			hasChange = true
		}
	}
	if d.HasChange("city") {
		if _v, ok := d.GetOkExists("city"); ok {
			v := _v.(string)
			optArgs.L = &v
			hasChange = true
		}
	}
	if d.HasChange("province") {
		if _v, ok := d.GetOkExists("province"); ok {
			v := _v.(string)
			optArgs.St = &v
			hasChange = true
		}
	}
	if d.HasChange("postal_code") {
		if _v, ok := d.GetOkExists("postal_code"); ok {
			v := _v.(string)
			optArgs.Postalcode = &v
			hasChange = true
		}
	}
	if d.HasChange("organisation_unit") {
		if _v, ok := d.GetOkExists("organisation_unit"); ok {
			v := _v.(string)
			optArgs.Ou = &v
			hasChange = true
		}
	}
	if d.HasChange("job_title") {
		if _v, ok := d.GetOkExists("job_title"); ok {
			v := _v.(string)
			optArgs.Title = &v
			hasChange = true
		}
	}
	if d.HasChange("employee_number") {
		if _v, ok := d.GetOkExists("employee_number"); ok {
			v := _v.(string)
			optArgs.Employeenumber = &v
			hasChange = true
		}
	}
	if d.HasChange("employee_type") {
		if _v, ok := d.GetOkExists("employee_type"); ok {
			v := _v.(string)
			optArgs.Employeetype = &v
			hasChange = true
		}
	}
	if d.HasChange("preferred_language") {
		if _v, ok := d.GetOkExists("preferred_language"); ok {
			v := _v.(string)
			optArgs.Preferredlanguage = &v
			hasChange = true
		}
	}
	if d.HasChange("account_disabled") {
		if _v, ok := d.GetOkExists("account_disabled"); ok {
			v := _v.(bool)
			optArgs.Nsaccountlock = &v
			hasChange = true
		}
	}
	if d.HasChange("telephone_numbers") {
		if _v, ok := d.GetOkExists("telephone_numbers"); ok {
			v := utilsGetArry(_v.([]interface{}))
			optArgs.Telephonenumber = &v
			hasChange = true
		}
	}
	if d.HasChange("mobile_numbers") {
		if _v, ok := d.GetOkExists("mobile_numbers"); ok {
			v := utilsGetArry(_v.([]interface{}))
			optArgs.Mobile = &v
			hasChange = true
		}
	}
	if d.HasChange("krb_principal_name") {
		if _v, ok := d.GetOkExists("krb_principal_name"); ok {
			v := utilsGetArry(_v.([]interface{}))
			optArgs.Krbprincipalname = &v
			hasChange = true
		}
	}
	if d.HasChange("ssh_public_key") {
		if _v, ok := d.GetOkExists("ssh_public_key"); ok {
			v := utilsGetArry(_v.([]interface{}))
			optArgs.Ipasshpubkey = &v
			hasChange = true
		}
	}
	if d.HasChange("car_license") {
		if _v, ok := d.GetOkExists("car_license"); ok {
			v := utilsGetArry(_v.([]interface{}))
			optArgs.Carlicense = &v
			hasChange = true
		}
	}
	if d.HasChange("email_address") {
		if _v, ok := d.GetOkExists("email_address"); ok {
			v := utilsGetArry(_v.([]interface{}))
			optArgs.Mail = &v
			hasChange = true
		}
	}
	if d.HasChange("krb_principal_expiration") {
		if _v, ok := d.GetOkExists("krb_principal_expiration"); ok {
			v := _v.(string)
			if v != "" {
				timestamp, err := time.Parse(time.RFC3339, v)
				if err != nil {
					return diag.Errorf("The krb_principal_expiration timestamp could not be parsed as RFC3339: %s", err)
				}
				optArgs.Krbprincipalexpiration = &timestamp
				hasChange = true
			}
		}
	}
	if d.HasChange("krb_password_expiration") {
		if _v, ok := d.GetOkExists("krb_password_expiration"); ok {
			v := _v.(string)
			if v != "" {
				timestamp, err := time.Parse(time.RFC3339, v)
				if err != nil {
					return diag.Errorf("The krb_password_expiration timestamp could not be parsed as RFC3339: %s", err)
				}
				optArgs.Krbpasswordexpiration = &timestamp
				hasChange = true
			}
		}
	}
	if d.HasChange("userclass") {
		if _v, ok := d.GetOkExists("userclass"); ok {
			v := utilsGetArry(_v.([]interface{}))
			optArgs.Userclass = &v
			hasChange = true
		}
	}

	if hasChange {
		_, err = client.UserMod(&ipa.UserModArgs{}, &optArgs)
		if err != nil {
			if strings.Contains(err.Error(), "EmptyModlist") {
				log.Printf("[DEBUG] EmptyModlist (4202): no modifications to be performed")
			} else {
				return diag.Errorf("Error update freeipa user: %s", err)
			}
		}
	}

	return resourceFreeIPADNSUserRead(ctx, d, meta)
}

func resourceFreeIPADNSUserDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Delete freeipa user")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}
	optArgs := ipa.UserDelOptionalArgs{}

	if _v, ok := d.GetOkExists("name"); ok {
		v := []string{_v.(string)}
		optArgs.UID = &v
	}
	_, err = client.UserDel(&ipa.UserDelArgs{}, &optArgs)
	if err != nil {
		return diag.Errorf("Error delete freeipa user: %s", err)
	}

	d.SetId("")

	return nil
}
