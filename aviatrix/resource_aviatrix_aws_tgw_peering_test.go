package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)
func TestAccAviatrixAWSTgwPeering_basic(t *testing.T) {
	t.Parallel()

	awsRegion1 := "us-east-1"
	awsRegion2 := "us-east-2"
	terraformDir := "../examples/aviatrix-aws-tgw-peering"

	accountName := fmt.Sprintf("tfa-%s", strings.ToLower(random.UniqueId()))

	// Skip the test if the SKIP_AWS_TGW_PEERING env var is set to "yes"
	if os.Getenv("SKIP_AWS_TGW_PEERING") == "yes" {
		t.Skip("Skipping Aviatrix AWS TGW peering tests as 'SKIP_AWS_TGW_PEERING' is set")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: terraformDir,

		// Variables to pass to Terraform
		Vars: map[string]interface{}{
			"account_name":       accountName,
			"aws_region1":        awsRegion1,
			"aws_region2":        awsRegion2,
			"aws_account_number": os.Getenv("AWS_ACCOUNT_NUMBER"),
			"aws_access_key":     os.Getenv("AWS_ACCESS_KEY"),
			"aws_secret_key":     os.Getenv("AWS_SECRET_KEY"),
		},
	}

	// Destroy the Terraform infrastructure at the end of the test
	defer terraform.Destroy(t, terraformOptions)

	// Deploy the Terraform infrastructure
	terraform.InitAndApply(t, terraformOptions)

	// Verify the Aviatrix AWS TGW peering exists
	verifyAwsTgwPeeringExists(t, terraformOptions)
}

// Verify the Aviatrix AWS TGW peering exists
func verifyAwsTgwPeeringExists(t *testing.T, terraformOptions *terraform.Options) {
	t.Helper()

	// Get the Terraform output
	tgwName1 := terraform.Output(t, terraformOptions, "tgw_name1")
	tgwName2 := terraform.Output(t, terraformOptions, "tgw_name2")

	// Create an Aviatrix client
	client, err := goaviatrix.NewClientWithConfig(goaviatrix.ClientConfig{
		APIUser:     os.Getenv("AVIATRIX_API_USER"),
		APIPass:     os.Getenv("AVIATRIX_API_PASSWORD"),
		APIEndpoint: os.Getenv("AVIATRIX_API_ENDPOINT"),
	})
	assert.NoError(t, err)

	// Get the AWS TGW peering
	awsTgwPeering := &goaviatrix.AwsTgwPeering{
		TgwName1: tgwName1,
		TgwName2: tgwName2,
	}
	err = client.GetAwsTgwPeering(awsTgwPeering)
	assert.NoError(t, err)
	assert.NotNil(t, awsTgwPeering)
}
func testAccCheckAWSTgwPeeringDestroy(t *testing.T, terraformOptions *terraform.Options) {
    // Retrieve the provider client from the Terraform options
    client := terraformProvider.Meta().(*goaviatrix.Client)

    // Destroyed resources should no longer exist
    for _, resource := range terraform.ListResources(t, terraformOptions) {
        if resource.Type != "aviatrix_aws_tgw_peering" {
            continue
        }

        foundAwsTgwPeering := &goaviatrix.AwsTgwPeering{
            TgwName1: resource.Primary.Attributes["tgw_name1"],
            TgwName2: resource.Primary.Attributes["tgw_name2"],
        }

        err := client.GetAwsTgwPeering(foundAwsTgwPeering)
        if err != goaviatrix.ErrNotFound {
            t.Errorf("Resource %s still exists", resource.Primary.ID)
        }
    }
}
