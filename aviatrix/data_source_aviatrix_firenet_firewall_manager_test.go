package test

import (
    "fmt"
    "os"
    "testing"

    "github.com/gruntwork-io/terratest/modules/random"
    "github.com/gruntwork-io/terratest/modules/terraform"
    "github.com/stretchr/testify/assert"
)

func TestAccDataSourceAviatrixFireNetFirewallManager_basic(t *testing.T) {
    t.Parallel()

    // Skip the test if the environment variable is set
    if skip, ok := os.LookupEnv("SKIP_DATA_FIRENET_FIREWALL_MANAGER"); ok && skip == "yes" {
        t.Skip("Skipping Data Source FireNet Firewall Manager test as SKIP_DATA_FIRENET_FIREWALL_MANAGER is set")
    }

    // Generate a random name for the resources to avoid naming conflicts
    rName := random.UniqueId()

    // Define Terraform options
    terraformOptions := &terraform.Options{
        TerraformDir: "../",
        Vars: map[string]interface{}{
            "aws_access_key":       os.Getenv("AWS_ACCESS_KEY"),
            "aws_secret_key":       os.Getenv("AWS_SECRET_KEY"),
            "aws_account_number":   os.Getenv("AWS_ACCOUNT_NUMBER"),
            "aws_region":           os.Getenv("AWS_REGION"),
            "firewall_name":        "my-firewall",
            "enable_firewall":      true,
            "vpc_name":             fmt.Sprintf("vpc-for-firenet-%s", rName),
            "transit_gateway_name": fmt.Sprintf("tftg-%s", rName),
            "enable_ha":            false,
            "vendor_type":          "Generic",
        },
        EnvVars: map[string]string{
            "AVIATRIX_API_USER":        os.Getenv("AVIATRIX_API_USER"),
            "AVIATRIX_API_PASSWORD":    os.Getenv("AVIATRIX_API_PASSWORD"),
            "AVIATRIX_CONTROLLER_IP_1": os.Getenv("AVIATRIX_CONTROLLER_IP_1"),
        },
    }

    // Clean up resources after the test is complete
    defer terraform.Destroy(t, terraformOptions)

    // Create resources needed for the test
    terraform.InitAndApply(t, terraformOptions)

    // Test that the data source returns expected results
    gatewayName := terraform.OutputRequired(t, terraformOptions, "gateway_name")
    assert.Equal(t, fmt.Sprintf("tftg-%s", rName), gatewayName)
}
