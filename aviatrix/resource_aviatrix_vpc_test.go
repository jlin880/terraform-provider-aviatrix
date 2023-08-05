package aviatrix

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

func TestAccAviatrixVpc(t *testing.T) {
	t.Parallel()

	terraformOptions := &terraform.Options{
		TerraformDir: "./",
		Vars: map[string]interface{}{
			"cloud_type":                1,
			"account_name":              "my-account",
			"name":                      "my-vpc",
			"region":                    "us-west-2",
			"cidr":                      "10.0.0.0/16",
			"subnet_size":               24,
			"num_of_subnet_pairs":       2,
			"enable_private_oob_subnet": false,
			"aviatrix_transit_vpc":      false,
			"aviatrix_firenet_vpc":      false,
			"enable_native_gwlb":        false,
			"private_mode_subnets":      false,
		},
	}

	// Clean up resources after the test
	defer terraform.Destroy(t, terraformOptions)

	// Deploy the resources with Terraform
	terraform.InitAndApply(t, terraformOptions)

	// Validate the VPC creation
	vpcID := terraform.Output(t, terraformOptions, "vpc_id")
	if vpcID == "" {
		t.Errorf("Failed to create VPC")
	}
}
