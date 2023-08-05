package aviatrix_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAviatrixGatewayDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAviatrixGatewayDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.aviatrix_gateway.test_gateway", "gw_name", "test-gateway"),
					resource.TestCheckResourceAttr("data.aviatrix_gateway.test_gateway", "cloud_type", "1"),
					resource.TestCheckResourceAttr("data.aviatrix_gateway.test_gateway", "account_name", "test-account"),
					resource.TestCheckResourceAttr("data.aviatrix_gateway.test_gateway", "vpc_id", "test-vpc"),
					resource.TestCheckResourceAttr("data.aviatrix_gateway.test_gateway", "vpc_reg", "us-west-2"),
					// Add more attribute checks here based on your schema definition
				),
			},
		},
	})
}

const testAccAviatrixGatewayDataSourceConfig = `
data "aviatrix_gateway" "test_gateway" {
  gw_name     = "test-gateway"
  account_name = "test-account"
  vpc_id       = "test-vpc"
  vpc_reg      = "us-west-2"
  // Set other required attributes here
}
`
