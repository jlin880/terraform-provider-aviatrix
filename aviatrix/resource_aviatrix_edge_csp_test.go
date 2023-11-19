package aviatrix_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAccAviatrixEdgeCSP_basic(t *testing.T) {
	if os.Getenv("SKIP_EDGE_CSP") == "yes" {
		t.Skip("Skipping Edge CSP test as SKIP_EDGE_CSP is set")
	}

	// Generate a random resource name to avoid collisions
	resourceName := fmt.Sprintf("aviatrix_edge_csp.test-%s", random.UniqueId())

	// Set up Terraform options
	terraformOptions := &terraform.Options{
		TerraformDir: "../examples/edge_csp",

		// Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			"account_name":      fmt.Sprintf("edge-csp-acc-%s", random.UniqueId()),
			"gw_name":           fmt.Sprintf("edge-csp-%s", random.UniqueId()),
			"site_id":           fmt.Sprintf("site-%s", random.UniqueId()),
			"project_uuid":      os.Getenv("EDGE_CSP_PROJECT_UUID"),
			"compute_node_uuid": os.Getenv("EDGE_CSP_COMPUTE_NODE_UUID"),
			"template_uuid":     os.Getenv("EDGE_CSP_TEMPLATE_UUID"),
			"username":          os.Getenv("EDGE_CSP_USERNAME"),
			"password":          os.Getenv("EDGE_CSP_PASSWORD"),
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEdgeCSPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEdgeCSPBasic(accountName, gwName, siteId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEdgeCSPExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "gw_name", gwName),
					resource.TestCheckResourceAttr(resourceName, "site_id", siteId),
					resource.TestCheckResourceAttr(resourceName, "interfaces.0.ip_address", "10.230.5.32/24"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.1.ip_address", "10.230.3.32/24"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.2.ip_address", "172.16.15.162/20"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	}

	// Import the Edge CSP into Terraform state
	terraform.Import(t, importedTerraformOptions, importedResourceName)
	importedOutput := terraform.Output(t, importedTerraformOptions, "interfaces.0.ip_address")
	assert.Equal(t, "10.230.5.32/24", importedOutput)
}

func testAccEdgeCSPBasic(accountName, gwName, siteId string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_account" {
 	account_name      = "%s"
	cloud_type        = 65536
	edge_csp_username = "%s"
	edge_csp_password = "%s"
}
resource "aviatrix_edge_csp" "test" {
	account_name      = aviatrix_account.test_account.account_name
	gw_name           = "%s"
	site_id           = "%s"
 	project_uuid      = "%s"
 	compute_node_uuid = "%s"
 	template_uuid     = "%s"

	interfaces {
		name          = "eth0"
		type          = "WAN"
		ip_address    = "10.230.5.32/24"
		gateway_ip    = "10.230.5.100"
		wan_public_ip = "64.71.24.221"
	}
	
	interfaces {
		name       = "eth1"
		type       = "LAN"
		ip_address = "10.230.3.32/24"
	}
	
	interfaces {
		name        = "eth2"
		type        = "MANAGEMENT"
		enable_dhcp = false
		ip_address  = "172.16.15.162/20"
		gateway_ip  = "172.16.0.1"
	}
}
 `, accountName, os.Getenv("EDGE_CSP_USERNAME"), os.Getenv("EDGE_CSP_PASSWORD"), gwName, siteId,
		os.Getenv("EDGE_CSP_PROJECT_UUID"), os.Getenv("EDGE_CSP_COMPUTE_NODE_UUID"), os.Getenv("EDGE_CSP_TEMPLATE_UUID"))
}

func testAccCheckEdgeCSPExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("edge csp not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no edge csp id is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		edgeSpoke, err := client.GetEdgeCSP(context.Background(), rs.Primary.Attributes["gw_name"])
		if err != nil {
			return err
		}
		if edgeSpoke.GwName != rs.Primary.ID {
			return fmt.Errorf("could not find edge csp")
		}
		return nil
	}
}

func testAccCheckEdgeCSPDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_edge_csp" {
			continue
		}

		_, err := client.GetEdgeCSP(context.Background(), rs.Primary.Attributes["gw_name"])
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("edge csp still exists")
		}
	}

	return nil
}
