package aviatrix

import (
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceAviatrixVPNUser(t *testing.T) {
	resourceName := "aviatrix_vpn_user.test_user"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAviatrixVPNUserDestroy,
		Steps:        testAccCheckAviatrixVPNUserSteps(resourceName),
	})
}

func testAccCheckAviatrixVPNUserDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_vpn_user" {
			continue
		}

		vpnUser := &goaviatrix.VPNUser{
			UserName: rs.Primary.ID,
		}

		err := client.GetVPNUser(vpnUser)
		if err == nil {
			return fmt.Errorf("VPN User still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckAviatrixVPNUserSteps(resourceName string) []*resource.TestStep {
	return []*resource.TestStep{
		{
			Config: testAccAviatrixVPNUserConfig,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr(resourceName, "user_name", "test_user"),
				resource.TestCheckResourceAttr(resourceName, "vpc_id", "test_vpc_id"),
				resource.TestCheckResourceAttr(resourceName, "gw_name", "test_gateway"),
				resource.TestCheckResourceAttr(resourceName, "user_email", "test@test.com"),
				resource.TestCheckResourceAttr(resourceName, "saml_endpoint", "test_saml_endpoint"),
				resource.TestCheckResourceAttr(resourceName, "profiles.#", "1"),
				resource.TestCheckResourceAttr(resourceName, "profiles.0", "test_profile"),
			),
		},
		{
			ResourceName:      resourceName,
			ImportState:       true,
			ImportStateVerify: true,
		},
		{
			Config: testAccAviatrixVPNUserUpdateConfig,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr(resourceName, "user_name", "test_user"),
				resource.TestCheckResourceAttr(resourceName, "vpc_id", "test_vpc_id_updated"),
				resource.TestCheckResourceAttr(resourceName, "gw_name", "test_gateway_updated"),
				resource.TestCheckResourceAttr(resourceName, "user_email", "test@test.com_updated"),
				resource.TestCheckResourceAttr(resourceName, "saml_endpoint", "test_saml_endpoint_updated"),
				resource.TestCheckResourceAttr(resourceName, "profiles.#", "1"),
				resource.TestCheckResourceAttr(resourceName, "profiles.0", "test_profile_updated"),
			),
		},
	}
}

const testAccAviatrixVPNUserConfig = `
resource "aviatrix_vpn_user" "test_user" {
  vpc_id    = "test_vpc_id"
  gw_name   = "test_gateway"
  user_name = "test_user"
  user_email = "test@test.com"
  saml_endpoint = "test_saml_endpoint"

  profiles = [
    "test_profile",
  ]
}`

const testAccAviatrixVPNUserUpdateConfig = `
resource "aviatrix_vpn_user" "test_user" {
  vpc_id    = "test_vpc_id_updated"
  gw_name   = "test_gateway_updated"
  user_name = "test_user"
  user_email = "test@test.com_updated"
  saml_endpoint = "test_saml_endpoint_updated"

  profiles = [
    "test_profile_updated",
  ]
}`
