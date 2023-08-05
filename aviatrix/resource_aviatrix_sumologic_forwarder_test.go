package aviatrix_test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

func TestAviatrixSumologicForwarder(t *testing.T) {
	t.Parallel()

	terraformOptions := &terraform.Options{
		TerraformDir: "./",
		Vars: map[string]interface{}{
			"access_id": random.UniqueId(),
			"access_key": random.UniqueId(),
			// Set other required variables here
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	// Run assertions or tests to verify the resources created by Terraform
	// For example, you can use the AWS SDK or other tools to check if the Sumologic Forwarder is enabled and configured correctly.
}
