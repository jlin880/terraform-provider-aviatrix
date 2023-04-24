package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAviatrixFQDN_basic(t *testing.T) {
	var fqdn goaviatrix.FQDN

	rName := acctest.RandString(5)

	skipAcc := os.Getenv("SKIP_FQDN")
	if skipAcc == "yes" {
		t.Skip("Skipping FQDN test as SKIP_FQDN is set")
	}

	resourceName := "aviatrix_fqdn.foo"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preGatewayCheck(t, ". Set SKIP_FQDN to yes to skip FQDN tests")
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFQDNDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFQDNConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFQDNExists(resourceName, &fqdn),
					resource.TestCheckResourceAttr(resourceName, "fqdn_tag", fmt.Sprintf("tff-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "fqdn_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "fqdn_mode", "white"),
					resource.TestCheckResourceAttr(resourceName, "gw_filter_tag_list.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "gw_filter_tag_list.0.gw_name", fmt.Sprintf("tfg-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "domain_names.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "domain_names.0.fqdn", "facebook.com"),
					resource.TestCheckResourceAttr(resourceName, "domain_names.0.proto", "tcp"),
					resource.TestCheckResourceAttr(resourceName, "domain_names.0.port", "443"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
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
		proto  = "tcp"
		port   = "443"
		action = "Allow"
	}
}

resource "aviatrix_fqdn_pass_through" "test_fqdn_pass_through" {
	gw_name            = aviatrix_gateway.test_gw_aws.gw_name
	pass_through_cidrs = [
		"10.0.0.0/24",
		"10.0.1.0/24",
		"10.0.2.0/24",
	]

	depends_on         = [aviatrix_fqdn.test_fqdn]
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		os.Getenv("AWS_REGION"))
}

func testAccCheckFQDNPassThroughExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("fqdn_pass_through Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no fqdn_pass_through ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		gw := &goaviatrix.Gateway{GwName: rs.Primary.Attributes["gw_name"]}

		_, err := client.GetFQDNPassThroughCIDRs(gw)
		if err != nil {
			return err
		}
		if gw.GwName != rs.Primary.ID {
			return fmt.Errorf("fqdn_pass_through not found")
		}

		return nil
	}
}

func testAccCheckFQDNPassThroughDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_fqdn_pass_through" {
			continue
		}
		gw := &goaviatrix.Gateway{GwName: rs.Primary.Attributes["gw_name"]}
		_, err := client.GetFQDNPassThroughCIDRs(gw)
		if err == nil {
			return fmt.Errorf("fqdn_pass_through still exists")
		}
	}

	return nil
}
