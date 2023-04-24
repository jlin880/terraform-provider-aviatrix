package test

import (
	"fmt"
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"

	"github.com/AviatrixSystems/terratest-aviatrix-nonet/api"
)

func TestAwsTgwVpcAttachment(t *testing.T) {
	t.Parallel()

	// Declare variables for the test
	rName := fmt.Sprintf("test-aws-tgw-vpc-attach-%s", strings.ToLower(random.UniqueId()))
	awsSideAsNumber := "64512"
	nDm := "test"

	// Deploy the Terraform code
	terraformOptions := &terraform.Options{
		TerraformDir: "../",
		Vars: map[string]interface{}{
			"aws_region":     os.Getenv("AWS_REGION"),
			"aws_vpc_id":     os.Getenv("AWS_VPC_ID"),
			"tgw_vpc_attach": rName,
		},
	}
	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApply(t, terraformOptions)

	// Check if the AWS TGW VPC Attachment exists and save its details to a variable
	var awsTgwVpcAttachment api.AwsTgwVpcAttachment
	resourceName := fmt.Sprintf("aviatrix_aws_tgw_vpc_attachment.test-%s", rName)
	assert.NoError(t, testAccCheckAwsTgwVpcAttachmentExists(resourceName, &awsTgwVpcAttachment))

	// Check if the Security Domain was correctly set
	assert.Equal(t, nDm, awsTgwVpcAttachment.SecurityDomainName)

	// Check if the AWS account number was correctly set
	accountName := fmt.Sprintf("tfa-%s", rName)
	assert.Equal(t, accountName, awsTgwVpcAttachment.VpcAccountName)

	// Check if the VPC ID was correctly set
	assert.Equal(t, os.Getenv("AWS_VPC_ID"), awsTgwVpcAttachment.VpcID)

	// Check if the TGW VPC Attachment was correctly attached
	assert.Equal(t, "attached", awsTgwVpcAttachment.AttachmentStatus)
}

func testAccCheckAwsTgwVpcAttachmentExists(n string, awsTgwVpcAttachment *api.AwsTgwVpcAttachment) error {
	terraformOptions := &terraform.Options{TerraformDir: "../"}
	output := terraform.Show(t, terraformOptions, "json")
	resources := output["values"].(map[string]interface{})["root_module"].(map[string]interface{})["resources"].([]interface{})

	for _, res := range resources {
		if res.(map[string]interface{})["type"].(string) != "aviatrix_aws_tgw_vpc_attachment" {
			continue
		}

		attributes := res.(map[string]interface{})["values"].(map[string]interface{})
		if attributes["tgw_name"].(string) != os.Getenv("AWS_TGW_NAME") || attributes["network_domain_name"].(string) != awsTgwVpcAttachment.SecurityDomainName || attributes["vpc_id"].(string) != os.Getenv("AWS_VPC_ID") {
			continue
		}

		awsTgwVpcAttachment.TgwName = attributes["tgw_name"].(string)
		awsTgwVpcAttachment.SecurityDomainName = attributes["network_domain_name"].(string)
		awsTgwVpcAttachment.VpcID = attributes["vpc_id"].(string)
		awsTgwVpcAttachment.AttachmentStatus = attributes["attachment_status"].(string)

		return nil
	}

	return fmt.Errorf("AWS TGW VPC ATTACH not found: %s", n)
}


func testAccCheckAwsTgwVpcAttachmentDestroy(t *testing.T, awsTgwVpcAttachment *goaviatrix.AwsTgwVpcAttachment) {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	err := retry.DoWithRetry(t, "Waiting for AWS TGW VPC attachment to be destroyed", 30, 5*time.Second, func() (string, error) {
		foundAwsTgwVpcAttachment, err := client.GetAwsTgwVpcAttachment(awsTgwVpcAttachment)
		if err != nil {
			if strings.Contains(err.Error(), "no such resource") {
				return "", nil
			}
			return "", fmt.Errorf("failed to get AWS TGW VPC attachment: %v", err)
		}
		if foundAwsTgwVpcAttachment != nil {
			return "", fmt.Errorf("AWS TGW VPC attachment still exists")
		}
		return "AWS TGW VPC attachment destroyed", nil
	})

	require.NoError(t, err)
}