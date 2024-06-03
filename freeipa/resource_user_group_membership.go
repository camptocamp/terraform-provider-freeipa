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

func resourceFreeIPAUserGroupMembership() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFreeIPAUserGroupMembershipCreate,
		ReadContext:   resourceFreeIPAUserGroupMembershipRead,
		DeleteContext: resourceFreeIPAUserGroupMembershipDelete,
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
			"user": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"group"},
				Description:   "User to add",
			},
			"group": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"user"},
				Description:   "Group to add",
			},
		},
	}
}

func resourceFreeIPAUserGroupMembershipCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Creating freeipa the user group membership")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	user_id := "u"

	optArgs := ipa.GroupAddMemberOptionalArgs{}

	args := ipa.GroupAddMemberArgs{
		Cn: d.Get("name").(string),
	}
	if _v, ok := d.GetOkExists("user"); ok {
		v := []string{_v.(string)}
		optArgs.User = &v
		user_id = "u"
	}
	if _v, ok := d.GetOkExists("group"); ok {
		v := []string{_v.(string)}
		optArgs.Group = &v
		user_id = "g"
	}

	_, err = client.GroupAddMember(&args, &optArgs)
	if err != nil {
		return diag.Errorf("Error creating freeipa the user group membership: %s", err)
	}

	switch user_id {
	case "g":
		id := fmt.Sprintf("%s/g/%s", d.Get("name").(string), d.Get("group").(string))
		d.SetId(id)
	case "u":
		id := fmt.Sprintf("%s/u/%s", d.Get("name").(string), d.Get("user").(string))
		d.SetId(id)
	}

	return resourceFreeIPAUserGroupMembershipRead(ctx, d, meta)
}

func resourceFreeIPAUserGroupMembershipRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Read freeipa the user group membership")

	name, typeId, userId, err := parseUserMembershipID(d.Id())

	if err != nil {
		return diag.Errorf("Error parsing ID of freeipa_user_group_membership: %s", err)
	}

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	optArgs := ipa.GroupFindOptionalArgs{
		Cn: &name,
	}

	switch typeId {
	case "g":
		v := []string{userId}
		optArgs.Group = &v
	case "u":
		v := []string{userId}
		optArgs.User = &v
	}

	res, err := client.GroupFind("", &ipa.GroupFindArgs{}, &optArgs)
	if err != nil {
		return diag.Errorf("Error find freeipa the user group membership: %s", err)
	}

	if strings.Contains(*res.Summary, "0 groups matched") {
		log.Printf("[DEBUG] Warning! Group or User membership not exist")
		d.Set("user", "")
		d.Set("group", "")
		d.SetId("")
	}

	return nil
}

func resourceFreeIPAUserGroupMembershipDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Delete freeipa the user group membership")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	optArgs := ipa.GroupRemoveMemberOptionalArgs{}

	nameId, typeId, userId, err := parseUserMembershipID(d.Id())

	if err != nil {
		return diag.Errorf("Error parsing ID of freeipa_user_group_membership: %s", err)
	}

	args := ipa.GroupRemoveMemberArgs{
		Cn: nameId,
	}

	switch typeId {
	case "g":
		v := []string{userId}
		optArgs.Group = &v
	case "u":
		v := []string{userId}
		optArgs.User = &v
	}

	_, err = client.GroupRemoveMember(&args, &optArgs)
	if err != nil {
		return diag.Errorf("Error delete freeipa the user group membership: %s", err)
	}

	d.SetId("")

	return nil
}

func parseUserMembershipID(id string) (string, string, string, error) {
	idParts := strings.Split(id, "/")
	if len(idParts) < 3 {
		return "", "", "", fmt.Errorf("Unable to determine user membership ID %s", id)
	}

	name := idParts[0]
	_type := idParts[1]
	user := idParts[2]

	return name, _type, user, nil
}
