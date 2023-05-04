package aviatrix

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)
import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/goaviatrix"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAccAviatrixRbacGroupPermissionAttachment_basic(t *testing.T) {
	var rbacGroupPermissionAttachment goaviatrix.RbacGroupPermissionAttachment

	rName := random.UniqueId()

	skipAcc := os.Getenv("SKIP_RBAC_GROUP_PERMISSION_ATTACHMENT")
	if skipAcc == "yes" {
		t.Skip("Skipping rbac group permission attachment tests as 'SKIP_RBAC_GROUP_PERMISSION_ATTACHMENT' is set")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "./",
		Vars: map[string]interface{}{
			"group_name":      fmt.Sprintf("tf-%s", rName),
			"permission_name": "all_write",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	resourceName := "aviatrix_rbac_group_permission_attachment.test"

	err := terraform.Import(t, terraformOptions, resourceName)
	assert.NoError(t, err)

	rbacGroupPermissionAttachment.GroupName = terraform.Output(t, terraformOptions, "group_name")
	rbacGroupPermissionAttachment.PermissionName = terraform.Output(t, terraformOptions, "permission_name")

	assert.Equal(t, fmt.Sprintf("tf-%s", rName), rbacGroupPermissionAttachment.GroupName)
	assert.Equal(t, "all_write", rbacGroupPermissionAttachment.PermissionName)

	err = testAccCheckRbacGroupPermissionAttachmentExists(t, &rbacGroupPermissionAttachment)
	assert.NoError(t, err)
}

func testAccCheckRbacGroupPermissionAttachmentExists(t *testing.T, rAttachment *goaviatrix.RbacGroupPermissionAttachment) error {
	client := goaviatrix.NewClientFromEnv()

	foundAttachment, err := client.GetRbacGroupPermissionAttachment(rAttachment)
	if err != nil {
		return err
	}

	assert.Equal(t, rAttachment.GroupName, foundAttachment.GroupName)
	assert.Equal(t, rAttachment.PermissionName, foundAttachment.PermissionName)

	*rAttachment = *foundAttachment

	return nil
}
func testAccCheckRbacGroupPermissionAttachmentExists(n string, rAttachment *goaviatrix.RbacGroupPermissionAttachment) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("RbacGroupPermissionAttachment Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no RbacGroupPermissionAttachment ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		foundAttachment := &goaviatrix.RbacGroupPermissionAttachment{
			GroupName:      rs.Primary.Attributes["group_name"],
			PermissionName: rs.Primary.Attributes["permission_name"],
		}

		foundAttachment2, err := client.GetRbacGroupPermissionAttachment(foundAttachment)
		if err != nil {
			return err
		}
		if foundAttachment2.GroupName != rs.Primary.Attributes["group_name"] {
			return fmt.Errorf("'group_name' Not found in created attributes")
		}
		if foundAttachment2.PermissionName != rs.Primary.Attributes["permission_name"] {
			return fmt.Errorf("'permission_name' Not found in created attributes")
		}

		*rAttachment = *foundAttachment2
		return nil
	}
}

func testAccCheckRbacGroupPermissionAttachmentDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_rbac_group_permission_attachment" {
			continue
		}
		foundAttachment := &goaviatrix.RbacGroupPermissionAttachment{
			GroupName:      rs.Primary.Attributes["group_name"],
			PermissionName: rs.Primary.Attributes["permission_name"],
		}

		_, err := client.GetRbacGroupPermissionAttachment(foundAttachment)
		if err != goaviatrix.ErrNotFound {
			if strings.Contains(err.Error(), "is invalid") {
				return nil
			}
			return fmt.Errorf("rbac group user attachment still exists")
		}
	}

	return nil
}
