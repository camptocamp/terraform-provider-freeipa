package freeipa

import (
	"context"
	"log"
	"strings"

	ipa "github.com/camptocamp/go-freeipa/freeipa"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceFreeIPAHostGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFreeIPADNSHostGroupCreate,
		ReadContext:   resourceFreeIPADNSHostGroupRead,
		UpdateContext: resourceFreeIPADNSHostGroupUpdate,
		DeleteContext: resourceFreeIPADNSHostGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Hostgroup's name",
			},
			"description": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Description: "A description of this hostgroup",
			},
		},
	}
}

func resourceFreeIPADNSHostGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Creating freeipa hostgroup")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	optArgs := ipa.HostgroupAddOptionalArgs{}

	args := ipa.HostgroupAddArgs{
		Cn: d.Get("name").(string),
	}

	if _v, ok := d.GetOkExists("description"); ok {
		v := _v.(string)
		optArgs.Description = &v
	}
	_, err = client.HostgroupAdd(&args, &optArgs)
	if err != nil {
		return diag.Errorf("Error creating freeipa hostgroup: %s", err)
	}

	d.SetId(d.Get("name").(string))

	return resourceFreeIPADNSHostGroupRead(ctx, d, meta)
}

func resourceFreeIPADNSHostGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Read freeipa hostgroup")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	all := true
	args := ipa.HostgroupShowArgs{
		Cn: d.Get("name").(string),
	}
	optArgs := ipa.HostgroupShowOptionalArgs{
		All: &all,
	}

	res, err := client.HostgroupShow(&args, &optArgs)
	if err != nil {
		if strings.Contains(err.Error(), "NotFound") {
			d.SetId("")
			log.Printf("[DEBUG] Hostgroup not found")
			return nil
		} else {
			return diag.Errorf("Error reading freeipa hostgroup: %s", err)
		}
	}

	log.Printf("[DEBUG] Read freeipa hostgroup %s", res.Result.Cn)

	return nil
}

func resourceFreeIPADNSHostGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Update freeipa hostgroup")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}
	var hasChange = false
	args := ipa.HostgroupModArgs{
		Cn: d.Get("name").(string),
	}
	optArgs := ipa.HostgroupModOptionalArgs{}

	if d.HasChange("description") {
		if _v, ok := d.GetOkExists("description"); ok {
			v := _v.(string)
			if v != "" {
				optArgs.Description = &v
				hasChange = true
			}
		}
	}
	if hasChange {
		_, err = client.HostgroupMod(&args, &optArgs)
		if err != nil {
			if strings.Contains(err.Error(), "EmptyModlist") {
				log.Printf("[DEBUG] EmptyModlist (4202): no modifications to be performed")
			} else {
				return diag.Errorf("Error update freeipa hostgroup: %s", err)
			}
		}
	}

	return resourceFreeIPADNSHostGroupRead(ctx, d, meta)
}

func resourceFreeIPADNSHostGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Delete freeipa hostgroup")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}
	args := ipa.HostgroupDelArgs{
		Cn: []string{d.Get("name").(string)},
	}
	optArgs := ipa.HostgroupDelOptionalArgs{}

	_, err = client.HostgroupDel(&args, &optArgs)
	if err != nil {
		return diag.Errorf("Error delete freeipa host: %s", err)
	}

	d.SetId("")

	return nil
}
