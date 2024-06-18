package freeipa

import (
	"context"
	"log"
	"strings"

	ipa "github.com/camptocamp/go-freeipa/freeipa"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceFreeIPASudocmdgroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFreeIPASudocmdgroupCreate,
		ReadContext:   resourceFreeIPASudocmdgroupRead,
		UpdateContext: resourceFreeIPASudocmdgroupUpdate,
		DeleteContext: resourceFreeIPASudocmdgroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the sudo command group",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Sudo command group description",
			},
		},
	}
}

func resourceFreeIPASudocmdgroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Creating freeipa sudo command group")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	optArgs := ipa.SudocmdgroupAddOptionalArgs{}

	args := ipa.SudocmdgroupAddArgs{
		Cn: d.Get("name").(string),
	}
	if _v, ok := d.GetOkExists("description"); ok {
		v := _v.(string)
		optArgs.Description = &v
	}
	_, err = client.SudocmdgroupAdd(&args, &optArgs)
	if err != nil {
		return diag.Errorf("Error creating freeipa sudo command group: %s", err)
	}

	d.SetId(d.Get("name").(string))

	return resourceFreeIPASudocmdgroupRead(ctx, d, meta)
}

func resourceFreeIPASudocmdgroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Read freeipa sudo command group")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	all := true
	optArgs := ipa.SudocmdgroupShowOptionalArgs{
		All: &all,
	}

	args := ipa.SudocmdgroupShowArgs{
		Cn: d.Id(),
	}

	res, err := client.SudocmdgroupShow(&args, &optArgs)
	if err != nil {
		if strings.Contains(err.Error(), "NotFound") {
			d.SetId("")
			log.Printf("[DEBUG] Sudo command group not found")
			return nil
		} else {
			return diag.Errorf("Error reading freeipa sudo command group: %s", err)
		}
	}

	log.Printf("[DEBUG] Read freeipa sudo command group %s", res.Result.Cn)
	return nil
}

func resourceFreeIPASudocmdgroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Update freeipa sudo command group")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	args := ipa.SudocmdgroupModArgs{
		Cn: d.Id(),
	}
	optArgs := ipa.SudocmdgroupModOptionalArgs{}

	var hasChange = false

	if d.HasChange("description") {
		if _v, ok := d.GetOkExists("description"); ok {
			v := _v.(string)
			optArgs.Description = &v
			hasChange = true
		}
	}

	// TODO: Change No-Posix, Posix, External

	if hasChange {
		_, err = client.SudocmdgroupMod(&args, &optArgs)
		if err != nil {
			if strings.Contains(err.Error(), "EmptyModlist") {
				log.Printf("[DEBUG] EmptyModlist (4202): no modifications to be performed")
			} else {
				return diag.Errorf("Error update freeipa sudo command group: %s", err)
			}
		}
	}

	d.SetId(d.Get("name").(string))

	return resourceFreeIPASudocmdgroupRead(ctx, d, meta)
}

func resourceFreeIPASudocmdgroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Delete freeipa sudo command group")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}
	args := ipa.SudocmdgroupDelArgs{
		Cn: []string{d.Id()},
	}
	_, err = client.SudocmdgroupDel(&args, &ipa.SudocmdgroupDelOptionalArgs{})
	if err != nil {
		return diag.Errorf("Error delete freeipa sudo command group: %s", err)
	}

	d.SetId("")
	return nil
}
