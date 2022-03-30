package freeipa

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	ipa "github.com/tehwalris/go-freeipa/freeipa"
)

func resourceFreeIPADNSZone() *schema.Resource {
	return &schema.Resource{
		Create: resourceFreeIPADNSZoneCreate,
		Read:   resourceFreeIPADNSZoneRead,
		Update: resourceFreeIPADNSZoneUpdate,
		Delete: resourceFreeIPADNSZoneDelete,
		Importer: &schema.ResourceImporter{
			State: resourceFreeIPADNSZoneImport,
		},

		Schema: map[string]*schema.Schema{
			"idnssoaserial": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"idnsname": {
				Type:     schema.TypeString,
				Required: true,
				// Optional: true,
			},
		},
	}
}

func resourceFreeIPADNSZoneCreate(d *schema.ResourceData, meta interface{}) error {

	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	idnssoaserial := d.Get("idnssoaserial").(int)
	idnsname := d.Get("idnsname")
	optArgs := ipa.DnszoneAddOptionalArgs{
		Idnsname: &idnsname,
	}

	client.DnszoneAdd(
		&ipa.DnszoneAddArgs{
			Idnssoaserial: idnssoaserial,
		},
		&optArgs,
	)

	d.SetId(fmt.Sprintf("%s.%d", idnsname, idnssoaserial))

	return resourceFreeIPADNSZoneRead(d, meta)
}

func resourceFreeIPADNSZoneUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating DNS Zone: %s", d.Id())

	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	idnsname := d.Get("idnsname")
	optArgs := ipa.DnszoneModOptionalArgs{
		Idnsname: &idnsname,
	}

	client.DnszoneMod(
		&ipa.DnszoneModArgs{},
		&optArgs,
	)
	return resourceFreeIPADNSZoneRead(d, meta)
}

func resourceFreeIPADNSZoneRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}
	idnsname := d.Get("idnsname")

	optArgs := ipa.DnszoneShowOptionalArgs{
		Idnsname: &idnsname,
	}

	client.DnszoneShow(
		&ipa.DnszoneShowArgs{},
		&optArgs,
	)

	return nil
}
func resourceFreeIPADNSZoneDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting Host: %s", d.Id())
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}
	idnsname := d.Get("idnsname")

	client.DnszoneDel(
		&ipa.DnszoneDelArgs{},
		&ipa.DnszoneDelOptionalArgs{
			Idnsname: &[]interface{}{idnsname},
		},
	)

	d.SetId("")
	return nil
}

func resourceFreeIPADNSZoneImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {

	d.SetId(d.Id())
	d.Set("idnsname", d.Id())
	err := resourceFreeIPADNSZoneRead(d, meta)
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	return []*schema.ResourceData{d}, nil
}
