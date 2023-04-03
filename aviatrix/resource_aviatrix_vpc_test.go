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

		terraform.InitAndApply(t, terraformOptions)

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
				"subnets.0.region":              "us-east1",
				"subnets.0.cidr":                "10.0

	} else {
		t.Log("Skipping VPC tests in GCP as 'SKIP_VPC_GCP' is set")
	}

	if os.Getenv("SKIP_VPC_AZURE") != "yes" {
		testAccCheckVpcExists(t, terraformOptions, resourceName, &vpc)
		testAccCheckVpcBasicAZURE(t, resourceName, &vpc)
	} else {
		t.Log("Skipping VPC tests in Azure as 'SKIP_VPC_AZURE' is set")
	}
}
func testAccCheckVpcExists(t *testing.T, terraformOptions *terraform.Options, resourceName string, vpc *goaviatrix.Vpc) {
	terraform.InitAndApply(t, terraformOptions)

	output := terraform.Show(t, terraformOptions)
	expectedOutput := fmt.Sprintf(`%s = {
  "account_name" = "%s"
  "cidr" = "10.0.0.0/16"
  "cloud_type" = "1"
  "id" = ""
  "name" = "%s"
  "region" = "%s"
  "subnets" = []
}`, resourceName, fmt.Sprintf("tfa-%s", terraformOptions.Vars["r_name"].(string)), fmt.Sprintf("tfg-%s", terraformOptions.Vars["r_name"].(string)), os.Getenv("AWS_REGION"))

	if output != expectedOutput {
		t.Fatalf("Unexpected output:\n%s\nExpected output:\n%s", output, expectedOutput)
	}

	client := goaviatrix.NewClient(os.Getenv("AVIATRIX_CONTROLLER_IP"), os.Getenv("AVIATRIX_USERNAME"), os.Getenv("AVIATRIX_PASSWORD"))
	foundVpc := &goaviatrix.Vpc{
		Name: fmt.Sprintf("tfg-%s", terraformOptions.Vars["r_name"].(string)),
	}

	err := client.GetVpc(foundVpc)
	if err != nil {
		t.Fatalf("Error getting VPC: %v", err)
	}

	if foundVpc.Name != fmt.Sprintf("tfg-%s", terraformOptions.Vars["r_name"].(string)) {
		t.Fatalf("VPC not found")
	}

	*vpc = *foundVpc
}

func testAccCheckVpcDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_vpc" {
			continue
		}

		foundVpc := &goaviatrix.Vpc{
			Name: rs.Primary.Attributes["name"],
		}

		_, err := client.GetVpc(foundVpc)
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("VPC still exists")
		}
	}

	return nil
}
