package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/random"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAccAviatrixTrafficClassifier_basic(t *testing.T) {
	if os.Getenv("SKIP_TRAFFIC_CLASSIFIER") == "yes" {
		t.Skip("Skipping traffic classifier test as SKIP_TRAFFIC_CLASSIFIER is set")
	}

	resourceName := "aviatrix_traffic_classifier.test"
	policyName := "policy-" + random.UniqueId()

	terraformOptions := &terraform.Options{
		TerraformDir: "./",
		Vars: map[string]interface{}{
			"traffic_classifier_name": policyName,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	assert.NoError(t, terraform.OutputStruct(resourceName, &struct{}{}))

	// Import the resource using the resource ID
	importedResource := resourceName + "-imported"
	importedTerraformOptions := &terraform.Options{
		TerraformDir: "./",
		Vars: map[string]interface{}{
			"traffic_classifier_name": policyName,
		},
	}
	importedTerraformOptions.ImportState = fmt.Sprintf("%s=%s", resourceName, terraform.Output(t, resourceName))
	importedTerraformOptions.ImportStateVerify = true

	terraform.InitAndApply(t, importedTerraformOptions)

	assert.NoError(t, terraform.OutputStruct(importedResource, &struct{}{}))
}

func TestAccAviatrixTrafficClassifier_import(t *testing.T) {
	if os.Getenv("SKIP_TRAFFIC_CLASSIFIER") == "yes" {
		t.Skip("Skipping traffic classifier test as SKIP_TRAFFIC_CLASSIFIER is set")
	}

	resourceName := "aviatrix_traffic_classifier.test"
	policyName := "policy-" + random.UniqueId()

	terraformOptions := &terraform.Options{
		TerraformDir: "./",
		Vars: map[string]interface{}{
			"traffic_classifier_name": policyName,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	assert.NoError(t, terraform.OutputStruct(resourceName, &struct{}{}))

	// Import the resource using the resource ID
	importedResource := resourceName + "-imported"
	importedTerraformOptions := &terraform.Options{
		TerraformDir: "./",
		Vars: map[string]interface{}{
			"traffic_classifier_name": policyName,
		},
	}
	importedTerraformOptions.ImportState = fmt.Sprintf("%s=%s", resourceName, terraform.Output(t, resourceName))
	importedTerraformOptions.ImportStateVerify = true

	terraform.InitAndApply(t, importedTerraformOptions)

	assert.NoError(t, terraform.OutputStruct(importedResource, &struct{}{}))
}
