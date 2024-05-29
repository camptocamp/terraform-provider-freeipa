FreeIPA Terraform Provider
==========================

[![Terraform Registry Version](https://img.shields.io/badge/dynamic/json?color=blue&label=registry&query=%24.version&url=https%3A%2F%2Fregistry.terraform.io%2Fv1%2Fproviders%2Fcamptocamp%2Ffreeipa)](https://registry.terraform.io/providers/camptocamp/freeipa)
[![Go Report Card](https://goreportcard.com/badge/github.com/camptocamp/terraform-provider-freeipa)](https://goreportcard.com/report/github.com/camptocamp/terraform-provider-freeipa)
[![Build Status](https://travis-ci.org/camptocamp/terraform-provider-freeipa.svg?branch=master)](https://travis-ci.org/camptocamp/terraform-provider-freeipa)
[![By Camptocamp](https://img.shields.io/badge/by-camptocamp-fb7047.svg)](http://www.camptocamp.com)

This provider adds integration between Terraform and FreeIPA.

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.12.x
-	[Go](https://golang.org/doc/install) 1.10


Building The Provider
---------------------

Download the provider source code

```sh
$ go get github.com/camptocamp/terraform-provider-freeipa
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/camptocamp/terraform-provider-freeipa
$ make build
```

Installing the provider
-----------------------

After building the provider, install it using the Terraform instructions for [installing a third party provider](https://www.terraform.io/docs/configuration/providers.html#third-party-plugins).

Example
----------------------

```hcl
provider freeipa {
  host = "ipa.example.test"   # or set $FREEIPA_HOST
  username = "admin"          # or set $FREEIPA_USERNAME
  password = "P@S5sw0rd"      # or set $FREEIPA_PASSWORD
}

resource freeipa_host "foo" {
  fqdn = "foo.example.test"
  description = "This is my foo host"
  force = true
  random = true
}

resource freeipa_host "bar" {
  fqdn = "bar.example.test"
  userpassword = "abcde"
}

resource freeipa_dns_record "foo" {
  dnszoneidnsname = "your.zone.name."
  idnsname = "foo"
  records = ["192.168.10.10"]
  type = "A"
}
```

Usage
----------------------


Import
------

DNS records can be imported using the record name and the zone name from <record_name>/<zone_name>/\<type\>

```
$ terraform import freeipa_dns_record.foo foo/example.tld./A
```
