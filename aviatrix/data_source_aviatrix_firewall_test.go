package aviatrix_test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAviatrixFirewallDataSource(t *testing.T) {
	// Set up the Terraform options for the Aviatrix provider
	terraformOptions := &terraform.Options{
		TerraformDir: "<path_to_terraform_module>", // Update with the path to your Terraform module
	}

	// Run Terraform apply and defer destroy
	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApply(t, terraformOptions)

	// Get the output values from Terraform
	gwName := terraform.Output(t, terraformOptions, "gw_name")

	// Set up the input data for the data source
	dataSourceInput := map[string]interface{}{
		"gw_name": gwName,
	}

	// Run the data source function
	data, err := dataSourceAviatrixFirewallRead(dataSourceInput)

	// Assert that no error occurred
	assert.NoError(t, err)

	// Assert that the data source result is as expected
	assert.Equal(t, "deny-all", data.Get("base_policy"))
	assert.False(t, data.Get("base_log_enabled").(bool))
	// ... assert other attributes and policies as needed
}
