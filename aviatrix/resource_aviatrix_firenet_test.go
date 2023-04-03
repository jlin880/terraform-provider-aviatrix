package test

import (
    "context"
    "fmt"
    "os"
    "reflect"
    "testing"

    "github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
    "github.com/gruntwork-io/terratest/modules/random"
    "github.com/gruntwork-io/terratest/modules/terraform"
)

func TestAccAviatrixFireNet(t *testing.T) {
    t.Parallel()

    terraformOptions := createTerraformOptions(t)

    defer terraform.Destroy(t, terraformOptions)

    terraform.InitAndApply(t, terraformOptions)

    fireNetID := terraform.Output(t, terraformOptions, "firenet_id")

    client, err := goaviatrix.NewClient(os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"), os.Getenv("AWS_REGION"), true)
    if err != nil {
        t.Fatalf("Failed to create Aviatrix client: %s", err)
    }

    var fireNet goaviatrix.FireNet

    err = terraform.ReadStateFromFile(t, terraformOptions.StatePath, &fireNet)
    if err != nil {
        t.Fatalf("Failed to read state file: %s", err)
    }

    foundFireNet, err := client.GetFireNet(&goaviatrix.FireNet{
        VpcID: fireNet.VpcID,
    })

    if err != nil {
        t.Fatalf("Failed to get firenet %s: %s", fireNetID, err)
    }

    if foundFireNet.VpcID != fireNetID {
        t.Errorf("FireNet not found")
    }
}

func createTerraformOptions(t *testing.T) *terraform.Options {
    uniqueID := random.UniqueId()

    terraformDir := "../path/to/aviatrix/module"

    return terraform.WithDefaultRetryableErrors(t, &terraform.Options{
        TerraformDir: terraformDir,
        Vars: map[string]interface{}{
            "aws_access_key":     os.Getenv("AWS_ACCESS_KEY"),
            "aws_secret_key":     os.Getenv("AWS_SECRET_KEY"),
            "aws_region":         os.Getenv("AWS_REGION"),
            "aws_account_number": os.Getenv("AWS_ACCOUNT_NUMBER"),
            "firenet_name":       fmt.Sprintf("terratest-firenet-%s", uniqueID),
            "vpc_cidr":           "10.0.0.0/16",
            "firewall_size":      "m5.xlarge",
            "firewall_image":     "Palo Alto Networks VM-Series Next-Generation Firewall Bundle 1",
        },
        StatePath: "terraform.tfstate",
    })
}

func TestAccFireNetConfigBasic(t *testing.T) {
    rName := random.UniqueId()

    terraformOptions := createTerraformOptions(t)
    terraformOptions.Vars["config"] = testAccFireNetConfigBasic(rName)

    defer terraform.Destroy(t, terraformOptions)

    terraform.InitAndApply(t, terraformOptions)

    fireNetID := terraform.Output(t, terraformOptions, "firenet_id")

    client, err := goaviatrix.NewClient(os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"), os.Getenv("AWS_REGION"), true)
    if err != nil {
        t.Fatalf("Failed to create Aviatrix client: %s", err)
    }

    var fireNet goaviatrix.FireNet

    err = terraform.ReadStateFromFile(t, terraformOptions.StatePath, &fireNet)
    if err != nil {
        t.Fatalf("Failed to read state file:

func testAccFireNetConfigBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_account" {
	account_name       = "tfa-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}
resource "aviatrix_vpc" "test_vpc" {
	cloud_type           = 1
	account_name         = aviatrix_account.test_account.account_name
	region               = "%s"
	name                 = "vpc-for-firenet"
	cidr                 = "10.10.0.0/24"
	aviatrix_firenet_vpc = true
}
resource "aviatrix_transit_gateway" "test_transit_gateway" {
	cloud_type               = aviatrix_vpc.test_vpc.cloud_type
	account_name             = aviatrix_account.test_account.account_name
	gw_name                  = "tftg-%s"
	vpc_id                   = aviatrix_vpc.test_vpc.vpc_id
	vpc_reg                  = aviatrix_vpc.test_vpc.region
	gw_size                  = "c5.xlarge"
	subnet                   = aviatrix_vpc.test_vpc.subnets[0].cidr
	enable_hybrid_connection = true
	enable_firenet           = true
}
resource "aviatrix_firewall_instance" "test_firewall_instance" {
	vpc_id            = aviatrix_vpc.test_vpc.vpc_id
	firenet_gw_name   = aviatrix_transit_gateway.test_transit_gateway.gw_name
	firewall_name     = "tffw-%s"
	firewall_image    = "Palo Alto Networks VM-Series Next-Generation Firewall Bundle 1"
	firewall_size     = "m5.xlarge"
	management_subnet = aviatrix_vpc.test_vpc.subnets[0].cidr
	egress_subnet     = aviatrix_vpc.test_vpc.subnets[1].cidr
}
resource "aviatrix_firewall_instance_association" "firewall_instance_association" {
	vpc_id               = aviatrix_firewall_instance.test_firewall_instance.vpc_id
	firenet_gw_name      = aviatrix_transit_gateway.test_transit_gateway.gw_name
	instance_id          = aviatrix_firewall_instance.test_firewall_instance.instance_id
	firewall_name        = aviatrix_firewall_instance.test_firewall_instance.firewall_name
	lan_interface        = aviatrix_firewall_instance.test_firewall_instance.lan_interface
	management_interface = aviatrix_firewall_instance.test_firewall_instance.management_interface
	egress_interface     = aviatrix_firewall_instance.test_firewall_instance.egress_interface
	attached             = true
}
resource "aviatrix_firenet" "test_firenet" {
	vpc_id             = aviatrix_vpc.test_vpc.vpc_id
	inspection_enabled = true
	egress_enabled     = false

	depends_on = [aviatrix_firewall_instance_association.firewall_instance_association]
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		os.Getenv("AWS_REGION"), rName, rName)
}

func testAccCheckFireNetExists(n string, fireNet *goaviatrix.FireNet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("fireNet Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no FireNet ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		foundFireNet := &goaviatrix.FireNet{
			VpcID: rs.Primary.Attributes["vpc_id"],
		}

		foundFireNet2, err := client.GetFireNet(foundFireNet)
		if err != nil {
			if err == goaviatrix.ErrNotFound {
				return fmt.Errorf("fireNet not found")
			}
			return err
		}
		if foundFireNet2.VpcID != rs.Primary.ID {
			return fmt.Errorf("fireNet not found")
		}

		*fireNet = *foundFireNet
		return nil
	}
}

func testAccCheckFireNetDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_firenet" {
			continue
		}

		foundFireNet := &goaviatrix.FireNet{
			VpcID: rs.Primary.Attributes["vpc_id"],
		}

		_, err := client.GetFireNet(foundFireNet)
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("fireNet still exists")
		}
	}

	return nil
}

func testResourceFireNetStateDataV0() map[string]interface{} {
	return map[string]interface{}{}
}

func testResourceFireNetStateDataV1() map[string]interface{} {
	return map[string]interface{}{
		"manage_firewall_instance_association": true,
	}
}

func TestResourceFireNetStateUpgradeV0(t *testing.T) {
	expected := testResourceFireNetStateDataV1()
	actual, err := resourceAviatrixFireNetStateUpgradeV0(context.Background(), testResourceFireNetStateDataV0(), nil)
	if err != nil {
		t.Fatalf("error migrating state: %s", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected:%#v\ngot:%#v\n", expected, actual)
	}
}
