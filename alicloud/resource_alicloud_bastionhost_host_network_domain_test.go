package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudBastionhostHostNetworkDomain_Basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_bastionhost_host_network_domains.default"
	ra := resourceAttrInit(resourceId, AlicloudBastionhostHostNetworkDomain0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &YundunBastionhostService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "ListBastionhostHostNetworkDomains")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sbastionhosthostnetworkdomain%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudBastionhostHostNetworkDomainBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id": "${data.alicloud_bastionhost_instances.default.ids.0}",
					"keyword":     "tf-testHostDomainName12345",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id": CHECKSET,
						"keyword":     "tf-testHostDomainName12345",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"comment": "tf-testHostDomainComment12345",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"comment": "tf-testHostDomainComment12345",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"keyword": "tf-testHostDomainpName12345-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"keyword": "tf-testHostDomainName12345-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"comment":             "tf-testHostDomainComment12345update",
					"network_domain_name": "tf-testHostDomainName12345",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"comment":             "tf-testHostDomainComment12345update",
						"network_domain_name": "tf-testHostDomainName12345",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var AlicloudBastionhostHostNetworkDomain0 = map[string]string{
	"instance_id": CHECKSET,
}

func AlicloudBastionhostHostNetworkDomainBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
data "alicloud_bastionhost_instances" "default" {}
`, name)
}
