package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAccDataSourceAviatrixFireNetVendorIntegration(t *testing.T) {
	rName := acctest.RandString(5)
	resourceName := "data.aviatrix_firenet_vendor_integration.test"

	skipAcc := os.Getenv("SKIP_DATA_FIRENET_VENDOR_INTEGRATION")
	if skipAcc == "yes" {
		t.Skip("Skipping Data Source FireNet Vendor Integration test as SKIP_DATA_FIRENET_VENDOR_INTEGRATION is set")
	}
	msg := ". Set SKIP_DATA_FIRENET_VENDOR_INTEGRATION to yes to skip Data Source FireNet Vendor Integration tests"

	terraformOptions, err := configureTerraformOptions(rName)
	require.NoError(t, err)

	// Destroy the resources after the test.
	defer terraform.Destroy(t, terraformOptions)

	// Deploy the resources with Terraform.
	terraform.InitAndApply(t, terraformOptions)

	// Check the resource state using Terraform output.
	err = testAccDataSourceAviatrixFireNetVendorIntegration(t, resourceName)(terraformOptions)
	require.NoError(t, err)
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

func TestAccDataSourceAviatrixFireNetVendorIntegration(t *testing.T) {
	rName := acctest.RandString(5)
	resourceName := "data.aviatrix_firenet_vendor_integration.test"

	skipAcc := os.Getenv("SKIP_DATA_FIRENET_VENDOR_INTEGRATION")
	if skipAcc == "yes" {
		t.Skip("Skipping Data Source FireNet Vendor Integration test as SKIP_DATA_FIRENET_VENDOR_INTEGRATION is set")
	}
	msg := ". Set SKIP_DATA_FIRENET_VENDOR_INTEGRATION to yes to skip Data Source FireNet Vendor Integration tests"

	terraformOptions, err := configureTerraformOptions(rName)
	require.NoError(t, err)

	// Destroy the resources after the test.
	defer terraform.Destroy(t, terraformOptions)

	// Deploy the resources with Terraform.
	terraform.InitAndApply(t, terraformOptions)

	// Check the resource state using Terraform output.
	err = testAccDataSourceAviatrixFireNetVendorIntegration(t, resourceName)(terraformOptions)
	require.NoError(t, err)
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


func testAccDataSourceAviatrixFireNetVendorIntegration(t *testing.T, resourceName string) func(*terraform.Options) error {
	return func(terraformOptions *terraform.Options) error {
		resourceState := terraform.ShowResource(t, terraformOptions, resourceName)

		if _, ok := resourceState.Attributes["firewall_name"]; !ok {
			return fmt.Errorf("Expected firewall name to be set but it was not")
		}

		return nil
	}
}
