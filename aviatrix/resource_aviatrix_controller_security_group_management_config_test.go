package aviatrix

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAviatrixControllerSecurityGroupManagementConfig(t *testing.T) {
	t.Parallel()

	// Skip the test if the environment variable is set
	if os.Getenv("SKIP_CONTROLLER_SECURITY_GROUP_MANAGEMENT_CONFIG") == "yes" {
		t.Skip("Skipping Controller Config test as SKIP_CONTROLLER_SECURITY_GROUP_MANAGEMENT_CONFIG is set")
	}

	// Generate a random name to avoid naming conflicts
	resourceName := fmt.Sprintf("aviatrix_controller_security_group_management_config.test-%s", random.UniqueId())

	terraformOptions := &terraform.Options{
		TerraformDir: "../path/to/terraform/config/dir",
		Vars: map[string]interface{}{
			"account_name":       fmt.Sprintf("tfa-%s", random.UniqueId()),
			"cloud_type":         1,
			"aws_account_number": os.Getenv("AWS_ACCOUNT_NUMBER"),
			"aws_iam":            false,
			"aws_access_key":     os.Getenv("AWS_ACCESS_KEY"),
			"aws_secret_key":     os.Getenv("AWS_SECRET_KEY"),
			"resource_name":      resourceName,
		},
	}

	// Destroy the infrastructure after the test is finished
	defer terraform.Destroy(t, terraformOptions)

	// Create the infrastructure using Terraform
	terraform.InitAndApply(t, terraformOptions)

	// Check that the security group management config exists
	checkControllerSecurityGroupManagementConfigExists(t, resourceName)

	// Import the security group management config and verify it
	terraform.Import(t, terraformOptions, resourceName)
	checkControllerSecurityGroupManagementConfigExists(t, resourceName)
}

func checkControllerSecurityGroupManagementConfigExists(t *testing.T, resourceName string) {
	client := goaviatrix.NewClient(os.Getenv("AVIATRIX_CONTROLLER_IP"), os.Getenv("AVIATRIX_USERNAME"), os.Getenv("AVIATRIX_PASSWORD"))
	status, err := client.GetSecurityGroupManagementStatus()
	assert.NoError(t, err)

	// Verify that the security group management config is enabled
	assert.False(t, status.EnableSecurityGroupManagement)

	// Verify that the security group management config resource exists
	rs, err := terraform.Show(t, terraformOptions, "aviatrix_controller_security_group_management_config.test")
	assert.NoError(t, err)
	assert.Equal(t, resourceName, rs.Primary.ID)

	// Verify that the security group management config resource is associated with the correct controller
	assert.Equal(t, strings.Replace(os.Getenv("AVIATRIX_CONTROLLER_IP"), ".", "-", -1), rs.Primary.ID)
}

func testAccCheckControllerSecurityGroupManagementConfigDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_controller_security_group_management_config" {
			continue
		}

		_, err := client.GetSecurityGroupManagementStatus()
		if err != nil {
			return fmt.Errorf("could not retrieve Controller Security Group Management Status due to err: %v", err)
		}
	}

	return nil
}
