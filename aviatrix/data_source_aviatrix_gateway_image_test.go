package aviatrix_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTerraformAviatrixDataSourceGatewayImage(t *testing.T) {
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: ".",
		EnvVars: map[string]string{
			"SKIP_DATA_GATEWAY_IMAGE": os.Getenv("SKIP_DATA_GATEWAY_IMAGE"),
		},
	})

	resourceName := "data.aviatrix_gateway_image.foo"

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	output, err := terraform.OutputE(t, terraformOptions, "image_version")
	require.NoError(t, err, "failed to get output")
	expectedOutput := "hvm-cloudx-aws-022021"

	assert.Equal(t, expectedOutput, output)

	// Check if resource exists
	res, err := terraform.Provider().GetSchema().ResourceTypes["aviatrix_gateway_image"].DataSourcesMap["foo"].StateFunc(terraform.NewResourceConfigRaw(nil), terraformOptions)
	require.NoError(t, err, "failed to get resource state")
	require.NotNil(t, res, "resource does not exist")
}

func TestTerraformAviatrixDataSourceGatewayImageSkip(t *testing.T) {
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: ".",
		EnvVars: map[string]string{
			"SKIP_DATA_GATEWAY_IMAGE": "yes",
		},
	})

	resourceName := "data.aviatrix_gateway_image.foo"

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	// Check if resource exists
	res, err := terraform.Provider().GetSchema().ResourceTypes["aviatrix_gateway_image"].DataSourcesMap["foo"].StateFunc(terraform.NewResourceConfigRaw(nil), terraformOptions)
	require.NoError(t, err, "failed to get resource state")
	assert.Nil(t, res, "resource exists even though 'SKIP_DATA_GATEWAY_IMAGE' is set")
}

func TestAccDataSourceAviatrixGatewayImage_basic(t *testing.T) {
	resourceName := "data.aviatrix_gateway_image.foo"

	skipAcc := os.Getenv("SKIP_DATA_GATEWAY_IMAGE")
	if skipAcc == "yes" {
		t.Skip("Skipping Data Source Gateway Image test as SKIP_DATA_GATEWAY_IMAGE is set")
	}

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: ".",
		Vars: map[string]interface{}{
			"rand": random.UniqueId(),
		},
	})

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	// Check the gateway image data source
	expectedImageVersion := "hvm-cloudx-aws-022021"
	actualImageVersion, err := terraform.OutputE(t, terraformOptions, "image_version")
	require.NoError(t, err, "failed to get output")
	assert.Equal(t, expectedImageVersion, actualImageVersion)

	// Check if resource exists
	res, err := terraform.Provider().GetSchema().ResourceTypes["aviatrix_gateway_image"].DataSourcesMap["foo"].StateFunc(terraform.NewResourceConfigRaw(nil), terraformOptions)
	require.NoError(t, err, "failed to get resource state")
	require.NotNil(t, res, "resource does not exist")
}

func TestAviatrixGatewayImageDataSource(t *testing.T) {
	// Generate a random name to avoid naming conflicts
	rName := random.UniqueId()

	// Set up Terraform options
	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../path/to/terraform/code",

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
	output := terraform.Output(t, terraformOptions, "image_version")
	expectedOutput := "hvm-cloudx-aws-022021"
	assert.Equal(t, expectedOutput, output)

	// Check if the data source resource exists in the Terraform state
	resourceName := "data.aviatrix_gateway_image.foo"
	res, ok := terraform.OutputMap(t, terraformOptions, "gateway_images")[resourceName].(map[string]interface{})
	assert.True(t, ok, fmt.Sprintf("resource %s not found in state", resourceName))
	assert.Equal(t, float64(1), res["cloud_type"])
	assert.Equal(t, "6.5", res["software_version"])
}
