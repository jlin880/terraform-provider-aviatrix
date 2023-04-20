package test

import (
	"fmt"
	"os"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func configureTerraformOptions(rName string) (*terraform.Options, error) {
	awsRegion := os.Getenv("AWS_REGION")
	awsAccountNumber := os.Getenv("AWS_ACCOUNT_NUMBER")
	awsAccessKey := os.Getenv("AWS_ACCESS_KEY")
	awsSecretKey := os.Getenv("AWS_SECRET_KEY")

	terraformOptions := &terraform.Options{
		TerraformDir: "../examples/",
		Vars: map[string]interface{}{
			"aws_region":         awsRegion,
			"aws_account_number": awsAccountNumber,
			"aws_access_key":     awsAccessKey,
			"aws_secret_key":     awsSecretKey,
		},
	}

	return terraformOptions, nil
}

func aviatrixClientFromResourceState(t *testing.T, resourceState map[string]interface{}) *goaviatrix.Client {
	cid, ok := resourceState["cid"].(string)
	assert.True(t, ok, fmt.Sprintf("Expected to get CID from resource state but did not get it: %v", resourceState))

	client := goaviatrix.NewClient(cid, "")
	err := client.Login()
	assert.NoError(t, err, "Failed to authenticate to Aviatrix Controller")

	return client
}
