package aviatrix
import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/acctest"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/assert"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/goaviatrix"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

func TestAccAviatrixEdgeSpoke_basic(t *testing.T) {
	if os.Getenv("SKIP_EDGE_SPOKE") == "yes" {
		t.Skip("Skipping Edge as a Spoke test as SKIP_EDGE_SPOKE is set")
	}

	edgeName := fmt.Sprintf("edge-%s", random.UniqueId())
	siteId := fmt.Sprintf("site-%s", random.UniqueId())
	path, _ := os.Getwd()

	terraformOptions := &terraform.Options{
		TerraformDir: "./",
		Vars: map[string]interface{}{
			"gw_name":                edgeName,
			"site_id":                siteId,
			"ztp_file_type":          "iso",
			"ztp_file_download_path": path,
			"interfaces": []interface{}{
				map[string]interface{}{
					"name":          "eth0",
					"type":          "WAN",
					"ip_address":    "10.230.5.32/24",
					"gateway_ip":    "10.230.5.100",
					"wan_public_ip": "64.71.24.221",
				},
				map[string]interface{}{
					"name":       "eth1",
					"type":       "LAN",
					"ip_address": "10.230.3.32/24",
				},
				map[string]interface{}{
					"name":        "eth2",
					"type":        "MANAGEMENT",
					"enable_dhcp": false,
					"ip_address":  "172.16.15.162/20",
					"gateway_ip":  "172.16.0.1",
				},
			},
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	// Check that the edge as a spoke resource exists
	edgeSpoke, err := getEdgeSpoke(t, terraformOptions, edgeName)
	assert.NoError(t, err)
	assert.NotNil(t, edgeSpoke)
	assert.Equal(t, edgeName, edgeSpoke.GwName)

	// Import the edge as a spoke resource
	importedTerraformOptions := &terraform.Options{
		TerraformDir: "./",
		ImportState:  fmt.Sprintf("aviatrix_edge_spoke.test %s", edgeName),
	}

	terraform.Import(t, importedTerraformOptions)

	// Verify that the imported resource exists
	importedEdgeSpoke, err := getEdgeSpoke(t, terraformOptions, edgeName)
	assert.NoError(t, err)
	assert.NotNil(t, importedEdgeSpoke)
	assert.Equal(t, edgeName, importedEdgeSpoke.GwName)
}

func testAccCheckEdgeSpokeExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("edge as a spoke not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no edge as a spoke id is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		edgeSpoke, err := client.GetEdgeSpoke(context.Background(), rs.Primary.Attributes["gw_name"])
		if err != nil {
			return err
		}
		if edgeSpoke.GwName != rs.Primary.ID {
			return fmt.Errorf("could not find edge as a spoke")
		}
		return nil
	}
}

func testAccCheckEdgeSpokeDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_edge_spoke" {
			continue
		}

		_, err := client.GetEdgeSpoke(context.Background(), rs.Primary.Attributes["gw_name"])
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("edge as a spoke still exists")
		}
	}

	return nil
}
