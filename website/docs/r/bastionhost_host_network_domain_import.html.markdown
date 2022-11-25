---
subcategory: "Bastion Host"
layout: "alicloud"
page_title: "Alicloud: alicloud_bastionhost_host_network_domain_import"
sidebar_current: "docs-alicloud-datasource-bastionhost-host-network-domain-import"
description: |-
Provides a Alicloud Bastionhost Host Network Domain import Resource.
---

# alicloud\_bastionhost\_host\_network\_domain\_import

Provides a Alicloud Bastionhost Host Network Domain import Resource.

-> **NOTE:** Available in v1.188.3+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_bastionhost_host_network_domain_import" "default" {
    instance_id = "bastionhost-cn-xxxxxxxxx"
    network_domain_id = "2"
    host_ids = [ "26", "33" ]
    database_ids = [ "271", "272" ]
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) You Want to Query the User the Bastion Host ID of.
* `network_domain_id` - (Required, Computed) The ID of Bastionhost Network Domain.
* `host_ids` - (Optional) The list of host ID which imported to Bastionhost.
* `database_ids` - (Optional) The list of host ID which imported to Bastionhost.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Bastionhost Network Domain. The value formats as `<instance_id>:<network_domain_id>`.

## Import

Bastion Host Network Domain can be imported using the id, e.g.

```
$ terraform import alicloud_bastionhost_host_network_domain_import.example <instance_id>:<network_domain_id>
```