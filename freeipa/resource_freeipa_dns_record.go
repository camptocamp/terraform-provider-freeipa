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
			"srv_part_priority": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"srv_part_weight": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"srv_part_port": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"srv_part_target": {
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

	if _aSrvPartPriority, ok := d.GetOkExists("srv_part_priority"); ok {
		aSrvPartPriority := _aSrvPartPriority.(int)
		optArgs.SrvPartPriority = &aSrvPartPriority
	}

	if _aSrvPartWeight, ok := d.GetOkExists("srv_part_weight"); ok {
		aSrvPartWeight := _aSrvPartWeight.(int)
		optArgs.SrvPartWeight = &aSrvPartWeight
	}

	if _aSrvPartPort, ok := d.GetOkExists("srv_part_port"); ok {
		aSrvPartPort := _aSrvPartPort.(int)
		optArgs.SrvPartPort = &aSrvPartPort
	}

	if _aSrvPartTarget, ok := d.GetOkExists("srv_part_target"); ok {
		aSrvPartTarget := _aSrvPartTarget
		optArgs.SrvPartTarget = &aSrvPartTarget
	}

	_, err = client.DnsrecordAdd(&args, &optArgs)
	if err != nil {
		return err
	}

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

	optArgs := ipa.DnsrecordModOptionalArgs{}

	if _dnszoneidnsname, ok := d.GetOkExists("dnszoneidnsname"); ok {
		dnszoneidnsname := _dnszoneidnsname
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

	if _aSrvPartPriority, ok := d.GetOkExists("srv_part_priority"); ok {
		aSrvPartPriority := _aSrvPartPriority.(int)
		optArgs.SrvPartPriority = &aSrvPartPriority
	}

	if _aSrvPartWeight, ok := d.GetOkExists("srv_part_weight"); ok {
		aSrvPartWeight := _aSrvPartWeight.(int)
		optArgs.SrvPartWeight = &aSrvPartWeight
	}

	if _aSrvPartPort, ok := d.GetOkExists("srv_part_port"); ok {
		aSrvPartPort := _aSrvPartPort.(int)
		optArgs.SrvPartPort = &aSrvPartPort
	}

	if _aSrvPartTarget, ok := d.GetOkExists("srv_part_target"); ok {
		aSrvPartTarget := _aSrvPartTarget
		optArgs.SrvPartTarget = &aSrvPartTarget
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
		dnszoneidnsname := _dnszoneidnsname
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

	if res.Result.SrvPartPriority != nil {
		d.Set("srv_part_priority", *res.Result.SrvPartPriority)
	}

	if res.Result.SrvPartWeight != nil {
		d.Set("srv_part_weight", *res.Result.SrvPartWeight)
	}

	if res.Result.SrvPartPort != nil {
		d.Set("srv_part_port", *res.Result.SrvPartPort)
	}

	if res.Result.SrvPartTarget != nil {
		d.Set("srv_part_target", *res.Result.SrvPartTarget)
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
	delAll := true
	optArgs := ipa.DnsrecordDelOptionalArgs{
		Dnszoneidnsname: &dnszoneidnsname,
		DelAll:          &delAll,
	}

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
