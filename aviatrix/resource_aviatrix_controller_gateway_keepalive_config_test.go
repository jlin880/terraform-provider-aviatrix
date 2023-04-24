package aviatrix_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAviatrixControllerGatewayKeepaliveConfig(t *testing.T) {
	t.Parallel()

	// Skip test if environment variable is set
	if os.Getenv("SKIP_CONTROLLER_GATEWAY_KEEPALIVE_CONFIG") == "yes" {
		t.Skip("Skipping Controller Gateway Keepalive Config test as SKIP_CONTROLLER_GATEWAY_KEEPALIVE_CONFIG is set")
	}

	terraformOptions := &terraform.Options{
		// Set the path to the Terraform code directory
		TerraformDir: "./terraform",

		// A unique name for the test to avoid naming conflicts
		// with resources that may already exist in the account
		// or with previous runs of the test
		Upgrade: true,
		Vars: map[string]interface{}{
			"resource_name": random.UniqueId(),
		},
	}

	// Run Terraform commands to apply and destroy the resources
	terraform.InitAndApply(t, terraformOptions)
	defer terraform.Destroy(t, terraformOptions)

	// Test that the resource was created successfully
	controllerIP := terraform.Output(t, terraformOptions, "controller_ip")
	client := NewAviatrixClient(t, controllerIP)
	speed, err := client.GetGatewayKeepaliveConfig(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, "slow", speed)

	// Import the resource into Terraform state and verify it
	resourceName := terraform.Options(t, terraformOptions).Vars["resource_name"].(string)
	importedTerraformOptions := terraform.ImportedResourceTerraformOptions(terraformOptions, resourceName, strings.Replace(controllerIP, ".", "-", -1))
	importedTerraformOptions.BackendConfig = terraformOptions.BackendConfig
	terraform.Import(t, importedTerraformOptions)
	terraform.Refresh(t, importedTerraformOptions)
	importedResourceState := terraform.Show(t, importedTerraformOptions, "-no-color")
	assert.Contains(t, importedResourceState, fmt.Sprintf(`resource "aviatrix_controller_gateway_keepalive_config" "%s"`, resourceName))
	assert.Contains(t, importedResourceState, fmt.Sprintf(`keepalive_speed = %q`, "slow"))
}

func NewAviatrixClient(t *testing.T, controllerIP string) *goaviatrix.Client {
	username := os.Getenv("AVIATRIX_USERNAME")
	password := os.Getenv("AVIATRIX_PASSWORD")
	if username == "" || password == "" {
		t.Fatal("AVIATRIX_USERNAME and/or AVIATRIX_PASSWORD environment variables are not set")
	}

	client, err := goaviatrix.NewClient(context.Background(), controllerIP, username, password, "")
	if err != nil {
		t.Fatalf("failed to create Aviatrix client: %v", err)
	}

	return client
}
func CheckControllerGatewayKeepaliveConfigExists(t *testing.T, resourceName string, client *goaviatrix.Client) {
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../terraform",
		Vars: map[string]interface{}{
			"resource_name": resourceName,
		},
	})
	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	// Verify the gateway keepalive config is set to slow
	speed, err := client.GetGatewayKeepaliveConfig(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, "slow", speed)

	// Import the resource into Terraform state and verify it
	importedTerraformOptions := terraform.ImportedResourceTerraformOptions(terraformOptions, resourceName, strings.Replace(client.ControllerIP, ".", "-", -1))
	terraform.Import(t, importedTerraformOptions)
	terraform.Refresh(t, importedTerraformOptions)
	importedResourceState := terraform.Show(t, importedTerraformOptions, "-no-color")
	assert.Contains(t, importedResourceState, fmt.Sprintf(`resource "aviatrix_controller_gateway_keepalive_config" "%s"`, resourceName))
	assert.Contains(t, importedResourceState, fmt.Sprintf(`keepalive_speed = %q`, "slow"))
}

func CheckControllerGatewayKeepaliveConfigDestroy(t *testing.T, resourceName string, client *goaviatrix.Client) {
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../terraform",
		Vars: map[string]interface{}{
			"resource_name": resourceName,
		},
	})
	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	// Verify the gateway keepalive config is set to medium
	speed, err := client.GetGatewayKeepaliveConfig(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, "medium", speed)

	// Refresh the Terraform state and verify that the resource is destroyed
	terraform.Refresh(t, terraformOptions)
	importedResourceState := terraform.Show(t, terraformOptions, "-no-color")
	assert.NotContains(t, importedResourceState, fmt.Sprintf(`resource "aviatrix_controller_gateway_keepalive_config" "%s"`, resourceName))
}
