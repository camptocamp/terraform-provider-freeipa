package freeipa

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"

	ipa "github.com/tehwalris/go-freeipa/freeipa"
)

func resourceFreeIpaHost() *schema.Resource {
	return &schema.Resource{
		Create: resourceFreeIpaHostCreate,
		Read:   resourceFreeIpaHostRead,
		Update: resourceFreeIpaHostUpdate,
		Delete: resourceFreeIpaHostDelete,
		/*
			Importer: &schema.ResourceImporter{
				State: resourceFreeIpaHostImport,
			},
		*/

		Schema: map[string]*schema.Schema{
			"fqdn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"random": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"userpassword": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"randompassword": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"force": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceFreeIpaHostCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO][freeipa] Creating Host: %s", d.Id())
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	fqdn := d.Get("fqdn").(string)
	random := d.Get("random").(bool)
	userpassword := d.Get("userpassword").(string)
	force := d.Get("force").(bool)

	optArgs := ipa.HostAddOptionalArgs{
		Random: &random,
		Force:  &force,
	}

	if userpassword != "" {
		optArgs.Userpassword = &userpassword
	}

	res, err := client.HostAdd(
		&ipa.HostAddArgs{
			Fqdn: fqdn,
		},
		&optArgs,
	)
	if err != nil {
		return err
	}

	d.SetId(fqdn)

	if random {
		d.Set("randompassword", *res.Result.Randompassword)
	}

	return nil
}

func resourceFreeIpaHostRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing Host: %s", d.Id())
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	fqdn := d.Get("fqdn").(string)

	// Only check that we can find the host
	_, err = client.HostShow(
		&ipa.HostShowArgs{
			Fqdn: fqdn,
		},
		&ipa.HostShowOptionalArgs{},
	)
	if err != nil {
		return err
	}

	return nil
}

func resourceFreeIpaHostUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating Host: %s", d.Id())
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	fqdn := d.Get("fqdn").(string)
	random := d.Get("random").(bool)
	userpassword := d.Get("userpassword").(string)

	optArgs := ipa.HostModOptionalArgs{
		Random: &random,
	}

	if userpassword != "" {
		optArgs.Userpassword = &userpassword
	}

	res, err := client.HostMod(
		&ipa.HostModArgs{
			Fqdn: fqdn,
		},
		&optArgs,
	)
	if err != nil {
		return err
	}

	if random {
		d.Set("randompassword", *res.Result.Randompassword)
	}

	return nil
}

func resourceFreeIpaHostDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting Host: %s", d.Id())
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	fqdn := d.Get("fqdn").(string)

	_, err = client.HostDel(
		&ipa.HostDelArgs{
			Fqdn: []string{fqdn},
		},
		&ipa.HostDelOptionalArgs{},
	)
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}
