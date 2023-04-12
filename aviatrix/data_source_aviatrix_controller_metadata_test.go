package aviatrix_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAccDataSourceAviatrixControllerMetadata_basic(t *testing.T) {
	t.Parallel()

	// Skips the test if the environment variable is set
	if skip, ok := os.LookupEnv("SKIP_DATA_CONTROLLER_METADATA"); ok && skip == "yes" {
		t.Skip("Skipping Data Source Controller Metadata test as SKIP_DATA_CONTROLLER_METADATA is set")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: ".",
	}

	// Clean up resources after the test is complete
	defer terraform.Destroy(t, terraformOptions)

	// Create resources needed for the test
	terraform.InitAndApply(t, terraformOptions)

	// Test that the data source returns expected results
	data := terraform.Output(t, terraformOptions, "metadata")
	assert.NotEmpty(t, data)

	// Additional assertions can be made here based on the metadata returned
	// Example:
	// assert.Contains(t, data, "aviatrix_version")
	// assert.Contains(t, data, "product_name")

	fmt.Println("Controller Metadata:\n", data)
}

func TestAccDataSourceAviatrixControllerMetadata(t *testing.T) {
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

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	err := testAccDataSourceAviatrixControllerMetadata(resourceName)(terraformOptions.State)
	assert.NoError(t, err)
}

func testAccDataSourceAviatrixControllerMetadata(name string) terraform.ResourceCheck {
	return terraform.ResourceCheck{
		Name: name,
		Exists: true,
		ExpectedOutput: "metadata",
	}
}
