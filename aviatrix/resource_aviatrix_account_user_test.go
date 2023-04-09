package aviatrix

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAccAviatrixAccountUser_basic(t *testing.T) {
	rInt := random.UniqueId()

	terraformOptions := &terraform.Options{
		TerraformDir: "./path/to/terraform/directory",
		Vars: map[string]interface{}{
			"username": fmt.Sprintf("tf-testing-%d", rInt),
			"email":    "abc@xyz.com",
			"password": "Password-1234^",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	var account goaviatrix.AccountUser
	err := json.Unmarshal([]byte(terraform.OutputJson(t, terraformOptions, "")), &account)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, fmt.Sprintf("tf-testing-%d", rInt), account.UserName)
	assert.Equal(t, "abc@xyz.com", account.Email)
	assert.Equal(t, "Password-1234^", account.Password)
}

func TestAccAviatrixAccountUser_import(t *testing.T) {
	rInt := random.UniqueId()

	terraformOptions := &terraform.Options{
		TerraformDir: "./path/to/terraform/directory",
		Vars: map[string]interface{}{
			"username": fmt.Sprintf("tf-testing-%d", rInt),
			"email":    "abc@xyz.com",
			"password": "Password-1234^",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	var account goaviatrix.AccountUser
	err := json.Unmarshal([]byte(terraform.OutputJson(t, terraformOptions, "")), &account)
	if err != nil {
		t.Fatal(err)
	}

	importedTerraformOptions := &terraform.Options{
		TerraformDir: "./path/to/terraform/directory",
		// Set the ID of the resource to import
		ImportState: fmt.Sprintf("aviatrix_account_user.foo %s", account.UserName),
	}

	// Verify the import worked
	terraform.Import(t, importedTerraformOptions)

	resourceState := terraform.Show(t, importedTerraformOptions, "-json")
	assert.NoError(t, json.Unmarshal([]byte(resourceState), &account))

	assert.Equal(t, fmt.Sprintf("tf-testing-%d", rInt), account.UserName)
	assert.Equal(t, "abc@xyz.com", account.Email)
	assert.Equal(t, "Password-1234^", account.Password)
}
func testAccCheckAccountUserExists(n string, account *goaviatrix.AccountUser) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("account not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no account ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		foundAccount := &goaviatrix.AccountUser{
			UserName: rs.Primary.Attributes["username"],
		}

		_, err := client.GetAccountUser(foundAccount)
		if err != nil {
			return fmt.Errorf("failed to get account: %s", err)
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
		if err == nil {
			return fmt.Errorf("account still exists")
		}
	}

	return nil
}