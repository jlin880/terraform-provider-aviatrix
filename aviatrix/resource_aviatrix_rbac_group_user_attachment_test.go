package test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAccAviatrixRbacGroupUserAttachment_basic(t *testing.T) {
	var rbacGroupUserAttachment goaviatrix.RbacGroupUserAttachment

	rName := random.UniqueId()

	skipAcc := os.Getenv("SKIP_RBAC_GROUP_USER_ATTACHMENT")
	if skipAcc == "yes" {
		t.Skip("Skipping rbac group user attachment tests as SKIP_RBAC_GROUP_USER_ATTACHMENT is set")
	}

	resourceName := "aviatrix_rbac_group_user_attachment.test"
	msgCommon := ". Set SKIP_RBAC_GROUP_USER_ATTACHMENT to 'yes' to skip rbac group user attachment tests"

	terraformOptions := &terraform.Options{
		TerraformDir: "./",
		Vars: map[string]interface{}{
			"group_name": fmt.Sprintf("tf-%s", rName),
			"user_name":  fmt.Sprintf("tf-user-%s", rName),
			"email":      "abc@xyz.com",
			"password":   "Password-1234",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	checkRbacGroupUserAttachmentExists(t, terraformOptions, resourceName, &rbacGroupUserAttachment)
	checkResourceAttrs(t, terraformOptions, resourceName, rName)

	terraform.Import(t, terraformOptions, resourceName)

	checkRbacGroupUserAttachmentExists(t, terraformOptions, resourceName, &rbacGroupUserAttachment)
	checkResourceAttrs(t, terraformOptions, resourceName, rName)
}

func checkRbacGroupUserAttachmentExists(t *testing.T, terraformOptions *terraform.Options, resourceName string, attachment *goaviatrix.RbacGroupUserAttachment) {
	client := getAviatrixClient(t)

	err := terraform.WithRetryableErrors(t, terraformOptions, func() error {
		groupName := terraformOptions.Vars["group_name"].(string)
		userName := terraformOptions.Vars["user_name"].(string)

		foundAttachment := &goaviatrix.RbacGroupUserAttachment{
			GroupName: groupName,
			UserName:  userName,
		}

		foundAttachment2, err := client.GetRbacGroupUserAttachment(foundAttachment)
		if err != nil {
			return err
		}
		if foundAttachment2.GroupName != groupName {
			return fmt.Errorf("'group_name' not found in created attributes")
		}
		if foundAttachment2.UserName != userName {
			return fmt.Errorf("'user_name' not found in created attributes")
		}

		*attachment = *foundAttachment2
		return nil
	})

	assert.NoError(t, err)
}

func checkResourceAttrs(t *testing.T, terraformOptions *terraform.Options, resourceName string, rName string) {
	expectedResourceAttrs := map[string]string{
		"group_name": fmt.Sprintf("tf-%s", rName),
		"user_name":  fmt.Sprintf("tf-user-%s", rName),
	}

	actualResourceAttrs := terraform.OutputAll(t, terraformOptions)

	for attrName, expectedAttrValue := range expectedResourceAttrs {
		actualAttrValue := actualResourceAttrs[attrName].Value

		assert.Equal(t, expectedAttrValue, actualAttrValue, "Attribute %s does not match", attrName)
	}
}
