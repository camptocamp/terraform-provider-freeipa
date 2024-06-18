package freeipa

import (
	"context"
	"fmt"
	"log"
	"strings"

	ipa "github.com/camptocamp/go-freeipa/freeipa"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceFreeIPAHostHostGroupMembership() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFreeIPAHostHostGroupMembershipCreate,
		ReadContext:   resourceFreeIPAHostHostGroupMembershipRead,
		DeleteContext: resourceFreeIPAHostHostGroupMembershipDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Group name",
			},
			"host": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"hostgroup"},
				Description:   "Host to add",
			},
			"hostgroup": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"host"},
				Description:   "HostGroup to add",
			},
		},
	}
}

func resourceFreeIPAHostHostGroupMembershipCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Creating freeipa the host group membership")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	host_id := "h"

	optArgs := ipa.HostgroupAddMemberOptionalArgs{}

	args := ipa.HostgroupAddMemberArgs{
		Cn: d.Get("name").(string),
	}
	if _v, ok := d.GetOkExists("host"); ok {
		v := []string{_v.(string)}
		optArgs.Host = &v
		host_id = "h"
	}
	if _v, ok := d.GetOkExists("hostgroup"); ok {
		v := []string{_v.(string)}
		optArgs.Hostgroup = &v
		host_id = "hg"
	}

	_, err = client.HostgroupAddMember(&args, &optArgs)
	if err != nil {
		return diag.Errorf("Error creating freeipa the host group membership: %s", err)
	}

	switch host_id {
	case "hg":
		id := fmt.Sprintf("%s/hg/%s", d.Get("name").(string), d.Get("hostgroup").(string))
		d.SetId(id)
	case "h":
		id := fmt.Sprintf("%s/h/%s", d.Get("name").(string), d.Get("host").(string))
		d.SetId(id)
	}

	return resourceFreeIPAHostHostGroupMembershipRead(ctx, d, meta)
}

func resourceFreeIPAHostHostGroupMembershipRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Read freeipa the host group membership")

	name, typeId, hostId, err := parseHostMembershipID(d.Id())

	if err != nil {
		return diag.Errorf("Error parsing ID of freeipa_host_group_membership: %s", err)
	}

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	optArgs := ipa.HostgroupFindOptionalArgs{
		Cn: &name,
	}

	switch typeId {
	case "hg":
		v := []string{hostId}
		optArgs.Hostgroup = &v
	case "h":
		v := []string{hostId}
		optArgs.Host = &v
	}

	res, err := client.HostgroupFind("", &ipa.HostgroupFindArgs{}, &optArgs)
	if err != nil {
		return diag.Errorf("Error find freeipa the host group membership: %s", err)
	}

	if strings.Contains(*res.Summary, "0 groups matched") {
		log.Printf("[DEBUG] Warning! Hostgroup or Host membership not exist")
		d.Set("host", "")
		d.Set("hostgroup", "")
		d.SetId("")
	}

	return nil
}

func resourceFreeIPAHostHostGroupMembershipDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Delete freeipa the host group membership")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	optArgs := ipa.HostgroupRemoveMemberOptionalArgs{}

	nameId, typeId, hostId, err := parseHostMembershipID(d.Id())

	if err != nil {
		return diag.Errorf("Error parsing ID of freeipa_user_group_membership: %s", err)
	}

	args := ipa.HostgroupRemoveMemberArgs{
		Cn: nameId,
	}

	switch typeId {
	case "hg":
		v := []string{hostId}
		optArgs.Hostgroup = &v
	case "h":
		v := []string{hostId}
		optArgs.Host = &v
	}

	_, err = client.HostgroupRemoveMember(&args, &optArgs)
	if err != nil {
		if strings.Contains(err.Error(), "NotFound") {
			d.SetId("")
			log.Printf("[DEBUG] Hostgroup not found")
			return nil
		} else {
			return diag.Errorf("Error delete freeipa the host group membership: %s", err)
		}
	}
	d.SetId("")

	return nil
}

func parseHostMembershipID(id string) (string, string, string, error) {
	idParts := strings.Split(id, "/")
	if len(idParts) < 3 {
		return "", "", "", fmt.Errorf("Unable to determine host membership ID %s", id)
	}

	name := idParts[0]
	_type := idParts[1]
	host := idParts[2]

	return name, _type, host, nil
}
