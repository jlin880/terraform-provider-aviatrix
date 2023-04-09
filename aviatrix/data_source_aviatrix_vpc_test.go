package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAccDataSourceAviatrixVpc_basic(t *testing.T) {
	t.Parallel()

	awsRegion := os.Getenv("AWS_REGION")
	awsAccountID := os.Getenv("AWS_ACCOUNT_ID")

	// Skip test if environment variables are not set
	if awsRegion == "" || awsAccountID == "" {
		t.Skip("Skipping test due to missing AWS_REGION and/or AWS_ACCOUNT_ID environment variables")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "./",
		Vars: map[string]interface{}{
			"region":       awsRegion,
			"account_name": fmt.Sprintf("tfa-%s", random.UniqueId()),
			"name":         fmt.Sprintf("tfv-%s", random.UniqueId()),
			"cidr":         "10.0.0.0/16",
			"aws_account":  awsAccountID,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	// Validate the data source
	resourceName := "data.aviatrix_vpc.vpc"
	expectedRegion := awsRegion
	expectedCIDR := "10.0.0.0/16"

	data := terraform.OutputMap(t, terraformOptions, "vpc")
	if len(data) == 0 {
		t.Fatalf("No VPC data returned")
	}

	assert.Equal(t, expectedRegion, data["region"], "Unexpected region.")
	assert.Equal(t, expectedCIDR, data["cidr"], "Unexpected CIDR.")

	// Check if the resource exists
	res := terraform.GetStateResource(t, terraformOptions, resourceName)
	if res == nil {
		t.Fatalf("resource '%s' does not exist", resourceName)
	}
}
