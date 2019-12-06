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
		//		Update: resourceFreeIpaHostUpdate,
		//Delete: resourceFreeIpaHostDelete,
		/*
			Importer: &schema.ResourceImporter{
				State: resourceFreeIpaHostImport,
			},
		*/

		Schema: map[string]*schema.Schema{
			"fqdn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"randompassword": {
				Type:     schema.TypeString,
				Computed: true,
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
	//description := d.Get("description").(string)
	truue := true

	res, err := client.HostAdd(
		&ipa.HostAddArgs{
			Fqdn: fqdn,
		},
		&ipa.HostAddOptionalArgs{
			Force:  &truue,
			Random: &truue,
		},
	)
	if err != nil {
		return err
	}

	d.Set("randompassword", *res.Result.Randompassword)

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
