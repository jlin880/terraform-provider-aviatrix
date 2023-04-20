package test

import (
    "context"
    "fmt"
    "os"
    "testing"

    "github.com/gruntwork-io/terratest/modules/terraform"
    "github.com/stretchr/testify/assert"

    goaviatrix "github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
)

func TestAccAviatrixPrivateModeMulticloudEndpoint_basic(t *testing.T) {
	rName := acctest.RandString(5)

	skipAcc := os.Getenv("SKIP_PRIVATE_MODE_MULTICLOUD_ENDPOINT")
	if skipAcc == "yes" {
		t.Skip("Skipping Controller Private Mode multicloud endpoint tests as SKIP_PRIVATE_MODE_MULTICLOUD_ENDPOINT is set")
	}
	msgCommon := "Set SKIP_PRIVATE_MODE_MULTICLOUD_ENDPOINT to yes to skip Controller Private Mode load balancer tests"
	resourceName := "aviatrix_private_mode_multicloud_endpoint.test"

	terraformOptions := &terraform.Options{
		TerraformDir: "./",
		Vars: map[string]interface{}{
			"account_name":         fmt.Sprintf("tfa-%s", rName),
			"aws_account_number":   os.Getenv("AWS_ACCOUNT_NUMBER"),
			"aws_iam":              false,
			"aws_access_key":       os.Getenv("AWS_ACCESS_KEY"),
			"aws_secret_key":       os.Getenv("AWS_SECRET_KEY"),
			"vpc_id":               os.Getenv("AWS_VPC_ID"),
			"region":               os.Getenv("AWS_REGION"),
			"controller_lb_vpc_id": os.Getenv("CONTROLLER_VPC_ID"),
			"lb_type":              "controller",
			"enable_private_mode":  true,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	assert.True(t, terraform.IsResourcePresent(t, terraformOptions, "aviatrix_account.test_account"))
	assert.True(t, terraform.IsResourcePresent(t, terraformOptions, "aviatrix_controller_private_mode_config.test"))
	assert.True(t, terraform.IsResourcePresent(t, terraformOptions, "aviatrix_private_mode_lb.test"))
	assert.True(t, terraform.IsResourcePresent(t, terraformOptions, resourceName))

	// Import the resource using the import ID
	importedResource := terraform.ImportState(t, terraformOptions, resourceName)
	assert.Equal(t, resourceName, importedResource)

	// Check the state using the import ID
	state := terraform.Show(t, terraformOptions)
	assert.Contains(t, state, resourceName)
}

func TestAccAviatrixPrivateModeMulticloudEndpoint_import(t *testing.T) {
	rName := acctest.RandString(5)

	skipAcc := os.Getenv("SKIP_PRIVATE_MODE_MULTICLOUD_ENDPOINT")
	if skipAcc == "yes" {
		t.Skip("Skipping Controller Private Mode multicloud endpoint tests as SKIP_PRIVATE_MODE_MULTICLOUD_ENDPOINT is set")
	}
	msgCommon := "Set SKIP_PRIVATE_MODE_MULTICLOUD_ENDPOINT to yes to skip Controller Private Mode load balancer tests"
	resourceName := "aviatrix_private_mode_multicloud_endpoint.test"

	terraformOptions := &terraform.Options{
		TerraformDir: "./",
		Vars: map[string]interface{}{
			"account_name":         fmt.Sprintf("tfa-%s", rName),
			"aws_account_number":   os.Getenv("AWS_ACCOUNT_NUMBER"),
			"aws_iam":              false,
			"aws_access_key":       os.Getenv("AWS_ACCESS_KEY"),
			"aws_secret_key":       os.Getenv("AWS_SECRET_KEY"),
			"vpc_id":               os.Getenv("AWS_VPC

func testAccAviatrixPrivateModeMulticloudEndpointBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_account" {
	account_name       = "tfa-%[1]s"
	cloud_type         = 1
	aws_account_number = "%[2]s"
	aws_iam            = false
	aws_access_key     = "%[3]s"
	aws_secret_key     = "%[4]s"
}

resource "aviatrix_controller_private_mode_config" "test" {
	enable_private_mode = true
}

resource "aviatrix_private_mode_lb" "test" {
	account_name = aviatrix_account.test_account.account_name
	vpc_id       = "%[5]s"
	region       = "%[6]s"
	lb_type      = "controller"

	depends_on = [aviatrix_controller_private_mode_config.test]
}

resource "aviatrix_private_mode_multicloud_endpoint" "test" {
	account_name         = "tfa-%[1]s"
	vpc_id               = "%[7]s"
	region               = "%[6]s"
	controller_lb_vpc_id = aviatrix_private_mode_lb.test.vpc_id
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"), os.Getenv("CONTROLLER_VPC_ID"), os.Getenv("AWS_REGION"), os.Getenv("AWS_VPC_ID"))
}

func testAccAviatrixPrivateModeMulticloudEndpointExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("private mode load multicloud endpoint Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no private mode multicloud endpoint ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		vpcId := rs.Primary.ID
		_, err := client.GetPrivateModeMulticloudEndpoint(context.Background(), vpcId)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccAviatrixPrivateModeMulticloudEndpointDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aviatrix_private_mode_lb" {
			vpcId := rs.Primary.ID
			_, err := client.GetPrivateModeLoadBalancer(context.Background(), vpcId)
			if err != nil {
				if err == goaviatrix.ErrNotFound {
					continue
				}
				return fmt.Errorf("failed to destroy Private Mode load balancer: %s", err)
			}
			return fmt.Errorf("failed to destroy Private Mode load balancer")
		} else if rs.Type == "aviatrix_private_mode_multicloud_load_balancer" {
			vpcId := rs.Primary.ID
			_, err := client.GetPrivateModeMulticloudEndpoint(context.Background(), vpcId)

			if err != nil {
				if err == goaviatrix.ErrNotFound {
					continue
				}
				return fmt.Errorf("error getting Private Mode multicloud endpoint after destroy: %s", err)
			}
			return fmt.Errorf("failed to destroy Private Mode multicloud endpoint")
		}
		if rs.Type != "aviatrix_private_mode_lb" {
			continue
		}

	}

	return nil
}
