package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceAviatrixFireNetVendorIntegration_basic(t *testing.T) {
	t.Parallel()

	rName := random.UniqueId()

	skipAcc := os.Getenv("SKIP_DATA_FIRENET_VENDOR_INTEGRATION")
	if skipAcc == "yes" {
		t.Skip("Skipping Data Source FireNet Vendor Integration test as SKIP_DATA_FIRENET_VENDOR_INTEGRATION is set")
	}

	terraformOptions, err := configureTerraformOptions(rName)
	if err != nil {
		t.Fatalf("Failed to configure Terraform options: %v", err)
	}

	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApply(t, terraformOptions)

	resourceName := "data.aviatrix_firenet_vendor_integration.test"
	check := testAccDataSourceAviatrixFireNetVendorIntegration(t, resourceName)

	if err := check(terraformOptions); err != nil {
		t.Fatalf("Failed test: %v", err)
	}
}

func configureTerraformOptions(rName string) (*terraform.Options, error) {
	awsRegion := os.Getenv("AWS_REGION")
	awsAccountNumber := os.Getenv("AWS_ACCOUNT_NUMBER")
	awsAccessKey := os.Getenv("AWS_ACCESS_KEY")
	awsSecretKey := os.Getenv("AWS_SECRET_KEY")

	accountName := fmt.Sprintf("tfa-%s", rName)
	gwName := fmt.Sprintf("tftg-%s", rName)
	firewallName := fmt.Sprintf("tffw-%s", rName)

	terraformOptions := &terraform.Options{
		TerraformDir: "../examples/firenet_vendor_integration",
		Vars: map[string]interface{}{
			"aws_region":         awsRegion,
			"aws_account_number": awsAccountNumber,
			"aws_access_key":     awsAccessKey,
			"aws_secret_key":     awsSecretKey,
			"account_name":       accountName,
			"gw_name":            gwName,
			"firewall_name":      firewallName,
		},
	}

	return terraformOptions, nil
}

func testAccDataSourceAviatrixFireNetVendorIntegration(t *testing.T, resourceName string) resource.TestCheckFunc {
	return func(terraformOptions *terraform.Options) error {
		resourceState := terraform.GetResourceState(t, terraformOptions, resourceName)

		if _, ok := resourceState.Attributes["firewall_name"]; !ok {
			return fmt.Errorf("Expected firewall name to be set but it was not")
		}

		return nil
	}
}
