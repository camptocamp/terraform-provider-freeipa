package freeipa

import (
	"fmt"
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
			"record": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
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
		Arecord:         &[]string{},
		Aaaarecord:      &[]string{},
		Cnamerecord:     &[]string{},
		Mxrecord:        &[]string{},
		Nsrecord:        &[]string{},
		Ptrrecord:       &[]string{},
		Srvrecord:       &[]string{},
		Txtrecord:       &[]string{},
		Sshfprecord:     &[]string{},
	}

	records := d.Get("record").(*schema.Set)
	for _, recordI := range records.List() {
		record := recordI.(map[string]interface{})

		_type := record["type"]
		_value := record["value"].(string)

		switch _type {
		case "A":
			*optArgs.Arecord = append(*optArgs.Arecord, _value)
		case "AAAA":
			*optArgs.Aaaarecord = append(*optArgs.Aaaarecord, _value)
		case "CNAME":
			*optArgs.Cnamerecord = append(*optArgs.Cnamerecord, _value)
		case "MX":
			*optArgs.Mxrecord = append(*optArgs.Mxrecord, _value)
		case "NS":
			*optArgs.Nsrecord = append(*optArgs.Nsrecord, _value)
		case "PTR":
			*optArgs.Ptrrecord = append(*optArgs.Ptrrecord, _value)
		case "SRV":
			*optArgs.Srvrecord = append(*optArgs.Srvrecord, _value)
		case "TXT":
			*optArgs.Txtrecord = append(*optArgs.Txtrecord, _value)
		case "SSHFP":
			*optArgs.Sshfprecord = append(*optArgs.Sshfprecord, _value)
		}
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
		Arecord:         &[]string{},
		Aaaarecord:      &[]string{},
		Cnamerecord:     &[]string{},
		Mxrecord:        &[]string{},
		Nsrecord:        &[]string{},
		Ptrrecord:       &[]string{},
		Srvrecord:       &[]string{},
		Txtrecord:       &[]string{},
		Sshfprecord:     &[]string{},
	}

	records := d.Get("record").(*schema.Set)
	for _, recordI := range records.List() {
		record := recordI.(map[string]interface{})

		_type := record["type"]
		_value := record["value"].(string)

		switch _type {
		case "A":
			*optArgs.Arecord = append(*optArgs.Arecord, _value)
		case "AAAA":
			*optArgs.Aaaarecord = append(*optArgs.Aaaarecord, _value)
		case "CNAME":
			*optArgs.Cnamerecord = append(*optArgs.Cnamerecord, _value)
		case "MX":
			*optArgs.Mxrecord = append(*optArgs.Mxrecord, _value)
		case "NS":
			*optArgs.Nsrecord = append(*optArgs.Nsrecord, _value)
		case "PTR":
			*optArgs.Ptrrecord = append(*optArgs.Ptrrecord, _value)
		case "SRV":
			*optArgs.Srvrecord = append(*optArgs.Srvrecord, _value)
		case "TXT":
			*optArgs.Txtrecord = append(*optArgs.Txtrecord, _value)
		case "SSHFP":
			*optArgs.Sshfprecord = append(*optArgs.Sshfprecord, _value)
		}
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

	var records []map[string]interface{}

	if res.Result.Arecord != nil {
		for _, record := range *res.Result.Arecord {
			records = append(records, map[string]interface{}{
				"type":  "A",
				"value": record,
			})
		}
	}

	if res.Result.Aaaarecord != nil {
		for _, record := range *res.Result.Aaaarecord {
			records = append(records, map[string]interface{}{
				"type":  "AAAA",
				"value": record,
			})
		}
	}

	if res.Result.Cnamerecord != nil {
		for _, record := range *res.Result.Cnamerecord {
			records = append(records, map[string]interface{}{
				"type":  "CNAME",
				"value": record,
			})
		}
	}

	if res.Result.Mxrecord != nil {
		for _, record := range *res.Result.Mxrecord {
			records = append(records, map[string]interface{}{
				"type":  "MX",
				"value": record,
			})
		}
	}

	if res.Result.Nsrecord != nil {
		for _, record := range *res.Result.Nsrecord {
			records = append(records, map[string]interface{}{
				"type":  "NS",
				"value": record,
			})
		}
	}

	if res.Result.Ptrrecord != nil {
		for _, record := range *res.Result.Ptrrecord {
			records = append(records, map[string]interface{}{
				"type":  "PTR",
				"value": record,
			})
		}
	}

	if res.Result.Srvrecord != nil {
		for _, record := range *res.Result.Srvrecord {
			records = append(records, map[string]interface{}{
				"type":  "SRV",
				"value": record,
			})
		}
	}

	if res.Result.Txtrecord != nil {
		for _, record := range *res.Result.Txtrecord {
			records = append(records, map[string]interface{}{
				"type":  "TXT",
				"value": record,
			})
		}
	}

	if res.Result.Sshfprecord != nil {
		for _, record := range *res.Result.Sshfprecord {
			records = append(records, map[string]interface{}{
				"type":  "SSHFP",
				"value": record,
			})
		}
	}

	d.Set("record", records)

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
		Arecord:         &[]string{},
		Aaaarecord:      &[]string{},
		Cnamerecord:     &[]string{},
		Mxrecord:        &[]string{},
		Nsrecord:        &[]string{},
		Ptrrecord:       &[]string{},
		Srvrecord:       &[]string{},
		Txtrecord:       &[]string{},
		Sshfprecord:     &[]string{},
	}

	records := d.Get("record").(*schema.Set)
	for _, recordI := range records.List() {
		record := recordI.(map[string]interface{})

		_type := record["type"]
		_value := record["value"].(string)

		switch _type {
		case "A":
			*optArgs.Arecord = append(*optArgs.Arecord, _value)
		case "AAAA":
			*optArgs.Aaaarecord = append(*optArgs.Aaaarecord, _value)
		case "CNAME":
			*optArgs.Cnamerecord = append(*optArgs.Cnamerecord, _value)
		case "MX":
			*optArgs.Mxrecord = append(*optArgs.Mxrecord, _value)
		case "NS":
			*optArgs.Nsrecord = append(*optArgs.Nsrecord, _value)
		case "PTR":
			*optArgs.Ptrrecord = append(*optArgs.Ptrrecord, _value)
		case "SRV":
			*optArgs.Srvrecord = append(*optArgs.Srvrecord, _value)
		case "TXT":
			*optArgs.Txtrecord = append(*optArgs.Txtrecord, _value)
		case "SSHFP":
			*optArgs.Sshfprecord = append(*optArgs.Sshfprecord, _value)
		}
	}

	_, err = client.DnsrecordDel(&args, &optArgs)
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}

func resourceFreeIPADNSRecordImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	idnsname, dnszoneidnsname := splitID(d.Id())

	d.SetId(fmt.Sprintf("%s.%s", idnsname, dnszoneidnsname))

	d.Set("idnsname", idnsname)
	d.Set("dnszoneidnsname", dnszoneidnsname)

	log.Printf("[INFO] Importing DNS Record `%s` in zone `%s`.", idnsname, dnszoneidnsname)

	err := resourceFreeIPADNSRecordRead(d, meta)
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	return []*schema.ResourceData{d}, nil
}
