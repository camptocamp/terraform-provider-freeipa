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

func resourceFreeIPASudoRuleAllowCommandMembership() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFreeIPASudoRuleAllowCommandMembershipCreate,
		ReadContext:   resourceFreeIPASudoRuleAllowCommandMembershipRead,
		DeleteContext: resourceFreeIPASudoRuleAllowCommandMembershipDelete,
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
				Description:   "Sudo command to allow by the sudo rule",
			},
			"sudocmd_group": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"sudocmd"},
				Description:   "Sudo command group to allow by the sudo rule",
			},
		},
	}
}

func resourceFreeIPASudoRuleAllowCommandMembershipCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Creating freeipa sudo rule allowed command membership")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	cmd_id := "srac"

	optArgs := ipa.SudoruleAddAllowCommandOptionalArgs{}

	args := ipa.SudoruleAddAllowCommandArgs{
		Cn: d.Get("name").(string),
	}
	if _v, ok := d.GetOkExists("sudocmd"); ok {
		v := []string{_v.(string)}
		optArgs.Sudocmd = &v
		cmd_id = "srac"
	}
	if _v, ok := d.GetOkExists("sudocmd_group"); ok {
		v := []string{_v.(string)}
		optArgs.Sudocmdgroup = &v
		cmd_id = "sracg"
	}

	_, err = client.SudoruleAddAllowCommand(&args, &optArgs)
	if err != nil {
		return diag.Errorf("Error creating freeipa sudo rule allowed command membership: %s", err)
	}

	switch cmd_id {
	case "srac":
		id := fmt.Sprintf("%s/%s/%s", d.Get("name").(string), cmd_id, d.Get("sudocmd").(string))
		d.SetId(id)
	case "sracg":
		id := fmt.Sprintf("%s/%s/%s", d.Get("name").(string), cmd_id, d.Get("sudocmd_group").(string))
		d.SetId(id)
	}

	return resourceFreeIPASudoRuleAllowCommandMembershipRead(ctx, d, meta)
}

func resourceFreeIPASudoRuleAllowCommandMembershipRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Read freeipa sudo rule allowed command membership")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	sudoruleId, typeId, cmdId, err := parseSudoRuleAllowCommandMembershipID(d.Id())

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
	case "srac":
		if res.Result.MemberallowcmdSudocmd == nil || !slices.Contains(*res.Result.MemberallowcmdSudocmd, cmdId) {
			log.Printf("[DEBUG] Warning! Sudo rule allowed command membership does not exist")
			d.Set("name", "")
			d.Set("sudocmd", "")
			d.Set("sudocmd_group", "")
			d.SetId("")
			return nil
		}
	case "sracg":
		if res.Result.MemberallowcmdSudocmdgroup == nil || !slices.Contains(*res.Result.MemberallowcmdSudocmdgroup, cmdId) {
			log.Printf("[DEBUG] Warning! Sudo rule allowed command membership does not exist")
			d.Set("name", "")
			d.Set("sudocmd", "")
			d.Set("sudocmd_group", "")
			d.SetId("")
			return nil
		}
	}

	if err != nil {
		return diag.Errorf("Error show freeipa sudo rule allowed command membership: %s", err)
	}

	log.Printf("[DEBUG] Read freeipa sudo rule allowed command membership %s", res.Result.Cn)
	return nil
}

func resourceFreeIPASudoRuleAllowCommandMembershipDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Delete freeipa sudo rule allowed command membership")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	sudoruleId, typeId, cmdId, err := parseSudoRuleAllowCommandMembershipID(d.Id())

	if err != nil {
		return diag.Errorf("Error parsing ID of freeipa_sudocmdgroup_membership: %s", err)
	}

	optArgs := ipa.SudoruleRemoveAllowCommandOptionalArgs{}

	args := ipa.SudoruleRemoveAllowCommandArgs{
		Cn: sudoruleId,
	}

	switch typeId {
	case "srac":
		v := []string{cmdId}
		optArgs.Sudocmd = &v
	case "sracg":
		v := []string{cmdId}
		optArgs.Sudocmdgroup = &v
	}

	_, err = client.SudoruleRemoveAllowCommand(&args, &optArgs)
	if err != nil {
		return diag.Errorf("Error delete freeipa sudo rule allowed command membership: %s", err)
	}

	d.SetId("")
	return nil
}

func parseSudoRuleAllowCommandMembershipID(id string) (string, string, string, error) {
	idParts := strings.SplitN(id, "/", 3)
	if len(idParts) < 3 {
		return "", "", "", fmt.Errorf("Unable to determine sudo rule allowed command membership ID %s", id)
	}

	name := idParts[0]
	_type := idParts[1]
	sudocmd := idParts[2]

	return name, _type, sudocmd, nil
}
