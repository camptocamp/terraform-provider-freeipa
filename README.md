FreeIPA Terraform Provider
==========================
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
  host = "ipa.example.test"
  username = "admin"
  password = "P@S5sw0rd"
  insecure = true
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
```

Usage
----------------------

