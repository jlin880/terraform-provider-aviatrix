package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

func TestAccDataSourceAviatrixFireNet_basic(t *testing.T) {
	t.Parallel()

	// Check if the test should be skipped
	skipAcc := os.Getenv("SKIP_DATA_FIRENET")
	if skipAcc == "yes" {
		t.Skip("Skipping Data Source FireNet tests as SKIP_DATA_FIRENET is set")
	}

	// Generate a random name to avoid naming conflicts
	rName := random.UniqueId()

	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../examples/data-sources/aviatrix_firenet",

		// Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			"prefix": rName,
		},
	}

	// Run `terraform init` and `terraform apply`
	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApply(t, terraformOptions)

	// Check that the data source is accessible and has the expected values
	data := terraform.OutputAll(t, terraformOptions, "aviatrix_firenet_foo")
	expectedFireNetAttributes := map[string]string{
		"vpc_id":                                          fmt.Sprintf("vpc-for-firenet-%s", rName),
		"inspection_enabled":                              "true",
		"egress_enabled":                                  "false",
		"firewall_instance_association.#":                 "1",
		"firewall_instance_association.0.attached":        "true",
		"firewall_instance_association.0.firenet_gw_name": fmt.Sprintf("tftg-%s", rName),
		"firewall_instance_association.0.firewall_name":   fmt.Sprintf("tffw-%s", rName),
	}
	for key, expectedValue := range expectedFireNetAttributes {
		if data[key].(string) != expectedValue {
			t.Errorf("Unexpected value for %s. Got %v but expected %v", key, data[key].(string), expectedValue)
		}
	}

	// Test importing the data source using the state file
	importedTerraformOptions := terraformOptions
	importedTerraformOptions.StatePath = terraformOptions.StatePath + ".backup"
	terraform.Import(t, importedTerraformOptions, "aviatrix_firenet_foo."+rName)
	terraform.Refresh(t, terraformOptions)
	importedData := terraform.OutputAll(t, terraformOptions, "aviatrix_firenet_foo")
	for key, expectedValue := range expectedFireNetAttributes {
		if importedData[key].(string) != expectedValue {
			t.Errorf("Unexpected value for %s after import. Got %v but expected %v", key, importedData[key].(string), expectedValue)
		}
	}
}

func testAccDataSourceFireNetConfigBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_account" {
	account_name       = "tfa-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}

resource "aviatrix_vpc" "test_vpc" {
	cloud_type           = 1
	account_name         = aviatrix_account.test_account.account_name
	region               = "%s"
	name                 = "vpc-for-firenet"
	cidr                 = "10.10.0.0/24"
	aviatrix_firenet_vpc = true
}

resource "aviatrix_transit_gateway" "test_transit_gateway" {
	cloud_type               = aviatrix_vpc.test_vpc.cloud_type
	account_name             = aviatrix_account.test_account.account_name
	gw_name                  = "tftg-%s"
	vpc_id                   = aviatrix_vpc.test_vpc.vpc_id
	vpc_reg                  = aviatrix_vpc.test_vpc.region
	gw_size                  = "c5.xlarge"
	subnet                   = aviatrix_vpc.test_vpc.subnets[0].cidr
	enable_hybrid_connection = true
	enable_firenet           = true
}

resource "aviatrix_firewall_instance" "test_firewall_instance" {
	vpc_id            = aviatrix_vpc.test_vpc.vpc_id
	firenet_gw_name   = aviatrix_transit_gateway.test_transit_gateway.gw_name
	firewall_name     = "tffw-%s"
	firewall_image    = "Palo Alto Networks VM-Series Next-Generation Firewall Bundle 1"
	firewall_size     = "m5.xlarge"
	management_subnet = aviatrix_vpc.test_vpc.subnets[0].cidr
	egress_subnet     = aviatrix_vpc.test_vpc.subnets[1].cidr
}

resource "aviatrix_firewall_instance_association" "firewall_instance_association" {
	vpc_id               = aviatrix_firewall_instance.test_firewall_instance.vpc_id
	firenet_gw_name      = aviatrix_transit_gateway.test_transit_gateway.gw_name
	instance_id          = aviatrix_firewall_instance.test_firewall_instance.instance_id
	firewall_name        = aviatrix_firewall_instance.test_firewall_instance.firewall_name
	lan_interface        = aviatrix_firewall_instance.test_firewall_instance.lan_interface
	management_interface = aviatrix_firewall_instance.test_firewall_instance.management_interface
	egress_interface     = aviatrix_firewall_instance.test_firewall_instance.egress_interface
	attached             = true
}

resource "aviatrix_firenet" "test_firenet" {
	vpc_id             = aviatrix_vpc.test_vpc.vpc_id
	inspection_enabled = true
	egress_enabled     = false

	depends_on = [aviatrix_firewall_instance_association.firewall_instance_association]
}

data "aviatrix_firenet" "foo" {
	vpc_id = aviatrix_firenet.test_firenet.vpc_id
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		os.Getenv("AWS_REGION"), rName, rName)
}

func TestAccDataSourceAviatrixFireNet(t *testing.T) {
	t.Parallel()

	// Check if the test should be skipped
	skipAcc := os.Getenv("SKIP_DATA_FIRENET")
	if skipAcc == "yes" {
		t.Skip("Skipping Data Source FireNet tests as SKIP_DATA_FIRENET is set")
	}

	// Generate a random name to avoid naming conflicts
	rName := random.UniqueId()

	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../examples/data-sources/aviatrix_firenet",

		// Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			"prefix": rName,
		},
	}

	// Clean up resources after the test is complete
	defer terraform.Destroy(t, terraformOptions)

	// Create resources needed for the test
	terraform.InitAndApply(t, terraformOptions)

	// Test that the data source returns expected results
	dataSourceName := "data.aviatrix_firenet.foo"
	expectedFireNetAttributes := map[string]string{
		"vpc_id":                                          fmt.Sprintf("vpc-for-firenet-%s", rName),
		"inspection_enabled":                              "true",
		"egress_enabled":                                  "false",
		"firewall_instance_association.#":                 "1",
		"firewall_instance_association.0.attached":        "true",
		"firewall_instance_association.0.firenet_gw_name": fmt.Sprintf("tftg-%s", rName),
		"firewall_instance_association.0.firewall_name":   fmt.Sprintf("tffw-%s", rName),
	}
	actualFireNetAttributes := terraform.OutputAll(t, terraformOptions, dataSourceName)
	for key, expectedValue := range expectedFireNetAttributes {
		if actualFireNetAttributes[key] != expectedValue {
			t.Errorf("Unexpected value for %s. Got %v but expected %v", key, actualFireNetAttributes[key], expectedValue)
		}
	}
}