package test

import (
	"fmt"
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestTerraformAviatrixTrafficClassifier(t *testing.T) {
	t.Parallel()

	if os.Getenv("SKIP_TRAFFIC_CLASSIFIER") == "yes" {
		t.Skip("Skipping traffic classifier test as SKIP_TRAFFIC_CLASSIFIER is set")
	}

	trafficClassifierName := "policy-" + random.UniqueId()
	terraformOptions := &terraform.Options{
		TerraformDir: ".",
		Vars: map[string]interface{}{
			"traffic_classifier_name": trafficClassifierName,
		},
	}

	// Clean up everything at the end of the test
	defer terraform.Destroy(t, terraformOptions)

	// Deploy the infrastructure with Terraform
	terraform.InitAndApply(t, terraformOptions)

	// Verify the infrastructure was created successfully
	resourceName := "aviatrix_traffic_classifier.test"
	assert.NoError(t, terraform.OutputStruct(resourceName, &struct{}{}))

	// Import the resource using the resource ID
	importedResource := resourceName + "-imported"
	importedTerraformOptions := &terraform.Options{
		TerraformDir: ".",
		Vars: map[string]interface{}{
			"traffic_classifier_name": trafficClassifierName,
		},
	}
	importedTerraformOptions.ImportState = fmt.Sprintf("%s=%s", resourceName, terraform.Output(t, resourceName))
	importedTerraformOptions.ImportStateVerify = true

	// Verify the imported resource
	terraform.InitAndApply(t, importedTerraformOptions)
	assert.NoError(t, terraform.OutputStruct(importedResource, &struct{}{}))
}

func TestTerraformAviatrixTrafficClassifier_import(t *testing.T) {
	t.Parallel()

	if os.Getenv("SKIP_TRAFFIC_CLASSIFIER") == "yes" {
		t.Skip("Skipping traffic classifier test as SKIP_TRAFFIC_CLASSIFIER is set")
	}

	trafficClassifierName := "policy-" + random.UniqueId()
	terraformOptions := &terraform.Options{
		TerraformDir: ".",
		Vars: map[string]interface{}{
			"traffic_classifier_name": trafficClassifierName,
		},
	}

	// Clean up everything at the end of the test
	defer terraform.Destroy(t, terraformOptions)

	// Deploy the infrastructure with Terraform
	terraform.InitAndApply(t, terraformOptions)

	// Verify the infrastructure was created successfully
	resourceName := "aviatrix_traffic_classifier.test"
	assert.NoError(t, terraform.OutputStruct(resourceName, &struct{}{}))

	// Import the resource using the resource ID
	importedResource := resourceName + "-imported"
	importedTerraformOptions := &terraform.Options{
		TerraformDir: ".",
		Vars: map[string]interface{}{
			"traffic_classifier_name": trafficClassifierName,
		},
	}
	importedTerraformOptions.ImportState = fmt.Sprintf("%s=%s", resourceName, terraform.Output(t, resourceName))
	importedTerraformOptions.ImportStateVerify = true

	// Verify the imported resource
	terraform.InitAndApply(t, importedTerraformOptions)
	assert.NoError(t, terraform.OutputStruct(importedResource, &struct{}{}))
}
