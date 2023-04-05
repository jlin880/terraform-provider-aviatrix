package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/gruntwork-io/terratest/modules/acctest"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAccAviatrixVPNUserAccelerator_basic(t *testing.T) {
	t.Parallel()

	testEnvironment := acctest.EnvironmentVariable(t, "SKIP_VPN_USER_ACCELERATOR")
	if testEnvironment == "yes" {
		t.Skip("Skipping VPN User Accelerator test as SKIP_VPN_USER_ACCELERATOR is set")
	}
	msgCommon := ". Set SKIP_VPN_USER_ACCELERATOR to skip VPN User Accelerator tests"

	awsRegion := os.Getenv("AWS_REGION")
	awsVpcId := os.Getenv("AWS_VPC_ID")
	awsSubnet := os.Getenv("AWS_SUBNET")
	awsAccountNumber := os.Getenv("AWS_ACCOUNT_NUMBER")
	awsAccessKey := os.Getenv("AWS_ACCESS_KEY")
	awsSecretKey := os.Getenv("AWS_SECRET_KEY")

	resourceName := "aviatrix_vpn_user_accelerator.test_elb"
	uniqueID := random.UniqueId()

	terraformOptions := &terraform.Options{
		TerraformDir: "../../examples/vpn_user_accelerator",
		Vars: map[string]interface{}{
			"aws_region":         awsRegion,
			"aws_vpc_id":         awsVpcId,
			"aws_subnet":         awsSubnet,
			"aws_account_number": awsAccountNumber,
			"aws_access_key":     awsAccessKey,
			"aws_secret_key":     awsSecretKey,
			"env":                uniqueID,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	err := testAccCheckVPNUserAcceleratorExists(t, resourceName)
	assert.NoError(t, err)
}

func testAccCheckVPNUserAcceleratorExists(t *testing.T, resourceName string) error {
	client := goaviatrix.NewClient(os.Getenv("AVIATRIX_CONTROLLER_IP"), os.Getenv("AVIATRIX_USERNAME"), os.Getenv("AVIATRIX_PASSWORD"))

	rs, err := terraform.ReadResourceFromRootFolder(t, terraformOptions, resourceName)
	if err != nil {
		return err
	}

	if rs.Primary.ID == "" {
		return fmt.Errorf("no vpn user accelerator ID is set")
	}

	elbList, err := client.GetVpnUserAccelerator()
	if err != nil {
		return err
	}
	if !goaviatrix.Contains(elbList, rs.Primary.ID) {
		return fmt.Errorf("vpn user accelerator ID not found")
	}

	return nil
}


func checkVPNUserAcceleratorDestroyed(t *testing.T, client *goaviatrix.Client, state *terraform.State) error {
    for _, rs := range state.RootModule().Resources {
        if rs.Type != "aviatrix_vpn_user_accelerator" {
            continue
        }

        elbList, err := client.GetVpnUserAccelerator()
        if err != nil {
            return fmt.Errorf("error retrieving vpn user accelerator: %s", err)
        }

        if goaviatrix.Contains(elbList, rs.Primary.ID) {
            return fmt.Errorf("vpn user accelerator still exists")
        }
    }

    return nil
}
