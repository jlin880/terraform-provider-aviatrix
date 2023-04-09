package aviatrix

import (
    "fmt"
    "os"
    "testing"

    "github.com/gruntwork-io/terratest/modules/random"
    "github.com/gruntwork-io/terratest/modules/terraform"
)

func TestAccDataSourceAviatrixVpcTracker_basic(t *testing.T) {
    t.Parallel()

    // Generate a random name to avoid naming conflicts
    rName := random.UniqueId()
    resourceName := "data.aviatrix_vpc_tracker.test"

    // Check if the test should be skipped
    skipAcc := os.Getenv("SKIP_DATA_VPC_TRACKER")
    if skipAcc == "yes" {
        t.Skip("Skipping data source vpc_tracker tests as 'SKIP_DATA_VPC_TRACKER' is set")
    }

    // Set up Terraform options
    terraformOptions := &terraform.Options{
        // The path to where our Terraform code is located
        TerraformDir: "../examples/data-sources/aviatrix_vpc_tracker",

        // Variables to pass to our Terraform code using -var options
        Vars: map[string]interface{}{
            "prefix": rName,
        },
    }

    // Run `terraform init` and `terraform apply`
    defer terraform.Destroy(t, terraformOptions)
    terraform.InitAndApply(t, terraformOptions)

    // Check that the data source is accessible and has the expected values
    data := terraform.OutputAll(t, terraformOptions, resourceName)
    expectedVpcTrackerAttributes := map[string]string{
        "vpc_list.0.account_name": "tfa-" + rName,
        "vpc_list.0.name":         "vpc-for-vpc-tracker-" + rName,
        "vpc_list.0.vpc_id":       "vpc-" + rName,
        "vpc_list.0.cloud_type":   "1",
        "cloud_type":              "1",
    }
    for key, expectedValue := range expectedVpcTrackerAttributes {
        if data[key].(string) != expectedValue {
            t.Errorf("Unexpected value for %s. Got %v but expected %v", key, data[key].(string), expectedValue)
        }
    }
}
