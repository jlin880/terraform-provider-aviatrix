Copy code
package aviatrix_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/network/mgmt/network"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAccAviatrixControllerPrivateModeConfig_basic(t *testing.T) {
	rName := randomString(5)

	skipAcc := os.Getenv("SKIP_CONTROLLER_PRIVATE_MODE_CONFIG")
	if skipAcc == "yes" {
		t.Skip("Skipping Controller Private Mode config tests as SKIP_CONTROLLER_PRIVATE_MODE_CONFIG is set")
	}
	msgCommon := ". Set SKIP_CONTROLLER_PRIVATE_MODE_CONFIG to yes to skip Controller Private Mode config tests"
	resourceName := "aviatrix_controller_private_mode_config.test"

	terraformOptions := &terraform.Options{
		TerraformDir: "./",
		Vars: map[string]interface{}{
			"resource_name": rName,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	// Verify the resource exists
	err := testAccControllerPrivateModeConfigExists(t, terraformOptions, resourceName)
	assert.NoError(t, err)

	// Import the resource
	importedResourceName := "imported_controller_private_mode_config"
	importedTerraformOptions := &terraform.Options{
		TerraformDir: "./",
		ImportState:  true,
		// Set the ID to the controller IP
		// Replace . with - in the controller IP since . is not a valid character in a Terraform resource name
		State: fmt.Sprintf("aviatrix_controller_private_mode_config.test,%s\n", strings.Replace(getControllerIP(t), ".", "-", -1)),
		// Use a different resource name to avoid conflicts
		Vars: map[string]interface{}{
			"resource_name": rName + "-imported",
		},
	}

	defer terraform.Destroy(t, importedTerraformOptions)

	terraform.InitAndApply(t, importedTerraformOptions)

	// Verify the imported resource exists
	err = testAccControllerPrivateModeConfigExists(t, importedTerraformOptions, importedResourceName)
	assert.NoError(t, err)
}
func testAccControllerPrivateModeConfigExists(t *testing.T, terraformOptions *terraform.Options) {
    // retrieve the controller IP from the terraform output
    controllerIP := terraform.Output(t, terraformOptions, "controller_ip")

    client := testAccProvider.Meta().(*goaviatrix.Client)

    // check if the controller private mode is enabled
    info, err := client.GetPrivateModeInfo(context.Background())
    assert.NoError(t, err)
    assert.True(t, info.EnablePrivateMode, "controller private mode is not enabled")
    assert.Equal(t, strings.Replace(controllerIP, ".", "-", -1), info.CID, "controller IP not matching with the CID")
}

func testAccControllerPrivateModeConfigDestroy(t *testing.T, terraformOptions *terraform.Options) {
    // retrieve the controller IP from the terraform output
    controllerIP := terraform.Output(t, terraformOptions, "controller_ip")

    client := testAccProvider.Meta().(*goaviatrix.Client)

    // destroy the controller private mode configuration
    terraform.Destroy(t, terraformOptions)

    // check if the controller private mode is disabled
    info, err := client.GetPrivateModeInfo(context.Background())
    assert.NoError(t, err)
    assert.False(t, info.EnablePrivateMode, "controller private mode is still enabled")
    assert.Equal(t, "", info.CID, "controller CID is not empty")
    assert.NotContains(t, info.PrivateIP, controllerIP, "controller IP not removed from private CIDR")
}