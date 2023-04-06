package aviatrix

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAccAviatrixCopilotAssociation_basic(t *testing.T) {
	skipAcc := os.Getenv("SKIP_COPILOT_ASSOCIATION")
	if skipAcc == "yes" {
		t.Skip("Skipping Copilot Association test as SKIP_COPILOT_ASSOCIATION is set")
	}
	resourceName := "aviatrix_copilot_association.test"
	testAccVersion := os.Getenv("TESTACC_AVIATRIX_VERSION")

	terraformOptions := &terraform.Options{
		TerraformDir: "../../_example/copilot_association",
		Upgrade:      true,
		Vars: map[string]interface{}{
			"copilot_address": "aviatrix.com",
		},
	}

	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApply(t, terraformOptions)

	assert.NoError(t, terraform.InitAndApplyE(t, terraformOptions))

	status, err := aviatrixClient.GetCopilotAssociationStatus(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, "Succeed", status.Status)
}

func TestAccAviatrixCopilotAssociation_update(t *testing.T) {
	skipAcc := os.Getenv("SKIP_COPILOT_ASSOCIATION")
	if skipAcc == "yes" {
		t.Skip("Skipping Copilot Association test as SKIP_COPILOT_ASSOCIATION is set")
	}
	resourceName := fmt.Sprintf("aviatrix_copilot_association.test_%s", random.UniqueId())

	terraformOptions := &terraform.Options{
		TerraformDir: "../../_example/copilot_association",
		Upgrade:      true,
		Vars: map[string]interface{}{
			"copilot_address": "aviatrix.com",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	terraformOptions.Vars["copilot_address"] = "updated.aviatrix.com"

	terraform.Apply(t, terraformOptions)

	status, err := aviatrixClient.GetCopilotAssociationStatus(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, "Succeed", status.Status)
}

func TestAccAviatrixCopilotAssociation_destroy(t *testing.T) {
	skipAcc := os.Getenv("SKIP_COPILOT_ASSOCIATION")
	if skipAcc == "yes" {
		t.Skip("Skipping Copilot Association test as SKIP_COPILOT_ASSOCIATION is set")
	}
	resourceName := "aviatrix_copilot_association.test"

	terraformOptions := &terraform.Options{
		TerraformDir: "../../_example/copilot_association",
		Upgrade:      true,
		Vars: map[string]interface{}{
			"copilot_address": "aviatrix.com",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	assert.NoError(t, terraform.DestroyE(t, terraformOptions))

	_, err := aviatrixClient.GetCopilotAssociationStatus(context.Background())
	assert.Equal(t, goaviatrix.ErrNotFound, err)
}
