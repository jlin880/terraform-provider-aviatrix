package aviatrix_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/acctest"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/goaviatrix"
)

func TestAviatrixSpokeExternalDeviceConn_basic(t *testing.T) {
	var externalDeviceConn goaviatrix.ExternalDeviceConn

	rName := random.UniqueId()
	resourceName := "aviatrix_spoke_external_device_conn.test"

	skipAcc := os.Getenv("SKIP_SPOKE_EXTERNAL_DEVICE_CONN")
	if skipAcc == "yes" {
		t.Skip("Skipping spoke external device connection tests as 'SKIP_SPOKE_EXTERNAL_DEVICE_CONN' is set")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: ".",
		Vars: map[string]interface{}{
			"connection_name":   rName,
			"aws_account_name":  fmt.Sprintf("tfa-%s", rName),
			"vpc_id":            os.Getenv("AWS_VPC_ID"),
			"vpc_reg":           os.Getenv("AWS_REGION"),
			"aws_subnet":        os.Getenv("AWS_SUBNET"),
			"bgp_local_as_num":  "123",
			"bgp_remote_as_num": "345",
			"remote_gateway_ip": "172.12.13.14",
			"gw_name":           fmt.Sprintf("tfg-%s", rName),
			"gw_size":           "t2.micro",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	if err := checkSpokeExternalDeviceConnExists(t, resourceName, &externalDeviceConn); err != nil {
		t.Fatal(err)
	}

	expected := map[string]string{
		"vpc_id":            os.Getenv("AWS_VPC_ID"),
		"connection_name":   rName,
		"gw_name":           fmt.Sprintf("tfg-%s", rName),
		"connection_type":   "bgp",
		"bgp_local_as_num":  "123",
		"bgp_remote_as_num": "345",
		"remote_gateway_ip": "172.12.13.14",
	}

	for key, value := range expected {
		if externalDeviceConn.Get(key) != value {
			t.Errorf("Output %s: expected %s, but got %s", key, value, externalDeviceConn.Get(key))
		}
	}
}

func checkSpokeExternalDeviceConnExists(t *testing.T, n string, externalDeviceConn *goaviatrix.ExternalDeviceConn) error {
	client := goaviatrix.NewClient(goaviatrix.ClientConfig{
		APIEndpoint:  os.Getenv("AVIATRIX_API_ENDPOINT"),
		APIUsername:  os.Getenv("AVIATRIX_USERNAME"),
		APIPassword:  os.Getenv("AVIATRIX_PASSWORD"),
		APIToken:     os.Getenv("AVIATRIX_API_TOKEN"),
		APIVersion:   os.Getenv("AVIATRIX_API_VERSION"),
		LogLevel:     os.Getenv("AVIATRIX_LOG_LEVEL"),
		LogFormatter: os.Getenv("AVIATRIX_LOG_FORMATTER"),
	})

	state := terraform.GetState(t, &terraform.Options{TerraformDir: "."})

	rs, ok := state.RootModule().Resources[n]
	if !ok {
		return fmt.Errorf("spoke external device connection Not found: %s", n

func testAccCheckSpokeGatewaySubnetGroupExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("spoke gateway subnet group not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no spoke gateway subnet group ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		spokeGatewaySubnetGroup := &goaviatrix.SpokeGatewaySubnetGroup{
			SubnetGroupName: rs.Primary.Attributes["name"],
			GatewayName:     rs.Primary.Attributes["gw_name"],
		}
		err := client.GetSpokeGatewaySubnetGroup(context.Background(), spokeGatewaySubnetGroup)
		if err != nil {
			return err
		}
		if spokeGatewaySubnetGroup.GatewayName+"~"+spokeGatewaySubnetGroup.SubnetGroupName != rs.Primary.ID {
			return fmt.Errorf("spoke gateway subnet group not found")
		}
		return nil
	}
}

func testAccCheckSpokeGatewaySubnetGroupSubnetsMatch(resourceName string, input []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("spoke gateway subnet group not found: %s", resourceName)
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		spokeGatewaySubnetGroup := &goaviatrix.SpokeGatewaySubnetGroup{
			SubnetGroupName: rs.Primary.Attributes["name"],
			GatewayName:     rs.Primary.Attributes["gw_name"],
		}
		err := client.GetSpokeGatewaySubnetGroup(context.Background(), spokeGatewaySubnetGroup)
		if err != nil {
			return err
		}

		if !goaviatrix.Equivalent(spokeGatewaySubnetGroup.SubnetList, input) {
			return fmt.Errorf("subnets don't match with the input")
		}
		return nil
	}
}

func testAccCheckSpokeGatewaySubnetGroupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_spoke_gateway_subnet_group" {
			continue
		}

		spokeGatewaySubnetGroup := &goaviatrix.SpokeGatewaySubnetGroup{
			SubnetGroupName: rs.Primary.Attributes["name"],
			GatewayName:     rs.Primary.Attributes["gw_name"],
		}

		err := client.GetSpokeGatewaySubnetGroup(context.Background(), spokeGatewaySubnetGroup)
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("spoke gateway subnet group still exists")
		}
	}
	return nil
}
