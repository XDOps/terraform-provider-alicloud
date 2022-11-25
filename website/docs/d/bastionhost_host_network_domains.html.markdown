---
subcategory: "Bastion Host"
layout: "alicloud"
page_title: "Alicloud: alicloud_bastionhost_host_network_domains"
sidebar_current: "docs-alicloud-datasource-bastionhost-host-network-domains"
description: |-
Provides a list of Bastionhost Host Network Domain.
---

# alicloud\_bastionhost\_host\_network\_domains

Provides a list of Bastionhost Host Network Domain.

-> **NOTE:** Available in v1.188.3+.

## Example Usage

Basic Usage

```terraform
data "alicloud_bastionhost_host_network_domains" "default" {
  instance_id  = "bastionhost-cn-xxxxx"
  keyword      = "proxy"
  connect_type = "Direct"
}
output "domains" {
  value = data.alicloud_bastionhost_host_network_domains.default
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) You Want to Query the User the Bastion Host ID of.
* `keyword` - (Optional) The keyword to be matched.
* `connect_type` - (Optional) The type of Bastionhost host connect Network Domain.

## Attributes Reference

The following attributes are exported:

* `network_domain_ids` - The list of Bastionhost Network Domain IDs.
* `network_domains` - The list of Bastionhost Network Domain
  - `comment` The comment of the Network Domain.
  - `database_count` The count of database which import to Network Domain.
  - `http_proxy_config` The proxy config of http.
  - `host_count` The count of host which import to Network Domain.
  - `is_build_in` A unknown result.
  - `network_domain_id` The ID of Bastionhost Network Domain.
  - `network_domain_name` The name of Bastionhost Network Domain.
  - `network_domain_type` The type of Bastionhost Network Domain. Types: Socks5Proxy、SSHProxy、HTTPProxy.
  - `proxies_state` The state of proxies.
  - `ssh_proxy_config` The proxy config of ssh.
  - `socks5_config` The proxy config of socks5.
  - `vpc_id` The id of the VPC.
  - `vpc_name` The name of the VPC.
  - `vpc_region_id` The Region of the VPC.
