package freeipa

import (
	"context"
	"fmt"
	"log"
	"strings"

	ipa "github.com/camptocamp/go-freeipa/freeipa"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/exp/slices"
)

func resourceFreeIPASudoRuleOption() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFreeIPASudoRuleOptionCreate,
		ReadContext:   resourceFreeIPASudoRuleOptionRead,
		DeleteContext: resourceFreeIPASudoRuleOptionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Sudo rule name",
			},
			"option": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Sudo option to add to the sudo rule.",
			},
		},
	}
}

func resourceFreeIPASudoRuleOptionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Creating freeipa sudo rule option")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	optArgs := ipa.SudoruleAddOptionOptionalArgs{}

	args := ipa.SudoruleAddOptionArgs{
		Cn:         d.Get("name").(string),
		Ipasudoopt: d.Get("option").(string),
	}

	_, err = client.SudoruleAddOption(&args, &optArgs)
	if err != nil {
		return diag.Errorf("Error creating freeipa sudo rule option: %s", err)
	}

	id := fmt.Sprintf("%s/%s/%s", d.Get("name").(string), "sro", d.Get("option").(string))
	d.SetId(id)

	return resourceFreeIPASudoRuleOptionRead(ctx, d, meta)
}

func resourceFreeIPASudoRuleOptionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Read freeipa sudo rule user membership")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	sudoruleId, typeId, opt_id, err := parseSudoRuleOptionID(d.Id())

	if err != nil {
		return diag.Errorf("Error parsing ID of freeipa_sudo_rule_option: %s", err)
	}

	all := true
	optArgs := ipa.SudoruleShowOptionalArgs{
		All: &all,
	}

	args := ipa.SudoruleShowArgs{
		Cn: sudoruleId,
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

	switch typeId {
	case "sro":
		if res.Result.Ipasudoopt == nil || !slices.Contains(*res.Result.Ipasudoopt, opt_id) {
			log.Printf("[DEBUG] Warning! Sudo rule option does not exist")
			d.Set("name", "")
			d.Set("option", "")
			d.SetId("")
			return nil
		}
	}

	log.Printf("[DEBUG] Read freeipa sudo rule option %s", res.Result.Cn)
	return nil
}

func resourceFreeIPASudoRuleOptionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Delete freeipa sudo rule option")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	sudoruleId, typeId, opt_id, err := parseSudoRuleOptionID(d.Id())

	if err != nil {
		return diag.Errorf("Error parsing ID of freeipa_rule_runasuser_membership: %s", err)
	}

	optArgs := ipa.SudoruleRemoveOptionOptionalArgs{}

	args := ipa.SudoruleRemoveOptionArgs{
		Cn: sudoruleId,
	}

	switch typeId {
	case "sro":
		args.Ipasudoopt = opt_id
	}

	_, err = client.SudoruleRemoveOption(&args, &optArgs)
	if err != nil {
		return diag.Errorf("Error delete freeipa sudo rule option: %s", err)
	}

	d.SetId("")
	return nil
}

func parseSudoRuleOptionID(id string) (string, string, string, error) {
	idParts := strings.SplitN(id, "/", 3)
	if len(idParts) < 3 {
		return "", "", "", fmt.Errorf("Unable to determine sudo rule option ID %s", id)
	}

	name := idParts[0]
	_type := idParts[1]
	user := idParts[2]

	return name, _type, user, nil
}
