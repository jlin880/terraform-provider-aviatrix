package test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/acctest"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/goaviatrix"
)

func prePrivateModeCheck(t *testing.T, msgEnd string) {
	for _, key := range []string{"CONTROLLER_VPC_ID", "AWS_REGION"} {
		if os.Getenv(key) == "" {
			t.Fatalf("%s must be set for Private Mode tests using load balancers. %s", key, msgEnd)
		}
	}
}

func TestAccAviatrixPrivateModeLb_basic(t *testing.T) {
	rName := acctest.RandString(5)

	skipAcc := os.Getenv("SKIP_PRIVATE_MODE_LB")
	if skipAcc == "yes" {
		t.Skip("Skipping Controller Private Mode load balancer tests as SKIP_PRIVATE_MODE_LB is set")
	}
	msgCommon := "Set SKIP_PRIVATE_MODE_LB to yes to skip Controller Private Mode load balancer tests"
	resourceName := "aviatrix_private_mode_lb.test"

	terraformOptions := &terraform.Options{
		TerraformDir: ".",
		Vars: map[string]interface{}{
			"rname":             rName,
			"aws_account_num":   os.Getenv("AWS_ACCOUNT_NUMBER"),
			"aws_access_key":    os.Getenv("AWS_ACCESS_KEY"),
			"aws_secret_key":    os.Getenv("AWS_SECRET_KEY"),
			"controller_vpc_id": os.Getenv("CONTROLLER_VPC_ID"),
			"region":            os.Getenv("AWS_REGION"),
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	assert.NotNil(t, terraformOptions, "Terraform options are nil")

	// Test whether the resource exists
	err := aviatrixPrivateModeLbExists(t, resourceName, terraformOptions)
	assert.NoError(t, err, "Error thrown while checking if resource exists")

	// Test that the correct attributes were set
	assert.Equal(t, "tfa-"+rName, terraform.Output(t, terraformOptions, "account_name"))
	assert.Equal(t, os.Getenv("AWS_REGION"), terraform.Output(t, terraformOptions, "region"))
	assert.Equal(t, "controller", terraform.Output(t, terraformOptions, "lb_type"))
}

func aviatrixPrivateModeLbExists(t *testing.T, resourceName string, terraformOptions *terraform.Options) error {
	resourceState := terraform.StateFromFile(t, terraformOptions.StatePath)

	client := aviatrix.NewClient(os.Getenv("AVIATRIX_API_ACCOUNT_NAME"), os.Getenv("AVIATRIX_API_KEY"), os.Getenv("AVIATRIX_CONTROLLER_IP"))

	vpcId := resourceState.RootModule().Outputs["this"].Value.(string)

	lb, err := client.GetPrivateModeLoadBalancer(context.Background(), vpcId)
	if err != nil {
		return err
	}

	assert.Equal(t, lb.VpcID, vpcId)
	assert.Equal(t, lb.AccountName, "tfa-"+terraformOptions.Vars["rname"].(string))
	assert.Equal(t, lb.Region, os.Getenv("AWS_REGION"))
	assert.Equal(t, lb.LBType, "controller")

	return nil
}
