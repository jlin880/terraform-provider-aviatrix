package aviatrix_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAviatrixControllerBgpMaxAsLimitConfig_basic(t *testing.T) {
	t.Parallel()

	skipAcc := os.Getenv("SKIP_CONTROLLER_BGP_MAX_AS_LIMIT_CONFIG")
	if skipAcc == "yes" {
		t.Skip("Skipping Controller BGP Max AS Limit Config test as SKIP_CONTROLLER_BGP_MAX_AS_LIMIT_CONFIG is set")
	}

	resourceName := fmt.Sprintf("aviatrix_controller_bgp_max_as_limit_config.test-%s", strings.ToLower(random.UniqueId()))

	terraformOptions := &terraform.Options{
		TerraformDir: "../examples/controller_bgp_max_as_limit_config/",
		Vars: map[string]interface{}{
			"max_as_limit": 1,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	checkControllerBgpMaxAsLimitConfigExists(t, terraformOptions, resourceName)
}

func checkControllerBgpMaxAsLimitConfigExists(t *testing.T, terraformOptions *terraform.Options, resourceName string) {
	client := goaviatrix.NewClient(
		terraform.Output(t, terraformOptions, "controller_endpoint"),
		terraform.Output(t, terraformOptions, "username"),
		terraform.Output(t, terraformOptions, "password"),
		"",
		terraform.Output(t, terraformOptions, "aws_account"),
		"",
		"",
		terraform.Output(t, terraformOptions, "aws_role_arn"),
		"",
		"",
		"",
		true,
	)

	bgpMaxAsLimit, err := client.GetControllerBgpMaxAsLimit(context.Background())
	if err != nil {
		t.Fatalf("Failed to get controller BGP Max AS Limit config status: %v", err)
	}

	assert.Equal(t, 1, bgpMaxAsLimit, "BGP Max AS Limit value is not as expected")
}
