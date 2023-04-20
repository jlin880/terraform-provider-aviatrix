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

func TestAccAviatrixS2C_basic(t *testing.T) {
	var s2c goaviatrix.Site2Cloud

	rName := RandomString(5)
	resourceName := "aviatrix_site2cloud.foo"

	skipAcc := os.Getenv("SKIP_S2C")
	if skipAcc == "yes" {
		t.Skip("Skipping Site2Cloud test as SKIP_S2C is set")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../examples/aviatrix-site2cloud-basic",
		Vars: map[string]interface{}{
			"prefix": rName,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	resource.Assert(t, terraformOptions, terraform.IsResourceUpToDate(resourceName))
	resourceState := terraform.StateFromFile(t, terraformOptions.StateFilePath)

	err := testAccCheckS2CExists(resourceName, &s2c)(resourceState)
	assert.NoError(t, err)

	assert.Equal(t, fmt.Sprintf("tfs-%s", rName), s2c.TunnelName)
	assert.Equal(t, os.Getenv("AWS_VPC_ID"), s2c.VpcID)
	assert.Equal(t, "policy", s2c.TunnelType)
	assert.Equal(t, fmt.Sprintf("tfg-%s", rName), s2c.PrimaryCloudGatewayName)
	assert.Equal(t, "8.8.8.8", s2c.RemoteGatewayIP)
	assert.Equal(t, "10.23.0.0/24", s2c.RemoteSubnetCIDR)
	assert.Equal(t, "generic", s2c.RemoteGatewayType)
	assert.Equal(t, "unmapped", s2c.ConnectionType)
}

func testAccS2CConfigBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test" {
	account_name       = "tfa-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}
resource "aviatrix_gateway" "test" {
	cloud_type   = 1
	account_name = aviatrix_account.test.account_name
	gw_name      = "tfg-%[1]s"
	vpc_id       = "%[5]s"
	vpc_reg      = "%[6]s"
	gw_size      = "t2.micro"
	subnet       = "%[7]s"
}
resource "aviatrix_site2cloud" "foo" {
	vpc_id                     = aviatrix_gateway.test.vpc_id
	connection_name            = "tfs-%[1]s"
	connection_type            = "unmapped"
	remote_gateway_type        = "generic"
	tunnel_type                = "policy"
	primary_cloud_gateway_name = aviatrix_gateway.test.gw_name
	remote_gateway_ip          = "8.8.8.8"
	remote_subnet_cidr         = "10.23.0.0/24"
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_REGION"), os.Getenv("AWS_SUBNET"))
}

func testAccCheckS2CExists(n string, s2c *goaviatrix.Site2Cloud) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("site2cloud Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no site2cloud ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		foundS2C := &goaviatrix.Site2Cloud{
			TunnelName: rs.Primary.Attributes["connection_name"],
			VpcID:      rs.Primary.Attributes["vpc_id"],
		}

		_, err := client.GetSite2Cloud(foundS2C)
		if err != nil {
			return err
		}
		if foundS2C.TunnelName+"~"+foundS2C.VpcID != rs.Primary.ID {
			return fmt.Errorf("site2cloud connection not found")
		}

		*s2c = *foundS2C
		return nil
	}
}

func testAccCheckS2CDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_site2cloud" {
			continue
		}

		foundS2C := &goaviatrix.Site2Cloud{
			TunnelName: rs.Primary.Attributes["connection_name"],
			VpcID:      rs.Primary.Attributes["vpc_id"],
		}

		_, err := client.GetSite2Cloud(foundS2C)
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("site2cloud still exists")
		}
	}

	return nil
}
