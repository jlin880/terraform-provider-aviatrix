package aviatrix_test

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
 
func TestAccAviatrixFQDNTagRule_basic(t *testing.T) {
    if os.Getenv("SKIP_FQDN_TAG_RULE") == "yes" {
        t.Skip("Skipping fqdn tag rule test as SKIP_FQDN_TAG_RULE is set")
    }
 
    testTagName := fmt.Sprintf("test-%s", strings.ToLower(random.UniqueId()))
    terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
        TerraformDir: "../examples/aviatrix-fqdn-tag-rule",
 
        Vars: map[string]interface{}{
            "fqdn_tag_name": testTagName,
            "fqdn": "*.aviatrix.com",
            "protocol": "tcp",
            "port": "443",
            "action": "Allow",
        },
    })
 
    defer terraform.Destroy(t, terraformOptions)
    terraform.InitAndApply(t, terraformOptions)
 
    // Create a client that will be used to test the resources
    client, err := goaviatrix.NewClientFromEnv()
    if err != nil {
        t.Fatal(err)
    }
 
    fqdn := &goaviatrix.FQDN{
        FQDNTag: testTagName,
        DomainList: []*goaviatrix.Filters{
            {
                FQDN:     "*.aviatrix.com",
                Protocol: "tcp",
                Port:     "443",
                Verdict:  "Allow",
            },
        },
    }
 
    t.Run("Check fqdn tag rule exists", func(t *testing.T) {
        fqdn, err := client.GetFQDNTagRule(fqdn)
        assert.NoError(t, err)
        assert.NotNil(t, fqdn)
    })
 
    t.Run("Check fqdn tag rule is imported correctly", func(t *testing.T) {
        importedResource := terraform.ImportState{
            ID: testTagName,
            Attributes: map[string]string{
                "fqdn_tag_name": testTagName,
                "fqdn": "*.aviatrix.com",
                "protocol": "tcp",
                "port": "443",
                "action": "Allow",
            },
        }
 
        terraformOptionsImport := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
            TerraformDir: "../examples/aviatrix-fqdn-tag-rule",
        })
        terraform.Import(t, terraformOptionsImport, "aviatrix_fqdn_tag_rule.test_fqdn_tag_rule", importedResource)
 
        fqdn, err := client.GetFQDNTagRule(fqdn)
        assert.NoError(t, err)
        assert.NotNil(t, fqdn)
    })
}

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
