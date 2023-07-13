---
layout: "freeipa"
page_title: "FreeIPA: freeipa_dns_record"
sidebar_current: "docs-freeipa-resource-dns-record
description: "Provides a FreeIPA DNS record resource. This cas be used to create DNS records in FreeIPA DNS zones."
---

# freeipa_dns_record

Provides a FreeIPA DNS record resource. This cas be used to create DNS records in FreeIPA DNS zones.

## Example Usage

```hcl
resource "freeipa_dns_record" "at" {
  idnsname = "@"
  dnszoneidnsname = "example.tld"
  dnsttl = 300
  type = "A"
  records = ["10.10.10.10"]
}
```

## Argument Reference

The following arguments are supported:

* `idnsname` - (Required) The record DNS without zone.
* `dnszoneidnsname` - (Required) The name of the DNS zone.
* `type` - (Required) The type of DNS record. The following types are supported: `A`, `AAAA`, `CNAME`, `MX`, `NS`, `PTR`, `SRV`, `TXT`, `SSHFP`.
* `records` - (Required) List of records of the same type.
* `dnsttl` - (Optional) Time-To-Live of the DNS record.

## Import

DNS records can be imported using the record name, the zone name and the record type from `<record_name>/<zone_name>/<type>`.

```shell
$ terraform import freeipa_dns_record.foo foo/example.tld./A
```
