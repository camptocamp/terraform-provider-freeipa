package freeipa

import (
	"context"
	"log"
	"strings"

	ipa "github.com/camptocamp/go-freeipa/freeipa"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceFreeIPAHBACPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFreeIPADNSHBACPolicyCreate,
		ReadContext:   resourceFreeIPADNSHBACPolicyRead,
		UpdateContext: resourceFreeIPADNSHBACPolicyUpdate,
		DeleteContext: resourceFreeIPADNSHBACPolicyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "HBAC policy name",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    false,
				Description: "HBAC policy description",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				ForceNew:    false,
				Description: "Enable this policy (Defaults to `true`)",
			},
			"usercategory": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    false,
				Description: "User category the policy is applied to (allowed value: `all`)",
			},
			"hostcategory": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    false,
				Description: "Host category the policy is applied to (allowed value: `all`)",
			},
			"servicecategory": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    false,
				Description: "Service category the policy is applied to (allowed value: `all`)",
			},
		},
	}
}

func resourceFreeIPADNSHBACPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Creating freeipa the HBAC policy")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	optArgs := ipa.HbacruleAddOptionalArgs{}

	args := ipa.HbacruleAddArgs{
		Cn: d.Get("name").(string),
	}
	if _v, ok := d.GetOkExists("description"); ok {
		v := _v.(string)
		optArgs.Description = &v
	}
	if _v, ok := d.GetOkExists("enabled"); ok {
		v := _v.(bool)
		optArgs.Ipaenabledflag = &v
	}
	if _v, ok := d.GetOkExists("usercategory"); ok {
		v := _v.(string)
		optArgs.Usercategory = &v
	}
	if _v, ok := d.GetOkExists("hostcategory"); ok {
		v := _v.(string)
		optArgs.Hostcategory = &v
	}
	if _v, ok := d.GetOkExists("servicecategory"); ok {
		v := _v.(string)
		optArgs.Servicecategory = &v
	}

	_, err = client.HbacruleAdd(&args, &optArgs)
	if err != nil {
		return diag.Errorf("Error creating freeipa the HBAC policy: %s", err)
	}

	d.SetId(d.Get("name").(string))

	return resourceFreeIPADNSHBACPolicyRead(ctx, d, meta)
}

func resourceFreeIPADNSHBACPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Read freeipa the HBAC policy")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	all := true
	optArgs := ipa.HbacruleShowOptionalArgs{
		All: &all,
	}

	args := ipa.HbacruleShowArgs{
		Cn: d.Id(),
	}

	res, err := client.HbacruleShow(&args, &optArgs)
	if err != nil {
		if strings.Contains(err.Error(), "NotFound") {
			d.SetId("")
			log.Printf("[DEBUG] HBAC policy not found")
			return nil
		} else {
			return diag.Errorf("Error reading freeipa HBAC policy: %s", err)
		}
	}

	log.Printf("[DEBUG] Read freeipa HBAC policy %s", res.Result.Cn)
	return nil
}

func resourceFreeIPADNSHBACPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Update freeipa HBAC policy")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	args := ipa.HbacruleModArgs{
		Cn: d.Id(),
	}
	optArgs := ipa.HbacruleModOptionalArgs{}

	var hasChange = false

	if d.HasChange("description") {
		if _v, ok := d.GetOkExists("description"); ok {
			v := _v.(string)
			optArgs.Description = &v
			hasChange = true
		}
	}
	if d.HasChange("enabled") {
		if _v, ok := d.GetOkExists("enabled"); ok {
			v := _v.(bool)
			optArgs.Ipaenabledflag = &v
			hasChange = true
		}
	}
	if d.HasChange("usercategory") {
		if _v, ok := d.GetOkExists("usercategory"); ok {
			v := _v.(string)
			optArgs.Usercategory = &v
			hasChange = true
		}
	}
	if d.HasChange("hostcategory") {
		if _v, ok := d.GetOkExists("hostcategory"); ok {
			v := _v.(string)
			optArgs.Hostcategory = &v
			hasChange = true
		}
	}
	if d.HasChange("servicecategory") {
		if _v, ok := d.GetOkExists("servicecategory"); ok {
			v := _v.(string)
			optArgs.Servicecategory = &v
			hasChange = true
		}
	}

	if hasChange {
		_, err = client.HbacruleMod(&args, &optArgs)
		if err != nil {
			if strings.Contains(err.Error(), "EmptyModlist") {
				log.Printf("[DEBUG] EmptyModlist (4202): no modifications to be performed")
			} else {
				return diag.Errorf("Error update freeipa HBAC policy: %s", err)
			}
		}
	}

	d.SetId(d.Get("name").(string))

	return resourceFreeIPADNSHBACPolicyRead(ctx, d, meta)
}

func resourceFreeIPADNSHBACPolicyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Delete freeipa the HBAC policy")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	args := ipa.HbacruleDelArgs{
		Cn: []string{d.Id()},
	}
	_, err = client.HbacruleDel(&args, &ipa.HbacruleDelOptionalArgs{})
	if err != nil {
		return diag.Errorf("Error delete freeipa the HBAC policy: %s", err)
	}

	d.SetId("")

	return nil
}
