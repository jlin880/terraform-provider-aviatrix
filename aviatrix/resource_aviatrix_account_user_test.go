package test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/aviatrix/goaviatrix"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAccAviatrixRbacGroupAccessAccountAttachment_basic(t *testing.T) {
	var rbacGroupAccessAccountAttachment goaviatrix.RbacGroupAccessAccountAttachment

	rName := RandomString(5)

	skipAcc := GetEnvVar(t, "SKIP_RBAC_GROUP_ACCESS_ACCOUNT_ATTACHMENT")
	if skipAcc == "yes" {
		t.Skip("Skipping rbac group access account attachment tests as SKIP_RBAC_GROUP_ACCESS_ACCOUNT_ATTACHMENT is set")
	}

	resourceName := "aviatrix_rbac_group_access_account_attachment.test"
	msgCommon := ". Set SKIP_RBAC_GROUP_ACCESS_ACCOUNT_ATTACHMENT to yes to skip rbac group access account attachment tests"

	terraformOptions := &terraform.Options{
		TerraformDir: "../path/to/terraform/dir",
		Vars: map[string]interface{}{
			"group_name":           fmt.Sprintf("tf-%s", rName),
			"access_account_name":  fmt.Sprintf("tf-acc-%s", rName),
			"aws_account_number":   GetEnvVar(t, "AWS_ACCOUNT_NUMBER"),
			"aws_access_key":       GetEnvVar(t, "AWS_ACCESS_KEY"),
			"aws_secret_key":       GetEnvVar(t, "AWS_SECRET_KEY"),
			"cloud_type":           1,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	// Check that the resource was created correctly
	groupName := terraform.Output(t, terraformOptions, "group_name")
	accessAccountName := terraform.Output(t, terraformOptions, "access_account_name")

	assert.Equal(t, fmt.Sprintf("tf-%s", rName), groupName)
	assert.Equal(t, fmt.Sprintf("tf-acc-%s", rName), accessAccountName)

	resourceImportID := fmt.Sprintf("%s/%s", groupName, accessAccountName)
	resourceImport := fmt.Sprintf(`
resource "aviatrix_rbac_group_access_account_attachment" "test" {
	group_name         = aviatrix_rbac_group.test.id
	access_account_name = aviatrix_account.test_account.account_name
}
	`)

	// Import the resource using terraform import
	terraform.Import(t, terraformOptions, resourceImportID)
	resource.Refresh(t, terraformOptions, resourceName)

	// Check that the imported resource matches the expected resource
	assert.NoError(t, testAccCheckRbacGroupAccessAccountAttachmentExists(resourceName, &rbacGroupAccessAccountAttachment))
	assert.Equal(t, groupName, rbacGroupAccessAccountAttachment.GroupName)
	assert.Equal(t, accessAccountName, rbacGroupAccessAccountAttachment.AccessAccountName)
}

func testAccCheckRbacGroupAccessAccountAttachmentExists(resourceName string, rAttachment *goaviatrix.RbacGroupAccessAccountAttachment) error {
	groupName := terraform.OutputRequired(resourceName, "group_name")
	accessAccountName := terraform.OutputRequired(resourceName, "access_account_name")

	client := NewAviatrixClient()

	foundAttachment := &goaviatrix.RbacGroupAccessAccountAttachment{
		GroupName:         groupName,
		AccessAccountName: accessAccountName,
	}

	foundAttachment2, err := client.GetRbacGroupAccessAccountAttachment(foundAttachment)
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
