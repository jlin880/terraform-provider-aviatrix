package aviatrix

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

	resourceName := "data.aviatrix_firenet_firewall_manager.test"
	rName := random.UniqueId()

	skipAcc := os.Getenv("SKIP_DATA_FIRENET_FIREWALL_MANAGER")
	if skipAcc == "yes" {
		t.Skip("Skipping Data Source FireNet Firewall Manager test as SKIP_DATA_FIRENET_FIREWALL_MANAGER is set")
	}
	msg := ". Set SKIP_DATA_FIRENET_FIREWALL_MANAGER to yes to skip Data Source FireNet FIREWALL MANAGER tests"

	terraformOptions := testAccDataSourceAviatrixFireNetFirewallManagerConfigBasic(rName)

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	gatewayName := terraform.OutputRequired(t, terraformOptions, "gateway_name")

	assert.Equal(t, fmt.Sprintf("tftg-%s", rName), gatewayName)

}

func testAccDataSourceAviatrixFireNetFirewallManagerConfigBasic(rName string) *terraform.Options {
	return terraform.WithDefaultRetryableErrors(t, &terraform.Options{
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
	})
}
