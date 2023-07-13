---
layout: "freeipa"
page_title: "Provider: FreeIPA"
description: |-
  This provider adds integration between Terraform and FreeIPA.
---

# FreeIPA Provider

This provider adds integration between Terraform and FreeIPA.

Use the navigation to the left to read about the available resources.

## Example Usage

Terraform 0.13 and later:

```hcl
provider "freeipa" {
  host = "ipa.example.test"   # or set $FREEIPA_HOST
  username = "admin"          # or set $FREEIPA_USERNAME
  password = "P@S5sw0rd"      # or set $FREEIPA_PASSWORD
  insecure = true
}
```

## Argument Reference

The following arguments are supported:

- `host` - (Required) FreeIPA API url. It must be provided, but it can also be sourced from the `FREEIPA_HOST` environment variable.
- `username` - (Optional) FreeIPA API username. It can also be sourced from the `FREEIPA_USERNAME` environment variable.
- `password` - (Optional) FreeIPA API password. It can also be sourced from the `FREEIPA_PASSWORD` environment variable.
- `insecure` - (Optional) Skip SSL/TLS verification.
