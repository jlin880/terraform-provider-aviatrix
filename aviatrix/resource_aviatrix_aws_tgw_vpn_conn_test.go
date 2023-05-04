package test

import (
    "context"
    "fmt"
    "os"
    "testing"
    "time"

    "github.com/gruntwork-io/terratest/modules/random"
    "github.com/gruntwork-io/terratest/modules/terraform"
)

func TestTerraformAwsTgwVpnConn(t *testing.T) {
    terraformOptions := &terraform.Options{
        TerraformDir: "./",
    }

    skipAcc := os.Getenv("SKIP_AWS_TGW_VPN_CONN")
    if skipAcc == "yes" {
        t.Skip("Skipping AWS TGW VPN CONN test as SKIP_AWS_TGW_VPN_CONN is set")
    }

    // Clean up resources after test is done
    defer terraform.Destroy(t, terraformOptions)

    // Provision resources using Terraform
    terraform.InitAndApply(t, terraformOptions)

    // Test the output
    vpnConnectionName := fmt.Sprintf("tfc-%s", random.UniqueId())
    awsTgwVpnConnConfig := fmt.Sprintf(`
        resource "aviatrix_aws_tgw_vpn_conn" "test" {
            tgw_name          = aviatrix_aws_tgw.test_aws_tgw.tgw_name
            route_domain_name = aviatrix_aws_tgw_network_domain.Default_Domain.name
            connection_name   = "%s"
            public_ip         = "40.0.0.0"
            remote_as_number  = "12"
        }
    `, vpnConnectionName)

    // Create the resource and verify it exists
    terraformOptions.Vars["aws_tgw_vpn_conn_config"] = awsTgwVpnConnConfig
    terraform.InitAndApply(t, terraformOptions)
    terraform.Output(t, terraformOptions, "aws_tgw_vpn_conn_id")

    // Import the resource
    vpnConnectionId := terraform.Output(t, terraformOptions, "aws_tgw_vpn_conn_id")
    importState := fmt.Sprintf("aviatrix_aws_tgw_vpn_conn.test %s", vpnConnectionId)
    terraform.Import(t, terraformOptions, importState)

    // Verify the imported resource
    awsTgwVpnConnResource := terraform.Output(t, terraformOptions, "aviatrix_aws_tgw_vpn_conn.test")
    expectedAwsTgwVpnConnResource := fmt.Sprintf(`{
        "connection_name": "%s",
        "tgw_name": "tft-%s",
        "route_domain_name": "Default_Domain",
        "public_ip": "40.0.0.0",
        "remote_as_number": "12"
    }`, vpnConnectionName, terraformOptions.Vars["name"].(string))

    if awsTgwVpnConnResource != expectedAwsTgwVpnConnResource {
        t.Errorf("expected %s but got %s", expectedAwsTgwVpnConnResource, awsTgwVpnConnResource)
    }
}


func testAccAwsTgwVpnConnConfigBasic(rName string, awsSideAsNumber string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_account" {
	account_name       = "tfa-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam	           = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}
resource "aviatrix_aws_tgw" "test_aws_tgw" {
	account_name       = aviatrix_account.test_account.account_name
	aws_side_as_number = "64512"
	region             = "%s"
	tgw_name           = "tft-%s"
}
resource "aviatrix_aws_tgw_network_domain" "Default_Domain" {
	name     = "Default_Domain"
	tgw_name = aviatrix_aws_tgw.test.tgw_name
}
resource "aviatrix_aws_tgw_network_domain" "Shared_Service_Domain" {
	name     = "Shared_Service_Domain"
	tgw_name = aviatrix_aws_tgw.test.tgw_name
}
resource "aviatrix_aws_tgw_network_domain" "Aviatrix_Edge_Domain" {
	name     = "Aviatrix_Edge_Domain"
	tgw_name = aviatrix_aws_tgw.test.tgw_name
}
resource "aviatrix_aws_tgw_peering_domain_conn" "default_nd_conn1" {
	tgw_name1    = aviatrix_aws_tgw.test.tgw_name
	domain_name1 = aviatrix_aws_tgw_network_domain.Aviatrix_Edge_Domain.name
	tgw_name2    = aviatrix_aws_tgw.test.tgw_name
	domain_name2 = aviatrix_aws_tgw_network_domain.Default_Domain.name
}
resource "aviatrix_aws_tgw_peering_domain_conn" "default_nd_conn2" {
	tgw_name1    = aviatrix_aws_tgw.test.tgw_name
	domain_name1 = aviatrix_aws_tgw_network_domain.Aviatrix_Edge_Domain.name
	tgw_name2    = aviatrix_aws_tgw.test.tgw_name
	domain_name2 = aviatrix_aws_tgw_network_domain.Shared_Service_Domain.name
}
resource "aviatrix_aws_tgw_peering_domain_conn" "default_nd_conn3" {
	tgw_name1    = aviatrix_aws_tgw.test.tgw_name
	domain_name1 = aviatrix_aws_tgw_network_domain.Default_Domain.name
	tgw_name2    = aviatrix_aws_tgw.test.tgw_name
	domain_name2 = aviatrix_aws_tgw_network_domain.Shared_Service_Domain.name
}
resource "aviatrix_aws_tgw_vpn_conn" "test" {
	tgw_name          = aviatrix_aws_tgw.test_aws_tgw.tgw_name
	route_domain_name = aviatrix_aws_tgw_network_domain.Default_Domain.name
	connection_name   = "tfc-%s"
	public_ip         = "40.0.0.0"
	remote_as_number  = "%s"
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		os.Getenv("AWS_REGION"), rName, rName, awsSideAsNumber)
}

func tesAccCheckAwsTgwVpnConnExists(n string, awsTgwVpnConn *goaviatrix.AwsTgwVpnConn) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("AWS TGW VPN CONN Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no AWS TGW VPN CONN ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		foundAwsTgwVpnConn := &goaviatrix.AwsTgwVpnConn{
			TgwName: rs.Primary.Attributes["tgw_name"],
			VpnID:   rs.Primary.Attributes["vpn_id"],
		}

		foundAwsTgwVpnConn2, err := client.GetAwsTgwVpnConn(foundAwsTgwVpnConn)
		if err != nil {
			return err
		}
		if foundAwsTgwVpnConn2.TgwName != rs.Primary.Attributes["tgw_name"] {
			return fmt.Errorf("tgw_name Not found in created attributes")
		}
		if foundAwsTgwVpnConn2.ConnName != rs.Primary.Attributes["connection_name"] {
			return fmt.Errorf("connection_name Not found in created attributes")
		}

		*awsTgwVpnConn = *foundAwsTgwVpnConn
		return nil
	}
}

func testAccCheckAwsTgwVpnConnDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_aws_tgw_vpn_conn" {
			continue
		}

		foundAwsTgwVpnConn := &goaviatrix.AwsTgwVpnConn{
			TgwName: rs.Primary.Attributes["tgw_name"],
			VpnID:   rs.Primary.Attributes["vpn_id"],
		}

		_, err := client.GetAwsTgwVpnConn(foundAwsTgwVpnConn)
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("aviatrix AWS TGW VPN CONN still exists")
		}
	}

	return nil
}
