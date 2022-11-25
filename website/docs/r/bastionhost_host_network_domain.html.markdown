---
subcategory: "Bastion Host"
layout: "alicloud"
page_title: "Alicloud: alicloud_bastionhost_host_network_domain"
sidebar_current: "docs-alicloud-datasource-bastionhost-host-network-domain"
description: |-
Provides a Alicloud Bastionhost Host Network Domain Resource.
---

# alicloud\_bastionhost\_host\_network\_domain

Provides a Alicloud Bastionhost Host Network Domain Resource.

-> **NOTE:** Available in v1.188.3+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_bastionhost_host_network_domain" "default" {
    instance_id = "bastionhost-cn-xxxxxxxxxx"
    network_domain_name = "Proxy"
    network_domain_type = "Proxy"
    comment = "Proxy"
    proxies = jsonencode([
        {
            ProxyId = "32"
            ProxyType = "Socks5Proxy"
            Weight = 100
            Socks5ProxyConfig = {
                Address         = "10.35.2.215"
                Port            = 10088
                User            = "ops"
                Password        = base64encode("123456")
            }
        }
    ])
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) You Want to Query the User the Bastion Host ID of.
* `comment` - (Optional) The comment of the entry, Supports up to 128 Characters.
* `proxies` - (Optional) The list of proxy.
* `network_domain_name` - (Optional) The name of Bastionhost Network Domain.
* `network_domain_type` - (Optional) The type of Bastionhost Network Domain. Types: Socks5Proxy、SSHProxy、HTTPProxy.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Bastionhost Network Domain. The value formats as `<instance_id>:<network_domain_id>`.
* `network_domain_id` - (Optional, Computed) The ID of Bastionhost Network Domain.

## Import

Bastion Host User can be imported using the id, e.g.

```
$ terraform import alicloud_bastionhost_host_network_domain.example <instance_id>:<network_domain_id>
```