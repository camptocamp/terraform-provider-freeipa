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

func resourceFreeIPASudoRuleDenyCommandMembership() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFreeIPASudoRuleDenyCommandMembershipCreate,
		ReadContext:   resourceFreeIPASudoRuleDenyCommandMembershipRead,
		DeleteContext: resourceFreeIPASudoRuleDenyCommandMembershipDelete,
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
			"sudocmd": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"sudocmd_group"},
				Description:   "Sudo command to deny by the sudo rule",
			},
			"sudocmd_group": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"sudocmd"},
				Description:   "Sudo command group to deny by the sudo rule",
			},
		},
	}
}

func resourceFreeIPASudoRuleDenyCommandMembershipCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Creating freeipa sudo rule denied command membership")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	cmd_id := "srdc"

	optArgs := ipa.SudoruleAddDenyCommandOptionalArgs{}

	args := ipa.SudoruleAddDenyCommandArgs{
		Cn: d.Get("name").(string),
	}
	if _v, ok := d.GetOkExists("sudocmd"); ok {
		v := []string{_v.(string)}
		optArgs.Sudocmd = &v
		cmd_id = "srdc"
	}
	if _v, ok := d.GetOkExists("sudocmd_group"); ok {
		v := []string{_v.(string)}
		optArgs.Sudocmdgroup = &v
		cmd_id = "srdcg"
	}

	_, err = client.SudoruleAddDenyCommand(&args, &optArgs)
	if err != nil {
		return diag.Errorf("Error creating freeipa sudo rule denied command membership: %s", err)
	}

	switch cmd_id {
	case "srdc":
		id := fmt.Sprintf("%s/%s/%s", d.Get("name").(string), cmd_id, d.Get("sudocmd").(string))
		d.SetId(id)
	case "srdcg":
		id := fmt.Sprintf("%s/%s/%s", d.Get("name").(string), cmd_id, d.Get("sudocmd_group").(string))
		d.SetId(id)
	}

	return resourceFreeIPASudoRuleDenyCommandMembershipRead(ctx, d, meta)
}

func resourceFreeIPASudoRuleDenyCommandMembershipRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Read freeipa sudo rule denied command membership")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	sudoruleId, typeId, cmdId, err := parseSudoRuleDenyCommandMembershipID(d.Id())

	if err != nil {
		return diag.Errorf("Error parsing ID of freeipa_sudocmdgroup_membership: %s", err)
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
	case "srdc":
		if res.Result.MemberdenycmdSudocmd == nil || !slices.Contains(*res.Result.MemberdenycmdSudocmd, cmdId) {
			log.Printf("[DEBUG] Warning! Sudo rule denied command membership does not exist")
			d.Set("name", "")
			d.Set("sudocmd", "")
			d.Set("sudocmd_group", "")
			d.SetId("")
			return nil
		}
	case "srdcg":
		if res.Result.MemberdenycmdSudocmdgroup == nil || !slices.Contains(*res.Result.MemberdenycmdSudocmdgroup, cmdId) {
			log.Printf("[DEBUG] Warning! Sudo rule denied command membership does not exist")
			d.Set("name", "")
			d.Set("sudocmd", "")
			d.Set("sudocmd_group", "")
			d.SetId("")
			return nil
		}
	}

	if err != nil {
		return diag.Errorf("Error show freeipa sudo rule denied command membership: %s", err)
	}

	log.Printf("[DEBUG] Read freeipa sudo rule denied command membership %s", res.Result.Cn)
	return nil
}

func resourceFreeIPASudoRuleDenyCommandMembershipDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Delete freeipa sudo rule denied command membership")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	sudoruleId, typeId, cmdId, err := parseSudoRuleDenyCommandMembershipID(d.Id())

	if err != nil {
		return diag.Errorf("Error parsing ID of freeipa_sudocmdgroup_membership: %s", err)
	}

	optArgs := ipa.SudoruleRemoveDenyCommandOptionalArgs{}

	args := ipa.SudoruleRemoveDenyCommandArgs{
		Cn: sudoruleId,
	}

	switch typeId {
	case "srdc":
		v := []string{cmdId}
		optArgs.Sudocmd = &v
	case "srdcg":
		v := []string{cmdId}
		optArgs.Sudocmdgroup = &v
	}

	_, err = client.SudoruleRemoveDenyCommand(&args, &optArgs)
	if err != nil {
		return diag.Errorf("Error delete freeipa sudo rule denied command membership: %s", err)
	}

	d.SetId("")
	return nil
}

func parseSudoRuleDenyCommandMembershipID(id string) (string, string, string, error) {
	idParts := strings.SplitN(id, "/", 3)
	if len(idParts) < 3 {
		return "", "", "", fmt.Errorf("Unable to determine sudo rule denied command membership ID %s", id)
	}

	name := idParts[0]
	_type := idParts[1]
	sudocmd := idParts[2]

	return name, _type, sudocmd, nil
}
