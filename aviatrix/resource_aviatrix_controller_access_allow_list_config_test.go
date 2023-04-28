package test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAccAviatrixControllerAccessAllowListConfig_basic(t *testing.T) {
	skipAcc := os.Getenv("SKIP_CONTROLLER_ACCESS_ALLOW_LIST_CONFIG")
	if skipAcc == "yes" {
		t.Skip("Skipping Controller Access Allow List Config test as SKIP_CONTROLLER_ACCESS_ALLOW_LIST_CONFIG is set")
	}

	testName := fmt.Sprintf("terratest-%s", strings.ToLower(random.UniqueId()))
	terraformOptions := &terraform.Options{
		TerraformDir: "./",
		Vars: map[string]interface{}{
			"allow_list": []interface{}{
				map[string]interface{}{
					"ip_address": "0.0.0.0",
				},
			},
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	id := strings.Replace(testName, "-", ".", -1)

	client := goaviatrix.NewClient(os.Getenv("AVIATRIX_CONTROLLER_USERNAME"), os.Getenv("AVIATRIX_CONTROLLER_PASSWORD"), os.Getenv("AVIATRIX_CONTROLLER_URL"), "", "", false)

	err := client.Login(context.Background())
	assert.Nil(t, err)

	defer client.Logout(context.Background())

	_, err = client.GetControllerAccessAllowList(context.Background())
	assert.Nil(t, err)

	// Assert that the allow list has been configured properly
	assert.Equal(t, id, client.ControllerIP)

	// Import the resource and verify that it matches the current state
	importedResource := terraform.ImportState(t, terraformOptions, fmt.Sprintf("aviatrix_controller_access_allow_list_config.test[%s]", id))

	assert.Equal(t, id, importedResource.Primary.ID)
	assert.Equal(t, "0.0.0.0", importedResource.Primary.Attributes["allow_list.0.ip_address"])
}

func testAccCheckControllerAccessAllowListConfigExists(ctx context.Context, client *goaviatrix.Client, resourceName string) error {
	rs, err := terraform.GetResource(ctx, resourceName)
	if err != nil {
		return err
	}

	if rs.Primary.ID == "" {
		return fmt.Errorf("no controller access allow list config ID is set")
	}

	if strings.Replace(client.ControllerIP, ".", "-", -1) != rs.Primary.ID {
		return fmt.Errorf("controller access allow list config ID not found")
	}

	return nil
}

func testAccCheckControllerAccessAllowListConfigDestroy(ctx context.Context, client *goaviatrix.Client, s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_controller_access_allow_list_config" {
			continue
		}

		_, err := client.GetControllerAccessAllowList(ctx)
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("controller access allow list config still exists")
		}
	}

	return nil
}
