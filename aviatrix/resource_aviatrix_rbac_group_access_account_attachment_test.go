package aviatrix_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/goaviatrix"
	"github.com/gruntwork-io/terratest/modules/acctest"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAccAviatrixRbacGroupAccessAccountAttachment_basic(t *testing.T) {
	t.Parallel()

	rName := random.UniqueId()
	groupName := fmt.Sprintf("tf-%s", rName)
	accountName := fmt.Sprintf("tf-acc-%s", rName)

	skipAcc := os.Getenv("SKIP_RBAC_GROUP_ACCESS_ACCOUNT_ATTACHMENT")
	if skipAcc == "yes" {
		t.Skip("Skipping rbac group access account attachment tests as SKIP_RBAC_GROUP_ACCESS_ACCOUNT_ATTACHMENT is set")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "./path/to/terraform/dir",
		Vars: map[string]interface{}{
			"group_name":         groupName,
			"access_account_name": accountName,
			"aws_account_number": os.Getenv("AWS_ACCOUNT_NUMBER"),
			"aws_access_key":     os.Getenv("AWS_ACCESS_KEY"),
			"aws_secret_key":     os.Getenv("AWS_SECRET_KEY"),
			"cloud_type":         1,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	// Check that the resource was created correctly
	client := goaviatrix.NewClient(goaviatrix.ClientConfig{
		Username: os.Getenv("AVIATRIX_USERNAME"),
		Password: os.Getenv("AVIATRIX_PASSWORD"),
		APIURL:   os.Getenv("AVIATRIX_API_URL"),
	})

	rbacGroupAccessAccountAttachment := goaviatrix.RbacGroupAccessAccountAttachment{
		GroupName:         groupName,
		AccessAccountName: accountName,
	}

	foundAttachment, err := client.GetRbacGroupAccessAccountAttachment(&rbacGroupAccessAccountAttachment)
	assert.NoError(t, err)
	assert.Equal(t, groupName, foundAttachment.GroupName)
	assert.Equal(t, accountName, foundAttachment.AccessAccountName)

	// Test import
	importedResourceName := "aviatrix_rbac_group_access_account_attachment.test"
	resourceType := "aviatrix_rbac_group_access_account_attachment"
	resourceID := fmt.Sprintf("%s:%s", groupName, accountName)

	err = terraform.Import(t, terraformOptions, importedResourceName, resourceID)
	assert.NoError(t, err)

	err = terraform.Refresh(t, terraformOptions)
	assert.NoError(t, err)

	// Check that the imported resource was created correctly
	foundAttachment, err = client.GetRbacGroupAccessAccountAttachment(&rbacGroupAccessAccountAttachment)
	assert.NoError(t, err)
	assert.Equal(t, groupName, foundAttachment.GroupName)
	assert.Equal(t, accountName, foundAttachment.AccessAccountName)

	// Check that the resource can be destroyed
	terraform.Destroy(t, terraformOptions)

	// Check that the resource was destroyed successfully
	_, err = client.GetRbacGroupAccessAccountAttachment(&rbacGroupAccessAccountAttachment)
	assert.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "is invalid"), "Expected error to contain 'is invalid'")
}

func testAccRbacGroupAccessAccountAttachmentConfigBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_rbac_group" "test" {
	group_name = "tf-%s"
}
resource "aviatrix_account" "test_account" {
	account_name       = "tf-acc-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}
resource "aviatrix_rbac_group_access_account_attachment" "test" {
	group_name 			= aviatrix_rbac_group.test.group_name
	access_account_name = aviatrix_account.test_account.account_name
}
	`, rName, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"))
}

func testAccCheckRbacGroupAccessAccountAttachmentExists(n string, rAttachment *goaviatrix.RbacGroupAccessAccountAttachment) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("RbacGroupAccessAccountAttachment Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no RbacGroupAccessAccountAttachment ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		foundAttachment := &goaviatrix.RbacGroupAccessAccountAttachment{
			GroupName:         rs.Primary.Attributes["group_name"],
			AccessAccountName: rs.Primary.Attributes["access_account_name"],
		}

		foundAttachment2, err := client.GetRbacGroupAccessAccountAttachment(foundAttachment)
		if err != nil {
			return err
		}
		if foundAttachment2.GroupName != rs.Primary.Attributes["group_name"] {
			return fmt.Errorf("'group_name' Not found in created attributes")
		}
		if foundAttachment2.AccessAccountName != rs.Primary.Attributes["access_account_name"] {
			return fmt.Errorf("'access_account_name' Not found in created attributes")
		}

		*rAttachment = *foundAttachment2
		return nil
	}
}

func testAccCheckRbacGroupAccessAccountAttachmentDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_rbac_group_access_account_attachment" {
			continue
		}
		foundAttachment := &goaviatrix.RbacGroupAccessAccountAttachment{
			GroupName:         rs.Primary.Attributes["group_name"],
			AccessAccountName: rs.Primary.Attributes["access_account_name"],
		}

		_, err := client.GetRbacGroupAccessAccountAttachment(foundAttachment)
		if err != goaviatrix.ErrNotFound {
			if strings.Contains(err.Error(), "is invalid") {
				return nil
			}
			return fmt.Errorf("rbac group access account attachment still exists")
		}
	}

	return nil
}
