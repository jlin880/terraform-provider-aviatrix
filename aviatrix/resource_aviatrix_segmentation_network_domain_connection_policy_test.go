package aviatrix_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAviatrixSegmentationNetworkDomainConnectionPolicy_basic(t *testing.T) {
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

	require.NoError(t, resource.RetryContext(ctx, retryTimeout, func() *resource.RetryError {
		_, err := client.GetSegmentationSecurityDomainConnectionPolicy(foundSegmentationNetworkDomainConnectionPolicy)
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("failed to get segmentation network domain connection policy: %v", err))
		}

		if foundSegmentationNetworkDomainConnectionPolicy.Domain1.DomainName+"~"+foundSegmentationNetworkDomainConnectionPolicy.Domain2.DomainName != resourceName {
			return resource.RetryableError(fmt.Errorf("segmentation network domain connection policy not found yet"))
		}

		return nil
	}))
}

func ResourceDestroyFunc(resourceType string) schema.ResourceTestStep {
	return schema.ResourceTestStep{
		// We assume that the resource id is formed by concatenating the values of its attributes
		// separated by the character "~"
		// For instance, a policy connecting domains "domain1" and "domain2" would have an id "domain1~domain2"
		Check: resource.ComposeTestCheckFunc(
			func(state *terraform.State) error {
				client := GetAviatrixClientFromProvider()

				for _, rs := range state.RootModule().Resources {
					if rs.Type != resourceType {
						continue
					}

					domain1Name := rs.Primary.Attributes["domain_name_1"]
					domain2Name := rs.Primary.Attributes["domain_name_2"]

					policy := &goaviatrix.SegmentationSecurityDomainConnectionPolicy{
						Domain1: &goaviatrix.SegmentationSecurityDomain{DomainName: domain1Name},
						Domain2: &goaviatrix.SegmentationSecurityDomain{DomainName: domain2Name},
					}

					_, err := client.GetSegmentationSecurityDomainConnectionPolicy(policy)
					if err == nil {
						return fmt.Errorf("Resource %s with id %s still exists", resourceType, rs.Primary.ID)
					}
				}

				return nil
			},
		),
	}
}