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

func resourceFreeIPASudoRuleRunAsGroupMembership() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFreeIPASudoRuleRunAsGroupMembershipCreate,
		ReadContext:   resourceFreeIPASudoRuleRunAsGroupMembershipRead,
		DeleteContext: resourceFreeIPASudoRuleRunAsGroupMembershipDelete,
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
			"runasgroup": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Run As Group to add to the sudo rule. Can be an external group (local group of ipa clients)",
			},
		},
	}
}

func resourceFreeIPASudoRuleRunAsGroupMembershipCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Creating freeipa sudo rule runasgroup membership")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	group_id := "srrau"

	optArgs := ipa.SudoruleAddRunasgroupOptionalArgs{}

	args := ipa.SudoruleAddRunasgroupArgs{
		Cn: d.Get("name").(string),
	}
	if _v, ok := d.GetOkExists("runasgroup"); ok {
		v := []string{_v.(string)}
		optArgs.Group = &v
		group_id = "srraug"
	}

	_, err = client.SudoruleAddRunasgroup(&args, &optArgs)
	if err != nil {
		return diag.Errorf("Error creating freeipa sudo rule runasgroup membership: %s", err)
	}

	switch group_id {
	case "srraug":
		id := fmt.Sprintf("%s/%s/%s", d.Get("name").(string), group_id, d.Get("runasgroup").(string))
		d.SetId(id)
	}

	return resourceFreeIPASudoRuleRunAsGroupMembershipRead(ctx, d, meta)
}

func resourceFreeIPASudoRuleRunAsGroupMembershipRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Read freeipa sudo rule user membership")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	sudoruleId, typeId, group_id, err := parseSudoRuleRunAsGroupMembershipID(d.Id())

	if err != nil {
		return diag.Errorf("Error parsing ID of freeipa_sudo_rule_runasgroup_membership: %s", err)
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
	case "srraug":
		if (res.Result.IpasudorunasgroupGroup == nil || !slices.Contains(*res.Result.IpasudorunasgroupGroup, group_id)) && (res.Result.Ipasudorunasextgroup == nil || !slices.Contains(*res.Result.Ipasudorunasextgroup, group_id)) {
			log.Printf("[DEBUG] Warning! Sudo rule user membership does not exist")
			d.Set("name", "")
			d.Set("runasgroup", "")
			d.Set("runasgroup", "")
			d.SetId("")
			return nil
		}
	}

	if err != nil {
		return diag.Errorf("Error show freeipa sudo rule runasgroup membership: %s", err)
	}

	log.Printf("[DEBUG] Read freeipa sudo rule user membership %s", res.Result.Cn)
	return nil
}

func resourceFreeIPASudoRuleRunAsGroupMembershipDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Delete freeipa sudo rule runasgroup membership")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	sudoruleId, typeId, group_id, err := parseSudoRuleRunAsGroupMembershipID(d.Id())

	if err != nil {
		return diag.Errorf("Error parsing ID of freeipa_rule_runasgroup_membership: %s", err)
	}

	optArgs := ipa.SudoruleRemoveRunasgroupOptionalArgs{}

	args := ipa.SudoruleRemoveRunasgroupArgs{
		Cn: sudoruleId,
	}

	switch typeId {
	case "srraug":
		v := []string{group_id}
		optArgs.Group = &v
	}

	_, err = client.SudoruleRemoveRunasgroup(&args, &optArgs)
	if err != nil {
		return diag.Errorf("Error delete freeipa sudo rule runasgroup membership: %s", err)
	}

	d.SetId("")
	return nil
}

func parseSudoRuleRunAsGroupMembershipID(id string) (string, string, string, error) {
	idParts := strings.SplitN(id, "/", 3)
	if len(idParts) < 3 {
		return "", "", "", fmt.Errorf("Unable to determine sudo rule runasgroup membership ID %s", id)
	}

	name := idParts[0]
	_type := idParts[1]
	group := idParts[2]

	return name, _type, group, nil
}
