package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

func TestAccAviatrixTransitExternalDeviceConn_basic(t *testing.T) {
	terraformOptions := &terraform.Options{
		TerraformDir: "./",
		Vars: map[string]interface{}{
			"aws_account_number": os.Getenv("AWS_ACCOUNT_NUMBER"),
			"aws_access_key":     os.Getenv("AWS_ACCESS_KEY"),
			"aws_secret_key":     os.Getenv("AWS_SECRET_KEY"),
			"aws_region":         os.Getenv("AWS_REGION"),
			"aws_vpc_id":         os.Getenv("AWS_VPC_ID"),
			"aws_subnet":         os.Getenv("AWS_SUBNET"),
		},
	}

	// Skip the test if the SKIP_TRANSIT_EXTERNAL_DEVICE_CONN environment variable is set
	skipAcc := os.Getenv("SKIP_TRANSIT_EXTERNAL_DEVICE_CONN")
	if skipAcc == "yes" {
		t.Skip("Skipping transit external device connection tests as 'SKIP_TRANSIT_EXTERNAL_DEVICE_CONN' is set")
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	// Get the name of the created resource from the Terraform state
	resourceName := terraform.Output(t, terraformOptions, "resource_name")

	// Get the connection details from the Aviatrix API
	client := aviatrix.NewClient(os.Getenv("AVIATRIX_USERNAME"), os.Getenv("AVIATRIX_PASSWORD"), os.Getenv("AVIATRIX_CONTROLLER"))
	externalDeviceConn, err := client.GetExternalDeviceConnDetail(resourceName)
	require.NoError(t, err)

	// Verify that the connection details match the expected values
	assert.Equal(t, os.Getenv("AWS_VPC_ID"), externalDeviceConn.VpcID)
	assert.Equal(t, fmt.Sprintf("tfg-%s", resourceName), externalDeviceConn.GWName)
	assert.Equal(t, "bgp", externalDeviceConn.ConnectionType)
	assert.Equal(t, "123", externalDeviceConn.BGPLocalASNum)
	assert.Equal(t, "345", externalDeviceConn.BGPRemoteASNum)
	assert.Equal(t, "172.12.13.14", externalDeviceConn.RemoteGatewayIP)

	// Import the resource and verify that the connection details match the expected values
	importedResource, err := terraform.ImportE(t, terraformOptions, "aviatrix_transit_external_device_conn", resourceName)
	require.NoError(t, err)

	assert.Equal(t, resourceName, importedResource.Id())
	assert.Equal(t, os.Getenv("AWS_VPC_ID"), importedResource.Get("vpc_id").(string))
	assert.Equal(t, fmt.Sprintf("tfg-%s", resourceName), importedResource.Get("gw_name").(string))
	assert.Equal(t, "bgp", importedResource.Get("connection_type").(string))
	assert.Equal(t, "123", importedResource.Get("bgp_local_as_num").(string))
	assert.Equal(t, "345", importedResource.Get("bgp_remote_as_num").(string))
	assert.Equal(t, "172.12.13.14", importedResource.Get("remote_gateway_ip").(string))
}

func testAccTransitExternalDeviceConnConfigBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test" {
	account_name       = "tfa-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}
resource "aviatrix_transit_gateway" "test" {
	cloud_type   = 1
	account_name = aviatrix_account.test.account_name
	gw_name      = "tfg-%s"
	vpc_id       = "%s"
	vpc_reg      = "%s"
	gw_size      = "t2.micro"
	subnet       = "%s"
}
resource "aviatrix_transit_external_device_conn" "test" {
	vpc_id            = aviatrix_transit_gateway.test.vpc_id
	connection_name   = "%s"
	gw_name           = aviatrix_transit_gateway.test.gw_name
	connection_type   = "bgp"
	bgp_local_as_num  = "123"
	bgp_remote_as_num = "345"
	remote_gateway_ip = "172.12.13.14"
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		rName, os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_REGION"), os.Getenv("AWS_SUBNET"), rName)
}

func checkTransitExternalDeviceConnExists(t *testing.T, resourceName string, externalDeviceConn *goaviatrix.ExternalDeviceConn) error {
	client := getAviatrixClient(t)

	output, err := client.GetTransitExternalDeviceConnList()
	if err != nil {
		return err
	}

	for _, conn := range output {
		if conn.ConnectionName == externalDeviceConn.ConnectionName && conn.VpcID == externalDeviceConn.VpcID {
			*externalDeviceConn = conn
			return nil
		}
	}

	return fmt.Errorf("transit external device connection %s not found", resourceName)
}

func testAccCheckTransitExternalDeviceConnDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_transit_external_device_conn" {
			continue
		}

		foundExternalDeviceConn := &goaviatrix.ExternalDeviceConn{
			VpcID:          rs.Primary.Attributes["vpc_id"],
			ConnectionName: rs.Primary.Attributes["connection_name"],
		}

		_, err := client.GetExternalDeviceConnDetail(foundExternalDeviceConn)
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("site2cloud still exists %s", err.Error())
		}
	}

	return nil
}
