package freeipa

import (
	"context"
	"log"
	"strings"

	ipa "github.com/camptocamp/go-freeipa/freeipa"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceFreeIPASudoRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFreeIPASudoRuleCreate,
		ReadContext:   resourceFreeIPASudoRuleRead,
		UpdateContext: resourceFreeIPASudoRuleUpdate,
		DeleteContext: resourceFreeIPASudoRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the sudo rule",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Sudo rule description",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				ForceNew:    false,
				Description: "Enable this sudo rule",
			},
			"usercategory": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    false,
				Description: "User category the sudo rule is applied to (allowed value: all)",
			},
			"hostcategory": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    false,
				Description: "Host category the sudo rule is applied to (allowed value: all)",
			},
			"commandcategory": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    false,
				Description: "Command category the sudo rule is applied to (allowed value: all)",
			},
			"runasusercategory": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    false,
				Description: "Run as user category the sudo rule is applied to (allowed value: all)",
			},
			"runasgroupcategory": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    false,
				Description: "Run as group category the sudo rule is applied to (allowed value: all)",
			},
			"order": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    false,
				Description: "Sudo rule order (must be unique)",
			},
		},
	}
}

func resourceFreeIPASudoRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Creating freeipa sudo rule")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	optArgs := ipa.SudoruleAddOptionalArgs{}

	args := ipa.SudoruleAddArgs{
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
	if _v, ok := d.GetOkExists("runasusercategory"); ok {
		v := _v.(string)
		optArgs.Ipasudorunasusercategory = &v
	}
	if _v, ok := d.GetOkExists("commandcategory"); ok {
		v := _v.(string)
		optArgs.Cmdcategory = &v
	}
	if _v, ok := d.GetOkExists("runasgroupcategory"); ok {
		v := _v.(string)
		optArgs.Ipasudorunasgroupcategory = &v
	}
	if _v, ok := d.GetOkExists("order"); ok {
		v := _v.(int)
		optArgs.Sudoorder = &v
	}
	_, err = client.SudoruleAdd(&args, &optArgs)
	if err != nil {
		return diag.Errorf("Error creating freeipa sudo rule: %s", err)
	}

	d.SetId(d.Get("name").(string))

	return resourceFreeIPASudoRuleRead(ctx, d, meta)
}

func resourceFreeIPASudoRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Read freeipa sudo rule")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	all := true
	optArgs := ipa.SudoruleShowOptionalArgs{
		All: &all,
	}

	args := ipa.SudoruleShowArgs{
		Cn: d.Id(),
	}

	res, err := client.SudoruleShow(&args, &optArgs)
	if err != nil {
		if strings.Contains(err.Error(), "NotFound") {
			d.SetId("")
			log.Printf("[DEBUG] Sudo rule not found")
			return nil
		} else {
			return diag.Errorf("Error reading freeipa sudo rule: %s", err)
		}
	}

	log.Printf("[DEBUG] Read freeipa sudo rule %s", res.Result.Cn)
	return nil
}

func resourceFreeIPASudoRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Update freeipa sudo rule")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	args := ipa.SudoruleModArgs{
		Cn: d.Id(),
	}
	optArgs := ipa.SudoruleModOptionalArgs{}

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
	if d.HasChange("runasusercategory") {
		if _v, ok := d.GetOkExists("runasusercategory"); ok {
			v := _v.(string)
			optArgs.Ipasudorunasusercategory = &v
			hasChange = true
		}
	}
	if d.HasChange("commandcategory") {
		if _v, ok := d.GetOkExists("commandcategory"); ok {
			v := _v.(string)
			optArgs.Cmdcategory = &v
			hasChange = true
		}
	}
	if d.HasChange("runasgroupcategory") {
		if _v, ok := d.GetOkExists("runasgroupcategory"); ok {
			v := _v.(string)
			optArgs.Ipasudorunasgroupcategory = &v
			hasChange = true
		}
	}
	if d.HasChange("order") {
		if _v, ok := d.GetOkExists("order"); ok {
			v := _v.(int)
			optArgs.Sudoorder = &v
			hasChange = true
		}
	}

	// TODO: Change No-Posix, Posix, External

	if hasChange {
		_, err = client.SudoruleMod(&args, &optArgs)
		if err != nil {
			if strings.Contains(err.Error(), "EmptyModlist") {
				log.Printf("[DEBUG] EmptyModlist (4202): no modifications to be performed")
			} else {
				return diag.Errorf("Error update freeipa sudo rule: %s", err)
			}
		}
	}

	d.SetId(d.Get("name").(string))

	return resourceFreeIPASudoRuleRead(ctx, d, meta)
}

func resourceFreeIPASudoRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Delete freeipa sudo rule")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}
	args := ipa.SudoruleDelArgs{
		Cn: []string{d.Id()},
	}
	_, err = client.SudoruleDel(&args, &ipa.SudoruleDelOptionalArgs{})
	if err != nil {
		return diag.Errorf("Error delete freeipa sudo rule: %s", err)
	}

	d.SetId("")
	return nil
}
