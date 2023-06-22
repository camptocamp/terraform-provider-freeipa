package freeipa

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	ipa "github.com/ccin2p3/go-freeipa/freeipa"
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
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"records": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
				// Set:      schema.HashString,
			},
			"dnsttl": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"dnsclass": {
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

	idnsname := d.Get("idnsname").(string)
	dnszoneidnsname := d.Get("dnszoneidnsname")

	args := ipa.DnsrecordAddArgs{
		Idnsname: idnsname,
	}

	optArgs := ipa.DnsrecordAddOptionalArgs{
		Dnszoneidnsname: &dnszoneidnsname,
	}

	_type := d.Get("type")
	_records := d.Get("records").(*schema.Set).List()
	records := make([]string, len(_records))
	for i, d := range _records {
		records[i] = d.(string)
	}
	switch _type {
	case "A":
		optArgs.Arecord = &records
	case "AAAA":
		optArgs.Aaaarecord = &records
	case "CNAME":
		optArgs.Cnamerecord = &records
	case "MX":
		optArgs.Mxrecord = &records
	case "NS":
		optArgs.Nsrecord = &records
	case "PTR":
		optArgs.Ptrrecord = &records
	case "SRV":
		optArgs.Srvrecord = &records
	case "TXT":
		optArgs.Txtrecord = &records
	case "SSHFP":
		optArgs.Sshfprecord = &records
	}

	if _dnsttl, ok := d.GetOkExists("dnsttl"); ok {
		dnsttl := _dnsttl.(int)
		optArgs.Dnsttl = &dnsttl
	}

	if _dnsclass, ok := d.GetOkExists("dnsclass"); ok {
		dnsclass := _dnsclass.(string)
		optArgs.Dnsclass = &dnsclass
	}

	_, err = client.DnsrecordAdd(&args, &optArgs)
	if err != nil {
		return err
	}

	// TODO: use aws_route53_records' way to generate ID
	d.SetId(fmt.Sprintf("%s.%s", idnsname, dnszoneidnsname))

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

	dnszoneidnsname := d.Get("dnszoneidnsname")
	optArgs := ipa.DnsrecordModOptionalArgs{
		Dnszoneidnsname: &dnszoneidnsname,
	}

	_type := d.Get("type")
	_records := d.Get("records").(*schema.Set).List()
	records := make([]string, len(_records))
	for i, d := range _records {
		records[i] = d.(string)
	}
	switch _type {
	case "A":
		optArgs.Arecord = &records
	case "AAAA":
		optArgs.Aaaarecord = &records
	case "CNAME":
		optArgs.Cnamerecord = &records
	case "MX":
		optArgs.Mxrecord = &records
	case "NS":
		optArgs.Nsrecord = &records
	case "PTR":
		optArgs.Ptrrecord = &records
	case "SRV":
		optArgs.Srvrecord = &records
	case "TXT":
		optArgs.Txtrecord = &records
	case "SSHFP":
		optArgs.Sshfprecord = &records
	}

	if _dnsttl, ok := d.GetOkExists("dnsttl"); ok {
		dnsttl := _dnsttl.(int)
		optArgs.Dnsttl = &dnsttl
	}

	if _dnsclass, ok := d.GetOkExists("dnsclass"); ok {
		dnsclass := _dnsclass.(string)
		optArgs.Dnsclass = &dnsclass
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

	dnszoneidnsname := d.Get("dnszoneidnsname")
	all := true
	optArgs := ipa.DnsrecordShowOptionalArgs{
		Dnszoneidnsname: &dnszoneidnsname,
		All:             &all,
	}

	res, err := client.DnsrecordShow(&args, &optArgs)
	if err != nil {
		return err
	}

	_type := d.Get("type")

	switch _type {
	case "A":
		if res.Result.Arecord != nil {
			d.Set("records", *res.Result.Arecord)
		}
	case "AAAA":
		if res.Result.Aaaarecord != nil {
			d.Set("records", *res.Result.Aaaarecord)
		}
	case "MX":
		if res.Result.Mxrecord != nil {
			d.Set("records", *res.Result.Mxrecord)
		}
	case "NS":
		if res.Result.Nsrecord != nil {
			d.Set("records", *res.Result.Nsrecord)
		}
	case "PTR":
		if res.Result.Ptrrecord != nil {
			d.Set("records", *res.Result.Ptrrecord)
		}
	case "SRV":
		if res.Result.Srvrecord != nil {
			d.Set("records", *res.Result.Srvrecord)
		}
	case "TXT":
		if res.Result.Txtrecord != nil {
			d.Set("records", *res.Result.Txtrecord)
		}
	case "SSHFP":
		if res.Result.Sshfprecord != nil {
			d.Set("records", *res.Result.Sshfprecord)
		}
	}

	if res.Result.Dnsttl != nil {
		d.Set("dnsttl", *res.Result.Dnsttl)
	}

	if res.Result.Dnsclass != nil {
		d.Set("dnsclass", *res.Result.Dnsclass)
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

	dnszoneidnsname := d.Get("dnszoneidnsname")
	optArgs := ipa.DnsrecordDelOptionalArgs{
		Dnszoneidnsname: &dnszoneidnsname,
	}

	_type := d.Get("type")
	_records := d.Get("records").(*schema.Set).List()
	records := make([]string, len(_records))
	for i, d := range _records {
		records[i] = d.(string)
	}
	switch _type {
	case "A":
		optArgs.Arecord = &records
	case "AAAA":
		optArgs.Aaaarecord = &records
	case "CNAME":
		optArgs.Cnamerecord = &records
	case "MX":
		optArgs.Mxrecord = &records
	case "NS":
		optArgs.Nsrecord = &records
	case "PTR":
		optArgs.Ptrrecord = &records
	case "SRV":
		optArgs.Srvrecord = &records
	case "TXT":
		optArgs.Ptrrecord = &records
	case "SSHFP":
		optArgs.Sshfprecord = &records
	}

	_, err = client.DnsrecordDel(&args, &optArgs)
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}

func resourceFreeIPADNSRecordImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	idnsname, dnszoneidnsname, _type := splitID(d.Id())

	d.SetId(fmt.Sprintf("%s.%s", idnsname, dnszoneidnsname))

	d.Set("idnsname", idnsname)
	d.Set("dnszoneidnsname", dnszoneidnsname)
	d.Set("type", _type)

	log.Printf("[INFO] Importing DNS Record `%s` in zone `%s`.", idnsname, dnszoneidnsname)

	err := resourceFreeIPADNSRecordRead(d, meta)
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	return []*schema.ResourceData{d}, nil
}
