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

func resourceFreeIPASudoRuleRunAsUserMembership() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFreeIPASudoRuleRunAsUserMembershipCreate,
		ReadContext:   resourceFreeIPASudoRuleRunAsUserMembershipRead,
		DeleteContext: resourceFreeIPASudoRuleRunAsUserMembershipDelete,
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
			"runasuser": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Run As User to add to the sudo rule. Can be an external user (local user of ipa clients)",
			},
		},
	}
}

func resourceFreeIPASudoRuleRunAsUserMembershipCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Creating freeipa sudo rule runasuser membership")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	user_id := "srrau"

	optArgs := ipa.SudoruleAddRunasuserOptionalArgs{}

	args := ipa.SudoruleAddRunasuserArgs{
		Cn: d.Get("name").(string),
	}
	if _v, ok := d.GetOkExists("runasuser"); ok {
		v := []string{_v.(string)}
		optArgs.User = &v
		user_id = "srrau"
	}

	_, err = client.SudoruleAddRunasuser(&args, &optArgs)
	if err != nil {
		return diag.Errorf("Error creating freeipa sudo rule runasuser membership: %s", err)
	}

	switch user_id {
	case "srrau":
		id := fmt.Sprintf("%s/%s/%s", d.Get("name").(string), user_id, d.Get("runasuser").(string))
		d.SetId(id)
	}

	return resourceFreeIPASudoRuleRunAsUserMembershipRead(ctx, d, meta)
}

func resourceFreeIPASudoRuleRunAsUserMembershipRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Read freeipa sudo rule user membership")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	sudoruleId, typeId, user_id, err := parseSudoRuleRunAsUserMembershipID(d.Id())

	if err != nil {
		return diag.Errorf("Error parsing ID of freeipa_sudo_rule_runasuser_membership: %s", err)
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
	case "srrau":
		if (res.Result.IpasudorunasUser == nil || !slices.Contains(*res.Result.IpasudorunasUser, user_id)) && (res.Result.Ipasudorunasextuser == nil || !slices.Contains(*res.Result.Ipasudorunasextuser, user_id)) {
			log.Printf("[DEBUG] Warning! Sudo rule user membership does not exist")
			d.Set("name", "")
			d.Set("runasuser", "")
			d.Set("runasgroup", "")
			d.SetId("")
			return nil
		}
	}

	if err != nil {
		return diag.Errorf("Error show freeipa sudo rule runasuser membership: %s", err)
	}

	log.Printf("[DEBUG] Read freeipa sudo rule user membership %s", res.Result.Cn)
	return nil
}

func resourceFreeIPASudoRuleRunAsUserMembershipDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Delete freeipa sudo rule runasuser membership")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	sudoruleId, typeId, user_id, err := parseSudoRuleRunAsUserMembershipID(d.Id())

	if err != nil {
		return diag.Errorf("Error parsing ID of freeipa_rule_runasuser_membership: %s", err)
	}

	optArgs := ipa.SudoruleRemoveRunasuserOptionalArgs{}

	args := ipa.SudoruleRemoveRunasuserArgs{
		Cn: sudoruleId,
	}

	switch typeId {
	case "srrau":
		v := []string{user_id}
		optArgs.User = &v
	}

	_, err = client.SudoruleRemoveRunasuser(&args, &optArgs)
	if err != nil {
		return diag.Errorf("Error delete freeipa sudo rule runasuser membership: %s", err)
	}

	d.SetId("")
	return nil
}

func parseSudoRuleRunAsUserMembershipID(id string) (string, string, string, error) {
	idParts := strings.SplitN(id, "/", 3)
	if len(idParts) < 3 {
		return "", "", "", fmt.Errorf("Unable to determine sudo rule runasuser membership ID %s", id)
	}

	name := idParts[0]
	_type := idParts[1]
	user := idParts[2]

	return name, _type, user, nil
}
