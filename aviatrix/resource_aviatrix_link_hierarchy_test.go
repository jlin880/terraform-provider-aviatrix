package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/goaviatrix"
)

func TestAviatrixLinkHierarchy_basic(t *testing.T) {
	t.Parallel()

	terraformOptions := &terraform.Options{
		TerraformDir: "./",
		Vars: map[string]interface{}{
			"name": random.UniqueId(),
			"link_name": fmt.Sprintf("link-%s", random.UniqueId()),
			"wan_tag": "wan3.10",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	resourceName := "aviatrix_link_hierarchy.test"
	linkHierarchyUUID := terraform.Output(t, terraformOptions, "aviatrix_link_hierarchy_uuid")
	linkName := terraform.Output(t, terraformOptions, "link_name")

	client := goaviatrix.NewClientFromEnv()

	linkHierarchy, err := client.GetLinkHierarchy(context.Background(), linkHierarchyUUID)
	assert.NoError(t, err)
	assert.NotNil(t, linkHierarchy)

	assert.Equal(t, terraformOptions.Vars["name"].(string), linkHierarchy.Name)

	links := linkHierarchy.Links
	assert.Len(t, links, 1)

	assert.Equal(t, linkName, links[0].Name)
	assert.Equal(t, terraformOptions.Vars["wan_tag"].(string), links[0].WANLink.WANTag)
}
func TestAccAviatrixLinkHierarchy_basic(t *testing.T) {
	t.Parallel()

	testName := fmt.Sprintf("link_hierarchy-%s", random.UniqueId())

	terraformOptions := &terraform.Options{
		TerraformDir: "../examples/link_hierarchy",
		Vars: map[string]interface{}{
			"name": testName,
			"link_name": fmt.Sprintf("link-%s", random.UniqueId()),
			"wan_tag": "wan3.10",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	// Check that the link hierarchy exists
	linkHierarchyID := terraform.Output(t, terraformOptions, "link_hierarchy_id")
	client := testAccProvider.Meta().(*goaviatrix.Client)
	_, err := client.GetLinkHierarchy(context.Background(), linkHierarchyID)
	assert.NoError(t, err)
}

func TestAccAviatrixLinkHierarchy_import(t *testing.T) {
	t.Parallel()

	testName := fmt.Sprintf("link_hierarchy-%s", random.UniqueId())

	terraformOptions := &terraform.Options{
		TerraformDir: "../examples/link_hierarchy",
		Vars: map[string]interface{}{
			"name": testName,
			"link_name": fmt.Sprintf("link-%s", random.UniqueId()),
			"wan_tag": "wan3.10",
		},
	}

	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApply(t, terraformOptions)

	resourceID := terraform.Output(t, terraformOptions, "link_hierarchy_id")

	importedTerraformOptions := &terraform.Options{
		TerraformDir: "../examples/link_hierarchy",
	}

	importedResourceName := "aviatrix_link_hierarchy.test"

	resource.Import(t, importedTerraformOptions, importedResourceName, resourceID)

	// Check that the link hierarchy exists
	client := testAccProvider.Meta().(*goaviatrix.Client)
	_, err := client.GetLinkHierarchy(context.Background(), resourceID)
	assert.NoError(t, err)
}
func testAccCheckLinkHierarchyExists(resourceName string) testutil.ResourceExists {
    return testutil.ResourceExists{
        ResourceType: "aviatrix_link_hierarchy",
        ResourceName: resourceName,
        Get:          getLinkHierarchy,
    }
}

func testAccCheckLinkHierarchyDestroy(s *terraform.State) error {
    client := testAccProvider.Meta().(*goaviatrix.Client)

    for _, rs := range s.RootModule().Resources {
        if rs.Type != "aviatrix_link_hierarchy" {
            continue
        }

        _, err := client.GetLinkHierarchy(context.Background(), rs.Primary.ID)
        if err == nil {
            return fmt.Errorf("link hierarchy still exists")
        }
    }

    return nil
}

func getLinkHierarchy(client *goaviatrix.Client, resourceName string) (interface{}, error) {
    rs, err := client.GetLinkHierarchy(context.Background(), resourceName)
    if err != nil {
        return nil, err
    }

    return rs, nil
}
