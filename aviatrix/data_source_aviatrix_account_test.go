package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceAviatrixAccount_basic(t *testing.T) {
	// Skip the test if the SKIP_DATA_ACCOUNT environment variable is set to "yes".
	skipAcc := os.Getenv("SKIP_DATA_ACCOUNT")
	if skipAcc == "yes" {
		t.Skip("Skipping Data Source Account test as SKIP_DATA_ACCOUNT is set")
	}

	// Set the AWS environment variables required for the test.
	awsRegion := os.Getenv("AWS_REGION")
	awsAccessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	awsAccountNumber := os.Getenv("AWS_ACCOUNT_NUMBER")

	// Create a random name to use for the resources.
	rName := random.UniqueId()

	// Define the Terraform options for the test.
	terraformOptions := &terraform.Options{
		TerraformDir: "./",
		EnvVars: map[string]string{
			"AWS_REGION":            awsRegion,
			"AWS_ACCESS_KEY_ID":     awsAccessKey,
			"AWS_SECRET_ACCESS_KEY": awsSecretKey,
			"AWS_ACCOUNT_NUMBER":    awsAccountNumber,
		},
		Vars: map[string]interface{}{
			"account_name":       fmt.Sprintf("tf-testing-%s", rName),
			"aws_account_number": awsAccountNumber,
			"aws_access_key":     awsAccessKey,
			"aws_secret_key":     awsSecretKey,
		},
	}

	// Destroy the Terraform resources at the end of the test.
	defer terraform.Destroy(t, terraformOptions)

	// Apply the Terraform configuration.
	terraform.InitAndApply(t, terraformOptions)

	// Check that the data source returns the expected result.
	expectedAccountName := fmt.Sprintf("tf-testing-%s", rName)
	dataSourceName := fmt.Sprintf("data.aviatrix_account.foo")
	expectedAttributes := map[string]string{
		"account_name": expectedAccountName,
	}
	terraform.OutputMap(t, terraformOptions, dataSourceName, expectedAttributes)
}

func testAccDataSourceAviatrixAccountConfigBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test" {
	account_name       = "tf-testing-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = "false"
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}
data "aviatrix_account" "foo" {
	account_name = aviatrix_account.test.id
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"))
}

func testAccDataSourceAviatrixAccount(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("root module has no resource called %s", name)
		}

		return nil
	}
}
