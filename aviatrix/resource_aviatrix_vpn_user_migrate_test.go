package aviatrix

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

func TestAccAviatrixVPNUserResourceV0(t *testing.T) {
	t.Parallel()

	terraformOptions := &terraform.Options{
		TerraformDir: "./",
		Vars: map[string]interface{}{
			"manage_user_attachment": false,
		},
	}

	// Clean up resources after the test
	defer terraform.Destroy(t, terraformOptions)

	// Deploy the resources with Terraform
	terraform.InitAndApply(t, terraformOptions)

	// Validate the VPN user resource
	vpnUserResource := terraform.Output(t, terraformOptions, "vpn_user_resource")
	if vpnUserResource != "true" {
		t.Errorf("Failed to create VPN user resource")
	}
}

func TestAccAviatrixVPNUserStateUpgradeV0(t *testing.T) {
	t.Parallel()

	terraformOptions := &terraform.Options{
		TerraformDir: "./",
		Vars: map[string]interface{}{
			"manage_user_attachment": false,
		},
	}

	// Clean up resources after the test
	defer terraform.Destroy(t, terraformOptions)

	// Deploy the resources with Terraform
	terraform.InitAndApply(t, terraformOptions)

	// Upgrade the state version to V0
	terraformOptionsV0 := terraformOptions
	terraformOptionsV0.Vars["manage_user_attachment"] = "false" // Simulating old state without the field

	// Upgrade the state to V0
	terraform.InitAndApply(t, terraformOptionsV0)

	// Validate the state upgrade
	vpnUserResourceV0 := terraform.Output(t, terraformOptionsV0, "vpn_user_resource")
	if vpnUserResourceV0 != "true" {
		t.Errorf("Failed to upgrade state to V0")
	}
}
