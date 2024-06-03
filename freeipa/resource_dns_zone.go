package freeipa

import (
	"context"
	"log"
	"strings"

	ipa "github.com/camptocamp/go-freeipa/freeipa"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceFreeIPADNSZone() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFreeIPADNSDNSZoneCreate,
		ReadContext:   resourceFreeIPADNSDNSZoneRead,
		UpdateContext: resourceFreeIPADNSDNSZoneUpdate,
		DeleteContext: resourceFreeIPADNSDNSZoneDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"zone_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Zone name (FQDN)",
			},
			"is_reverse_zone": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Allow create the reverse zone",
			},
			"disable_zone": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Allow disabled the zone",
			},
			"skip_overlap_check": { // Force DNS zone creation even if it will overlap with an existing zone
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Force DNS zone creation even if it will overlap with an existing zone",
			},
			"authoritative_nameserver": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Authoritative nameserver domain name",
			},
			"skip_nameserver_check": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Force DNS zone creation even if nameserver is not resolvable",
			},
			"admin_email_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Administrator e-mail address",
			},
			"soa_serial_number": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "SOA record serial number",
			},
			"soa_refresh": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     3600,
				Description: "SOA record refresh time",
			},
			"soa_retry": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     900,
				Description: "SOA record retry time",
			},
			"soa_expire": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1209600,
				Description: "SOA record expire time",
			},
			"soa_minimum": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     3600,
				Description: "How long should negative responses be cached",
			},
			"ttl": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Time to live for records at zone apex",
			},
			"default_ttl": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Time to live for records without explicit TTL definition",
			},
			"dynamic_updates": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Allow dynamic updates",
			},
			"bind_update_policy": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "BIND update policy",
			},
			"allow_query": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "any",
				Description: "Semicolon separated list of IP addresses or networks which are allowed to issue queries",
			},
			"allow_transfer": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "none",
				Description: "Semicolon separated list of IP addresses or networks which are allowed to transfer the zone",
			},
			"zone_forwarders": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Per-zone forwarders. A custom port can be specified for each forwarder using a standard format IP_ADDRESS port PORT",
			},
			"allow_prt_sync": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Allow synchronization of forward (A, AAAA) and reverse (PTR) records in the zone",
			},
			"allow_inline_dnssec_signing": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Allow inline DNSSEC signing of records in the zone",
			},
			"nsec3param_record": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "NSEC3PARAM record for zone in format: hash_algorithm flags iterations salt",
			},
		},
	}
}

func resourceFreeIPADNSDNSZoneCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Creating freeipa dns zone")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	optArgs := ipa.DnszoneAddOptionalArgs{}
	args := ipa.DnszoneAddArgs{}

	// if _v, ok := d.GetOkExists("soa_serial_number"); ok {
	// 	args.Idnssoaserial = _v.(int)
	// }

	if d.Get("is_reverse_zone").(bool) {
		if _v, ok := d.GetOkExists("zone_name"); ok {
			v := _v.(string)
			optArgs.NameFromIP = &v
		}
	} else {
		if _v, ok := d.GetOkExists("zone_name"); ok {
			v := _v.(any)
			optArgs.Idnsname = &v
		}
	}

	if _v, ok := d.GetOkExists("skip_overlap_check"); ok {
		v := _v.(bool)
		optArgs.SkipOverlapCheck = &v
	}
	if _v, ok := d.GetOkExists("skip_nameserver_check"); ok {
		v := _v.(bool)
		optArgs.SkipNameserverCheck = &v
	}
	if _v, ok := d.GetOkExists("authoritative_nameserver"); ok {
		v := _v.(any)
		optArgs.Idnssoamname = &v
	}
	if _v, ok := d.GetOkExists("admin_email_address"); ok {
		v := _v.(any)
		optArgs.Idnssoarname = &v
	}
	if _v, ok := d.GetOkExists("soa_refresh"); ok {
		v := _v.(int)
		optArgs.Idnssoarefresh = &v
	}
	if _v, ok := d.GetOkExists("soa_retry"); ok {
		v := _v.(int)
		optArgs.Idnssoaretry = &v
	}
	if _v, ok := d.GetOkExists("soa_expire"); ok {
		v := _v.(int)
		optArgs.Idnssoaexpire = &v
	}
	if _v, ok := d.GetOkExists("soa_minimum"); ok {
		v := _v.(int)
		optArgs.Idnssoaminimum = &v
	}
	if _v, ok := d.GetOkExists("ttl"); ok {
		v := _v.(int)
		optArgs.Dnsttl = &v
	}
	if _v, ok := d.GetOkExists("default_ttl"); ok {
		v := _v.(int)
		optArgs.Dnsdefaultttl = &v
	}
	if _v, ok := d.GetOkExists("dynamic_updates"); ok {
		v := _v.(bool)
		optArgs.Idnsallowdynupdate = &v
	}
	if _v, ok := d.GetOkExists("bind_update_policy"); ok {
		v := _v.(string)
		optArgs.Idnsupdatepolicy = &v
	}
	if _v, ok := d.GetOkExists("allow_query"); ok {
		v := _v.(string)
		optArgs.Idnsallowquery = &v
	}
	if _v, ok := d.GetOkExists("allow_transfer"); ok {
		v := _v.(string)
		optArgs.Idnsallowtransfer = &v
	}
	if _v, ok := d.GetOkExists("zone_forwarders"); ok {
		v := utilsGetArry(_v.([]interface{}))
		optArgs.Idnsforwarders = &v
	}
	if _v, ok := d.GetOkExists("allow_prt_sync"); ok {
		v := _v.(bool)
		optArgs.Idnsallowsyncptr = &v
	}
	if _v, ok := d.GetOkExists("allow_inline_dnssec_signing"); ok {
		v := _v.(bool)
		optArgs.Idnssecinlinesigning = &v
	}
	if _v, ok := d.GetOkExists("nsec3param_record"); ok {
		v := _v.(string)
		optArgs.Nsec3paramrecord = &v
	}
	res, err := client.DnszoneAdd(&args, &optArgs)
	if err != nil {
		return diag.Errorf("Error creating freeipa dns zone: %s", err)
	}

	d.SetId(res.Result.Idnsname.(string))

	if _v, ok := d.GetOkExists("disable_zone"); ok {
		v := _v.(bool)
		if v {
			var id any = d.Id()
			_, err = client.DnszoneDisable(&ipa.DnszoneDisableArgs{}, &ipa.DnszoneDisableOptionalArgs{Idnsname: &id})
			if err != nil {
				log.Printf("[DEBUG] DNS zone disable/enable. Something went wrong: %s", err)
			}
		}
	}

	return resourceFreeIPADNSDNSZoneRead(ctx, d, meta)
}

func resourceFreeIPADNSDNSZoneRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Read freeipa dns zone")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	all := true
	var id any = d.Id()
	optArgs := ipa.DnszoneShowOptionalArgs{
		All:      &all,
		Idnsname: &id,
	}

	res, err := client.DnszoneShow(&ipa.DnszoneShowArgs{}, &optArgs)
	if err != nil {
		if strings.Contains(err.Error(), "NotFound") {
			d.SetId("")
			log.Printf("[DEBUG] DNS Zone not found")
			return nil
		} else {
			return diag.Errorf("Error reading freeipa DNS zone: %s", err)
		}
	}

	d.Set("disable_zone", !*res.Result.Idnszoneactive)

	log.Printf("[DEBUG] Read freeipa dns zone %s", res.Result.Idnsname)
	return nil
}

func resourceFreeIPADNSDNSZoneUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Update freeipa dns zone")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	var id any = d.Id()
	optArgs := ipa.DnszoneModOptionalArgs{
		Idnsname: &id,
	}

	var hasChange = false

	if d.HasChange("authoritative_nameserver") {
		if _v, ok := d.GetOkExists("authoritative_nameserver"); ok {
			v := _v.(interface{})
			if v != "" {
				optArgs.Idnssoamname = &v
				hasChange = true
			}
		}
	}
	if d.HasChange("admin_email_address") {
		if _v, ok := d.GetOkExists("admin_email_address"); ok {
			v := _v.(interface{})
			if v != "" {
				optArgs.Idnssoarname = &v
				hasChange = true
			}
		}
	}
	if d.HasChange("soa_serial_number") {
		if _v, ok := d.GetOkExists("soa_serial_number"); ok {
			v := _v.(int)
			if v != 0 {
				optArgs.Idnssoaserial = &v
				hasChange = true
			}
		}
	}
	if d.HasChange("soa_refresh") {
		if _v, ok := d.GetOkExists("soa_refresh"); ok {
			v := _v.(int)
			optArgs.Idnssoarefresh = &v
			hasChange = true
		}
	}
	if d.HasChange("soa_retry") {
		if _v, ok := d.GetOkExists("soa_retry"); ok {
			v := _v.(int)
			optArgs.Idnssoaretry = &v
			hasChange = true
		}
	}
	if d.HasChange("soa_expire") {
		if _v, ok := d.GetOkExists("soa_expire"); ok {
			v := _v.(int)
			optArgs.Idnssoaexpire = &v
			hasChange = true
		}
	}
	if d.HasChange("soa_minimum") {
		if _v, ok := d.GetOkExists("soa_minimum"); ok {
			v := _v.(int)
			optArgs.Idnssoaminimum = &v
			hasChange = true
		}
	}
	if d.HasChange("ttl") {
		if _v, ok := d.GetOkExists("ttl"); ok {
			v := _v.(int)
			optArgs.Dnsttl = &v
			hasChange = true
		}
	}
	if d.HasChange("default_ttl") {
		if _v, ok := d.GetOkExists("default_ttl"); ok {
			v := _v.(int)
			optArgs.Dnsdefaultttl = &v
			hasChange = true
		}
	}
	if d.HasChange("dynamic_updates") {
		if _v, ok := d.GetOkExists("dynamic_updates"); ok {
			v := _v.(bool)
			optArgs.Idnsallowdynupdate = &v
			hasChange = true
		}
	}
	if d.HasChange("allow_prt_sync") {
		if _v, ok := d.GetOkExists("allow_prt_sync"); ok {
			v := _v.(bool)
			optArgs.Idnsallowsyncptr = &v
			hasChange = true
		}
	}
	if d.HasChange("allow_inline_dnssec_signing") {
		if _v, ok := d.GetOkExists("allow_inline_dnssec_signing"); ok {
			v := _v.(bool)
			optArgs.Idnssecinlinesigning = &v
			hasChange = true
		}
	}
	if d.HasChange("bind_update_policy") {
		if _v, ok := d.GetOkExists("bind_update_policy"); ok {
			v := _v.(string)
			optArgs.Idnsupdatepolicy = &v
			hasChange = true
		}
	}
	if d.HasChange("allow_query") {
		if _v, ok := d.GetOkExists("allow_query"); ok {
			v := _v.(string)
			optArgs.Idnsallowquery = &v
			hasChange = true
		}
	}
	if d.HasChange("allow_transfer") {
		if _v, ok := d.GetOkExists("allow_transfer"); ok {
			v := _v.(string)
			optArgs.Idnsallowtransfer = &v
			hasChange = true
		}
	}
	if d.HasChange("zone_forwarders") {
		if _v, ok := d.GetOkExists("zone_forwarders"); ok {
			v := utilsGetArry(_v.([]interface{}))
			optArgs.Idnsforwarders = &v
			hasChange = true
		}
	}
	if d.HasChange("nsec3param_record") {
		if _v, ok := d.GetOkExists("nsec3param_record"); ok {
			v := _v.(string)
			optArgs.Nsec3paramrecord = &v
			hasChange = true
		}
	}

	if hasChange {
		_, err = client.DnszoneMod(&ipa.DnszoneModArgs{}, &optArgs)
		if err != nil {
			if strings.Contains(err.Error(), "EmptyModlist") {
				log.Printf("[DEBUG] EmptyModlist (4202): no modifications to be performed")
			} else {
				return diag.Errorf("Error update freeipa dns zone: %s", err)
			}
		}
	}

	if d.HasChange("disable_zone") {
		if _v, ok := d.GetOkExists("disable_zone"); ok {
			v := _v.(bool)
			if v {
				_, err = client.DnszoneDisable(&ipa.DnszoneDisableArgs{}, &ipa.DnszoneDisableOptionalArgs{Idnsname: &id})
			} else {
				_, err = client.DnszoneEnable(&ipa.DnszoneEnableArgs{}, &ipa.DnszoneEnableOptionalArgs{Idnsname: &id})
			}
			if err != nil {
				log.Printf("[DEBUG] DNS zone disable/enable. Something went wrong: %s", err)
			}
		}
	}

	return resourceFreeIPADNSDNSZoneRead(ctx, d, meta)
}

func resourceFreeIPADNSDNSZoneDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Delete freeipa dns zone")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	strIds := []string{d.Id()}
	id := make([]interface{}, len(strIds))
	for i, v := range strIds {
		id[i] = v
	}
	optArgs := ipa.DnszoneDelOptionalArgs{
		Idnsname: &id,
	}
	_, err = client.DnszoneDel(&ipa.DnszoneDelArgs{}, &optArgs)
	if err != nil {
		return diag.Errorf("Error delete freeipa dns zone: %s", err)
	}

	d.SetId("")
	return nil
}
