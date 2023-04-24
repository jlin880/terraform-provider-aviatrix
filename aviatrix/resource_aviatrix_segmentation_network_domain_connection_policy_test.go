package aviatrix

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAccAviatrixSegmentationNetworkDomainConnectionPolicy_basic(t *testing.T) {
	t.Parallel()

	if os.Getenv("SKIP_SEGMENTATION_NETWORK_DOMAIN_CONNECTION_POLICY") == "yes" {
		t.Skip("Skipping segmentation network domain conn policy test as SKIP_SEGMENTATION_NETWORK_DOMAIN_CONNECTION_POLICY is set")
	}

	ctx := context.Background()
	rName := random.UniqueId()

	segmentationNetworkDomain1 := fmt.Sprintf("segmentation-sd-1-%s", rName)
	segmentationNetworkDomain2 := fmt.Sprintf("segmentation-sd-2-%s", rName)

	resourceOptions := &terraform.Options{
		TerraformDir: "./",
		VarFiles:     []string{"test-fixtures/segmentation_network_domain_connection_policy.tfvars"},
		Vars: map[string]interface{}{
			"segmentation_network_domain1": segmentationNetworkDomain1,
			"segmentation_network_domain2": segmentationNetworkDomain2,
		},
	}

	defer terraform.Destroy(ctx, resourceOptions)

	terraform.InitAndApply(ctx, resourceOptions)

	resourceID := fmt.Sprintf("aviatrix_segmentation_network_domain_connection_policy.test_segmentation_network_domain_connection_policy")
	verifySegmentationNetworkDomainConnectionPolicyExists(ctx, t, resourceID, segmentationNetworkDomain1, segmentationNetworkDomain2)
}

func verifySegmentationNetworkDomainConnectionPolicyExists(ctx context.Context, t *testing.T, resourceName string, segmentationNetworkDomain1 string, segmentationNetworkDomain2 string) {
	client := getAviatrixClientFromProvider()

	foundSegmentationNetworkDomainConnectionPolicy := &goaviatrix.SegmentationSecurityDomainConnectionPolicy{
		Domain1: &goaviatrix.SegmentationSecurityDomain{
			DomainName: segmentationNetworkDomain1,
		},
		Domain2: &goaviatrix.SegmentationSecurityDomain{
			DomainName: segmentationNetworkDomain2,
		},
	}

	err := resource.RetryContext(ctx, retryTimeout, func() *resource.RetryError {
		_, err := client.GetSegmentationSecurityDomainConnectionPolicy(foundSegmentationNetworkDomainConnectionPolicy)
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("failed to get segmentation network domain connection policy: %v", err))
		}

		if foundSegmentationNetworkDomainConnectionPolicy.Domain1.DomainName+"~"+foundSegmentationNetworkDomainConnectionPolicy.Domain2.DomainName != resourceName {
			return resource.RetryableError(fmt.Errorf("segmentation network domain connection policy not found yet"))
		}

		return nil
	})
	assert.Nil(t, err)
}

func getAviatrixClientFromProvider() *goaviatrix.Client {
	return testAccProvider.Meta().(*goaviatrix.Client)
}
func testAccCheckSegmentationNetworkDomainConnectionPolicyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_segmentation_network_domain_connection_policy" {
			continue
		}
		foundSegmentationNetworkDomainConnectionPolicy := &goaviatrix.SegmentationSecurityDomainConnectionPolicy{
			Domain1: &goaviatrix.SegmentationSecurityDomain{
				DomainName: rs.Primary.Attributes["domain_name_1"],
			},
			Domain2: &goaviatrix.SegmentationSecurityDomain{
				DomainName: rs.Primary.Attributes["domain_name_2"],
			},
		}
		_, err := client.GetSegmentationSecurityDomainConnectionPolicy(foundSegmentationNetworkDomainConnectionPolicy)
		if err == nil {
			return fmt.Errorf("segmentation_network_domain_connection_policy still exists")
		}
	}

	return nil
}

