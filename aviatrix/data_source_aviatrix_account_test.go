package aviatrix

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/aviatrix-systems/terraform-provider-aviatrix/aviatrix"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

func TestAccAviatrixDataSourceAccount_basic(t *testing.T) {
	testAccTerratestEnv := os.Getenv("TESTACC_TERRATEST_ENV")
	if testAccTerratestEnv == "" {
		t.Fatal("TESTACC_TERRATEST_ENV must be set for acceptance tests")
	}

	aviatrixAccountName := fmt.Sprintf("tf-testing-%s", random.UniqueId())

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "./",
		Vars: map[string]interface{}{
			"account_name":        aviatrixAccountName,
			"cloud_type":          1,
			"aws_account_number":  os.Getenv("AWS_ACCOUNT_NUMBER"),
			"aws_iam":             "false",
			"aws_access_key":      os.Getenv("AWS_ACCESS_KEY"),
			"aws_secret_key":      os.Getenv("AWS_SECRET_KEY"),
		},
	})

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	dataSourceName := "data.aviatrix_account.foo"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		Providers:         testAccProviders,
		CheckDestroy:      testAccCheckDataSourceAviatrixAccountDestroy,
		ExpectNonEmptyPlan: true,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAviatrixAccountConfigBasic(aviatrixAccountName),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceAviatrixAccount(dataSourceName),
				),
			},
		},
	})
}

func testAccCheckDataSourceAviatrixAccountDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_account" && rs.Type != "data.aviatrix_account" {
			continue
		}

		if _, err := testAccProviders["aviatrix"].Meta().(*aviatrix.Client).GetAccount(rs.Primary.ID); err == nil {
			return fmt.Errorf("account %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccDataSourceAviatrixAccountConfigBasic(accountName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test" {
	account_name       = "%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = "false"
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}

data "aviatrix_account" "foo" {
	account_name = aviatrix_account.test.id
}
`, accountName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"))
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
