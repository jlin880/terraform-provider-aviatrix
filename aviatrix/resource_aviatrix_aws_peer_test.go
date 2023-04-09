package aviatrix_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/goaviatrix"
	"github.com/gruntwork-io/terratest/modules/acctest"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAviatrixAWSPeer_basic(t *testing.T) {
	t.Parallel()

	// Get AWS region and VPC IDs from environment variables
	awsRegion1 := os.Getenv("AWS_REGION")
	awsRegion2 := os.Getenv("AWS_REGION2")
	vpcID1 := os.Getenv("AWS_VPC_ID")
	vpcID2 := os.Getenv("AWS_VPC_ID2")

	// Check if the necessary environment variables are set
	if awsRegion1 == "" || awsRegion2 == "" || vpcID1 == "" || vpcID2 == "" {
		t.Fatal("Missing environment variables: AWS_REGION, AWS_REGION2, AWS_VPC_ID, or AWS_VPC_ID2")
	}

	// Generate a random name to avoid collisions with existing resources
	randName := random.UniqueId()

	terraformOptions := &terraform.Options{
		// The path to the Terraform code to test
		TerraformDir: "../examples/aws_peer",

		// Variables to pass to the Terraform code during the test
		Vars: map[string]interface{}{
			"account_name1": fmt.Sprintf("tf-testing-%s-1", randName),
			"account_name2": fmt.Sprintf("tf-testing-%s-2", randName),
			"vpc_id1":       vpcID1,
			"vpc_id2":       vpcID2,
			"vpc_reg1":      awsRegion1,
			"vpc_reg2":      awsRegion2,
		},
	}

	// Delete the resources at the end of the test
	defer terraform.Destroy(t, terraformOptions)

	// Deploy the Terraform code
	terraform.InitAndApply(t, terraformOptions)

	// Check if the AWS peer exists
	awsPeer := goaviatrix.AWSPeer{
		VpcID1: vpcID1,
		VpcID2: vpcID2,
	}

	client := goaviatrix.NewClient(os.Getenv("AVIATRIX_CONTROLLER_IP"), os.Getenv("AVIATRIX_USERNAME"), os.Getenv("AVIATRIX_PASSWORD"))

	err := client.GetAWSPeer(&awsPeer)
	assert.NoError(t, err)

	// Check if the AWS peer attributes match the Terraform configuration
	assert.Equal(t, terraformOptions.Vars["account_name1"].(string), awsPeer.AccountName1)
	assert.Equal(t, terraformOptions.Vars["account_name2"].(string), awsPeer.AccountName2)
	assert.Equal(t, terraformOptions.Vars["vpc_id1"].(string), awsPeer.VpcID1)
	assert.Equal(t, terraformOptions.Vars["vpc_id2"].(string), awsPeer.VpcID2)
	assert.Equal(t, terraformOptions.Vars["vpc_reg1"].(string), awsPeer.VpcReg1)
	assert.Equal(t, terraformOptions.Vars["vpc_reg2"].(string), awsPeer.VpcReg2)
}


func testAccAWSPeerConfigBasic(rInt int, vpcID1 string, vpcID2 string, region1 string, region2 string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_account" {
	account_name       = "tf-testing-%d"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}
resource "aviatrix_aws_peer" "test_aws_peer" {
	account_name1 = aviatrix_account.test_account.account_name
	account_name2 = aviatrix_account.test_account.account_name
	vpc_id1       = "%s"
	vpc_id2       = "%s"
	vpc_reg1      = "%s"
	vpc_reg2      = "%s"
}
	`, rInt, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		vpcID1, vpcID2, region1, region2)
}

func tesAccCheckAWSPeerExists(n string, awsPeer *goaviatrix.AWSPeer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("awsPeer Not Created: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no AWSPeer ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		foundPeer := &goaviatrix.AWSPeer{
			VpcID1: rs.Primary.Attributes["vpc_id1"],
			VpcID2: rs.Primary.Attributes["vpc_id2"],
		}

		_, err := client.GetAWSPeer(foundPeer)
		if err != nil {
			return err
		}
		if foundPeer.VpcID1 != rs.Primary.Attributes["vpc_id1"] {
			return fmt.Errorf("vpc_id1 Not found in created attributes")
		}
		if foundPeer.VpcID2 != rs.Primary.Attributes["vpc_id2"] {
			return fmt.Errorf("vpc_id2 Not found in created attributes")
		}

		*awsPeer = *foundPeer
		return nil
	}
}

func testAccCheckAWSPeerDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_aws_peer" {
			continue
		}

		foundPeer := &goaviatrix.AWSPeer{
			VpcID1: rs.Primary.Attributes["vpc_id1"],
			VpcID2: rs.Primary.Attributes["vpc_id2"],
		}

		_, err := client.GetAWSPeer(foundPeer)
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("awsPeer still exists")
		}
	}

	return nil
}
