package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestTerraformAviatrixVGWConn(t *testing.T) {
	var vgwConn goaviatrix.VGWConn
	vpcID := os.Getenv("AWS_VPC_ID")
	bgpVGWId := os.Getenv("AWS_BGP_VGW_ID")

	rName := acctest.RandString(5)

	terraformOptions := &terraform.Options{
		TerraformDir: "./",
		Vars: map[string]interface{}{
			"account_name":      fmt.Sprintf("tfa-%s", rName),
			"aws_account_number": os.Getenv("AWS_ACCOUNT_NUMBER"),
			"aws_access_key":     os.Getenv("AWS_ACCESS_KEY"),
			"aws_secret_key":     os.Getenv("AWS_SECRET_KEY"),
			"gw_name":            fmt.Sprintf("tfg-%s", rName),
			"vpc_id":             vpcID,
			"vpc_reg":            os.Getenv("AWS_REGION"),
			"gw_size":            "t2.micro",
			"subnet":             os.Getenv("AWS_SUBNET"),
			"conn_name":          fmt.Sprintf("tfc-%s", rName),
			"bgp_vgw_id":         bgpVGWId,
			"bgp_vgw_region":     os.Getenv("AWS_REGION2"),
			"bgp_local_as_num":   "6451",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	resourceName := "aviatrix_vgw_conn.test_vgw_conn"
	importStateVerifyIgnore := []string{"enable_learned_cidrs_approval"}

	// Check if VGW connection resource exists
	assert.True(t, terraform.OutputExists(t, terraformOptions, "aviatrix_vgw_conn.test_vgw_conn_id"))
	assert.True(t, terraform.OutputExists(t, terraformOptions, "aviatrix_vgw_conn.test_vgw_conn_name"))

	// Check if VGW connection resource is created with correct attributes
	client := goaviatrix.NewClient(os.Getenv("AVIATRIX_CONTROLLER_IP"), os.Getenv("AVIATRIX_USERNAME"), os.Getenv("AVIATRIX_PASSWORD"))

	foundVGWConn := &goaviatrix.VGWConn{
		ConnName:      terraformOptions.Vars["conn_name"].(string),
		GwName:        terraformOptions.Vars["gw_name"].(string),
		VPCId:         terraformOptions.Vars["vpc_id"].(string),
		BgpVGWId:      terraformOptions.Vars["bgp_vgw_id"].(string),
		BgpVGWAccount: terraformOptions.Vars["account_name"].(string),
		BgpVGWRegion:  terraformOptions.Vars["bgp_vgw_region"].(string),
		BgpLocalAsNum: terraformOptions.Vars["bgp_local_as_num"].(string),
	}

	foundVGWConn2, err := client.GetVGWConnDetail(foundVGWConn)
	require.NoError(t, err)
	assert.Equal(t, foundVGWConn.ConnName, foundVGWConn2.ConnName)
	assert.Equal(t, foundVGWConn.GwName, foundVGWConn2.GwName)
	assert.Equal(t, foundVGWConn.V


		foundVGWConn2, err := client.GetVGWConnDetail(foundVGWConn)
		if err != nil {
			return err
		}
		if foundVGWConn2.ConnName != rs.Primary.Attributes["conn_name"] {
			return fmt.Errorf("conn_name Not found in created attributes")
		}

		*vgwConn = *foundVGWConn
		return nil
	}
}

func testAccCheckVGWConnDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_vgw_conn" {
			continue
		}

		foundVGWConn := &goaviatrix.VGWConn{
			ConnName:      rs.Primary.Attributes["conn_name"],
			GwName:        rs.Primary.Attributes["gw_name"],
			VPCId:         rs.Primary.Attributes["vpc_id"],
			BgpVGWId:      rs.Primary.Attributes["bgp_vgw_id"],
			BgpVGWAccount: rs.Primary.Attributes["bgp_vgw_account"],
			BgpVGWRegion:  rs.Primary.Attributes["bgp_vgw_region"],
			BgpLocalAsNum: rs.Primary.Attributes["bgp_local_as_num"],
		}

		_, err := client.GetVGWConnDetail(foundVGWConn)
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("vgw connection still exists")
		}
	}

	return nil
}
