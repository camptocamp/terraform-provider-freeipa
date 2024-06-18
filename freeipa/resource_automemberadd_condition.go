package freeipa

import (
	"context"
	"log"
	"strings"

	ipa "github.com/camptocamp/go-freeipa/freeipa"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceFreeIPAAutomemberaddCondition() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFreeIPAAutomemberaddConditionCreate,
		ReadContext:   resourceFreeIPAAutomemberaddConditionRead,
		DeleteContext: resourceFreeIPAAutomemberaddConditionDelete,
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
				ForceNew: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"inclusiveregex": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"exclusiveregex": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceFreeIPAAutomemberaddConditionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Creating freeipa membership condition")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	optArgs := ipa.AutomemberAddConditionOptionalArgs{}

	args := ipa.AutomemberAddConditionArgs{
		Cn:   d.Get("name").(string),
		Type: d.Get("type").(string),
		Key:  d.Get("key").(string),
	}
	if _v, ok := d.GetOkExists("description"); ok {
		v := _v.(string)
		optArgs.Description = &v
	}
	if _v, ok := d.GetOk("inclusiveregex"); ok {
		v := make([]string, len(_v.([]interface{})))
		for i, value := range _v.([]interface{}) {
			v[i] = value.(string)
		}
		optArgs.Automemberinclusiveregex = &v
	}
	if _v, ok := d.GetOk("exclusiveregex"); ok {
		v := make([]string, len(_v.([]interface{})))
		for i, value := range _v.([]interface{}) {
			v[i] = value.(string)
		}
		optArgs.Automemberexclusiveregex = &v
	}
	_, err = client.AutomemberAddCondition(&args, &optArgs)
	if err != nil {
		return diag.Errorf("Error creating freeipa automember condition: %s", err)
	}

	d.SetId(d.Get("name").(string))
	d.Set("type", d.Get("type").(string))
	d.Set("key", d.Get("key").(string))

	return resourceFreeIPAAutomemberaddConditionRead(ctx, d, meta)
}

func resourceFreeIPAAutomemberaddConditionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Read freeipa membership condition")

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
			log.Printf("[DEBUG] AutomemberaddCondition not found")
			return nil
		} else {
			return diag.Errorf("Error reading freeipa group: %s", err)
		}
	}

	log.Printf("[DEBUG] Read freeipa automember %s", res.Result.Cn)
	return nil
}

func resourceFreeIPAAutomemberaddConditionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Delete freeipa automember condition")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}
	optArgs := ipa.AutomemberRemoveConditionOptionalArgs{}
	if _v, ok := d.GetOk("inclusiveregex"); ok {
		v := make([]string, len(_v.([]interface{})))
		for i, value := range _v.([]interface{}) {
			v[i] = value.(string)
		}
		optArgs.Automemberinclusiveregex = &v
	}
	if _v, ok := d.GetOk("exclusiveregex"); ok {
		v := make([]string, len(_v.([]interface{})))
		for i, value := range _v.([]interface{}) {
			v[i] = value.(string)
		}
		optArgs.Automemberexclusiveregex = &v
	}
	args := ipa.AutomemberRemoveConditionArgs{
		Cn:   d.Id(),
		Key:  d.Get("key").(string),
		Type: d.Get("type").(string),
	}
	_, err = client.AutomemberRemoveCondition(&args, &optArgs)
	if err != nil {
		return diag.Errorf("Error delete freeipa automember: %s", err)
	}

	d.SetId("")
	return nil
}
