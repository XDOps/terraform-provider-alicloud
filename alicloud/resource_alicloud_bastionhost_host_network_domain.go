package alicloud

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudBastionhostHostNetworkDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudBastionhostHostNetworkDomainCreate,
		Read:   resourceAlicloudBastionhostHostNetworkDomainRead,
		Update: resourceAlicloudBastionhostHostNetworkDomainUpdate,
		Delete: resourceAlicloudBastionhostHostNetworkDomainDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"comment": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 128),
			},
			"proxies": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: BastionhostNetworkDomainDiffSuppressFunc,
			},
			"network_domain_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"network_domain_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"network_domain_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudBastionhostHostNetworkDomainCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateNetworkDomain"
	request := make(map[string]interface{})
	conn, err := client.NewBastionhostClient()
	if err != nil {
		return WrapError(err)
	}

	request["InstanceId"] = d.Get("instance_id")
	request["NetworkDomainName"] = d.Get("network_domain_name")
	request["NetworkDomainType"] = d.Get("network_domain_type")

	if request["NetworkDomainType"] == "Proxy" {
		var proxies []ProxyStruct

		p := d.Get("proxies")
		if len(p.(string)) > 0 {
			e := json.Unmarshal([]byte(p.(string)), &proxies)
			if e != nil {
				return WrapError(e)
			}

			var data []interface{}
			for i := 0; i < len(proxies); i++ {
				proxy := proxies[i]
				tmp := map[string]interface{}{}

				tmp["ProxyType"] = proxy.ProxyType
				tmp["Weight"] = proxy.Weight
				if proxy.ProxyType == "HTTPProxy" {
					tmp["HTTPProxyConfig"] = map[string]interface{}{
						"Address":  proxy.HTTPProxyConfig.Address,
						"Port":     proxy.HTTPProxyConfig.Port,
						"User":     proxy.HTTPProxyConfig.User,
						"Password": proxy.HTTPProxyConfig.Password,
					}
					data = append(data, tmp)
				} else if proxy.ProxyType == "Socks5Proxy" {
					tmp["Socks5ProxyConfig"] = map[string]interface{}{
						"Address":  proxy.Socks5ProxyConfig.Address,
						"Port":     proxy.Socks5ProxyConfig.Port,
						"User":     proxy.Socks5ProxyConfig.User,
						"Password": proxy.Socks5ProxyConfig.Password,
					}
					data = append(data, tmp)
				} else if proxy.ProxyType == "SSHProxy" {
					tmp["SSHProxyConfig"] = map[string]interface{}{
						"Address":  proxy.SSHProxyConfig.Address,
						"Port":     proxy.SSHProxyConfig.Port,
						"User":     proxy.SSHProxyConfig.User,
						"Password": proxy.SSHProxyConfig.Password,
					}
					data = append(data, tmp)
				}
			}

			r, e := json.Marshal(data)
			if e == nil {
				request["Proxies"] = string(r)
			}
		}
	}

	if v, ok := d.GetOk("comment"); ok {
		request["Comment"] = v.(string)
	}

	//wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-30"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			//if NeedRetry(err) {
			//	wait()
			//	return resource.RetryableError(err)
			//}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_bastionhost_host_network_domain", action, AlibabaCloudSdkGoERROR)
	}

	networkDomainId, err := jsonpath.Get("$.NetworkDomainId", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, d.Get("instance_id"), "$.NetworkDomainId", response)
	}

	d.SetId(fmt.Sprint(request["InstanceId"], ":", networkDomainId))

	return resourceAlicloudBastionhostHostNetworkDomainRead(d, meta)
}

func resourceAlicloudBastionhostHostNetworkDomainRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "GetNetworkDomain"
	conn, err := client.NewBastionhostClient()
	if err != nil {
		return WrapError(err)
	}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"InstanceId":      parts[0],
		"NetworkDomainId": parts[1],
	}

	log.Printf("-- %v --", request)

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-30"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_bastionhost_host_network_domain", action, AlibabaCloudSdkGoERROR)
	}

	object, err := jsonpath.Get("$.NetworkDomain", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, parts[0], "$.NetworkDomains", response)
	}

	networkDomain := object.(map[string]interface{})
	proxies := make([]interface{}, 0)

	var proxyStruct []ProxyStruct
	if v, ok := d.GetOk("proxies"); ok {
		err := json.Unmarshal([]byte(v.(string)), &proxyStruct)
		if err != nil {
			return WrapError(err)
		}
	}

	for i := 0; i < len(networkDomain["Proxies"].([]interface{})); i++ {
		item := networkDomain["Proxies"].([]interface{})[i].(map[string]interface{})
		var proxyKey = ""
		if item["ProxyType"] == "HTTPProxy" {
			proxyKey = "HTTPProxyConfig"
		} else if item["ProxyType"] == "Socks5Proxy" {
			proxyKey = "Socks5ProxyConfig"
		} else if item["ProxyType"] == "SSHProxy" {
			proxyKey = "SSHProxyConfig"
		} else {
			break
		}

		proxyConfig := item[proxyKey].(map[string]interface{})
		password := getPasswordFromProxies(proxyStruct, item)
		proxyValue := map[string]interface{}{
			"Address":  proxyConfig["Address"],
			"Port":     proxyConfig["Port"],
			"User":     proxyConfig["User"],
			"Password": password,
		}

		proxy := map[string]interface{}{
			"ProxyId":   item["ProxyId"],
			"ProxyType": item["ProxyType"],
			proxyKey:    proxyValue,
			"Weight":    item["Weight"],
		}
		proxies = append(proxies, proxy)
	}

	data, err := json.Marshal(proxies)
	if err != nil {
		return WrapError(err)
	}

	d.Set("instance_id", request["InstanceId"])
	d.Set("comment", networkDomain["Comment"].(string))
	d.Set("network_domain_id", networkDomain["NetworkDomainId"].(string))
	d.Set("network_domain_name", networkDomain["NetworkDomainName"].(string))
	d.Set("network_domain_type", networkDomain["NetworkDomainType"].(string))

	if err := d.Set("proxies", string(data)); err != nil {
		return WrapError(err)
	}

	return nil
}

func resourceAlicloudBastionhostHostNetworkDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	conn, err := client.NewBastionhostClient()
	if err != nil {
		return WrapError(err)
	}
	var response map[string]interface{}
	request := make(map[string]interface{})
	d.Partial(true)
	update := false

	request["InstanceId"] = d.Get("instance_id")
	oId, nId := d.GetChange("network_domain_id")
	if nId != nil && len(oId.(string)) > 0 {
		request["NetworkDomainId"] = nId
	} else {
		request["NetworkDomainId"] = oId
	}

	if d.HasChange("proxies") {
		var oldObjects, newObjects []ProxyStruct

		o, n := d.GetChange("proxies")

		e := json.Unmarshal([]byte(o.(string)), &oldObjects)
		if e != nil {
			return WrapError(e)
		}

		e = json.Unmarshal([]byte(n.(string)), &newObjects)
		if e != nil {
			return WrapError(e)
		}

		addObjects, removeObjects, updateObjects := compareProxies(oldObjects, newObjects)

		if len(removeObjects) > 0 {
			for _, removeProxy := range removeObjects {
				removeRequest := map[string]interface{}{
					"InstanceId":           request["InstanceId"],
					"NetworkDomainProxyId": removeProxy.ProxyId,
				}

				action := "DeleteNetworkDomainProxy"
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-30"), StringPointer("AK"), nil, removeRequest, &util.RuntimeOptions{})
				addDebug(action, response, request)
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, "alicloud_bastionhost_host_network_domain", action, AlibabaCloudSdkGoERROR)
				}

				d.SetPartial("proxies")
			}
		}

		if len(addObjects) > 0 {

			for _, addProxy := range addObjects {
				addRequest := map[string]interface{}{
					"InstanceId":           request["InstanceId"],
					"NetworkDomainId":      request["NetworkDomainId"],
					"NetworkDomainProxyId": addProxy.ProxyId,
					"ProxyType":            addProxy.ProxyType,
					"Weight":               addProxy.Weight,
					"HTTPProxyConfig":      "{\"Address\":\"\",\"Port\":0,\"User\":\"\",\"IsSetPassword\":false}",
					"Socks5ProxyConfig":    "{\"Address\":\"\",\"Port\":0,\"User\":\"\",\"IsSetPassword\":false}",
					"SSHProxyConfig":       "{\"Address\":\"\",\"Port\":0,\"User\":\"\",\"IsSetPassword\":false}",
				}

				var keyName = ""
				var keyValue map[string]interface{}
				if addProxy.ProxyType == "HTTPProxy" {
					keyName = "HTTPProxyConfig"
					keyValue = map[string]interface{}{
						"Address":  addProxy.HTTPProxyConfig.Address,
						"Port":     addProxy.HTTPProxyConfig.Port,
						"User":     addProxy.HTTPProxyConfig.User,
						"Password": addProxy.HTTPProxyConfig.Password,
					}
				} else if addProxy.ProxyType == "Socks5Proxy" {
					keyName = "Socks5ProxyConfig"
					keyValue = map[string]interface{}{
						"Address":  addProxy.Socks5ProxyConfig.Address,
						"Port":     addProxy.Socks5ProxyConfig.Port,
						"User":     addProxy.Socks5ProxyConfig.User,
						"Password": addProxy.Socks5ProxyConfig.Password,
					}
				} else if addProxy.ProxyType == "SSHProxy" {
					keyName = "SSHProxyConfig"
					keyValue = map[string]interface{}{
						"Address":  addProxy.SSHProxyConfig.Address,
						"Port":     addProxy.SSHProxyConfig.Port,
						"User":     addProxy.SSHProxyConfig.User,
						"Password": addProxy.SSHProxyConfig.Password,
					}
				}

				r, e := json.Marshal(keyValue)
				if e != nil {
					return WrapError(e)
				}

				addRequest[keyName] = string(r)

				action := "CreateNetworkDomainProxy"
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-30"), StringPointer("AK"), nil, addRequest, &util.RuntimeOptions{})
				addDebug(action, response, request)
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, "alicloud_bastionhost_host_network_domain", action, AlibabaCloudSdkGoERROR)
				}
				d.SetPartial("proxies")
			}
		}

		if len(updateObjects) > 0 {
			for _, updateProxy := range updateObjects {
				updateRequest := map[string]interface{}{
					"InstanceId":           request["InstanceId"],
					"NetworkDomainProxyId": updateProxy.ProxyId,
					"ProxyType":            updateProxy.ProxyType,
					"Weight":               updateProxy.Weight,
				}

				var keyName = ""
				var keyValue map[string]interface{}
				if updateProxy.ProxyType == "HTTPProxy" {
					keyName = "HTTPProxyConfig"
					keyValue = map[string]interface{}{
						"Address":  updateProxy.HTTPProxyConfig.Address,
						"Port":     updateProxy.HTTPProxyConfig.Port,
						"User":     updateProxy.HTTPProxyConfig.User,
						"Password": updateProxy.HTTPProxyConfig.Password,
					}
				} else if updateProxy.ProxyType == "Socks5Proxy" {
					keyName = "Socks5ProxyConfig"
					keyValue = map[string]interface{}{
						"Address":  updateProxy.Socks5ProxyConfig.Address,
						"Port":     updateProxy.Socks5ProxyConfig.Port,
						"User":     updateProxy.Socks5ProxyConfig.User,
						"Password": updateProxy.Socks5ProxyConfig.Password,
					}
				} else if updateProxy.ProxyType == "SSHProxy" {
					keyName = "SSHProxyConfig"
					keyValue = map[string]interface{}{
						"Address":  updateProxy.SSHProxyConfig.Address,
						"Port":     updateProxy.SSHProxyConfig.Port,
						"User":     updateProxy.SSHProxyConfig.User,
						"Password": updateProxy.SSHProxyConfig.Password,
					}
				}

				r, e := json.Marshal(keyValue)
				if e != nil {
					return WrapError(e)
				}

				updateRequest[keyName] = string(r)

				action := "ModifyNetworkDomainProxy"
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-30"), StringPointer("AK"), nil, updateRequest, &util.RuntimeOptions{})
				addDebug(action, response, request)
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, "alicloud_bastionhost_host_network_domain", action, AlibabaCloudSdkGoERROR)
				}
				d.SetPartial("proxies")
			}
		}
	}

	if d.HasChanges("comment", "network_domain_name", "network_domain_type") {
		update = true
		request["Comment"] = d.Get("comment")
		request["NetworkDomainName"] = d.Get("network_domain_name")
		request["NetworkDomainType"] = d.Get("network_domain_type")
	}

	if update {
		action := "ModifyNetworkDomain"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-30"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("proxies")
	}

	d.SetId(fmt.Sprint(request["InstanceId"], ":", request["NetworkDomainId"]))
	d.Partial(false)

	return resourceAlicloudBastionhostHostNetworkDomainRead(d, meta)
}

func resourceAlicloudBastionhostHostNetworkDomainDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	action := "DeleteNetworkDomain"
	var response map[string]interface{}
	conn, err := client.NewBastionhostClient()
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"InstanceId":      parts[0],
		"NetworkDomainId": parts[1],
	}
	request["RegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-30"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		if IsExpectedErrors(err, []string{"OBJECT_NOT_FOUND"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
