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
