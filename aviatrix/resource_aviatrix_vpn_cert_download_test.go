package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTerraformAviatrixVPNCertDownload(t *testing.T) {
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

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	// Check if VPN cert download resource exists
	vpnCertDownloadResource := "aviatrix_vpn_cert_download.test_vpn_cert_download"
	assert.True(t, terraform.OutputExists(t, terraformOptions, "aviatrix_vpn_cert_download.test_vpn_cert_download_id"))
	assert.True(t, terraform.OutputExists(t, terraformOptions, "aviatrix_vpn_cert_download.test_vpn_cert_download_name"))

	// Check if VPN cert download resource is enabled
	client := goaviatrix.NewClient(os.Getenv("AVIATRIX_CONTROLLER_IP"), os.Getenv("AVIATRIX_USERNAME"), os.Getenv("AVIATRIX_PASSWORD"))
	vpnCertDownloadStatus, err := client.GetVPNCertDownloadStatus()
	require.NoError(t, err)
	assert.True(t, vpnCertDownloadStatus.Results.Status)

	// Check if SAML endpoint is associated with VPN cert download resource
	vpnCertDownloadSamlEndpoints, err := client.GetVPNCertDownloadSamlEndpoints(vpnCertDownloadStatus.Results.ActionStatus, vpnCertDownloadStatus.Results.Info)
	require.NoError(t, err)
	assert.Equal(t, 1, len(vpnCertDownloadSamlEndpoints))
	assert.Equal(t, samlEndpointName, vpnCertDownloadSamlEndpoints[0])

	// Check if VPN user is associated with SAML endpoint
	samlUser, err := client.GetSamlUser(samlEndpointName, vpnUserName)
	require.NoError(t, err)
	assert.NotNil(t, samlUser)
}

func testAccVPNCertDownloadConfigBasic(rName string) string {
	idpMetadata := os.Getenv("IDP_METADATA")
	idpMetadataType := os.Getenv("IDP_METADATA_TYPE")
	vpnUserConfig := testAccVPNUserConfigBasic(rName, "true", rName)
	samlConfig := testAccSamlEndpointConfigBasic(rName, idpMetadata, idpMetadataType)
	return vpnUserConfig + samlConfig + `
resource "aviatrix_vpn_cert_download" "test_vpn_cert_download" {
    download_enabled = true
    saml_endpoints = [aviatrix_saml_endpoint.foo.endpoint_name]
	depends_on = [
    aviatrix_vpn_user.test_vpn_user, 
    aviatrix_saml_endpoint.foo
  ]
}
`
}

func testAccCheckVPNCertDownloadExists(n string) resource.TestCheckFunc {
    return func(s *terraform.State) error {
        rs, ok := s.RootModule().Resources[n]
        if !ok {
            return fmt.Errorf("aviatrix VPN Cert Download Resource is Not Created: %s", n)
        }
        if rs.Primary.ID == "" {
            return fmt.Errorf("no aviatrix VPN Cert Download Resource ID is set")
        }

        // Here, you would need to initialize an Aviatrix client using the credentials specified in the Terraform
        // variables, which you would pass into this test function as a parameter. I'm using a placeholder value here
        // to illustrate the code structure.
        client := initializeAviatrixClient()

        vpnCertDownloadStatus, err := client.GetVPNCertDownloadStatus()
        if err != nil {
            return err
        }
        if !vpnCertDownloadStatus.Results.Status {
            return fmt.Errorf("VPN Cert Download doesnt seem to be enabled")
        }
        return nil
    }
}

func testAccCheckVPNCertDownloadDestroy(s *terraform.State) error {
    // Here, you would need to initialize an Aviatrix client using the credentials specified in the Terraform
    // variables, which you would pass into this test function as a parameter. I'm using a placeholder value here
    // to illustrate the code structure.
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
            return fmt.Errorf("VPN Cert Download doesnt seem to be disabled")
        }
    }
    return nil
}

// Here's an example function that initializes an Aviatrix client. You would need to update this function to use
// the credentials specified in the Terraform variables.
func initializeAviatrixClient() *goaviatrix.Client {
    return goaviatrix.NewClient("AVIATRIX_API_KEY", "AVIATRIX_API_SECRET", "AVIATRIX_CONTROLLER_IP")
}