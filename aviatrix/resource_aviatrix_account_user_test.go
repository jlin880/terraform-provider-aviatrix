package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAccAviatrixAccountUser_basic(t *testing.T) {
	var account goaviatrix.AccountUser

	skipAcc := os.Getenv("SKIP_ACCOUNT_USER")
	if skipAcc == "yes" {
		t.Skip("Skipping Account User test as SKIP_ACCOUNT_USER is set")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "./path/to/terraform/directory",
		Vars: map[string]interface{}{
			"username": fmt.Sprintf("tf-testing-%d", random.Random(1000)),
			"email":    "abc@xyz.com",
			"password": "Password-1234^",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	resourceState := terraform.Show(t, terraformOptions, "-json")
	assert.NoError(t, json.Unmarshal([]byte(resourceState), &account))

	assert.Equal(t, fmt.Sprintf("tf-testing-%d", rInt), account.UserName)
	assert.Equal(t, "abc@xyz.com", account.Email)
	assert.Equal(t, "Password-1234^", account.Password)
}


func testAccAccountUserConfigBasic(rInt int) string {
	return fmt.Sprintf(`
resource "aviatrix_account_user" "foo" {
	username = "tf-testing-%d"
	email    = "abc@xyz.com"
	password = "Password-1234^"
}
	`, rInt)
}

func testAccCheckAccountUserExists(n string, account *goaviatrix.AccountUser) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("account Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no Account ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		foundAccount := &goaviatrix.AccountUser{
			UserName: rs.Primary.Attributes["username"],
		}

		_, err := client.GetAccountUser(foundAccount)
		if err == goaviatrix.ErrNotFound {
			return fmt.Errorf("account not found in REST response")
		}
		if foundAccount.UserName != rs.Primary.ID {
			return fmt.Errorf("account not found")
		}

		*account = *foundAccount
		return nil
	}
}

func testAccCheckAccountUserDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_account" {
			continue
		}

		foundAccount := &goaviatrix.AccountUser{
			UserName: rs.Primary.Attributes["username"],
		}

		_, err := client.GetAccountUser(foundAccount)
		if err != nil {
			return fmt.Errorf("account still exists")
		}
	}
	return nil
}
