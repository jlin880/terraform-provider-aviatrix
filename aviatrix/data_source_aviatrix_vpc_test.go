package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
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
	data := terraform.OutputMap(t, terraformOptions, "vpc")
	if len(data) == 0 {
		t.Fatalf("No VPC data returned")
	}

	if data["region"] != awsRegion {
		t.Fatalf("Unexpected region. Expected %s but got %s", awsRegion, data["region"])
	}

	if data["cidr"] != "10.0.0.0/16" {
		t.Fatalf("Unexpected CIDR. Expected 10.0.0.0/16 but got %s", data["cidr"])
	}
}
