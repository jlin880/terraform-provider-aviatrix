package aviatrix

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/avianto/go-terraform-aviatrix/helper"
	"github.com/avianto/go-terraform-aviatrix/provider"
	"github.com/avianto/go-terraform-aviatrix/terraform"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAccAviatrixEdgeSpokeTransitAttachment_basic(t *testing.T) {
	resourceName := fmt.Sprintf("aviatrix_edge_spoke_transit_attachment.test-%s", strings.ToLower(random.UniqueId()))

	skipAcc := os.Getenv("SKIP_EDGE_SPOKE_TRANSIT_ATTACHMENT")
	if skipAcc == "yes" {
		t.Skip("Skipping Edge as a Spoke transit attachment tests as 'SKIP_EDGE_SPOKE_TRANSIT_ATTACHMENT' is set")
	}

	terraformOptions, err := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/edge_spoke_transit_attachment/",
		Vars: map[string]interface{}{
			"acc_name":         fmt.Sprintf("tfa-%s", strings.ToLower(random.UniqueId())),
			"aws_account":      os.Getenv("AWS_ACCOUNT_NUMBER"),
			"aws_access_key":   os.Getenv("AWS_ACCESS_KEY"),
			"aws_secret_key":   os.Getenv("AWS_SECRET_KEY"),
			"aws_region":       os.Getenv("AWS_REGION"),
			"aws_vpc_id":       os.Getenv("AWS_VPC_ID"),
			"aws_subnet":       os.Getenv("AWS_SUBNET"),
			"spoke_gw_name":    os.Getenv("EDGE_SPOKE_NAME"),
			"transit_gw_name":  fmt.Sprintf("tft-%s", strings.ToLower(random.UniqueId())),
			"vpc_region":       os.Getenv("AWS_REGION"),
			"vpc_account_name": fmt.Sprintf("tfa-%s", strings.ToLower(random.UniqueId())),
			"gw_size":          "t2.micro",
		},
		NoColor: true,
	})
	if err != nil {
		t.Fatal(err)
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	// Verify the Edge as a Spoke Transit Attachment exists
	err = helper.Retry(func() error {
		attachment := &provider.SpokeTransitAttachment{
			SpokeGwName:   terraformOptions.Vars["spoke_gw_name"].(string),
			TransitGwName: terraformOptions.Vars["transit_gw_name"].(string),
		}
		if err := provider.GetEdgeSpokeTransitAttachment(context.Background(), attachment); err != nil {
			return err
		}

		return nil
	})

	assert.NoError(t, err)

	// Verify the Terraform state
	testAccCheckEdgeSpokeTransitAttachmentExists(resourceName, terraformOptions, t)
}

func preEdgeSpokeTransitAttachmentCheck(t *testing.T) {
	if os.Getenv("EDGE_SPOKE_NAME") == "" {
		t.Fatal("Environment variable EDGE_SPOKE_NAME is not set")
	}
}

func testAccEdgeSpokeTransitAttachmentConfigBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test" {
	account_name       = "tfa-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}
resource "aviatrix_transit_gateway" "test" {
	cloud_type   = 1
	account_name = aviatrix_account.test.account_name
	gw_name      = "tft-%s"
	vpc_id       = "%s"
	vpc_reg      = "%s"
	gw_size      = "t2.micro"
	subnet       = "%s"
}
resource "aviatrix_edge_spoke_transit_attachment" "test" {
	spoke_gw_name   = "%s"
	transit_gw_name = aviatrix_transit_gateway.test.gw_name
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		rName, os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_REGION"), os.Getenv("AWS_SUBNET"),
		os.Getenv("EDGE_SPOKE_NAME"))
}

func testAccCheckEdgeSpokeTransitAttachmentExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("edge as a spoke transit attachment not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no edge as a spoke transit attachment ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		attachment := &goaviatrix.SpokeTransitAttachment{
			SpokeGwName:   rs.Primary.Attributes["spoke_gw_name"],
			TransitGwName: rs.Primary.Attributes["transit_gw_name"],
		}
		attachment, err := client.GetEdgeSpokeTransitAttachment(context.Background(), attachment)
		if err != nil {
			return err
		}
		if attachment.SpokeGwName+"~"+attachment.TransitGwName != rs.Primary.ID {
			return fmt.Errorf("edge as a spoke transit attachment not found")
		}

		return nil
	}
}

func testAccCheckEdgeSpokeTransitAttachmentDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_spoke_transit_attachment" {
			continue
		}

		attachment := &goaviatrix.SpokeTransitAttachment{
			SpokeGwName:   rs.Primary.Attributes["spoke_gw_name"],
			TransitGwName: rs.Primary.Attributes["transit_gw_name"],
		}

		_, err := client.GetEdgeSpokeTransitAttachment(context.Background(), attachment)
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("edge as a spoke transit attachment still exists %s", err.Error())
		}
	}

	return nil
}
