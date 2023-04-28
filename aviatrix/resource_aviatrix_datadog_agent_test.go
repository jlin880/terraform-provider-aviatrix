package aviatrix_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAviatrixDatadogAgent(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "./",
		Vars: map[string]interface{}{
			"datadog_api_key": os.Getenv("DATADOG_API_KEY"),
			"site":            "datadoghq.com",
			"excluded_gateways": []string{"a", "b"},
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	client := goaviatrix.NewClient(os.Getenv("AVIATRIX_CONTROLLER_IP"), os.Getenv("AVIATRIX_USERNAME"), os.Getenv("AVIATRIX_PASSWORD"), os.Getenv("AVIATRIX_ADMIN_EMAIL"))

	datadogAgentStatus, err := client.GetDatadogAgentStatus()
	assert.NoError(t, err)

	assert.NotEmpty(t, datadogAgentStatus.ApiKey)

	expectedExcludedGateways := []string{"a", "b"}
	assert.ElementsMatch(t, expectedExcludedGateways, datadogAgentStatus.ExcludedGateways)

	importedResourceName := fmt.Sprintf("aviatrix_datadog_agent.%s", random.UniqueId())
	importState := terraform.ImportState{
		Resources: []terraform.ImportedResource{
			{
				Type:        "aviatrix_datadog_agent",
				ImportState: fmt.Sprintf("%s %s", datadogAgentStatus.ApiKey, "datadoghq.com"),
				Name:        importedResourceName,
			},
		},
	}

	terraform.Import(t, terraformOptions, &importState)

	importedResource := terraform.Instance(importedResourceName)
	assert.Equal(t, datadogAgentStatus.ApiKey, importedResource.Attr("api_key"))
	assert.Equal(t, "datadoghq.com", importedResource.Attr("site"))
	assert.ElementsMatch(t, expectedExcludedGateways, importedResource.Get("excluded_gateways").([]interface{}))
}


func testAccDatadogAgentBasic() string {
	return fmt.Sprintf(`
resource "aviatrix_datadog_agent" "test_datadog_agent" {
	api_key           = "%s"
	site              = "datadoghq.com"
	excluded_gateways = ["a", "b"]
}
`, os.Getenv("DATADOG_API_KEY"))
}

func testAccCheckDatadogAgentExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("datadog agent not found: %s", resourceName)
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		_, err := client.GetDatadogAgentStatus()
		if err == goaviatrix.ErrNotFound {
			return fmt.Errorf("datadog agent not found")
		}

		return nil
	}
}

func testAccCheckDatadogAgentExcludedGatewaysMatch(input []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*goaviatrix.Client)

		resp, _ := client.GetDatadogAgentStatus()
		if !goaviatrix.Equivalent(resp.ExcludedGateways, input) {
			return fmt.Errorf("excluded gateways don't match with the input")
		}
		return nil
	}
}

func testAccCheckDatadogAgentDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_datadog_agent" {
			continue
		}

		_, err := client.GetDatadogAgentStatus()
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("datadog_agent still exists")
		}
	}

	return nil
}
