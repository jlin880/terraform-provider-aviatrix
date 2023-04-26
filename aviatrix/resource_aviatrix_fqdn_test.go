package aviatrix_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/gruntwork-io/terratest/modules/acctest"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAccAviatrixFQDN_basic(t *testing.T) {
	var fqdn goaviatrix.FQDN

	rName := random.UniqueId()

	skipAcc := os.Getenv("SKIP_FQDN")
	if skipAcc == "yes" {
		t.Skip("Skipping FQDN test as SKIP_FQDN is set")
	}

	tfOptions := &terraform.Options{
		TerraformDir: "./",
		Vars: map[string]interface{}{
			"resource_name": "tff-" + rName,
			"aws_account_number": os.Getenv("AWS_ACCOUNT_NUMBER"),
			"aws_access_key": os.Getenv("AWS_ACCESS_KEY"),
			"aws_secret_key": os.Getenv("AWS_SECRET_KEY"),
			"aws_vpc_id": os.Getenv("AWS_VPC_ID"),
			"aws_region": os.Getenv("AWS_REGION"),
			"aws_subnet": os.Getenv("AWS_SUBNET"),
			"gateway_name": "tfg-" + rName,
		},
	}

	defer terraform.Destroy(t, tfOptions)

	terraform.InitAndApply(t, tfOptions)

	resourceName := "aviatrix_fqdn.foo"

	err := terraform.OutputStruct("fqdn", &fqdn)
	assert.NoError(t, err)

	expectedFqdn := fmt.Sprintf("facebook.com;tcp;443")

	assert.Equal(t, fqdn.FQDNTag, "tff-"+rName)
	assert.Equal(t, fqdn.FQDNMode, "white")
	assert.Equal(t, fqdn.FQDNEnabled, true)
	assert.Equal(t, fqdn.GWFilterTagList[0].GWName, "tfg-"+rName)
	assert.Equal(t, fqdn.DomainNames[0].FQDN, expectedFqdn)

	// Import the FQDN resource and check that the imported state matches the current state.
	importedResource := terraform.ImportState(t, tfOptions, resourceName)
	assert.Equal(t, fqdn.FQDNTag, importedResource.Primary.ID)
}


func testAccFQDNConfigBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test" {
	account_name       = "tfa-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}
resource "aviatrix_gateway" "test" {
	cloud_type     = 1
	account_name   = aviatrix_account.test.account_name
	gw_name        = "tfg-%[1]s"
	vpc_id         = "%[5]s"
	vpc_reg        = "%[6]s"
	gw_size        = "t2.micro"
	subnet         = "%[7]s"
	single_ip_snat = true
}
resource "aviatrix_fqdn" "foo" {
	fqdn_tag     = "tff-%[1]s"
	fqdn_enabled = true
	fqdn_mode    = "white"

	gw_filter_tag_list {
		gw_name        = aviatrix_gateway.test.gw_name
		source_ip_list = []
	}

	domain_names {
		fqdn  = "facebook.com"
		proto = "tcp"
		port  = "443"
	}
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_REGION"), os.Getenv("AWS_SUBNET"))
}

func testAccCheckFQDNExists(n string, fqdn *goaviatrix.FQDN) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("FQDN Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no FQDN ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		foundFQDN := &goaviatrix.FQDN{
			FQDNTag: rs.Primary.Attributes["fqdn_tag"],
		}

		_, err := client.GetFQDNTag(foundFQDN)
		if err != nil {
			return err
		}
		if foundFQDN.FQDNTag != rs.Primary.ID {
			return fmt.Errorf("FQDN not found")
		}

		*fqdn = *foundFQDN
		return nil
	}
}

func testAccCheckFQDNDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_fqdn" {
			continue
		}
		foundFQDN := &goaviatrix.FQDN{
			FQDNTag: rs.Primary.Attributes["fqdn_tag"],
		}

		_, err := client.GetFQDNTag(foundFQDN)
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("FQDN still exists")
		}
	}

	return nil
}
