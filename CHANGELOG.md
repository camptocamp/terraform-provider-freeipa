## 0.9.0 (May 22, 2024)

IMPROVEMENTS:

* Honour `HTTP_PROXY`, `HTTPS_PROXY` and `NO_PROXY` environment variables

## 0.8.1 (May 21, 2024)

BUG FIXES:

* Fix provider login for newer FreeIPA versions requiring the `Referer` HTTP request header

## 0.8.0 (June 23, 2023)

IMPROVEMENTS:

* Add “managed by hosts” support to `freeipa_host` resource
* Use Go 1.20
* Replace archived FreeIPA Go library with new maintained fork
* Migrate to Terraform plugin framework

## 0.7.0 (October 29, 2020)

* `freeipa_dns_record`: improve import function

## 0.6.0 (July 24, 2020)

* Release on Terraform registry

## 0.5.0 (April 10, 2020)

BREAKING CHANGES:

* resource/freeipa_dns_record: Add `type` and `records` arguments and remove `arecord`, `a_part_ip_address`, `srvrecord`, `srv_part_priority`, `srv_part_weight`, `srv_part_port` and `srv_part_target` arguments

BUG FIXES:

* resource/freeipa_dns_record: DNS records with multiple values where not always properly created.
* resource/freeipa_dns_record: DNS records with multiple values where not always properly deleted.

## 0.4.0 (April 9, 2020)

IMPROVEMENTS:

* resource/freeipa_dns_record: Add `arecord`, `srvrecord`, `srv_part_priority`, `srv_part_weight`, `srv_part_port` and `srv_part_target` attributes

BUG FIXES:

* Fix multi-valued records when using `arecord` or `srvrecord`

## 0.3.1 (April 8, 2020)

BUG FIXES:

* Fix `freeipa_dns_record` resource

## 0.3.0 (April 8, 2020)

IMPROVEMENTS:

* Add `freeipa_dns_record` resource

## 0.2.0 (February 26, 2020)

IMPROVEMENTS:

* Use new Terraform plugin SDK

BUG FIXES:

* Loop on read after create

## 0.1.0 (February 19, 2020)

* Initial release
