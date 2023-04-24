package test

import (
	"context"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAviatrixSmartGroup(t *testing.T) {
	terraformOptions := &terraform.Options{
		TerraformDir: "./",
		Vars: map[string]interface{}{
			"name": "test-smart-group",
			"selector": []map[string]interface{}{
				{
					"match_expressions": []map[string]interface{}{
						{
							"cidr": "11.0.0.0/16",
						},
						{
							"type":         "vm",
							"account_name": "devops",
							"region":       "us-west-2",
							"tags": map[string]string{
								"k3": "v3",
							},
						},
					},
				},
			},
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	client := goaviatrix.NewClient("AVIATRIX_CONTROLLER_IP", "AVIATRIX_USERNAME", "AVIATRIX_PASSWORD", true)

	smartGroup, err := client.GetSmartGroup(context.Background(), "test-smart-group")
	assert.NoError(t, err)
	assert.NotNil(t, smartGroup)

	assert.Equal(t, "test-smart-group", smartGroup.Name)
	assert.Equal(t, 2, len(smartGroup.Selector[0].MatchExpressions))
	assert.Equal(t, "11.0.0.0/16", smartGroup.Selector[0].MatchExpressions[0].Cidr)
	assert.Equal(t, "vm", smartGroup.Selector[0].MatchExpressions[1].Type)
	assert.Equal(t, "devops", smartGroup.Selector[0].MatchExpressions[1].AccountName)
	assert.Equal(t, "us-west-2", smartGroup.Selector[0].MatchExpressions[1].Region)
	assert.Equal(t, map[string]string{"k3": "v3"}, smartGroup.Selector[0].MatchExpressions[1].Tags)
}

func TestAviatrixSmartGroup_update(t *testing.T) {
	terraformOptions := &terraform.Options{
		TerraformDir: "./",
		Vars: map[string]interface{}{
			"name": "test-smart-group",
			"selector": []map[string]interface{}{
				{
					"match_expressions": []map[string]interface{}{
						{
							"cidr": "11.0.0.0/16",
						},
					},
				},
			},
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	terraformOptions.Vars["selector"] = []map[string]interface{}{
		{
			"match_expressions": []map[string]interface{}{
				{
					"cidr": "11.0.0.0/16",
				},
				{
					"type":         "vm",
					"account_name": "devops",
					"region":       "us-west-2",
					"tags": map[string]string{
						"k3": "v3",
					},
				},
			},
		},
	}

	terraform.Apply(t, terraformOptions)

	smartGroupName := terraform.Output(t, terraformOptions, "name")
	smartGroupID := terraform.Output(t, terraformOptions, "id")

	client := goaviatrix.NewClient("AVIATRIX_CONTROLLER_IP", "AVIATRIX_USERNAME", "AVIATRIX_PASSWORD", true)

	smartGroup, err := client.GetSmartGroup(context.Background(), smartGroupID)
	assert.NoError(t, err)
	assert.NotNil(t, smartGroup)

	assert.Equal(t, smartGroupName, smartGroup.Name)
	assert.Equal(t, 2, len(smartGroup.Selector[0].MatchExpressions))
	assert.Equal(t, "11.0.0.0/16", smartGroup.Selector[0].MatchExpressions[0].Cidr)
	assert.Equal(t, "vm", smartGroup.Selector[0].MatchExpressions[1].Type)
	assert.Equal(t, "devops", smartGroup.Selector[0].MatchExpressions[1].AccountName)
	assert.Equal(t, "us-west-2", smartGroup.Selector[0].MatchExpressions[1].Region)
	assert.Equal(t, map[string]string{"k3": "v3"}, smartGroup.Selector[0].MatchExpressions[1].Tags)
}


func testAccCheckSmartGroupExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("no Smart Group resource found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no Smart Group ID is set")
		}

		client := testAccProviderVersionValidation.Meta().(*goaviatrix.Client)

		smartGroup, err := client.GetSmartGroup(context.Background(), rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed to get Smart Group status: %v", err)
		}

		if smartGroup.UUID != rs.Primary.ID {
			return fmt.Errorf("smart Group ID not found")
		}

		return nil
	}
}

func testAccSmartGroupDestroy(s *terraform.State) error {
	client := testAccProviderVersionValidation.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_smart_group" {
			continue
		}

		_, err := client.GetSmartGroup(context.Background(), rs.Primary.ID)
		if err == nil || err != goaviatrix.ErrNotFound {
			return fmt.Errorf("smart group configured when it should be destroyed")
		}
	}

	return nil
}
