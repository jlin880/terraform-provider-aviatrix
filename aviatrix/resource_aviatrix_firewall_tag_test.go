package test

import (
	"fmt"
	"testing"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"

	goaviatrix "github.com/AviatrixSystems/terraform-provider-aviatrix/goaviatrix"
)

func TestAccAviatrixFirewallTag_basic(t *testing.T) {
	t.Parallel()

	terraformOptions := &terraform.Options{
		TerraformDir: "../examples/aviatrix_firewall_tag",
		Vars: map[string]interface{}{
			"firewall_tag": fmt.Sprintf("tft-%s", random.UniqueId()),
			"cidr_list": []map[string]string{
				{
					"cidr_tag_name": "a1",
					"cidr":          "10.1.0.0/24",
				},
				{
					"cidr_tag_name": "b1",
					"cidr":          "10.2.0.0/24",
				},
			},
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	var firewallTag goaviatrix.FirewallTag
	assert.NoError(t, goaviatrix.GetResourceByName(&firewallTag, terraformOptions.Var("firewall_tag").(string)))

	assert.Equal(t, terraformOptions.Var("firewall_tag").(string), firewallTag.Name)
	assert.Equal(t, "10.1.0.0/24", firewallTag.CIDRList[0].CIDR)
	assert.Equal(t, "a1", firewallTag.CIDRList[0].CIDRTagName)
	assert.Equal(t, "10.2.0.0/24", firewallTag.CIDRList[1].CIDR)
	assert.Equal(t, "b1", firewallTag.CIDRList[1].CIDRTagName)
}


func testAccFirewallTagConfigBasic(rInt int) string {
	return fmt.Sprintf(`
resource "aviatrix_firewall_tag" "foo" {
	firewall_tag = "tft-%d"
	cidr_list {
		cidr_tag_name = "a1"
		cidr          = "10.1.0.0/24"
	}
	cidr_list {
		cidr_tag_name = "b1"
		cidr          = "10.2.0.0/24"
	}
}
	`, rInt)
}

func testAccCheckFirewallTagExists(n string, firewallTag *goaviatrix.FirewallTag) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("firewall tag Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no tag ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		foundTag := &goaviatrix.FirewallTag{
			Name: rs.Primary.Attributes["firewall_tag"],
		}

		_, err := client.GetFirewallTag(foundTag)
		if err != nil {
			return err
		}
		if foundTag.Name != rs.Primary.ID {
			return fmt.Errorf("firewall tag not found")
		}

		*firewallTag = *foundTag
		return nil
	}
}

func testAccCheckFirewallTagDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_firewall_tag" {
			continue
		}
		foundTag := &goaviatrix.FirewallTag{
			Name: rs.Primary.Attributes["firewall_tag"],
		}

		_, err := client.GetFirewallTag(foundTag)
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("firewall tag still exists after destroy")
		}
	}

	return nil
}
