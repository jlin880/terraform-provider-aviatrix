package aviatrix_test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAccDataSourceAviatrixVpcTracker_basic(t *testing.T) {
	t.Parallel()

	// Generate a random name to avoid naming conflicts
	uniqueID := random.UniqueId()
	resourceName := "data.aviatrix_vpc_tracker.test"

	// Set up Terraform options
	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../examples/data-sources/aviatrix_vpc_tracker",

		// Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			"prefix": uniqueID,
		},
	}

	// Clean up resources after test is complete
	defer terraform.Destroy(t, terraformOptions)

	// Create resources using Terraform
	terraform.InitAndApply(t, terraformOptions)

	// Check that the data source is accessible and has the expected values
	output := terraform.OutputAll(t, terraformOptions)
	expectedVpcTrackerAttributes := map[string]string{
		"vpc_list.0.account_name": "tfa-" + uniqueID,
		"vpc_list.0.name":         "vpc-for-vpc-tracker-" + uniqueID,
		"vpc_list.0.vpc_id":       "vpc-" + uniqueID,
		"vpc_list.0.cloud_type":   "1",
		"cloud_type":              "1",
	}
	for key, expectedValue := range expectedVpcTrackerAttributes {
		assert.Equal(t, expectedValue, output[key])
	}
}
