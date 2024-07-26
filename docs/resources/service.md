---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "freeipa_service Resource - freeipa"
subcategory: ""
description: |-
  
---

# freeipa_service (Resource)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `krb_hostname` (String) Principal name Service principal. Format: <service_type>/<hostname>

### Optional

- `force` (Boolean) Force force principal name even if host not in DNS
- `skip_host_check` (Boolean) Skip host check force service to be created even when host object does not exist to manage it