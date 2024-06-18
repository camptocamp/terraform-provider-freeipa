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

func resourceFreeIPASudoRuleHostMembership() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFreeIPASudoRuleHostMembershipCreate,
		ReadContext:   resourceFreeIPASudoRuleHostMembershipRead,
		DeleteContext: resourceFreeIPASudoRuleHostMembershipDelete,
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
			"host": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"hostgroup"},
				Description:   "Host to add to the sudo rule",
			},
			"hostgroup": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"host"},
				Description:   "Hostgroup to add to the sudo rule",
			},
			// Hostmask not implemented yet. Maybe one day but I don't see the need.
			// "hostmask": {
			// 	Type:          schema.TypeString,
			// 	Optional:      true,
			// 	ForceNew:      true,
			// 	ConflictsWith: []string{"host", "hostgroup"},
			// },
		},
	}
}

func resourceFreeIPASudoRuleHostMembershipCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Creating freeipa sudo rule host membership")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	host_id := "srh"

	optArgs := ipa.SudoruleAddHostOptionalArgs{}

	args := ipa.SudoruleAddHostArgs{
		Cn: d.Get("name").(string),
	}
	if _v, ok := d.GetOkExists("host"); ok {
		v := []string{_v.(string)}
		optArgs.Host = &v
		host_id = "srh"
	}
	if _v, ok := d.GetOkExists("hostgroup"); ok {
		v := []string{_v.(string)}
		optArgs.Hostgroup = &v
		host_id = "srhg"
	}
	// Hostmask not implemented yet. Maybe one day but I don't see the need.
	// if _v, ok := d.GetOkExists("hostmask"); ok {
	// 	v := []string{_v.(string)}
	// 	optArgs.Hostmask = &v
	// 	host_id = "srhm"
	// }

	_, err = client.SudoruleAddHost(&args, &optArgs)
	if err != nil {
		return diag.Errorf("Error creating freeipa sudo rule host membership: %s", err)
	}

	switch host_id {
	case "srh":
		id := fmt.Sprintf("%s/%s/%s", d.Get("name").(string), host_id, d.Get("host").(string))
		d.SetId(id)
	case "srhg":
		id := fmt.Sprintf("%s/%s/%s", d.Get("name").(string), host_id, d.Get("hostgroup").(string))
		d.SetId(id)
		// Hostmask not implemented yet. Maybe one day but I don't see the need.
		// case "srhm":
		// 	id := fmt.Sprintf("%s/%s/%s", d.Get("name").(string), host_id, d.Get("hostmask").(string))
		// 	d.SetId(id)
	}

	return resourceFreeIPASudoRuleHostMembershipRead(ctx, d, meta)
}

func resourceFreeIPASudoRuleHostMembershipRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Read freeipa sudo rule host membership")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	sudoruleId, typeId, host_id, err := parseSudoRuleHostMembershipID(d.Id())

	if err != nil {
		return diag.Errorf("Error parsing ID of freeipa_sudo_rule_host_membership: %s", err)
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
	case "srh":
		if res.Result.MemberhostHost == nil || !slices.Contains(*res.Result.MemberhostHost, host_id) {
			log.Printf("[DEBUG] Warning! Sudo rule host membership does not exist")
			d.Set("name", "")
			d.Set("host", "")
			d.Set("hostgroup", "")
			d.SetId("")
			return nil
		}
	case "srhg":
		if res.Result.MemberhostHostgroup == nil || !slices.Contains(*res.Result.MemberhostHostgroup, host_id) {
			log.Printf("[DEBUG] Warning! Sudo rule host membership does not exist")
			d.Set("name", "")
			d.Set("host", "")
			d.Set("hostgroup", "")
			d.SetId("")
			return nil
		}
		// Hostmask not implemented yet. Maybe one day but I don't see the need.
		// case "srhm":
		// 	if res.Result.Hostmask == nil || !slices.Contains(*res.Result.Hostmask, host_id) {
		// 		log.Printf("[DEBUG] Warning! Sudo rule host membership does not exist")
		// 		d.Set("name", "")
		// 		d.Set("sudocmd", "")
		// 		d.Set("sudocmd_group", "")
		// 		d.SetId("")
		// 		return diag.Errorf("Error configuring freeipa Sudo rule host, hostmask not assigned: %s", host_id)
		// 	}
	}

	if err != nil {
		return diag.Errorf("Error show freeipa sudo rule host membership: %s", err)
	}

	log.Printf("[DEBUG] Read freeipa sudo rule host membership %s", res.Result.Cn)
	return nil
}

func resourceFreeIPASudoRuleHostMembershipDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Delete freeipa sudo rule host membership")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	sudoruleId, typeId, host_id, err := parseSudoRuleHostMembershipID(d.Id())

	if err != nil {
		return diag.Errorf("Error parsing ID of freeipa_rule_host_membership: %s", err)
	}

	optArgs := ipa.SudoruleRemoveHostOptionalArgs{}

	args := ipa.SudoruleRemoveHostArgs{
		Cn: sudoruleId,
	}

	switch typeId {
	case "srh":
		v := []string{host_id}
		optArgs.Host = &v
	case "srhg":
		v := []string{host_id}
		optArgs.Hostgroup = &v
		// Hostmask not implemented yet. Maybe one day but I don't see the need.
		// case "srhm":
		// 	v := []string{host_id}
		// 	optArgs.Hostmask = &v
	}

	_, err = client.SudoruleRemoveHost(&args, &optArgs)
	if err != nil {
		return diag.Errorf("Error delete freeipa sudo rule host membership: %s", err)
	}

	d.SetId("")
	return nil
}

func parseSudoRuleHostMembershipID(id string) (string, string, string, error) {
	idParts := strings.SplitN(id, "/", 3)
	if len(idParts) < 3 {
		return "", "", "", fmt.Errorf("Unable to determine sudo rule host membership ID %s", id)
	}

	name := idParts[0]
	_type := idParts[1]
	host := idParts[2]

	return name, _type, host, nil
}
