package test

import (
	"fmt"
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

func TestAccDataSourceAviatrixSpokeGateways_basic(t *testing.T) {
	rName := random.UniqueId()
	resourceName := "data.aviatrix_spoke_gateways.foo"

	skipAcc := os.Getenv("SKIP_DATA_SPOKE_GATEWAYS")
	if skipAcc == "yes" {
		t.Skip("Skipping Data Source All Spoke Gateways tests as SKIP_DATA_SPOKE_GATEWAYS is set")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../examples/data-sources/aviatrix_spoke_gateways",
		Upgrade:      true,
		Vars: map[string]interface{}{
			"rname": rName,
			"aws_account_number": os.Getenv("AWS_ACCOUNT_NUMBER"),
			"aws_access_key":     os.Getenv("AWS_ACCESS_KEY"),
			"aws_secret_key":     os.Getenv("AWS_SECRET_KEY"),
			"aws_vpc_id":         os.Getenv("AWS_VPC_ID"),
			"aws_region":         os.Getenv("AWS_REGION"),
			"aws_subnet":         os.Getenv("AWS_SUBNET"),
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	gwName := terraform.Output(t, terraformOptions, "gw_name")

	expected := map[string]interface{}{
		"gateway_list.#":      "1",
		"gateway_list.0.gw_name": fmt.Sprintf("aa-tfg-aws-%s", rName),
		"gateway_list.0.vpc_id":   os.Getenv("AWS_VPC_ID"),
		"gateway_list.0.vpc_reg":  os.Getenv("AWS_REGION"),
		"gateway_list.0.gw_size":  "t2.micro",
	}

	terraform.OutputStruct(t, terraformOptions, "foo", &expected)

	if gwName != expected["gateway_list.0.gw_name"].(string) {
		t.Errorf("Expected gateway name %s but got %s", expected["gateway_list.0.gw_name"], gwName)
	}
}


func testAccDataSourceAviatrixSpokeGatewaysConfigBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_account" {
	account_name 	   = "aa-tfa-%[1]s"
	cloud_type         = 1
	aws_account_number = "%[2]s"
	aws_iam            = "false"
	aws_access_key     = "%[3]s"
	aws_secret_key     = "%[4]s"
}
resource "aviatrix_transit_gateway" "test" {
	cloud_type   = 1
	account_name = aviatrix_account.test_account.account_name
	gw_name      = "aa-tfg-aws-%[1]s"
	vpc_id       = "%[5]s"
	vpc_reg      = "%[6]s"
	gw_size      = "t2.micro"
	subnet       = "%[7]s"
}
data "aviatrix_transit_gateways" "foo" {
    depends_on = [
		aviatrix_transit_gateway.test,
    ]
}
`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_REGION"), os.Getenv("AWS_SUBNET"))
}

func testAccDataSourceAviatrixSpokeGateways(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("root module has no data source called %s", name)
		}

		return nil
	}
}
