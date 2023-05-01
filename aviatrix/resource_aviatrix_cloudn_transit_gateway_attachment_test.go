package aviatrix

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAccAviatrixCloudnTransitGatewayAttachment_basic(t *testing.T) {
	skipAcc := os.Getenv("SKIP_CLOUDN_TRANSIT_GATEWAY_ATTACHMENT")
	if skipAcc == "yes" {
		t.Skip("Skipping transit gateway and cloudn attachment test as SKIP_CLOUDN_TRANSIT_GATEWAY_ATTACHMENT is set")
	}

	testName := fmt.Sprintf("terratest-%s", strings.ToLower(random.UniqueId()))
	terraformOptions := &terraform.Options{
		TerraformDir: "./",
		Vars: map[string]interface{}{
			"device_name":                           os.Getenv("CLOUDN_DEVICE_NAME"),
			"transit_gateway_name":                  os.Getenv("TRANSIT_GATEWAY_NAME"),
			"connection_name":                       fmt.Sprintf("connection-%s", testName),
			"transit_gateway_bgp_asn":               "65707",
			"cloudn_bgp_asn":                        os.Getenv("CLOUDN_BGP_ASN"),
			"cloudn_lan_interface_neighbor_ip":      os.Getenv("CLOUDN_LAN_INTERFACE_NEIGHBOR_IP"),
			"cloudn_lan_interface_neighbor_bgp_asn": os.Getenv("CLOUDN_LAN_INTERFACE_NEIGHBOR_BGP_ASN"),
			"enable_over_private_network":           true,
			"enable_jumbo_frame":                    false,
			"enable_dead_peer_detection":            true,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	resourceName := "aviatrix_cloudn_transit_gateway_attachment.test_cloudn_transit_gateway_attachment"

	err := testAccCheckCloudnTransitGatewayAttachmentExists(resourceName)
	assert.Nil(t, err)

	// Import the resource and verify that it matches the current state
	importedResource := terraform.ImportState(t, terraformOptions, resourceName)

	assert.Equal(t, resourceName, importedResource.Primary.ID)
	assert.Equal(t, "connection-"+testName, importedResource.Primary.Attributes["connection_name"])
}

func testAccCheckCloudnTransitGatewayAttachmentExists(resourceName string) error {
	client := goaviatrix.NewClientFromEnvParams()

	rs, err := terraform.GetResource(resourceName)
	if err != nil {
		return err
	}

	if rs.Primary.ID == "" {
		return fmt.Errorf("no cloudn_transit_gateway_attachment ID is set")
	}

	attachment := &goaviatrix.CloudnTransitGatewayAttachment{
		ConnectionName: rs.Primary.Attributes["connection_name"],
	}

	_, err = client.GetCloudnTransitGatewayAttachment(context.Background(), attachment.ConnectionName)
	if err != nil {
		return err
	}
	if attachment.ConnectionName != rs.Primary.ID {
		return fmt.Errorf("cloudn_transit_gateway_attachment not found")
	}

	return nil
}

func testCloudnTransitGatewayAttachmentDestroy(t *testing.T, terraformOptions *terraform.Options) {
	client := goaviatrix.NewClientFromEnvParams()

	err := terraform.DestroyE(t, terraformOptions)
	require.NoError(t, err)

	for _, rs := range terraformOptions.State().RootModule().Resources {
		if rs.Type != "aviatrix_cloudn_transit_gateway_attachment" {
			continue
		}

		attachment := &goaviatrix.CloudnTransitGatewayAttachment{
			ConnectionName: rs.Primary.Attributes["connection_name"],
		}
		_, err := client.GetCloudnTransitGatewayAttachment(context.Background(), attachment.ConnectionName)
		require.Error(t, err)
	}
}

func testAviatrixCloudnTransitGatewayAttachmentPreCheck(t *testing.T) {
	required := []string{
		"TRANSIT_GATEWAY_NAME",
		"CLOUDN_DEVICE_NAME",
		"CLOUDN_BGP_ASN",
		"CLOUDN_LAN_INTERFACE_NEIGHBOR_IP",
		"CLOUDN_LAN_INTERFACE_NEIGHBOR_BGP_ASN",
	}
	for _, v := range required {
		if os.Getenv(v) == "" {
			t.Fatalf("%s must be set for aviatrix_cloudn_transit_gateway_attachment acceptance test.", v)
		}
	}
}