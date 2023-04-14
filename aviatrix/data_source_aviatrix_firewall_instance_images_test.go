package aviatrix

import (
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/stretchr/testify/assert"
)

func TestAviatrixFirewallInstanceImagesDataSource_basic(t *testing.T) {
	t.Parallel()

	// Generate a random name to avoid naming conflicts
	rName := random.UniqueId()

	// Check if the test should be skipped
	skipAcc := os.Getenv("SKIP_DATA_FIREWALL_INSTANCE_IMAGES")
	if skipAcc == "yes" {
		t.Skip("Skipping Data Source Firewall Instance Images tests as SKIP_DATA_FIREWALL_INSTANCE_IMAGES is set")
	}

	// Set up Terraform options
	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../examples/data-sources/aviatrix_firewall_instance_images",

		// Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			"prefix": rName,
		},
	}

	// Clean up resources at the end of the test
	defer terraform.Destroy(t, terraformOptions)

	// Deploy the Terraform code
	terraform.InitAndApply(t, terraformOptions)

	// Check that the data source is accessible and has the expected values
	output := terraform.Output(t, terraformOptions, "firewall_images")
	assert.NotEmpty(t, output)
}

func TestAviatrixFirewallInstanceImagesDataSource_fail(t *testing.T) {
	t.Parallel()

	// Generate a random name to avoid naming conflicts
	rName := acctest.RandString(5)

	// Check if the test should be skipped
	skipAcc := os.Getenv("SKIP_DATA_FIREWALL_INSTANCE_IMAGES")
	if skipAcc == "yes" {
		t.Skip("Skipping Data Source Firewall Instance Images tests as SKIP_DATA_FIREWALL_INSTANCE_IMAGES is set")
	}

	// Set up Terraform options
	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../examples/data-sources/aviatrix_firewall_instance_images",

		// Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			"prefix": rName,
		},
	}

	// Clean up resources at the end of the test
	defer terraform.Destroy(t, terraformOptions)

	// Deploy the Terraform code and expect an error
	err := terraform.InitAndApplyE(t, terraformOptions)
	assert.Error(t, err)

	// Verify that the data source does not exist
	output := terraform.Output(t, terraformOptions, "firewall_images")
	assert.Empty(t, output)
}
