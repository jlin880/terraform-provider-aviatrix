go
Copy code
package test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
)

func preGateway2Check(t *testing.T, msgCommon string) {
	preAccountCheck(t, msgCommon)

	vpcID2 := os.Getenv("AWS_VPC_ID2")
	if vpcID2 == "" {
		t.Fatal("Environment variable AWS_VPC_ID2 is not set" + msgCommon)
	}

	region2 := os.Getenv("AWS_REGION2")
	if region2 == "" {
		t.Fatal("Environment variable AWS_REGION2 is not set" + msgCommon)
	}

	vpcNet2 := os.Getenv("AWS_SUBNET2")
	if vpcNet2 == "" {
		t.Fatal("Environment variable AWS_SUBNET2 is not set" + msgCommon)
	}
}

func TestAccAviatrixTunnel_basic(t *testing.T) {
    terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
        TerraformDir: "../path/to/terraform/code",
        Vars: map[string]interface{}{
            "account_name":       fmt.Sprintf("tfa-%s", random.UniqueId()),
            "cloud_type":         1,
            "aws_account_number": os.Getenv("AWS_ACCOUNT_NUMBER"),
            "aws_iam":            false,
            "aws_access_key":     os.Getenv("AWS_ACCESS_KEY"),
            "aws_secret_key":     os.Getenv("AWS_SECRET_KEY"),
            "gw_name1":           fmt.Sprintf("tfg-%s", random.UniqueId()),
            "vpc_id1":            os.Getenv("AWS_VPC_ID"),
            "vpc_reg1":           os.Getenv("AWS_REGION"),
            "gw_size1":           "t2.micro",
            "subnet1":            os.Getenv("AWS_SUBNET"),
            "gw_name2":           fmt.Sprintf("tfg2-%s", random.UniqueId()),
            "vpc_id2":            os.Getenv("AWS_VPC_ID2"),
            "vpc_reg2":           os.Getenv("AWS_REGION2"),
            "gw_size2":           "t2.micro",
            "subnet2":            os.Getenv("AWS_SUBNET2"),
        },
    })

    defer terraform.Destroy(t, terraformOptions)

    terraform.InitAndApply(t, terraformOptions)

    aviatrixClient, err := goaviatrix.NewClient()
    assert.NoError(t, err)

    tunnel := &goaviatrix.Tunnel{
        VpcName1: terraformOptions.Vars["gw_name1"].(string),
        VpcName2: terraformOptions.Vars["gw_name2"].(string),
    }

    err = aviatrixClient.GetTunnel(tunnel)
    assert.NoError(t, err)
}
func testAccTunnelConfigBasic(rName string, vpcID1 string, vpcID2 string, region1 string, region2 string,
	subnet1 string, subnet2 string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test" {
	account_name       = "tfa-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}
resource "aviatrix_gateway" "gw1" {
	cloud_type   = 1
	account_name = aviatrix_account.test.account_name
	gw_name      = "tfg-%[1]s"
	vpc_id       = "%[5]s"
	vpc_reg      = "%[7]s"
	gw_size      = "t2.micro"
	subnet       = "%[9]s"
}
resource "aviatrix_gateway" "gw2" {
	cloud_type   = 1
	account_name = aviatrix_account.test.account_name
	gw_name      = "tfg2-%[1]s"
	vpc_id       = "%[6]s"
	vpc_reg      = "%[8]s"
	gw_size      = "t2.micro"
	subnet       = "%[10]s"
}
resource "aviatrix_tunnel" "foo" {
	gw_name1 = aviatrix_gateway.gw1.gw_name
	gw_name2 = aviatrix_gateway.gw2.gw_name
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		vpcID1, vpcID2, region1, region2, subnet1, subnet2)
}

func tesAccCheckTunnelExists(n string, tunnel *goaviatrix.Tunnel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("aviatrix tunnel Not Created: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no aviatrix tunnel ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		foundTunnel := &goaviatrix.Tunnel{
			VpcName1: rs.Primary.Attributes["gw_name1"],
			VpcName2: rs.Primary.Attributes["gw_name2"],
		}

		_, err := client.GetTunnel(foundTunnel)
		if err != nil {
			return err
		}
		if foundTunnel.VpcName1 != rs.Primary.Attributes["gw_name1"] {
			return fmt.Errorf("gw_name1 Not found in created attributes")
		}
		if foundTunnel.VpcName2 != rs.Primary.Attributes["gw_name2"] {
			return fmt.Errorf("gw_name2 Not found in created attributes")
		}

		*tunnel = *foundTunnel
		return nil
	}
}

func testAccCheckTunnelDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_tunnel" {
			continue
		}

		foundTunnel := &goaviatrix.Tunnel{
			VpcName1: rs.Primary.Attributes["gw_name1"],
			VpcName2: rs.Primary.Attributes["gw_name2"],
		}

		_, err := client.GetTunnel(foundTunnel)
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("aviatrix tunnel still exists")
		}
	}

	return nil
}
