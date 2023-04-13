package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/assert"
)

func TestAccAviatrixSpokeHaGateway_basic(t *testing.T) {
	var gateway goaviatrix.Gateway

	rName := acctest.RandString(5)
	resourceName := "aviatrix_spoke_ha_gateway.test"

	skipGw := os.Getenv("SKIP_SPOKE_HA_GATEWAY")
	if skipGw == "yes" {
		t.Skip("Skipping Spoke HA Gateway test as SKIP_SPOKE_HA_GATEWAY is set")
	}

	// Setting default values for AWS_GW_SIZE and GCP_GW_SIZE
	awsGwSize := os.Getenv("AWS_GW_SIZE")
	if awsGwSize == "" {
		awsGwSize = "t2.micro"
	}

	if skipGw == "yes" {
		t.Log("Skipping AWS Spoke HA Gateway test as SKIP_SPOKE_HA_GATEWAY_AWS is set")
	} else {
		terraformOptions := &terraform.Options{
			TerraformDir: "../examples/spoke_ha_gateway/aws",
			Vars: map[string]interface{}{
				"account_name":        fmt.Sprintf("tfa-aws-%s", rName),
				"aws_account_number":  os.Getenv("AWS_ACCOUNT_NUMBER"),
				"aws_access_key":      os.Getenv("AWS_ACCESS_KEY"),
				"aws_secret_key":      os.Getenv("AWS_SECRET_KEY"),
				"region":              os.Getenv("AWS_REGION"),
				"aws_gw_size":         awsGwSize,
				"aws_spoke_cidr":      os.Getenv("AWS_SUBNET"),
				"aws_spoke_ha_subnet": os.Getenv("AWS_SPK_HA_SUBNET"),
			},
		}

		defer terraform.Destroy(t, terraformOptions)

		terraform.InitAndApply(t, terraformOptions)

		gwName := fmt.Sprintf("tfg-aws-%s-hagw", rName)

		gateway = goaviatrix.Gateway{
			GwName:      gwName,
			AccountName: fmt.Sprintf("tfa-aws-%s", rName),
		}

		err := aviatrixClient.GetGateway(&gateway)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, gwName, gateway.GwName)
		assert.Equal(t, fmt.Sprintf("tfa-aws-%s", rName), gateway.AccountName)
	}
}

func testAccSpokeHaGatewayConfigAWS(rName string) string {
	awsGwSize := os.Getenv("AWS_GW_SIZE")
	if awsGwSize == "" {
		awsGwSize = "t2.micro"
	}
	return fmt.Sprintf(`
resource "aviatrix_account" "test" {
	account_name       = "tfa-aws-%[1]s"
	cloud_type         = 1
	aws_account_number = "%[2]s"
	aws_iam            = false
	aws_access_key     = "%[3]s"
	aws_secret_key     = "%[4]s"
}
resource "aviatrix_vpc" "test" {
	account_name = aviatrix_account.test.account_name
	cidr         = "%[7]s"
	cloud_type   = 1
	name         = "aws-vpc-%[1]s"
	region       = "%[5]s"
}
resource "aviatrix_spoke_gateway" "test" {
	cloud_type        = 1
	account_name      = aviatrix_vpc.test.account_name
	gw_name           = "tfg-aws-%[1]s"
	gw_size           = "%[6]s"
	vpc_id            = aviatrix_vpc.test.vpc_id
	vpc_reg           = aviatrix_vpc.test.region
	subnet            = aviatrix_vpc.test.public_subnets[0].cidr
	manage_ha_gateway = false
}
resource "aviatrix_spoke_ha_gateway" "test" {
	primary_gw_name = aviatrix_spoke_gateway.test.gw_name
	gw_name         = "tfg-aws-%[1]s-hagw"
	gw_size         = "%[6]s"
	subnet          = aviatrix_vpc.test.public_subnets[1].cidr
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		os.Getenv("AWS_REGION"), awsGwSize, os.Getenv("AWS_SUBNET"))
}

func testAccCheckSpokeHaGatewayExists(n string, gateway *goaviatrix.Gateway) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("spoke gateway Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no spoke gateway ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		foundGateway := &goaviatrix.Gateway{
			GwName:      rs.Primary.Attributes["gw_name"],
			AccountName: rs.Primary.Attributes["account_name"],
		}

		_, err := client.GetGateway(foundGateway)
		if err != nil {
			return err
		}
		if foundGateway.GwName != rs.Primary.ID {
			return fmt.Errorf("spoke ha gateway not found")
		}

		*gateway = *foundGateway
		return nil
	}
}

func testAccCheckSpokeHaGatewayDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_spoke_ha_gateway" {
			continue
		}
		foundGateway := &goaviatrix.Gateway{
			GwName:      rs.Primary.Attributes["gw_name"],
			AccountName: rs.Primary.Attributes["account_name"],
		}

		_, err := client.GetGateway(foundGateway)
		if err == nil {
			return fmt.Errorf("spoke ha gateway still exists")
		}
	}

	return nil
}
