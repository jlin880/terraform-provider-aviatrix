
import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAccAviatrixCloudwatchAgent_basic(t *testing.T) {
	if os.Getenv("SKIP_CLOUDWATCH_AGENT") == "yes" {
		t.Skip("Skipping cloudwatch agent test as SKIP_CLOUDWATCH_AGENT is set")
	}

	terraformOptions := prepareCloudwatchAgentTest(t)
	resourceName := "aviatrix_cloudwatch_agent.test_cloudwatch_agent"

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	checkCloudwatchAgentResource(t, terraformOptions, resourceName)
	importedResourceName := importCloudwatchAgentResource(t, terraformOptions, resourceName)

	assert.Equal(t, resourceName, importedResourceName)
}

func prepareCloudwatchAgentTest(t *testing.T) *terraform.Options {
	terraformDir := "../../examples/cloudwatch_agent"

	region := "us-east-1"
	excludedGateways := []string{"a", "b"}

	terraformOptions := &terraform.Options{
		TerraformDir: terraformDir,
		Vars: map[string]interface{}{
			"cloudwatch_role_arn": "arn:aws:iam::469550033836:role/aviatrix-role-cloudwatch",
			"region":              region,
			"excluded_gateways":   excludedGateways,
		},
	}

	return terraformOptions
}

func checkCloudwatchAgentResource(t *testing.T, terraformOptions *terraform.Options, resourceName string) {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	resp, err := client.GetCloudwatchAgentStatus()
	assert.NoError(t, err)
	assert.True(t, resp.ExcludedGateways != nil)
	assert.Equal(t, terraformOptions.Vars["region"], resp.Region)
	assert.True(t, goaviatrix.Equivalent(resp.ExcludedGateways, terraformOptions.Vars["excluded_gateways"].([]string)))
}

func importCloudwatchAgentResource(t *testing.T, terraformOptions *terraform.Options, resourceName string) string {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	resp, err := client.GetCloudwatchAgentStatus()
	assert.NoError(t, err)

	importedTerraformOptions := terraformOptions
	importedTerraformOptions.ImportState = fmt.Sprintf(`{"excluded_gateways":["%v"],"cloudwatch_role_arn":"%v","region":"%v"}`,
		resp.ExcludedGateways, terraformOptions.Vars["cloudwatch_role_arn"], resp.Region)
	importedResourceName := terraform.Import(t, importedTerraformOptions, resourceName)

	return importedResourceName
}

func TestAccAviatrixCloudwatchAgent_import(t *testing.T) {
	if os.Getenv("SKIP_CLOUDWATCH_AGENT") == "yes" {
		t.Skip("Skipping cloudwatch agent test as SKIP_CLOUDWATCH_AGENT is set")
	}

	terraformOptions := prepareCloudwatchAgentTest(t)
	resourceName := "aviatrix_cloudwatch_agent.test_cloudwatch_agent"

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	importedResourceName := importCloudwatchAgentResource(t, terraformOptions, resourceName)

	assert.Equal(t, resourceName, importedResourceName)
}


func testAccCloudwatchAgentBasic() string {
	return `
resource "aviatrix_cloudwatch_agent" "test_cloudwatch_agent" {
	cloudwatch_role_arn = "arn:aws:iam::469550033836:role/aviatrix-role-cloudwatch"
	region              = "us-east-1"
	excluded_gateways   = ["a", "b"]
}
`
}

func testAccCheckCloudwatchAgentExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("cloudwatch agent not found: %s", resourceName)
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		_, err := client.GetCloudwatchAgentStatus()
		if err == goaviatrix.ErrNotFound {
			return fmt.Errorf("cloudwatch agent not found")
		}

		return nil
	}
}

func testAccCheckCloudwatchAgentExcludedGatewaysMatch(input []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*goaviatrix.Client)

		resp, _ := client.GetCloudwatchAgentStatus()
		if !goaviatrix.Equivalent(resp.ExcludedGateways, input) {
			return fmt.Errorf("excluded gateways don't match with the input")
		}
		return nil
	}
}

func testAccCheckCloudwatchAgentDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_cloudwatch_agent" {
			continue
		}

		_, err := client.GetCloudwatchAgentStatus()
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("cloudwatch_agent still exists")
		}
	}

	return nil
}
