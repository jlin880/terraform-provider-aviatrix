package aviatrix

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)
import (
	"context"
	"fmt"
	"testing"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
)

func TestAccAviatrixDNSProfile_basic(t *testing.T) {
	t.Parallel()

	resourceName := "aviatrix_dns_profile.test"
	profileName := fmt.Sprintf("test-dns-profile-%s", random.UniqueId())

	terraformOptions := &terraform.Options{
		TerraformDir: "./",
		Vars: map[string]interface{}{
			"name":                 profileName,
			"global_dns_servers":   []string{"8.8.8.8", "8.8.3.4"},
			"local_domain_names":   []string{"avx.internal.com", "avx.media.com"},
			"lan_dns_servers":      []string{"1.2.3.4", "5.6.7.8"},
			"wan_dns_servers":      []string{"2.3.4.5", "6.7.8.9"},
		},
	}

	// Run Terraform apply and capture the resource ID output
	terraformOutput := terraform.InitAndApplyAndGetOutput(t, terraformOptions, "test_dns_profile")

	// Verify the DNS profile exists
	client := aviatrix.NewClientFromEnv()
	dnsProfile, err := client.GetDNSProfile(context.Background(), profileName)
	assert.NoError(t, err)
	assert.Equal(t, profileName, dnsProfile.Name)
	assert.Equal(t, terraformOutput["global_dns_servers"].([]interface{}), dnsProfile.GlobalDNSServers)
	assert.Equal(t, terraformOutput["local_domain_names"].([]interface{}), dnsProfile.LocalDomainNames)
	assert.Equal(t, terraformOutput["lan_dns_servers"].([]interface{}), dnsProfile.LANDNSServers)
	assert.Equal(t, terraformOutput["wan_dns_servers"].([]interface{}), dnsProfile.WANDNSServers)

	// Import the DNS profile and verify it matches the Terraform state
	importedResourceName := "aviatrix_dns_profile.import"
	importedResource := terraform.Import(t, terraformOptions, importedResourceName)
	assert.Equal(t, resourceName, importedResource)

	terraform.OutputAll(t, terraformOptions)
}
func testAccDNSProfileBasic(profileName string) string {
	return fmt.Sprintf(`
resource "aviatrix_dns_profile" "test" {
	name               = "%s"
	global_dns_servers = ["8.8.8.8", "8.8.3.4"]
	local_domain_names = ["avx.internal.com", "avx.media.com"]
	lan_dns_servers    = ["1.2.3.4", "5.6.7.8"]
	wan_dns_servers    = ["2.3.4.5", "6.7.8.9"]
}
 `, profileName)
}

func testAccCheckDNSProfileExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("dns progile not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no dns profile id is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		_, err := client.GetDNSProfile(context.Background(), rs.Primary.Attributes["name"])
		if err != nil {
			if err == goaviatrix.ErrNotFound {
				return fmt.Errorf("could not find dns profile")
			}
			return err
		}

		return nil
	}
}

func testAccCheckDNSProfileDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_dns_profile" {
			continue
		}

		_, err := client.GetDNSProfile(context.Background(), rs.Primary.Attributes["name"])
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("dns profile still exists")
		}
	}

	return nil
}
