package aviatrix_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/acctest"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAccDataSourceAviatrixVpc_basic(t *testing.T) {
	t.Parallel()

	awsRegion := os.Getenv("AWS_REGION")
	awsAccountID := os.Getenv("AWS_ACCOUNT_ID")

	// Skip test if environment variables are not set
	if awsRegion == "" || awsAccountID == "" {
		t.Skip("Skipping test due to missing AWS_REGION and/or AWS_ACCOUNT_ID environment variables")
	}

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "./",
		Vars: map[string]interface{}{
			"region":       awsRegion,
			"account_name": fmt.Sprintf("tfa-%s", random.UniqueId()),
			"name":         fmt.Sprintf("tfv-%s", random.UniqueId()),
			"cidr":         "10.0.0.0/16",
			"aws_account":  awsAccountID,
		},
	})

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	// Validate the data source
	data := terraform.OutputMap(t, terraformOptions, "vpc")
	assert.NotEmpty(t, data, "No VPC data returned")

	assert.Equal(t, awsRegion, data["region"])
	assert.Equal(t, "10.0.0.0/16", data["cidr"])
}

func TestAccDataSourceAviatrixTransitGateways_basic(t *testing.T) {
	t.Parallel()

	testName := fmt.Sprintf("aviatrix-transit-gateways-%s", random.UniqueId())
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/data-sources/transit-gateways",
		Vars: map[string]interface{}{
			"test_name":          testName,
			"aws_account_number": os.Getenv("AWS_ACCOUNT_NUMBER"),
			"aws_access_key":     os.Getenv("AWS_ACCESS_KEY"),
			"aws_secret_key":     os.Getenv("AWS_SECRET_KEY"),
			"aws_region":         os.Getenv("AWS_REGION"),

			"aws_subnet":                os.Getenv("AWS_SUBNET"),
			"gcp_project_id":            os.Getenv("GCP_ID"),
			"gcp_credentials_file_path": os.Getenv("GCP_CREDENTIALS_FILEPATH"),
			"gcp_subnet":                os.Getenv("GCP_SUBNET"),
			"gcp_zone":                  os.Getenv("GCP_ZONE"),
			"aviatrix_version":          "3.5",
		},
	})

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	terraform.Output(t, terraformOptions, "all_transit_gateways")
}


func TestAccDataSourceAviatrixTransitGateways(t *testing.T) {
			t.Parallel()

			rName := acctest.RandString(5)

			awsVpcId := os.Getenv("AWS_VPC_ID")
			awsRegion := os.Getenv("AWS_REGION")
			awsSubnet := os.Getenv("AWS_SUBNET")
			gcpProjectId := os.Getenv("GCP_ID")
			gcpZone := os.Getenv("GCP_ZONE")
			gcpSubnet := os.Getenv("GCP_SUBNET")

			"gcp_zone":                  os.Getenv("GCP_ZONE"),
			"aviatrix_version":          "3.5",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	terraform.Output(t, terraformOptions, "all_transit_gateways")
}
func TestAccDataSourceAviatrixTransitGateways(t *testing.T) {
	rName := acctest.RandString(5)
	resourceName := "data.aviatrix_transit_gateways.foo"

	skipAcc := os.Getenv("SKIP_DATA_TRANSIT_GATEWAYS")
	if skipAcc == "yes" {
		t.Skip("Skipping Data Source All Transit Gateway tests as SKIP_DATA_TRANSIT_GATEWAYS is set")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "./",
		Vars: map[string]interface{}{
			"aviatrix_account_name":     fmt.Sprintf("aa-tfa-%s", rName),
			"aws_account_number":        os.Getenv("AWS_ACCOUNT_NUMBER"),
			"aws_access_key":            os.Getenv("AWS_ACCESS_KEY"),
			"aws_secret_key":            os.Getenv("AWS_SECRET_KEY"),
			"aws_vpc_id":                os.Getenv("AWS_VPC_ID"),
			"aws_region":                os.Getenv("AWS_REGION"),
			"aws_subnet":                os.Getenv("AWS_SUBNET"),
			"gcp_project_id":            os.Getenv("GCP_ID"),
			"gcp_credentials_file_path": os.Getenv("GCP_CREDENTIALS_FILEPATH"),
			"gcp_vpc_id":                os.Getenv("GCP_VPC_ID"),
			"gcp_zone":                  os.Getenv("GCP_ZONE"),
			"gcp_subnet":                os.Getenv("GCP_SUBNET"),
			"aviatrix_gw_name":          fmt.Sprintf("aa-tfg-aws-%s", rName),
			"aviatrix_gw_size":          "t2.micro",
			"aviatrix_gcp_gw_name":      fmt.Sprintf("aa-tfg-gcp-%s", rName),
			"aviatrix_gcp_gw_size":      "n1-standard-1",
		},
	}

	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApply(t, terraformOptions)

	// Check that the data source exists
	terraform.OutputRequired(t, terraformOptions, "data_source_exists")

	// Check the attributes of the data source
	terraform.OutputRequired(t, terraformOptions, "gateway_list.0.gw_name")
	terraform.OutputRequired(t, terraformOptions, "gateway_list.0.vpc_id")
	terraform.OutputRequired(t, terraformOptions, "gateway_list.0.vpc_reg")
	terraform.OutputRequired(t, terraformOptions, "gateway_list.0.gw_size")
	terraform.OutputRequired(t, terraformOptions, "gateway_list.1.gw_name")
	terraform.OutputRequired(t, terraformOptions, "gateway_list.1.gw_size")
	terraform.OutputRequired(t, terraformOptions, "gateway_list.1.account_name")
	terraform.OutputRequired(t, terraformOptions, "gateway_list.1.subnet")
	terraform.OutputRequired(t, terraformOptions, "gateway_list.1.vpc_reg")
}
