package test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
)

func TestAccAviatrixDistributedFirewallingConfig_basic(t *testing.T) {
	skipAcc := os.Getenv("SKIP_DISTRIBUTED_FIREWALLING_CONFIG")
	if skipAcc == "yes" {
		t.Skip("Skipping Distributed-firewalling config tests as SKIP_DISTRIBUTED_FIREWALLING_CONFIG is set")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../examples/aviatrix-distributed-firewalling-config",
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	distributedFirewallingStatus, err := aviatrix.GetDistributedFirewallingStatus(context.Background())
	assert.NoError(t, err)
	assert.True(t, distributedFirewallingStatus.EnableDistributedFirewalling)
}

func TestAccAviatrixDistributedFirewallingPolicyList_basic(t *testing.T) {
	skipAcc := os.Getenv("SKIP_DISTRIBUTED_FIREWALLING_POLICY_LIST")
	if skipAcc == "yes" {
		t.Skip("Skipping Distributed Firewalling Policy List test as SKIP_DISTRIBUTED_FIREWALLING_POLICY_LIST is set")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../examples/aviatrix-distributed-firewalling-policy-list",
		Vars: map[string]interface{}{
			"policies": []map[string]interface{}{
				{
					"name":             "test-distributed-firewalling-policy",
					"action":           "PERMIT",
					"logging":          true,
					"priority":         0,
					"protocol":         "TCP",
					"src_smart_groups": []string{"${aviatrix_smart_group.ad1.uuid}"},
					"dst_smart_groups": []string{"${aviatrix_smart_group.ad2.uuid}"},
					"port_ranges": []map[string]interface{}{
						{
							"hi": 10,
							"lo": 1,
						},
					},
				},
			},
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	policies, err := aviatrix.GetDistributedFirewallingPolicyList(context.Background())
	assert.NoError(t, err)

	assert.Equal(t, 1, len(policies))

	assert.Equal(t, "test-distributed-firewalling-policy", policies[0].Name)
	assert.Equal(t, "PERMIT", policies[0].Action)
	assert.Equal(t, true, policies[0].Logging)
	assert.Equal(t, 0, policies[0].Priority)
	assert.Equal(t, "TCP", policies[0].Protocol)
	assert.Equal(t, []string{"${aviatrix_smart_group.ad1.uuid}"}, policies[0].SrcSmartGroups)
	assert.Equal(t, []string{"${aviatrix_smart_group.ad2.uuid}"}, policies[0].DstSmartGroups)
	assert.Equal(t, 10, policies[0].PortRanges[0].Hi)
	assert.Equal(t, 1, policies[0].PortRanges[0].Lo)
}

func TestAccAviatrixDistributedFirewallingPolicyList_basic(t *testing.T) {
	skipAcc := os.Getenv("SKIP_DISTRIBUTED_FIREWALLING_POLICY_LIST")
	if skipAcc == "yes" {
		t.Skip("Skipping Distributed Firewalling Policy List test as SKIP_DISTRIBUTED_FIREWALLING_POLICY_LIST is set")
	}

	resourceName := "aviatrix_distributed_firewalling_policy_list.test"
	r := random.UniqueId()

	terraformOptions := &terraform.Options{
		TerraformDir: "./",
		Vars: map[string]interface{}{
			"policies": []map[string]interface{}{
				{
					"name":             fmt.Sprintf("test-distributed-firewalling-policy-%s", r),
					"action":           "PERMIT",
					"logging":          true,
					"priority":         0,
					"protocol":         "TCP",
					"src_smart_groups": []string{"test-smart_group-1.uuid"},
					"dst_smart_groups": []string{"test-smart_group-2.uuid"},
					"port_ranges": []map[string]interface{}{
						{
							"hi": 10,
							"lo": 1,
						},
					},
				},
			},
		},
	}

	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApply(t, terraformOptions)

	resource.Assert(t, terraformOptions, terraform.OutputExists(resourceName))
	resource.Assert(t, terraformOptions, terraform.OutputEquals(resourceName, "policies.#", "1"))
	resource.Assert(t, terraformOptions, terraform.OutputEquals(resourceName, "policies.0.name", fmt.Sprintf("test-distributed-firewalling-policy-%s", r)))
	resource.Assert(t, terraformOptions, terraform.OutputEquals(resourceName, "policies.0.action", "PERMIT"))
	resource.Assert(t, terraformOptions, terraform.OutputEquals(resourceName, "policies.0.priority", "0"))
	resource.Assert(t, terraformOptions, terraform.OutputEquals(resourceName, "policies.0.protocol", "TCP"))
	resource.Assert(t, terraformOptions, terraform.OutputEquals(resourceName, "policies.0.logging", "true"))
	resource.Assert(t, terraformOptions, terraform.OutputContains(resourceName, "policies.0.src_smart_groups", "test-smart_group-1"))
	resource.Assert(t, terraformOptions, terraform.OutputContains(resourceName, "policies.0.dst_smart_groups", "test-smart_group-2"))
}

func testAccDistributedFirewallingPolicyListDestroy(t *testing.T, client *goaviatrix.Client, resourceName string) {
	err := resource.RetryContext(context.Background(), 3*time.Minute, func() *resource.RetryError {
		rs, err := client.GetDistributedFirewallingPolicyList(context.Background())
		if err != nil {
			if err == goaviatrix.ErrNotFound {
				// The distributed firewalling policy list is not found, meaning it has been deleted.
				return nil
			}
			// The client encountered an error trying to retrieve the distributed firewalling policy list.
			return resource.NonRetryableError(fmt.Errorf("failed to get Distributed-firewalling Policy List status: %v", err))
		}

		if rs.EnableDistributedFirewalling {
			// The distributed firewalling policy list still exists, meaning it has not been deleted yet.
			return resource.RetryableError(fmt.Errorf("distributed-firewalling policy list still exists"))
		}

		// The distributed firewalling policy list has been deleted.
		return nil
	})

	require.NoError(t, err)
}

