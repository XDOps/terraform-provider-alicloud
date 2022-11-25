package alicloud

import (
	"fmt"
	"github.com/PaesslerAG/jsonpath"
	"log"
	"strconv"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudBastionhostHostAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudBastionhostHostAccountCreate,
		Read:   resourceAlicloudBastionhostHostAccountRead,
		Update: resourceAlicloudBastionhostHostAccountUpdate,
		Delete: resourceAlicloudBastionhostHostAccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"host_account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"host_account_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"host_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"pass_phrase": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("protocol_name"); ok && v.(string) == "SSH" {
						return false
					}
					return true
				},
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"private_key": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("protocol_name"); ok && v.(string) == "SSH" {
						return false
					}
					return true
				},
			},
			"protocol_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"RDP", "SSH"}, false),
			},
			"port": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudBastionhostHostAccountCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateHostAccount"
	request := make(map[string]interface{})
	conn, err := client.NewBastionhostClient()
	if err != nil {
		return WrapError(err)
	}
	request["HostAccountName"] = d.Get("host_account_name")
	request["HostId"] = d.Get("host_id")
	request["InstanceId"] = d.Get("instance_id")
	if v, ok := d.GetOk("pass_phrase"); ok {
		request["PassPhrase"] = v
	}
	if v, ok := d.GetOk("password"); ok {
		request["Password"] = v
	}
	if v, ok := d.GetOk("private_key"); ok {
		request["PrivateKey"] = v
	}
	request["ProtocolName"] = d.Get("protocol_name")
	request["RegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-12-09"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_bastionhost_host_account", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["InstanceId"], ":", response["HostAccountId"]))

	if v, ok := d.GetOk("port"); ok && (v != 22 && v != 3389) {
		req := map[string]interface{}{
			"InstanceId":   request["InstanceId"],
			"RegionId":     client.RegionId,
			"HostIds":      fmt.Sprintf("[\"%v\"]", request["HostId"]),
			"ProtocolName": request["ProtocolName"],
			"Port":         fmt.Sprintf("%v", v),
		}

		action = "ModifyHostsPort"
		err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-12-09"), StringPointer("AK"), nil, req, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, req)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "alicloud_bastionhost_host_account", action, AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAlicloudBastionhostHostAccountRead(d, meta)
}
func resourceAlicloudBastionhostHostAccountRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	conn, err := client.NewBastionhostClient()
	if err != nil {
		return WrapError(err)
	}
	yundunBastionhostService := YundunBastionhostService{client}
	object, err := yundunBastionhostService.DescribeBastionhostHostAccount(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_bastionhost_host_account yundunBastionhostService.DescribeBastionhostHostAccount Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	port := 22
	if object["ProtocolName"] == "RDP" {
		port = 3389
	}

	if v, ok := d.GetOk("port"); ok {
		port = v.(int)
	} else {
		var response map[string]interface{}
		request := map[string]interface{}{
			"InstanceId": parts[0],
			"HostId":     object["HostId"],
			"RegionId":   client.RegionId,
		}

		action := "GetHost"
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-12-09"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "alicloud_bastionhost_host_account", action, AlibabaCloudSdkGoERROR)
		}
		obj, e := jsonpath.Get("$.Host", response)
		if e != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, d.Id(), "$.Host", response)
		}
		host := obj.(map[string]interface{})
		for _, p := range host["Protocols"].([]interface{}) {
			protocol := p.(map[string]interface{})
			if protocol["ProtocolName"].(string) == object["ProtocolName"].(string) {
				res, er := strconv.Atoi(fmt.Sprintf("%v", protocol["Port"]))
				if er == nil {
					port = res
				}
			}
		}
	}

	d.Set("host_account_id", parts[1])
	d.Set("instance_id", parts[0])
	d.Set("host_account_name", object["HostAccountName"])
	d.Set("host_id", object["HostId"])
	d.Set("protocol_name", object["ProtocolName"])
	d.Set("port", port)

	return nil
}
func resourceAlicloudBastionhostHostAccountUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	conn, err := client.NewBastionhostClient()
	if err != nil {
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	update := false

	if d.HasChange("port") {
		request := map[string]interface{}{
			"InstanceId":   parts[0],
			"HostIds":      fmt.Sprintf("[\"%v\"]", d.Get("host_id")),
			"ProtocolName": d.Get("protocol_name"),
			"Port":         fmt.Sprintf("%v", d.Get("port")),
			"RegionId":     client.RegionId,
		}

		action := "ModifyHostsPort"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-12-09"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	}

	request := map[string]interface{}{
		"HostAccountId": parts[1],
		"InstanceId":    parts[0],
	}
	if d.HasChange("host_account_name") {
		update = true
		request["HostAccountName"] = d.Get("host_account_name")
	}
	if d.HasChange("pass_phrase") {
		update = true
		if v, ok := d.GetOk("pass_phrase"); ok {
			request["PassPhrase"] = v
		}
	}
	if d.HasChange("password") {
		update = true
		if v, ok := d.GetOk("password"); ok {
			request["Password"] = v
		}
	}
	if d.HasChange("private_key") {
		update = true
		if v, ok := d.GetOk("private_key"); ok {
			request["PrivateKey"] = v
		}
	}
	request["RegionId"] = client.RegionId
	if update {
		action := "ModifyHostAccount"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-12-09"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	}
	return resourceAlicloudBastionhostHostAccountRead(d, meta)
}
func resourceAlicloudBastionhostHostAccountDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	action := "DeleteHostAccount"
	var response map[string]interface{}
	conn, err := client.NewBastionhostClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"HostAccountId": parts[1],
		"InstanceId":    parts[0],
	}

	request["RegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-12-09"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		if IsExpectedErrors(err, []string{"Commodity.BizError.InvalidStatus", "OBJECT_NOT_FOUND"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
