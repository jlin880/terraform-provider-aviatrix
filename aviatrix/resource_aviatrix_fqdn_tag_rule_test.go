package test

import (
	"fmt"
	"testing"

	"github.com/Azure/go-autorest/autorest/to"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAviatrixFQDNTagRule(t *testing.T) {
	t.Parallel()

	// Define input variables
	resourceName := fmt.Sprintf("test_fqdn_tag_rule_%s", random.UniqueId())
	fqdnTagName := fmt.Sprintf("fqdn-%s", random.UniqueId())
	fqdn := "*.aviatrix.com"
	protocol := "tcp"
	port := "443"
	action := "Allow"

	// Specify the Terraform module folder and variables
	terraformOptions := &terraform.Options{
		TerraformDir: "./",
		Vars: map[string]interface{}{
			"resource_name":    resourceName,
			"fqdn_tag_name":    fqdnTagName,
			"fqdn":             fqdn,
			"protocol":         protocol,
			"port":             port,
			"action":           action,
			"manage_domain_names": to.BoolPtr(false),
			"fqdn_enabled":     to.BoolPtr(true),
			"fqdn_mode":        "white",
		},
	}

	// Call terraform init and apply
	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApply(t, terraformOptions)

	// Verify the FQDN tag rule exists
	err := verifyFQDNTagRule(t, terraformOptions, resourceName, fqdnTagName, fqdn, protocol, port, action)
	assert.NoError(t, err)
}

func verifyFQDNTagRule(t *testing.T, terraformOptions *terraform.Options, resourceName string, fqdnTagName string, fqdn string, protocol string, port string, action string) error {
	// Get the FQDN tag rule ID from Terraform output
	id := terraform.Output(t, terraformOptions, "fqdn_tag_rule_id")

	// Create the Aviatrix client using environment variables
	aviatrixClient, err := createAviatrixClientFromEnvironment()
	if err != nil {
		return err
	}

	// Get the FQDN tag rule
	fqdnTagRule, err := aviatrixClient.GetFQDNTagRule(fqdnTagName, fqdn, protocol, port, action)
	if err != nil {
		return err
	}

	// Verify the FQDN tag rule exists and has the correct ID
	if fqdnTagRule == nil {
		return fmt.Errorf("FQDN tag rule not found: %s", fqdnTagName)
	}
	if fqdnTagRule.ID != id {
		return fmt.Errorf("FQDN tag rule ID does not match: %s != %s", fqdnTagRule.ID, id)
	}

	// Verify the FQDN tag rule has the correct domain name
	if len(fqdnTagRule.DomainList) != 1 {
		return fmt.Errorf("FQDN tag rule does not have exactly 1 domain name: %d", len(fqdnTagRule.DomainList))
	}
	domain := fqdnTagRule.DomainList[0]
	if domain.FQDN != fqdn {
		return fmt.Errorf("FQDN tag rule domain name does not match: %s != %s", domain.FQDN, fqdn)
	}
	if domain.Protocol != protocol {
		return fmt.Errorf("FQDN tag rule

func testAccFQDNDomainNameBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_fqdn" "foo" {
	fqdn_tag            = "fqdn-%s"
	fqdn_enabled        = true
	fqdn_mode           = "white"
	manage_domain_names = false
}

resource "aviatrix_fqdn_tag_rule" "test_fqdn_tag_rule" {
	fqdn_tag_name = aviatrix_fqdn.foo.fqdn_tag
	fqdn          = "*.aviatrix.com"
	protocol      = "tcp"
	port          = "443"
	action        = "Allow"
}
`, rName)
}

func testAccCheckFQDNDomainNameExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("fqdn_tag_rule Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no fqdn_tag_rule ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		fqdn := &goaviatrix.FQDN{
			FQDNTag: rs.Primary.Attributes["fqdn_tag_name"],
			DomainList: []*goaviatrix.Filters{
				{
					FQDN:     rs.Primary.Attributes["fqdn"],
					Protocol: rs.Primary.Attributes["protocol"],
					Port:     rs.Primary.Attributes["port"],
					Verdict:  rs.Primary.Attributes["action"],
				},
			},
		}

		fqdn, err := client.GetFQDNTagRule(fqdn)
		if err != nil {
			return err
		}
		if getFQDNTagRuleID(fqdn) != rs.Primary.ID {
			return fmt.Errorf("fqdn_tag_rule not found")
		}

		return nil
	}
}

func testAccCheckFQDNDomainNameDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_fqdn_tag_rule" {
			continue
		}
		fqdn := &goaviatrix.FQDN{
			FQDNTag: rs.Primary.Attributes["fqdn_tag_name"],
			DomainList: []*goaviatrix.Filters{
				{
					FQDN:     rs.Primary.Attributes["fqdn"],
					Protocol: rs.Primary.Attributes["protocol"],
					Port:     rs.Primary.Attributes["port"],
					Verdict:  rs.Primary.Attributes["action"],
				},
			},
		}
		_, err := client.GetFQDNTagRule(fqdn)
		if err == nil {
			return fmt.Errorf("fqdn_tag_rule still exists")
		}
	}

	return nil
}
