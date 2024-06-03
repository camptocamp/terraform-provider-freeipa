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

func resourceFreeIPASudocmdgroupMembership() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFreeIPASudocmdgroupMembershipCreate,
		ReadContext:   resourceFreeIPASudocmdgroupMembershipRead,
		DeleteContext: resourceFreeIPASudocmdgroupMembershipDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the sudo command group",
			},
			"sudocmd": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Sudo command to add to the group",
			},
		},
	}
}

func resourceFreeIPASudocmdgroupMembershipCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Creating freeipa sudo command group membership")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	optArgs := ipa.SudocmdgroupAddMemberOptionalArgs{}

	args := ipa.SudocmdgroupAddMemberArgs{
		Cn: d.Get("name").(string),
	}
	if _v, ok := d.GetOkExists("sudocmd"); ok {
		v := []string{_v.(string)}
		optArgs.Sudocmd = &v
	}

	_, err = client.SudocmdgroupAddMember(&args, &optArgs)
	if err != nil {
		return diag.Errorf("Error creating freeipa sudo command group membership: %s", err)
	}

	id := fmt.Sprintf("%s/sc/%s", d.Get("name").(string), d.Get("sudocmd").(string))
	d.SetId(id)

	return resourceFreeIPASudocmdgroupMembershipRead(ctx, d, meta)
}

func resourceFreeIPASudocmdgroupMembershipRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Read freeipa sudo command group membership")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	cmdgrpId, typeId, cmdId, err := parseSudocmdgroupMembershipID(d.Id())

	if err != nil {
		return diag.Errorf("Error parsing ID of freeipa_sudocmdgroup_membership: %s", err)
	}

	all := true
	optArgs := ipa.SudocmdgroupShowOptionalArgs{
		All: &all,
	}

	args := ipa.SudocmdgroupShowArgs{
		Cn: cmdgrpId,
	}

	res, err := client.SudocmdgroupShow(&args, &optArgs)
	if err != nil {
		if strings.Contains(err.Error(), "NotFound") {
			d.SetId("")
			log.Printf("[DEBUG] Sudo command group not found")
			return nil
		} else {
			return diag.Errorf("Error reading freeipa Sudo command group: %s", err)
		}
	}

	switch typeId {
	case "sc":
		if res.Result.MemberSudocmd == nil || !slices.Contains(*res.Result.MemberSudocmd, cmdId) {
			log.Printf("[DEBUG] Warning! Sudo command group membership does not exist")
			d.Set("name", "")
			d.Set("sudocmd", "")
			d.SetId("")
			return nil
		}
	}

	if err != nil {
		return diag.Errorf("Error show freeipa sudo command group membership: %s", err)
	}

	log.Printf("[DEBUG] Read freeipa sudo command group membership %s", res.Result.Cn)
	return nil
}

func resourceFreeIPASudocmdgroupMembershipDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Delete freeipa sudo command group membership")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	cmdgrpId, typeId, cmdId, err := parseSudocmdgroupMembershipID(d.Id())

	if err != nil {
		return diag.Errorf("Error parsing ID of freeipa_sudocmdgroup_membership: %s", err)
	}

	optArgs := ipa.SudocmdgroupRemoveMemberOptionalArgs{}

	args := ipa.SudocmdgroupRemoveMemberArgs{
		Cn: cmdgrpId,
	}

	switch typeId {
	case "sc":
		v := []string{cmdId}
		optArgs.Sudocmd = &v
	}

	_, err = client.SudocmdgroupRemoveMember(&args, &optArgs)
	if err != nil {
		return diag.Errorf("Error delete freeipa sudo command group membership: %s", err)
	}

	d.SetId("")
	return nil
}

func parseSudocmdgroupMembershipID(id string) (string, string, string, error) {
	idParts := strings.SplitN(id, "/", 3)
	if len(idParts) < 3 {
		return "", "", "", fmt.Errorf("Unable to determine sudo command group membership ID %s", id)
	}

	name := idParts[0]
	_type := idParts[1]
	sudocmd := idParts[2]

	return name, _type, sudocmd, nil
}
