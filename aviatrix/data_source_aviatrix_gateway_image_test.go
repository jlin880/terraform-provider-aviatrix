package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/assert"
)

func TestTerraformAviatrixDataSourceGatewayImage(t *testing.T) {
	terraformOptions := &terraform.Options{
		TerraformDir: ".",
		EnvVars: map[string]string{
			"SKIP_DATA_GATEWAY_IMAGE": os.Getenv("SKIP_DATA_GATEWAY_IMAGE"),
		},
	}

	resourceName := "data.aviatrix_gateway_image.foo"

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	output := terraform.Output(t, terraformOptions, "image_version")
	expectedOutput := "hvm-cloudx-aws-022021"

	assert.Equal(t, expectedOutput, output)

	// check if resource exists
	res := terraform.GetStateResource(t, terraformOptions, resourceName)
	if res == nil {
		t.Fatalf("resource '%s' does not exist", resourceName)
	}
}

func TestTerraformAviatrixDataSourceGatewayImageSkip(t *testing.T) {
	terraformOptions := &terraform.Options{
		TerraformDir: ".",
		EnvVars: map[string]string{
			"SKIP_DATA_GATEWAY_IMAGE": "yes",
		},
	}

	resourceName := "data.aviatrix_gateway_image.foo"

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	// check if resource exists
	res := terraform.GetStateResource(t, terraformOptions, resourceName)
	if res != nil {
		t.Fatalf("resource '%s' exists even though 'SKIP_DATA_GATEWAY_IMAGE' is set", resourceName)
	}
}
func TestAccDataSourceAviatrixGatewayImage_basic(t *testing.T) {
	resourceName := "data.aviatrix_gateway_image.foo"

	skipAcc := os.Getenv("SKIP_DATA_GATEWAY_IMAGE")
	if skipAcc == "yes" {
		t.Skip("Skipping Data Source Gateway Image test as SKIP_DATA_GATEWAY_IMAGE is set")
	}

	test_structure.RunTestSuite(t, &test_structure.TestSuite{
		PreTestSteps: []test_structure.TestStep{
			{
				Config: testAccProviderVersionCheck,
			},
			{
				Config: testAccProviderLogin,
			},
			{
				Config: testAccSetAviatrixCIDR,
			},
			{
				Config: testAccSetAWSRegion,
			},
			{
				Config: testAccSetAzureRegion,
			},
		},
		TestSteps: []test_structure.TestStep{
			{
				Config: testAccDataSourceAviatrixGatewayImageConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceAviatrixGatewayImage(resourceName),
					resource.TestCheckResourceAttr(resourceName, "image_version", "hvm-cloudx-aws-022021"),
				),
			},
		},
		PostTestSteps: []test_structure.TestStep{
			{
				Config: testAccProviderLogout,
			},
		},
	})
}

func testAccDataSourceAviatrixGatewayImageConfigBasic() string {
	return `
data "aviatrix_gateway_image" "foo" {
	cloud_type       = 1
	software_version = "6.5" 
}
	`
}

func testAccDataSourceAviatrixGatewayImage(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("root module has no data source called %s", name)
		}

		return nil
	}
}
