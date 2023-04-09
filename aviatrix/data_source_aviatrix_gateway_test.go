package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

func TestAccDataSourceAviatrixGateway_basic(t *testing.T) {
	t.Parallel()

	rName := random.UniqueId()
	resourceName := "data.aviatrix_gateway.foo"

	skipAcc := os.Getenv("SKIP_DATA_GATEWAY")
	if skipAcc == "yes" {
		t.Skip("Skipping Data Source Gateway test as SKIP_DATA_GATEWAY is set")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../examples/data-sources/gateway",
		Vars: map[string]interface{}{
			"test_name":      rName,
			"aws_account":    os.Getenv("AWS_ACCOUNT_NUMBER"),
			"aws_access":     os.Getenv("AWS_ACCESS_KEY"),
			"aws_secret":     os.Getenv("AWS_SECRET_KEY"),
			"aws_vpc_id":     os.Getenv("AWS_VPC_ID"),
			"aws_subnet":     os.Getenv("AWS_SUBNET"),
			"aws_region":     os.Getenv("AWS_REGION"),
			"gw_name":        fmt.Sprintf("tfg-%s", rName),
			"account_name":   fmt.Sprintf("tfa-%s", rName),
			"cloud_type":     "1",
			"gw_size":        "t2.micro",
			"gw_interface":   "0",
			"public_ip":      "AUTO_ALLOCATE",
			"allocate_eip":   "true",
			"disable_srcdst": "false",
			"tags": map[string]string{
				"Automation": "Terraform",
			},
		},
		EnvVars: map[string]string{
			"SKIP_BACKEND": "true",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	gwName := terraform.Output(t, terraformOptions, "gw_name")

	if err := terraform.OutputStruct(resourceName, &GatewayData{}); err != nil {
		t.Fatalf("Failed to decode Terraform output: %v", err)
	}

	expected := GatewayData{
		AccountName: "tfa-" + rName,
		GwName:      "tfg-" + rName,
		VpcID:       os.Getenv("AWS_VPC_ID"),
		VpcReg:      os.Getenv("AWS_REGION"),
		GwSize:      "t2.micro",
	}

	if gwName != expected.GwName {
		t.Errorf("Expected gateway name %s but got %s", expected.GwName, gwName)
	}

	if err := terraform.OutputStruct(resourceName, &expected); err != nil {
		t.Fatalf("Failed to decode Terraform output: %v", err)
	}

	// Verify the output in Terratest format
	expectedOutput := map[string]string{
		"gw_name":       expected.GwName,
		"vpc_id":        expected.VpcID,
		"vpc_reg":       expected.VpcReg,
		"gw_size":       expected.GwSize,
		"account_name":  expected.AccountName,
		"cloud_type":    "1",
		"gw_interface":  "0",
		"public_ip":     "AUTO_ALLOCATE",
		"allocate_eip":  "true",
			"disable_srcdst": "false",
			"tags": map[string]string{
				"Automation": "Terraform",
			},
		},
		EnvVars: map[string]string{
			"SKIP_BACKEND": "true",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	gwName := terraform.Output(t, terraformOptions, "gw_name")

	if err := terraform.OutputStruct(resourceName, &GatewayData{}); err != nil {
		t.Fatalf("Failed to decode Terraform output: %v", err)
	}

	expected := GatewayData{
		AccountName: "tfa-" + rName,
		GwName:      "tfg-" + rName,
		VpcID:       os.Getenv("AWS_VPC_ID"),
		VpcReg:      os.Getenv("AWS_REGION"),
		GwSize:      "t2.micro",
	}

	if gwName != expected.GwName {
		t.Errorf("Expected gateway name %s but got %s", expected.GwName, gwName)
	}

	if err := terraform.OutputStruct(resourceName, &expected); err != nil {
		t.Fatalf("Failed to decode Terraform output: %v", err)
	}
}

type GatewayData struct {
	AccountName string `json:"account_name"`
	GwName      string `json:"gw_name"`
	VpcID       string `json:"vpc_id"`
	VpcReg      string `json:"vpc_reg"`
	GwSize      string `json:"gw_size"`
}

func testAccDataSourceAviatrixGatewayConfigBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test" {
	account_name 	   = "tfa-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = "false"
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}
resource "aviatrix_gateway" "test" {
	cloud_type   = 1
	account_name = aviatrix_account.test.account_name
	gw_name      = "tfg-%s"
	vpc_id       = "%s"
	vpc_reg      = "%s"
	gw_size      = "t2.micro"
	subnet       = "%s"
}
data "aviatrix_gateway" "foo" {
	gw_name = aviatrix_gateway.test.gw_name
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		rName, os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_REGION"), os.Getenv("AWS_SUBNET"))
}

func testAccDataSourceAviatrixGateway(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("root module has no data source called %s", name)
		}

		return nil
	}
}
