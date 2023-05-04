package test

import (
	"fmt"
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAviatrixControllerMetadataDataSource(t *testing.T) {
	t.Parallel()

	// Skip test if SKIP_DATA_CONTROLLER_METADATA environment variable is set to "yes"
	skipAcc := os.Getenv("SKIP_DATA_CONTROLLER_METADATA")
	if skipAcc == "yes" {
		t.Skip("Skipping Data Source Controller Metadata test as SKIP_DATA_CONTROLLER_METADATA is set")
	}

	// Set Terraform options
	terraformOptions := &terraform.Options{
		TerraformDir: ".",
		EnvVars: map[string]string{
			"SKIP_DATA_CONTROLLER_METADATA": skipAcc,
			"TF_ACC":                        "1",
		},
	}

	// Run Terraform Init and Apply
	terraform.InitAndApply(t, terraformOptions)

	// Test that the data source returns expected results
	err := terraform.ProviderDiagnosticsValidateExitCodeE(t, terraformOptions)
	assert.NoError(t, err)

	// Get the resource state
	state := terraform.ShowE(t, terraformOptions)
	resourceState := state.EvalForResource("data.aviatrix_controller_metadata.foo")

	// Test that the metadata attribute is not empty
	data := resourceState.Primary.Attributes["metadata"]
	assert.NotEmpty(t, data)

	// Additional assertions can be made here based on the metadata returned
	// Example:
	// assert.Contains(t, data, "aviatrix_version")
	// assert.Contains(t, data, "product_name")

	fmt.Println("Controller Metadata:\n", data)
}
