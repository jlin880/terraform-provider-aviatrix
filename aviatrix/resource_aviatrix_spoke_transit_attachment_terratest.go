package aviatrix_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAccAviatrixSpokeTransitAttachment_basic(t *testing.T) {
	var spokeTransitAttachment goaviatrix.SpokeTransitAttachment

	rName := randomUniqueName("tfs")
	resourceName := "aviatrix_spoke_transit_attachment.test"

	skipAcc := os.Getenv("SKIP_SPOKE_TRANSIT_ATTACHMENT")
	if skipAcc == "yes" {
		t.Skip("Skipping spoke transit attachment tests as 'SKIP_SPOKE_TRANSIT_ATTACHMENT' is set")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: ".",
		Vars: map[string]interface{}{
			"account_name":       fmt.Sprintf("tfa-%s", rName),
			"cloud_type":         "1",
			"aws_account_number": os.Getenv("AWS_ACCOUNT_NUMBER"),
			"aws_iam":            false,
			"aws_access_key":     os.Getenv("AWS_ACCESS_KEY"),
			"aws_secret_key":     os.Getenv("AWS_SECRET_KEY"),
			"gw_name1":           fmt.Sprintf("tfs-%s", rName),
			"gw_name2":           fmt.Sprintf("tft-%s", rName),
			"vpc_id1":            os.Getenv("AWS_VPC_ID"),
			"vpc_id2":            os.Getenv("AWS_VPC_ID2"),
			"vpc_reg1":           os.Getenv("AWS_REGION"),
			"vpc_reg2":           os.Getenv("AWS_REGION2"),
			"gw_size":            "t2.micro",
			"subnet1":            os.Getenv("AWS_SUBNET"),
			"subnet2":            os.Getenv("AWS_SUBNET2"),
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	assert.NoError(t, terraform.OutputStruct(resourceName, &struct{}{}))

	// Import the resource using the resource ID
	importedResource := resourceName + "-imported"
	importedTerraformOptions := &terraform.Options{
		TerraformDir: ".",
		Vars: map[string]interface{}{
			"account_name":       fmt.Sprintf("tfa-%s", rName),
			"cloud_type":         "1",
			"aws_account_number": os.Getenv("AWS_ACCOUNT_NUMBER"),
			"aws_iam":            false,
			"aws_access_key":     os.Getenv("AWS_ACCESS_KEY"),
			"aws_secret_key":     os.Getenv("AWS_SECRET_KEY"),
			"gw_name1":           fmt.Sprintf("tfs-%s", rName),
			"gw_name2":           fmt.Sprintf("tft-%s", rName),
			"vpc_id1":            os.Getenv("AWS_VPC_ID"),
			"vpc_id2":            os.Getenv("AWS_VPC_ID2"),
			"vpc_reg1":           os.Getenv("AWS_REGION"),
			"vpc_reg2":           os.Getenv("AWS_REGION2"),
			"gw_size":            "t2.micro",
			"subnet1":            os.Getenv("AWS_SUBNET"),
			"subnet2":            os.Getenv("AWS_SUBNET2"),
		},
		ImportState:       fmt.Sprintf("%s=%s", resourceName, terraform.Output(t, resourceName)),
		ImportStateVerify: true,
	}
}
func testAccCheckSpokeTransitAttachmentExists(n string, spokeTransitAttachment *goaviatrix.SpokeTransitAttachment) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("spoke transit attachment Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no spoke transit attachment ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		foundSpokeTransitAttachment := &goaviatrix.SpokeTransitAttachment{
			SpokeGwName:   rs.Primary.Attributes["spoke_gw_name"],
			TransitGwName: rs.Primary.Attributes["transit_gw_name"],
		}
		foundSpokeTransitAttachment2, err := client.GetSpokeTransitAttachment(foundSpokeTransitAttachment)
		if err != nil {
			return err
		}
		if foundSpokeTransitAttachment2.SpokeGwName+"~"+foundSpokeTransitAttachment2.TransitGwName != rs.Primary.ID {
			return fmt.Errorf("spoke transit attachment not found")
		}

		*spokeTransitAttachment = *foundSpokeTransitAttachment2
		return nil
	}
}

func testAccCheckSpokeTransitAttachmentDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_spoke_transit_attachment" {
			continue
		}

		foundSpokeTransitAttachment := &goaviatrix.SpokeTransitAttachment{
			SpokeGwName:   rs.Primary.Attributes["spoke_gw_name"],
			TransitGwName: rs.Primary.Attributes["transit_gw_name"],
		}

		_, err := client.GetSpokeTransitAttachment(foundSpokeTransitAttachment)
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("spoke transit attachment still exists %s", err.Error())
		}
	}

	return nil
}
