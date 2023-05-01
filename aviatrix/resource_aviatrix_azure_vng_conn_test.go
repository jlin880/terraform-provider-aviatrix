package test

import (
    "fmt"
    "os"
    "testing"

    "github.com/gruntwork-io/terratest/modules/random"
    "github.com/gruntwork-io/terratest/modules/terraform"
    "github.com/stretchr/testify/assert"

    goaviatrix "github.com/AviatrixSystems/terraform-provider-aviatrix/goaviatrix"
)

func TestTerraformAzureVngConn(t *testing.T) {
    terraformOptions := &terraform.Options{
        TerraformDir: "../",
        Vars: map[string]interface{}{
            "subscription_id":    os.Getenv("ARM_SUBSCRIPTION_ID"),
            "client_id":          os.Getenv("ARM_APPLICATION_ID"),
            "client_secret":      os.Getenv("ARM_APPLICATION_KEY"),
            "tenant_id":          os.Getenv("ARM_DIRECTORY_ID"),
            "resource_group":     os.Getenv("AZURE_RESOURCE_GROUP"),
            "location":           os.Getenv("AZURE_REGION"),
            "vnet_name":          os.Getenv("AZURE_VNG_VNET_NAME"),
            "vnet_subnet_name":   os.Getenv("AZURE_VNG_SUBNET_NAME"),
            "vnet_subnet_prefix": os.Getenv("AZURE_VNG_SUBNET_PREFIX"),
            "vng_name":           os.Getenv("AZURE_VNG_NAME"),
            "connection_name":    "test-connection",
        },
    }

    defer terraform.Destroy(t, terraformOptions)

    terraform.InitAndApply(t, terraformOptions)

    output := terraform.Output(t, terraformOptions, "connection_status")

    assert.Equal(t, "connected", output)
}

func TestAccAviatrixAzureVngConn_basic(t *testing.T) {
    if os.Getenv("SKIP_AZURE_VNG_CONN") == "yes" {
        t.Skip("Skipping azure vng conn test as SKIP_AZURE_VNG_CONN is set")
    }

    connectionName := fmt.Sprintf("test-%s", random.UniqueId())
    terraformOptions := &terraform.Options{
        TerraformDir: ".",
        Vars: map[string]interface{}{
            "connection_name": connectionName,
        },
    }

    defer terraform.Destroy(t, terraformOptions)

    terraform.InitAndApply(t, terraformOptions)

    client, err := goaviatrix.NewClient(os.Getenv("AVIATRIX_API_USER"), os.Getenv("AVIATRIX_API_KEY"), "", "")
    if err != nil {
        t.Fatalf("Failed to create Aviatrix client: %v", err)
    }

    azureVngConn, err := client.GetAzureVngConn(connectionName)
    if err != nil {
        t.Fatalf("Failed to get Azure VNG conn: %v", err)
    }

    assert.Equal(t, "test-tgw-azure", azureVngConn.PrimaryGatewayName)
    assert.Equal(t, connectionName, azureVngConn.ConnectionName)
    assert.Equal(t, os.Getenv("AZURE_VNG_VNET_ID"), azureVngConn.VpcID)
    assert.Equal(t, os.Getenv("AZURE_VNG"), azureVngConn.VngName)
}

func TestAccCheckAzureVngConnDestroy(t *testing.T) {
    terraformOptions := &terraform.Options{
        TerraformDir: ".",
    }

    terraform.Destroy(t, terraformOptions)

    err := terraform.InitAndPlanE(t, terraformOptions)
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "No changes. Infrastructure is up-to-date.")

    client, err := goaviatrix.NewClient(os.Getenv("AVIATRIX_API_USER"), os.Getenv("AVIATRIX_API_KEY

func testGetAzureVngConn(t *testing.T, terraformOptions *terraform.Options, connectionName string) *goaviatrix.AzureVngConn {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	err := client.Login(os.Getenv("AVIATRIX_API_USER"), os.Getenv("AVIATRIX_API_KEY"))
	if err != nil {
		t.Fatalf("Failed to login to Aviatrix API: %v", err)
	}

	azureVngConn, err := client.GetAzureVngConn(connectionName)
	if err != nil {
		t.Fatalf("Failed to get Azure VNG conn: %v", err)
	}

	return azureVngConn
}
func testAccCheckAzureVngConnDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_azure_vng_conn" {
			continue
		}

		connectionName := rs.Primary.Attributes["connection_name"]

		resp, err := client.GetAzureVngConnStatus(connectionName)
		if err != nil {
			if err != goaviatrix.ErrNotFound {
				return fmt.Errorf("failed to retrieve azure vng conn status: %s", err)
			}
			// If the error is ErrNotFound, the resource was successfully destroyed.
			continue
		}

		if resp.Attached {
			return fmt.Errorf("azure vng conn %s still exists", connectionName)
		}
	}

	return nil
}

