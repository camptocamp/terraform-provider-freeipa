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

func resourceFreeIPAHBACPolicyServiceMembership() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFreeIPADNSHBACPolicyServiceMembershipCreate,
		ReadContext:   resourceFreeIPADNSHBACPolicyServiceMembershipRead,
		DeleteContext: resourceFreeIPADNSHBACPolicyServiceMembershipDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "HBAC policy name",
			},
			"service": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"servicegroup"},
				Description:   "Service name the policy is applied to",
			},
			"servicegroup": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"service"},
				Description:   "Service group name the policy is applied to",
			},
		},
	}
}

func resourceFreeIPADNSHBACPolicyServiceMembershipCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Creating freeipa the HBAC policy service membership")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	svcmember_id := "s"
	optArgs := ipa.HbacruleAddServiceOptionalArgs{}

	args := ipa.HbacruleAddServiceArgs{
		Cn: d.Get("name").(string),
	}

	if _v, ok := d.GetOkExists("service"); ok {
		v := []string{_v.(string)}
		optArgs.Hbacsvc = &v
		svcmember_id = "s"
	}
	if _v, ok := d.GetOkExists("servicegroup"); ok {
		v := []string{_v.(string)}
		optArgs.Hbacsvcgroup = &v
		svcmember_id = "sg"
	}

	_, err = client.HbacruleAddService(&args, &optArgs)
	if err != nil {
		return diag.Errorf("Error creating freeipa the HBAC policy service membership: %s", err)
	}
	switch svcmember_id {
	case "sg":
		id := fmt.Sprintf("%s/sg/%s", d.Get("name").(string), d.Get("servicegroup").(string))
		d.SetId(id)
	case "s":
		id := fmt.Sprintf("%s/s/%s", d.Get("name").(string), d.Get("service").(string))
		d.SetId(id)
	}

	return resourceFreeIPADNSHBACPolicyServiceMembershipRead(ctx, d, meta)
}

func resourceFreeIPADNSHBACPolicyServiceMembershipRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Read freeipa the HBAC policy service membership")

	all := true
	name, typeId, svcId, err := parseHBACPolicyServiceMembershipID(d.Id())

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	args := ipa.HbacruleShowArgs{
		Cn: name,
	}
	optArgs := ipa.HbacruleShowOptionalArgs{
		All: &all,
	}
	res, err := client.HbacruleShow(&args, &optArgs)
	if err != nil {
		if strings.Contains(err.Error(), "NotFound") {
			d.Set("service", "")
			d.Set("servicegroup", "")
			d.SetId("")
			log.Printf("[DEBUG] HBAC policy not found")
			return nil
		} else {
			return diag.Errorf("Error reading freeipa HBAC policy: %s", err)
		}
	}

	switch typeId {
	case "sg":
		if res.Result.MemberserviceHbacsvcgroup == nil || !slices.Contains(*res.Result.MemberserviceHbacsvcgroup, svcId) {
			log.Printf("[DEBUG] Warning! Servicegroup membership does not exist")
			d.Set("service", "")
			d.Set("servicegroup", "")
			d.SetId("")
			return nil
		}
	case "s":
		if res.Result.MemberserviceHbacsvc == nil || !slices.Contains(*res.Result.MemberserviceHbacsvc, svcId) {
			log.Printf("[DEBUG] Warning! Service membership does not exist")
			d.Set("service", "")
			d.Set("servicegroup", "")
			d.SetId("")
			return nil
		}
	}

	return nil
}

func resourceFreeIPADNSHBACPolicyServiceMembershipDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Delete freeipa the HBAC policy service membership")

	name, typeId, svcId, err := parseHBACPolicyServiceMembershipID(d.Id())

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	args := ipa.HbacruleRemoveServiceArgs{
		Cn: name,
	}
	optArgs := ipa.HbacruleRemoveServiceOptionalArgs{}

	if typeId == "s" {
		v := []string{svcId}
		optArgs.Hbacsvc = &v
	}
	if typeId == "sg" {
		v := []string{svcId}
		optArgs.Hbacsvcgroup = &v
	}

	_, err = client.HbacruleRemoveService(&args, &optArgs)
	if err != nil {
		return diag.Errorf("Error delete freeipa the HBAC policy service membership: %s", err)
	}

	d.SetId("")

	return nil
}

func parseHBACPolicyServiceMembershipID(id string) (string, string, string, error) {
	idParts := strings.Split(id, "/")
	if len(idParts) < 3 {
		return "", "", "", fmt.Errorf("Unable to determine service membership ID %s", id)
	}

	name := idParts[0]
	_type := idParts[1]
	svc := idParts[2]

	return name, _type, svc, nil
}
