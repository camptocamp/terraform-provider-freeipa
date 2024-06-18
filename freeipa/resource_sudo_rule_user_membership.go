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

func resourceFreeIPASudoRuleUserMembership() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFreeIPASudoRuleUserMembershipCreate,
		ReadContext:   resourceFreeIPASudoRuleUserMembershipRead,
		DeleteContext: resourceFreeIPASudoRuleUserMembershipDelete,
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
			"user": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"group"},
				Description:   "User to add to the sudo rule",
			},
			"group": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"user"},
				Description:   "Group to add to the sudo rule",
			},
		},
	}
}

func resourceFreeIPASudoRuleUserMembershipCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Creating freeipa sudo rule user membership")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	user_id := "sru"

	optArgs := ipa.SudoruleAddUserOptionalArgs{}

	args := ipa.SudoruleAddUserArgs{
		Cn: d.Get("name").(string),
	}
	if _v, ok := d.GetOkExists("user"); ok {
		v := []string{_v.(string)}
		optArgs.User = &v
		user_id = "sru"
	}
	if _v, ok := d.GetOkExists("group"); ok {
		v := []string{_v.(string)}
		optArgs.Group = &v
		user_id = "srug"
	}

	_, err = client.SudoruleAddUser(&args, &optArgs)
	if err != nil {
		return diag.Errorf("Error creating freeipa sudo rule user membership: %s", err)
	}

	switch user_id {
	case "sru":
		id := fmt.Sprintf("%s/%s/%s", d.Get("name").(string), user_id, d.Get("user").(string))
		d.SetId(id)
	case "srug":
		id := fmt.Sprintf("%s/%s/%s", d.Get("name").(string), user_id, d.Get("group").(string))
		d.SetId(id)
	}

	return resourceFreeIPASudoRuleUserMembershipRead(ctx, d, meta)
}

func resourceFreeIPASudoRuleUserMembershipRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Read freeipa sudo rule user membership")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	sudoruleId, typeId, user_id, err := parseSudoRuleUserMembershipID(d.Id())

	if err != nil {
		return diag.Errorf("Error parsing ID of freeipa_sudo_rule_user_membership: %s", err)
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
	case "sru":
		if res.Result.MemberuserUser == nil || !slices.Contains(*res.Result.MemberuserUser, user_id) {
			log.Printf("[DEBUG] Warning! Sudo rule user membership does not exist")
			d.Set("name", "")
			d.Set("user", "")
			d.Set("group", "")
			d.SetId("")
			return nil
		}
	case "srug":
		if res.Result.MemberuserGroup == nil || !slices.Contains(*res.Result.MemberuserGroup, user_id) {
			log.Printf("[DEBUG] Warning! Sudo rule group membership does not exist")
			d.Set("name", "")
			d.Set("user", "")
			d.Set("group", "")
			d.SetId("")
			return nil
		}
	}

	if err != nil {
		return diag.Errorf("Error show freeipa sudo rule user membership: %s", err)
	}

	log.Printf("[DEBUG] Read freeipa sudo rule user membership %s", res.Result.Cn)
	return nil
}

func resourceFreeIPASudoRuleUserMembershipDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Delete freeipa sudo rule user membership")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	sudoruleId, typeId, user_id, err := parseSudoRuleUserMembershipID(d.Id())

	if err != nil {
		return diag.Errorf("Error parsing ID of freeipa_rule_user_membership: %s", err)
	}

	optArgs := ipa.SudoruleRemoveUserOptionalArgs{}

	args := ipa.SudoruleRemoveUserArgs{
		Cn: sudoruleId,
	}

	switch typeId {
	case "sru":
		v := []string{user_id}
		optArgs.User = &v
	case "srug":
		v := []string{user_id}
		optArgs.Group = &v
	}

	_, err = client.SudoruleRemoveUser(&args, &optArgs)
	if err != nil {
		return diag.Errorf("Error delete freeipa sudo rule user membership: %s", err)
	}

	d.SetId("")
	return nil
}

func parseSudoRuleUserMembershipID(id string) (string, string, string, error) {
	idParts := strings.SplitN(id, "/", 3)
	if len(idParts) < 3 {
		return "", "", "", fmt.Errorf("Unable to determine sudo rule user membership ID %s", id)
	}

	name := idParts[0]
	_type := idParts[1]
	user := idParts[2]

	return name, _type, user, nil
}
