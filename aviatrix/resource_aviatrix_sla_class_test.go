package aviatrix_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/aviatrix"
	"github.com/AviatrixSystems/terraform-provider-aviatrix/aviatrix/test"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAccAviatrixSLAClass_basic(t *testing.T) {
	t.Parallel()

	testStage := fmt.Sprintf("acc-%s", strings.ToLower(random.UniqueId()))

	terraformOptions := test.GetTerraformOptions(testStage, "../examples/sla_class_basic")
	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	resourceName := "aviatrix_sla_class.test"
	slaClassName := fmt.Sprintf("sla-%s", strings.ToLower(random.UniqueId()))

	expectedSLAClass := aviatrix.SLAClass{
		Name:           slaClassName,
		Latency:        43,
		Jitter:         1,
		PacketDropRate: 3,
	}

	// Check the SLA Class resource exists
	assert.NoError(t, test.WaitForResource(t, terraformOptions, resourceName))

	actualSLAClass, err := aviatrix.GetSLAClass(context.Background(), terraformOptions, slaClassName)
	assert.NoError(t, err)

	// Check the SLA Class attributes are correct
	assert.Equal(t, expectedSLAClass.Name, actualSLAClass.Name)
	assert.Equal(t, expectedSLAClass.Latency, actualSLAClass.Latency)
	assert.Equal(t, expectedSLAClass.Jitter, actualSLAClass.Jitter)
	assert.Equal(t, expectedSLAClass.PacketDropRate, actualSLAClass.PacketDropRate)

	// Import the SLA Class resource
	importedResourceName := "imported-sla-class"
	importedResource := terraform.ImportState(resourceName, importedResourceName, t, terraformOptions)
	importedSLAClass, err := aviatrix.GetSLAClass(context.Background(), terraformOptions, slaClassName)
	assert.NoError(t, err)

	// Verify the imported SLA Class resource
	assert.Equal(t, expectedSLAClass.Name, importedSLAClass.Name)
	assert.Equal(t, expectedSLAClass.Latency, importedSLAClass.Latency)
	assert.Equal(t, expectedSLAClass.Jitter, importedSLAClass.Jitter)
	assert.Equal(t, expectedSLAClass.PacketDropRate, importedSLAClass.PacketDropRate)
}


func testAccSLAClassBasic(slaClassName string) string {
	return fmt.Sprintf(`
resource "aviatrix_sla_class" "test" {
	name             = "%s"
	latency          = 43	
	jitter           = 1
	packet_drop_rate = 3
}
 `, slaClassName)
}

func testAccCheckSLAClassExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("sla class not found: %s", resourceName)
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		_, err := client.GetSLAClass(context.Background(), rs.Primary.Attributes["uuid"])
		if err == goaviatrix.ErrNotFound {
			return fmt.Errorf("sla class not found")
		}

		return nil
	}
}

func testAccCheckSLAClassDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_sla_class" {
			continue
		}

		_, err := client.GetSLAClass(context.Background(), rs.Primary.Attributes["uuid"])
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("sla class still exists")
		}
	}

	return nil
}

