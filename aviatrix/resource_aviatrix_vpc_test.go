package test

import (
	"fmt"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAccAviatrixVpc_basic(t *testing.T) {
	var vpc goaviatrix.Vpc

	rName := RandomString(5)
	resourceName := "aviatrix_vpc.test_vpc"

	skipAcc := os.Getenv("SKIP_VPC")
	if skipAcc == "yes" {
		t.Skip("Skipping VPC tests as 'SKIP_VPC' is set")
	}

	skipAccAWS := os.Getenv("SKIP_VPC_AWS")
	skipAccAZURE := os.Getenv("SKIP_VPC_AZURE")
	skipAccGCP := os.Getenv("SKIP_VPC_GCP")
	if skipAccAWS == "yes" && skipAccAZURE == "yes" && skipAccGCP == "yes" {
		t.Skip("Skipping VPC tests as 'SKIP_VPC_AWS', 'SKIP_VPC_AZURE' and 'SKIP_VPC_GCP' are all set")
	}

	if skipAccAWS != "yes" {
		msgCommon := ". Set 'SKIP_VPC_AWS' to 'yes' to skip VPC tests in AWS"
		terraformOptions := &terraform.Options{
			TerraformDir: ".",
			Vars: map[string]interface{}{
				"account_name":       fmt.Sprintf("tfa-%s", rName),
				"aws_account_number": os.Getenv("AWS_ACCOUNT_NUMBER"),
				"aws_iam":            false,
				"aws_access_key":     os.Getenv("AWS_ACCESS_KEY"),
				"aws_secret_key":     os.Getenv("AWS_SECRET_KEY"),
				"vpc_name":           fmt.Sprintf("tfg-%s", rName),
				"region":             os.Getenv("AWS_REGION"),
				"cidr":               "10.0.0.0/16",
			},
		}

		defer terraform.Destroy(t, terraformOptions)

		testAccCheckVpcExists(t, terraformOptions, resourceName, &vpc)

		assert.Equal(t, fmt.Sprintf("tfg-%s", rName), vpc.Name)
		assert.Equal(t, fmt.Sprintf("tfa-%s", rName), vpc.AccountName)
		assert.Equal(t, "1", vpc.CloudType)
		assert.Equal(t, "10.0.0.0/16", vpc.CIDR)
		assert.Equal(t, os.Getenv("AWS_REGION"), vpc.Region)
	} else {
		t.Log("Skipping VPC tests in AWS as 'SKIP_VPC_AWS' is set")
	}

	if skipAccGCP != "yes" {
		msgCommon := ". Set 'SKIP_VPC_GCP' to 'yes' to skip VPC tests in GCP"
		terraformOptions := &terraform.Options{
			TerraformDir: ".",
			Vars: map[string]interface{}{
				"account_name":                  fmt.Sprintf("tfa-%s", rName),
				"gcloud_project_id":             os.Getenv("GCP_ID"),
				"gcloud_project_credentials":    os.Getenv("GCP_CREDENTIALS_FILEPATH"),
				"vpc_name":                      fmt.Sprintf("tfg-%s", rName),
				"subnets.0.region":              "
