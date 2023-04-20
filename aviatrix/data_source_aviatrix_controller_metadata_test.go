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

	resourceName := "data.aviatrix_controller_metadata.foo"

	skipAcc := os.Getenv("SKIP_DATA_CONTROLLER_METADATA")
	if skipAcc == "yes" {
		t.Skip("Skipping Data Source Controller Metadata test as SKIP_DATA_CONTROLLER_METADATA is set")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: ".",
		EnvVars: map[string]string{
			"SKIP_DATA_CONTROLLER_METADATA": skipAcc,
			"TF_ACC":                        "1",
		},
	}

	// Clean up resources after the test is complete
	defer terraform.Destroy(t, terraformOptions)

	// Create resources needed for the test
	terraform.InitAndApply(t, terraformOptions)

	// Test that the data source returns expected results
	err := terraform.ProviderDiagnosticsValidateExitCode(t, terraformOptions)
	assert.NoError(t, err)

	err = terraform.ApplyAndIdempotent(t, terraformOptions)
	assert.NoError(t, err)

	state := terraform.Show(t, terraformOptions)
	resourceState := state.EvalForResource(resourceName)

	err = testAccDataSourceAviatrixControllerMetadata(resourceState)
	assert.NoError(t, err)
}

func testAccDataSourceAviatrixControllerMetadata(resourceState *terraform.StateResource) error {
	data := resourceState.Primary.Attributes["metadata"]
	assert.NotEmpty(t, data)

	// Additional assertions can be made here based on the metadata returned
	// Example:
	// assert.Contains(t, data, "aviatrix_version")
	// assert.Contains(t, data, "product_name")

	fmt.Println("Controller Metadata:\n", data)

	return nil
}
