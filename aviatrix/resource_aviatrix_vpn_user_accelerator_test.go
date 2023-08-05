package aviatrix

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

func TestAccAviatrixVPNUserAccelerator(t *testing.T) {
	t.Parallel()

	terraformOptions := &terraform.Options{
		TerraformDir: "./",
		Vars: map[string]interface{}{
			"elb_name": random.UniqueId(),
		},
	}

	// Clean up resources after the test
	defer terraform.Destroy(t, terraformOptions)

	// Deploy the resources with Terraform
	terraform.InitAndApply(t, terraformOptions)

	// Validate the VPN user accelerator resource
	vpnUserAccelerator := terraform.Output(t, terraformOptions, "vpn_user_accelerator")
	if vpnUserAccelerator != "true" {
		t.Errorf("Failed to create VPN user accelerator")
	}
}
