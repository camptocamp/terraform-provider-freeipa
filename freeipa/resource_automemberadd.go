package freeipa

import (
	"context"
	"log"
	"strings"

	ipa "github.com/camptocamp/go-freeipa/freeipa"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceFreeIPAAutomemberadd() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFreeIPAAutomemberaddCreate,
		ReadContext:   resourceFreeIPAAutomemberaddRead,
		UpdateContext: resourceFreeIPAAutomemberaddUpdate,
		DeleteContext: resourceFreeIPAAutomemberaddDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"addattr": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"setattr": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceFreeIPAAutomemberaddCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Creating freeipa automember")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	optArgs := ipa.AutomemberAddOptionalArgs{}

	args := ipa.AutomemberAddArgs{
		Cn:   d.Get("name").(string),
		Type: d.Get("type").(string),
	}
	if _v, ok := d.GetOkExists("description"); ok {
		v := _v.(string)
		optArgs.Description = &v
	}
	if _v, ok := d.GetOk("addattr"); ok {
		v := make([]string, len(_v.([]interface{})))
		for i, value := range _v.([]interface{}) {
			v[i] = value.(string)
		}
		optArgs.Addattr = &v
	}
	if _v, ok := d.GetOk("setattr"); ok {
		v := make([]string, len(_v.([]interface{})))
		for i, value := range _v.([]interface{}) {
			v[i] = value.(string)
		}
		optArgs.Setattr = &v
	}
	_, err = client.AutomemberAdd(&args, &optArgs)
	if err != nil {
		return diag.Errorf("Error creating freeipa automember: %s", err)
	}

	d.SetId(d.Get("name").(string))
	d.Set("type", d.Get("type").(string))

	return resourceFreeIPAAutomemberaddRead(ctx, d, meta)
}

func resourceFreeIPAAutomemberaddRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Read freeipa automember")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	all := true
	optArgs := ipa.AutomemberShowOptionalArgs{
		All: &all,
	}

	args := ipa.AutomemberShowArgs{
		Cn:   d.Id(),
		Type: d.Get("type").(string),
	}

	log.Printf("[DEBUG] Read freeipa automember %s", d.Id())
	res, err := client.AutomemberShow(&args, &optArgs)
	if err != nil {
		if strings.Contains(err.Error(), "NotFound") {
			d.SetId("")
			log.Printf("[DEBUG] Automemberadd not found")
			return nil
		} else {
			return diag.Errorf("Error reading freeipa group: %s", err)
		}
	}

	log.Printf("[DEBUG] Read freeipa automember %s", res.Result.Cn)
	return nil
}

func resourceFreeIPAAutomemberaddUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Update freeipa automember")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	args := ipa.AutomemberModArgs{
		Cn:   d.Id(),
		Type: d.Get("type").(string),
	}
	optArgs := ipa.AutomemberModOptionalArgs{}

	var hasChange = false

	if d.HasChange("description") {
		if _v, ok := d.GetOkExists("description"); ok {
			v := _v.(string)
			optArgs.Description = &v
			hasChange = true
		}
	}
	if d.HasChange("addattr") {
		if _v, ok := d.GetOkExists("addattr"); ok {
			v := make([]string, len(_v.([]interface{})))
			for i, value := range _v.([]interface{}) {
				v[i] = value.(string)
			}
			optArgs.Addattr = &v
			hasChange = true
		}
	}
	if d.HasChange("setattr") {
		if _v, ok := d.GetOkExists("setattr"); ok {
			v := make([]string, len(_v.([]interface{})))
			for i, value := range _v.([]interface{}) {
				v[i] = value.(string)
			}
			optArgs.Setattr = &v
			hasChange = true
		}
	}

	if hasChange {
		_, err = client.AutomemberMod(&args, &optArgs)
		if err != nil {
			if strings.Contains(err.Error(), "EmptyModlist") {
				log.Printf("[DEBUG] EmptyModlist (4202): no modifications to be performed")
			} else {
				return diag.Errorf("Error update freeipa automember: %s", err)
			}
		}
	}

	d.SetId(d.Get("name").(string))
	d.Set("type", d.Get("type").(string))

	return resourceFreeIPAAutomemberaddRead(ctx, d, meta)
}

func resourceFreeIPAAutomemberaddDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Delete freeipa automember")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}
	args := ipa.AutomemberDelArgs{
		Cn:   []string{d.Id()},
		Type: d.Get("type").(string),
	}
	_, err = client.AutomemberDel(&args, &ipa.AutomemberDelOptionalArgs{})
	if err != nil {
		return diag.Errorf("Error delete freeipa automember: %s", err)
	}

	d.SetId("")
	return nil
}
