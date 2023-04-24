package aviatrix

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)
func TestAccAviatrixCopilotSecurityGroupManagementConfig_basic(t *testing.T) {
	if os.Getenv("SKIP_COPILOT_SECURITY_GROUP_MANAGEMENT_CONFIG") == "yes" {
		t.Skip("Skipping copilot security group management config test as SKIP_COPILOT_SECURITY_GROUP_MANAGEMENT_CONFIG is set")
	}

	resourceName := "aviatrix_copilot_security_group_management_config.test"
	rName := acctest.RandString(5)

	ctx := context.Background()
	testAccProviderVersionValidation := testAccProvider

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCopilotSecurityGroupManagementConfigDestroy(ctx, testAccProvider),
		Steps: []resource.TestStep{
			{
				Config: testAccCopilotSecurityGroupManagementConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCopilotSecurityGroupManagementConfigExists(ctx, resourceName, testAccProvider),
					resource.TestCheckResourceAttr(resourceName, "account_name", fmt.Sprintf("tfa-aws-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", os.Getenv("AWS_VPC_ID")),
					resource.TestCheckResourceAttr(resourceName, "region", os.Getenv("AWS_REGION")),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckCopilotSecurityGroupManagementConfigExists(ctx context.Context, resourceName string, provider *schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("could not find copilot security group management config: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("copilot security group management config id is not set")
		}

		client := provider.Meta().(*goaviatrix.Client)

		if strings.Replace(client.ControllerIP, ".", "-", -1) != rs.Primary.ID {
			return fmt.Errorf("could not find copilot security group management id")
		}
		return nil
	}
}

func testAccCheckCopilotSecurityGroupManagementConfigDestroy(ctx context.Context, provider *schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := provider.Meta().(*goaviatrix.Client)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aviatrix_copilot_security_group_management_config" {
				continue
			}

			copilotSecurityGroupManagementConfig, err := client.GetCopilotSecurityGroupManagementConfig(ctx)
			if err != nil {
				return fmt.Errorf("could not read copilot security group management config due to err: %v", err)
			}
			if copilotSecurityGroupManagementConfig.State == "Enabled" {
				return fmt.Errorf("copilot security group management is still enabled")
			}
		}
		return nil
	}
}
