package freeipa

import (
	"context"
	"log"
	"strings"

	ipa "github.com/camptocamp/go-freeipa/freeipa"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceFreeIPASudocmd() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFreeIPASudocmdCreate,
		ReadContext:   resourceFreeIPASudocmdRead,
		UpdateContext: resourceFreeIPASudocmdUpdate,
		DeleteContext: resourceFreeIPASudocmdDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Absolute path of the sudo command",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Sudocmd description",
			},
		},
	}
}

func resourceFreeIPASudocmdCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Creating freeipa sudo command")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	optArgs := ipa.SudocmdAddOptionalArgs{}

	args := ipa.SudocmdAddArgs{
		Sudocmd: d.Get("name").(string),
	}
	if _v, ok := d.GetOkExists("description"); ok {
		v := _v.(string)
		optArgs.Description = &v
	}
	_, err = client.SudocmdAdd(&args, &optArgs)
	if err != nil {
		return diag.Errorf("Error creating freeipa sudo command: %s", err)
	}

	d.SetId(d.Get("name").(string))

	return resourceFreeIPASudocmdRead(ctx, d, meta)
}

func resourceFreeIPASudocmdRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Read freeipa sudo command")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	all := true
	optArgs := ipa.SudocmdShowOptionalArgs{
		All: &all,
	}

	args := ipa.SudocmdShowArgs{
		Sudocmd: d.Id(),
	}

	res, err := client.SudocmdShow(&args, &optArgs)
	if err != nil {
		if strings.Contains(err.Error(), "NotFound") {
			d.SetId("")
			log.Printf("[DEBUG] Sudo command not found")
			return nil
		} else {
			return diag.Errorf("Error reading freeipa Sudo command: %s", err)
		}
	}

	log.Printf("[DEBUG] Read freeipa sudo command %s", res.Result.Sudocmd)
	return nil
}

func resourceFreeIPASudocmdUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Update freeipa sudo command")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	args := ipa.SudocmdModArgs{
		Sudocmd: d.Id(),
	}
	optArgs := ipa.SudocmdModOptionalArgs{}

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
		_, err = client.SudocmdMod(&args, &optArgs)
		if err != nil {
			if strings.Contains(err.Error(), "EmptyModlist") {
				log.Printf("[DEBUG] EmptyModlist (4202): no modifications to be performed")
			} else {
				return diag.Errorf("Error update freeipa sudo command: %s", err)
			}
		}
	}

	d.SetId(d.Get("name").(string))

	return resourceFreeIPASudocmdRead(ctx, d, meta)
}

func resourceFreeIPASudocmdDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Delete freeipa sudo command")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}
	args := ipa.SudocmdDelArgs{
		Sudocmd: []string{d.Id()},
	}
	_, err = client.SudocmdDel(&args, &ipa.SudocmdDelOptionalArgs{})
	if err != nil {
		return diag.Errorf("Error delete freeipa sudo command: %s", err)
	}

	d.SetId("")
	return nil
}
