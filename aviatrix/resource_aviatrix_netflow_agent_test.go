package aviatrix_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

func TestAviatrixNetflowAgent_basic(t *testing.T) {
	if os.Getenv("SKIP_NETFLOW_AGENT") == "yes" {
		t.Skip("Skipping netflow agent test as SKIP_NETFLOW_AGENT is set")
	}

	// Generate a random resource name to avoid any naming conflicts
	resourceName := fmt.Sprintf("aviatrix_netflow_agent.test-%s", random.UniqueId())

	terraformOptions := &terraform.Options{
		TerraformDir: "../../path/to/terraform/directory",
		Vars: map[string]interface{}{
			"server_ip":         "1.2.3.4",
			"port":              10,
			"excluded_gateways": []string{"a", "b"},
		},
	}

	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApply(t, terraformOptions)

	// Verify that the netflow agent resource exists
	expectedExcludedGateways := []string{"a", "b"}
	assertNetflowAgentExists(t, resourceName, expectedExcludedGateways)

	// Import the netflow agent resource and verify it
	importedTerraformOptions := terraformOptions
	importedTerraformOptions.ImportState = true
	importedTerraformOptions.ImportStateVerify = true

	terraform.Import(t, importedTerraformOptions)
	assertNetflowAgentExists(t, resourceName, expectedExcludedGateways)
}

func assertNetflowAgentExists(t *testing.T, resourceName string, expectedExcludedGateways []string) {
	client := goaviatrix.NewClient("", os.Getenv("AVIATRIX_CONTROLLER_IP"), os.Getenv("AVIATRIX_USERNAME"), os.Getenv("AVIATRIX_PASSWORD"), true)

	resp, err := client.GetNetflowAgentStatus()
	if err != nil {
		t.Fatalf("Error getting netflow agent status: %s", err)
	}

	if len(resp.ExcludedGateways) != len(expectedExcludedGateways) {
		t.Fatalf("Expected %d excluded gateways, but got %d", len(expectedExcludedGateways), len(resp.ExcludedGateways))
	}

	for i, eg := range expectedExcludedGateways {
		if resp.ExcludedGateways[i] != eg {
			t.Fatalf("Expected excluded gateway %s, but got %s", eg, resp.ExcludedGateways[i])
		}
	}
}
