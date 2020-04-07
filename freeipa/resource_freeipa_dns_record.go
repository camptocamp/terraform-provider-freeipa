package freeipa

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	ipa "github.com/tehwalris/go-freeipa/freeipa"
)

func resourceFreeIPADNSRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceFreeIPADNSRecordCreate,
		Read:   resourceFreeIPADNSRecordRead,
		Update: resourceFreeIPADNSRecordUpdate,
		Delete: resourceFreeIPADNSRecordDelete,
		Importer: &schema.ResourceImporter{
			State: resourceFreeIPADNSRecordImport,
		},

		Schema: map[string]*schema.Schema{
			"idnsname": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"dnszoneidnsname": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"dnsttl": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"dnsclass": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"a_part_ip_address": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceFreeIPADNSRecordCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO][freeipa] Creating DNS Record: %s", d.Id())

	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	args := ipa.DnsrecordAddArgs{
		Idnsname: d.Get("idnsname").(string),
	}

	optArgs := ipa.DnsrecordAddOptionalArgs{}

	if _dnszoneidnsname, ok := d.GetOkExists("dnszoneidnsname"); ok {
		dnszoneidnsname := _dnszoneidnsname.(string)
		optArgs.Dnszoneidnsname = &dnszoneidnsname
	}

	if _dnsttl, ok := d.GetOkExists("dnsttl"); ok {
		dnsttl := _dnsttl.(int)
		optArgs.Dnsttl = &dnsttl
	}

	if _dnsclass, ok := d.GetOkExists("dnsclass"); ok {
		dnsclass := _dnsclass.(string)
		optArgs.Dnsclass = &dnsclass
	}

	if _aPartIPAddress, ok := d.GetOkExists("a_part_ip_address"); ok {
		aPartIPAddress := _aPartIPAddress.(string)
		optArgs.APartIPAddress = &aPartIPAddress
	}

	_, err = client.DnsrecordAdd(&args, &optArgs)
	if err != nil {
		return err
	}

	d.SetId(args.Idnsname)

	return resourceFreeIPADNSRecordRead(d, meta)
}

func resourceFreeIPADNSRecordUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating DNS Record: %s", d.Id())

	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	args := ipa.DnsrecordModArgs{
		Idnsname: d.Get("idnsname").(string),
	}

	optArgs := ipa.DnsrecordModOptionalArgs{}

	if _dnszoneidnsname, ok := d.GetOkExists("dnszoneidnsname"); ok {
		dnszoneidnsname := _dnszoneidnsname.(string)
		optArgs.Dnszoneidnsname = &dnszoneidnsname
	}

	if _dnsttl, ok := d.GetOkExists("dnsttl"); ok {
		dnsttl := _dnsttl.(int)
		optArgs.Dnsttl = &dnsttl
	}

	if _dnsclass, ok := d.GetOkExists("dnsclass"); ok {
		dnsclass := _dnsclass.(string)
		optArgs.Dnsclass = &dnsclass
	}

	if _aPartIPAddress, ok := d.GetOkExists("a_part_ip_address"); ok {
		aPartIPAddress := _aPartIPAddress.(string)
		optArgs.APartIPAddress = &aPartIPAddress
	}

	_, err = client.DnsrecordMod(&args, &optArgs)
	if err != nil {
		return err
	}

	return resourceFreeIPADNSRecordRead(d, meta)
}

func resourceFreeIPADNSRecordRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing DNS Record: %s", d.Id())

	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	args := ipa.DnsrecordShowArgs{
		Idnsname: d.Get("idnsname").(string),
	}

	optArgs := ipa.DnsrecordShowOptionalArgs{}

	if _dnszoneidnsname, ok := d.GetOkExists("dnszoneidnsname"); ok {
		dnszoneidnsname := _dnszoneidnsname.(string)
		optArgs.Dnszoneidnsname = &dnszoneidnsname
	}

	res, err := client.DnsrecordShow(&args, &optArgs)
	if err != nil {
		return err
	}

	if res.Result.Dnsttl != nil {
		d.Set("dnsttl", *res.Result.Dnsttl)
	}

	if res.Result.Dnsclass != nil {
		d.Set("dnsclass", *res.Result.Dnsclass)
	}

	if res.Result.APartIPAddress != nil {
		d.Set("a_part_ip_address", *res.Result.APartIPAddress)
	}

	return nil
}

func resourceFreeIPADNSRecordDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting DNS Record: %s", d.Id())

	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	args := ipa.DnsrecordDelArgs{
		Idnsname: d.Get("idnsname").(string),
	}

	optArgs := ipa.DnsrecordDelOptionalArgs{}

	_, err = client.DnsrecordDel(&args, &optArgs)
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}

func resourceFreeIPADNSRecordImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[INFO] Importing DNS Record: %s", d.Id())

	d.SetId(d.Id())
	d.Set("idnsname", d.Id())

	err := resourceFreeIPADNSRecordRead(d, meta)
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	return []*schema.ResourceData{d}, nil
}
