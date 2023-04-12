package aviatrix_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/goaviatrix"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAccDataSourceAviatrixCallerIdentity_basic(t *testing.T) {
	t.Parallel()

	rName := random.UniqueId()
	resourceName := "data.aviatrix_caller_identity.foo"

	skipAcc := os.Getenv("SKIP_DATA_CALLER_IDENTITY")
	if skipAcc == "yes" {
		t.Skip("Skipping Data Source Caller Identity test as SKIP_DATA_CALLER_IDENTITY is set")
	}

	terraformOptions, err := configureTerraformOptions(rName)
	if err != nil {
		t.Fatal(err)
	}
	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	resourceState := terraform.OutputAll(t, terraformOptions)

	client := aviatrixClientFromResourceState(t, resourceState)

	version, _, err := client.GetCurrentVersion()
	assert.NoError(t, err)
	assert.Contains(t, version, ".")

}

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
