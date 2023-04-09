package aviatrix_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAccAviatrixAwsGuardDuty_basic(t *testing.T) {
	// Skip the test if the environment variable is set
	skipGuardDuty := os.Getenv("SKIP_AWS_GUARD_DUTY")
	if skipGuardDuty == "yes" {
		t.Skip("Skipping AWS GuardDuty test as SKIP_AWS_GUARD_DUTY is set")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "./path/to/terraform/directory",
		Vars: map[string]interface{}{
			"account_name": fmt.Sprintf("tf-testing-%s", random.UniqueId()),
			"region":       "us-west-1",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	resourceState := terraform.Show(t, terraformOptions, "-json")
	accountName := terraformOptions.Vars["account_name"].(string)

	assert.Equal(t, accountName, resourceState.Primary.Attributes["account_name"])
	assert.Equal(t, "us-west-1", resourceState.Primary.Attributes["region"])
}

func testAccAwsGuardDutyBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_acc" {
	account_name       = "tf-testing-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = "false"
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}
resource "aviatrix_aws_guard_duty" "test_aws_guard_duty" {
	account_name = aviatrix_account.test_acc.account_name
	region = "us-west-1"
}
`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"))
}
