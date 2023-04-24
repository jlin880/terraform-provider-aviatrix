package aviatrix_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/network/mgmt/network"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAviatrixControllerEmailExceptionNotificationConfig(t *testing.T) {
	t.Parallel()

	// Skip test if environment variable is set
	if os.Getenv("SKIP_CONTROLLER_EMAIL_EXCEPTION_NOTIFICATION_CONFIG") == "yes" {
		t.Skip("Skipping Controller Email Exception Notification Config test as SKIP_CONTROLLER_EMAIL_EXCEPTION_NOTIFICATION_CONFIG is set")
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
	status, err := client.GetEmailExceptionNotificationStatus(context.Background())
	assert.NoError(t, err)
	assert.False(t, status)

	// Import the resource into Terraform state and verify it
	resourceName := terraform.Options(t, terraformOptions).Vars["resource_name"].(string)
	importedTerraformOptions := terraform.ImportedResourceTerraformOptions(terraformOptions, resourceName, strings.Replace(controllerIP, ".", "-", -1))
	importedTerraformOptions.BackendConfig = terraformOptions.BackendConfig
	terraform.Import(t, importedTerraformOptions)
	terraform.Refresh(t, importedTerraformOptions)
	importedResourceStatus := terraform.Show(t, importedTerraformOptions, "-no-color")
	assert.Contains(t, importedResourceStatus, fmt.Sprintf("enable_email_exception_notification = %q", false))
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
func CheckControllerEmailExceptionNotificationExists(t *testing.T, resourceName string, client *goaviatrix.Client) {
	status, err := client.GetEmailExceptionNotificationStatus(context.Background())
	assert.NoError(t, err)
	assert.False(t, status)

	// Get the resource from Terraform state
	resourceState := terraform.Show(t, terraformOptions, "-no-color")
	assert.Contains(t, resourceState, fmt.Sprintf(`resource "aviatrix_controller_email_exception_notification_config" "%s"`, resourceName))
	assert.Contains(t, resourceState, fmt.Sprintf(`enable_email_exception_notification = %q`, false))
}

func CheckControllerEmailExceptionNotificationConfigDestroy(t *testing.T, terraformOptions *terraform.Options, resourceName string, client *goaviatrix.Client) {
	terraform.Refresh(t, terraformOptions)
	CheckControllerEmailExceptionNotificationExists(t, resourceName, client)
}
