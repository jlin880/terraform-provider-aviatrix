package aviatrix_test

import (
	"fmt"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/assert"
)

func TestAccAviatrixVGWConn(t *testing.T) {
	connName := "test-connection"
	gwName := "test-gateway"
	vpcID := "test-vpc"
	bgpVGWID := "test-bgp-vgw-id"
	bgpVGWAccount := "test-bgp-vgw-account"
	bgpVGWRegion := "test-bgp-vgw-region"
	bgpLocalAsNum := "65000"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckAviatrixVGWConnDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "aviatrix_vgw_conn" "test" {
						conn_name = "%s"
						gw_name = "%s"
						vpc_id = "%s"
						bgp_vgw_id = "%s"
						bgp_vgw_account = "%s"
						bgp_vgw_region = "%s"
						bgp_local_as_num = "%s"
					}
				`, connName, gwName, vpcID, bgpVGWID, bgpVGWAccount, bgpVGWRegion, bgpLocalAsNum),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aviatrix_vgw_conn.test", "conn_name", connName),
					resource.TestCheckResourceAttr("aviatrix_vgw_conn.test", "gw_name", gwName),
					resource.TestCheckResourceAttr("aviatrix_vgw_conn.test", "vpc_id", vpcID),
					resource.TestCheckResourceAttr("aviatrix_vgw_conn.test", "bgp_vgw_id", bgpVGWID),
					resource.TestCheckResourceAttr("aviatrix_vgw_conn.test", "bgp_vgw_account", bgpVGWAccount),
					resource.TestCheckResourceAttr("aviatrix_vgw_conn.test", "bgp_vgw_region", bgpVGWRegion),
					resource.TestCheckResourceAttr("aviatrix_vgw_conn.test", "bgp_local_as_num", bgpLocalAsNum),
				),
			},
		},
	})
}

func testAccCheckAviatrixVGWConnDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_vgw_conn" {
			continue
		}

		connName := rs.Primary.Attributes["conn_name"]
		vpcID := rs.Primary.Attributes["vpc_id"]

		vgwConn := &goaviatrix.VGWConn{
			ConnName: connName,
			VPCId:    vpcID,
		}

		_, err := client.GetVGWConnDetail(vgwConn)
		if err == nil {
			return fmt.Errorf("aviatrix_vgw_conn %s still exists", connName)
		}
	}

	return nil
}
