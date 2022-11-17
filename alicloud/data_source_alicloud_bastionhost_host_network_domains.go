package alicloud

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudBastionhostHostNetworkDomains() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudBastionhostHostNetworkDomainsRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"keyword": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
			},
			"network_domain_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"connect_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"network_domains": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"comment": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"database_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"http_proxy_config": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"address": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"is_set_password": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"user": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"host_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"is_build_in": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"network_domain_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_domain_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_domain_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"proxies_state": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"proxy_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"proxy_state": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"weight": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"ssh_proxy_config": {
							Type:     schema.TypeMap,
							Computed: true,
						},
						"socks5_config": {
							Type:     schema.TypeMap,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudBastionhostHostNetworkDomainsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListNetworkDomains"
	request := make(map[string]interface{})
	request["InstanceId"] = d.Get("instance_id")
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	if v, ok := d.GetOk("keyword"); ok {
		request["Keyword"] = v.(string)
	} else {
		request["Keyword"] = ""
	}
	if v, ok := d.GetOk("connect_type"); ok {
		request["NetworkDomainType"] = v.(string)
	} else {
		request["NetworkDomainType"] = ""
	}

	var result []interface{}
	var response map[string]interface{}
	conn, err := client.NewBastionhostClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-30"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_bastionhost_host_network_domians", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.NetworkDomains", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.NetworkDomains", response)
		}
		result = resp.([]interface{})
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range result {
		network_domain := object.(map[string]interface{})
		database_count, _ := network_domain["DatabaseCount"].(json.Number).Int64()
		host_count, _ := network_domain["HostCount"].(json.Number).Int64()

		http_proxy_config := make(map[string]interface{}, 0)
		for k, v := range network_domain["HTTPProxyConfig"].(map[string]interface{}) {
			if k == "Address" {
				http_proxy_config["address"] = v.(string)
			} else if k == "IsSetPassword" {
				http_proxy_config["is_set_password"] = fmt.Sprintf("%t", v)
			} else if k == "Port" {
				http_proxy_config["port"] = fmt.Sprintf("%v", v)
			} else if k == "User" {
				http_proxy_config["user"] = v.(string)
			}
		}

		proxies_state := make([]map[string]interface{}, 0)
		for _, v := range network_domain["ProxiesState"].([]interface{}) {
			item := v.(map[string]interface{})
			weight, _ := item["Weight"].(json.Number).Int64()
			state := map[string]interface{}{
				"proxy_id":    item["ProxyId"].(string),
				"proxy_state": item["ProxyState"].(string),
				"weight":      weight,
			}
			proxies_state = append(proxies_state, state)
		}

		ssh_proxy_config := make(map[string]interface{}, 0)
		for k, v := range network_domain["SSHProxyConfig"].(map[string]interface{}) {
			if k == "Address" {
				ssh_proxy_config["address"] = v.(string)
			} else if k == "IsSetPassword" {
				ssh_proxy_config["is_set_password"] = fmt.Sprintf("%t", v)
			} else if k == "Port" {
				ssh_proxy_config["port"] = fmt.Sprintf("%v", v)
			} else if k == "User" {
				ssh_proxy_config["user"] = v.(string)
			}
		}

		socks5_config := make(map[string]interface{}, 0)
		for k, v := range network_domain["Socks5Config"].(map[string]interface{}) {
			if k == "Address" {
				socks5_config["address"] = v.(string)
			} else if k == "IsSetPassword" {
				socks5_config["is_set_password"] = fmt.Sprintf("%t", v)
			} else if k == "Port" {
				socks5_config["port"] = fmt.Sprintf("%v", v)
			} else if k == "User" {
				socks5_config["user"] = v.(string)
			}
		}

		mapping := map[string]interface{}{
			"comment":             network_domain["Comment"].(string),
			"database_count":      database_count,
			"host_count":          host_count,
			"http_proxy_config":   http_proxy_config,
			"is_build_in":         network_domain["IsBuiltIn"].(bool),
			"network_domain_id":   network_domain["NetworkDomainId"].(string),
			"network_domain_name": network_domain["NetworkDomainName"].(string),
			"network_domain_type": network_domain["NetworkDomainType"].(string),
			"proxies_state":       proxies_state,
			"ssh_proxy_config":    ssh_proxy_config,
			"socks5_config":       socks5_config,
			"vpc_id":              network_domain["VpcId"].(string),
			"vpc_name":            network_domain["VpcName"].(string),
			"vpc_region_id":       network_domain["VpcRegionId"].(string),
		}
		ids = append(ids, fmt.Sprint(mapping["network_domain_id"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))

	if err := d.Set("keyword", request["Keyword"]); err != nil {
		return WrapError(err)
	}

	if err := d.Set("connect_type", request["NetworkDomainType"]); err != nil {
		return WrapError(err)
	}

	if err := d.Set("network_domain_ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("network_domains", s); err != nil {
		return WrapError(err)
	}

	return nil
}
