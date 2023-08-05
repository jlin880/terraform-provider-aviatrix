package aviatrix

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

func TestAccAviatrixVPNProfileStateUpgradeV0(t *testing.T) {
	t.Parallel()

	terraformOptions := &terraform.Options{
		TerraformDir: "./",
		Vars: map[string]interface{}{
			"manage_user_attachment": random.Bool(),
		},
	}

	// Clean up resources after the test
	defer terraform.Destroy(t, terraformOptions)

	// Deploy the resources with Terraform
	terraform.InitAndApply(t, terraformOptions)

	// Validate the VPN profile state upgrade
	vpnProfile := terraform.Output(t, terraformOptions, "vpn_profile")
	if vpnProfile != "true" {
		t.Errorf("Failed to perform VPN profile state upgrade")
	}
}
