package aviatrix_test

import (
    "fmt"
    "os"
    "testing"

    "github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
    "github.com/gruntwork-io/terratest/modules/random"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"

    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
    "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccTerraformAviatrixVPNCertDownload(t *testing.T) {
    t.Parallel()

    skipVPNCertDownload := os.Getenv("SKIP_VPN_CERT_DOWNLOAD")
    if skipVPNCertDownload == "true" {
        t.Skip("Skipping test as SKIP_VPN_CERT_DOWNLOAD is set to true")
    }

    resourceGroupName := fmt.Sprintf("aviatrix-vpn-cert-download-test-%s", random.UniqueId())
    samlEndpointName := fmt.Sprintf("aviatrix-saml-endpoint-test-%s", random.UniqueId())
    vpnUserName := fmt.Sprintf("aviatrix-vpn-user-test-%s", random.UniqueId())

    terraformOptions := &terraform.Options{
        TerraformDir: "./",
        Vars: map[string]interface{}{
            "resource_group_name": resourceGroupName,
            "saml_endpoint_name":  samlEndpointName,
            "vpn_user_name":       vpnUserName,
        },
    }

    // Make sure to clean up resources after the test completes
    defer terraform.Destroy(t, terraformOptions)

    // Deploy the Terraform code
    terraform.InitAndApply(t, terraformOptions)

    // Check if the VPN cert download resource exists
    vpnCertDownloadResource := "aviatrix_vpn_cert_download.test_vpn_cert_download"
    resource.Test(t, resource.TestCase{
        PreCheck:          func() { testAccPreCheck(t) },
        Providers:         testAccProviders,
        CheckDestroy:      testAccCheckVPNCertDownloadDestroy,
        Steps: []resource.TestStep{
            {
                Config: testAccVPNCertDownloadConfigBasic(resourceGroupName),
                Check: resource.ComposeTestCheckFunc(
                    testAccCheckVPNCertDownloadExists(vpnCertDownloadResource),
                ),
            },
        },
    })
}

func testAccCheckVPNCertDownloadExists(n string) resource.TestCheckFunc {
    return func(s *terraform.State) error {
        rs, ok := s.RootModule().Resources[n]
        if !ok {
            return fmt.Errorf("aviatrix VPN Cert Download Resource is not created: %s", n)
        }
        if rs.Primary.ID == "" {
            return fmt.Errorf("no aviatrix VPN Cert Download Resource ID is set")
        }

        // Initialize the Aviatrix client
        client := initializeAviatrixClient()

        vpnCertDownloadStatus, err := client.GetVPNCertDownloadStatus()
        if err != nil {
            return err
        }
        if !vpnCertDownloadStatus.Results.Status {
            return fmt.Errorf("VPN Cert Download doesn't seem to be enabled")
        }

        return nil
    }
}
func testAccCheckVPNCertDownloadDestroy(s *terraform.State) error {
    // Initialize the Aviatrix client
    client := initializeAviatrixClient()

    for _, rs := range s.RootModule().Resources {
        if rs.Type != "aviatrix_vpn_cert_download" {
            continue
        }

        vpnCertDownloadStatus, err := client.GetVPNCertDownloadStatus()
        if err != nil {
            return err
        }
        if vpnCertDownloadStatus.Results.Status {
            return fmt.Errorf("VPN Cert Download is not disabled")
        }
    }
    return nil
}

// Initialize the Aviatrix client using the credentials specified in the Terraform variables
func initializeAviatrixClient() *goaviatrix.Client {
    apiKey := os.Getenv("AVIATRIX_API_KEY")
    apiSecret := os.Getenv("AVIATRIX_API_SECRET")
    controllerIP := os.Getenv("AVIATRIX_CONTROLLER_IP")
    return goaviatrix.NewClient(apiKey, apiSecret, controllerIP)
}
