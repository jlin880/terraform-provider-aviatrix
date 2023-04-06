package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
func TestAccDataSourceAviatrixControllerMetadata_basic(t *testing.T) {
	resourceName := "data.aviatrix_controller_metadata.foo"

	skipAcc := os.Getenv("SKIP_DATA_CONTROLLER_METADATA")
	if skipAcc == "yes" {
		t.Skip("Skipping Data Source Controller Metadata test as SKIP_DATA_CONTROLLER_METADATA is set")
	}

	testAccProvider := aviatrixProvider()
	testAccProviders := map[string]terraform.ResourceProvider{
		"aviatrix": testAccProvider,
	}

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: ".",
		EnvVars: map[string]string{
			"SKIP_DATA_CONTROLLER_METADATA": skipAcc,
		},
		Providers: testAccProviders,
	})

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	err := testAccDataSourceAviatrixControllerMetadata(resourceName)(terraformOptions.State)
	if err != nil {
		t.Fatalf("failed to verify data source: %s", err)
	}
}

func testAccDataSourceAviatrixControllerMetadata(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("root module has no data source called %s", name)
		}

		return nil
	}
}
