package freeipa

import (
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	ipa "github.com/tehwalris/go-freeipa/freeipa"
)

func resourceFreeIPAHost() *schema.Resource {
	return &schema.Resource{
		Create: resourceFreeIPAHostCreate,
		Read:   resourceFreeIPAHostRead,
		Update: resourceFreeIPAHostUpdate,
		Delete: resourceFreeIPAHostDelete,
		Importer: &schema.ResourceImporter{
			State: resourceFreeIPAHostImport,
		},

		Schema: map[string]*schema.Schema{
			"fqdn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"random": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"userpassword": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"force": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"randompassword": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceFreeIPAHostCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO][freeipa] Creating Host: %s", d.Id())

	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	args := ipa.HostAddArgs{
		Fqdn: d.Get("fqdn").(string),
	}

	optArgs := ipa.HostAddOptionalArgs{}

	if _description, ok := d.GetOkExists("description"); ok {
		description := _description.(string)
		optArgs.Description = &description
	}

	if _random, ok := d.GetOkExists("random"); ok {
		random := _random.(bool)
		optArgs.Random = &random
	}

	if _userpassword, ok := d.GetOkExists("userpassword"); ok {
		userpassword := _userpassword.(string)
		optArgs.Userpassword = &userpassword
	}

	if _force, ok := d.GetOkExists("force"); ok {
		force := _force.(bool)
		optArgs.Force = &force
	}

	res, err := client.HostAdd(&args, &optArgs)
	if err != nil {
		return err
	}

	d.SetId(args.Fqdn)

	// randompassword is not returned by HostShow
	if d.Get("random").(bool) {
		d.Set("randompassword", *res.Result.Randompassword)
	}

	// FIXME: When using a LB in front of a FreeIPA cluster, sometime the record
	// is not replicated on the server where the read is done, so we have to
	// retry to not have "Error: NotFound (4001)".
	// Maybe we should use resource.StateChangeConf instead...
	sleepDelay := 1 * time.Second
	for {
		err := resourceFreeIPAHostRead(d, meta)
		if err == nil {
			return nil
		}
		time.Sleep(sleepDelay)
		sleepDelay = sleepDelay * 2
	}
}

func resourceFreeIPAHostUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating Host: %s", d.Id())

	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	args := ipa.HostModArgs{
		Fqdn: d.Get("fqdn").(string),
	}

	optArgs := ipa.HostModOptionalArgs{}

	if _description, ok := d.GetOkExists("description"); ok {
		description := _description.(string)
		optArgs.Description = &description
	}

	if _random, ok := d.GetOkExists("random"); ok {
		random := _random.(bool)
		optArgs.Random = &random
	}

	if _userpassword, ok := d.GetOkExists("userpassword"); ok {
		userpassword := _userpassword.(string)
		optArgs.Userpassword = &userpassword
	}

	res, err := client.HostMod(&args, &optArgs)
	if err != nil {
		return err
	}

	// randompassword is not returned by HostShow
	if d.Get("random").(bool) {
		d.Set("randompassword", *res.Result.Randompassword)
	}

	return resourceFreeIPAHostRead(d, meta)
}

func resourceFreeIPAHostRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing Host: %s", d.Id())

	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	args := ipa.HostShowArgs{
		Fqdn: d.Get("fqdn").(string),
	}

	optArgs := ipa.HostShowOptionalArgs{}

	res, err := client.HostShow(&args, &optArgs)
	if err != nil {
		return err
	}

	if res.Result.Description != nil {
		d.Set("description", *res.Result.Description)
	}

	if res.Result.Userpassword != nil {
		d.Set("userpassword", *res.Result.Userpassword)
	}

	return nil
}

func resourceFreeIPAHostDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting Host: %s", d.Id())

	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	args := ipa.HostDelArgs{
		Fqdn: []string{
			d.Get("fqdn").(string),
		},
	}

	optArgs := ipa.HostDelOptionalArgs{}

	_, err = client.HostDel(&args, &optArgs)
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}

func resourceFreeIPAHostImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[INFO] Importing Host: %s", d.Id())

	d.SetId(d.Id())
	d.Set("fqdn", d.Id())

	err := resourceFreeIPAHostRead(d, meta)
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	return []*schema.ResourceData{d}, nil
}
