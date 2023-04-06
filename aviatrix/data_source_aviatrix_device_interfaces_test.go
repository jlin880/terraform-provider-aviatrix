package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/assert"
)

func TestAccDataSourceAviatrixDeviceInterfaces_basic(t *testing.T) {
	t.Parallel()

	rName := random.UniqueId()

	skipAcc := os.Getenv("SKIP_DATA_DEVICE_INTERFACES")
	if skipAcc == "yes" {
		t.Skip("Skipping Data Source Device Interfaces tests as SKIP_DATA_DEVICE_INTERFACES is set")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../examples/device_interfaces",
		Vars: map[string]interface{}{
			"device_name": os.Getenv("CLOUDN_DEVICE_NAME"),
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	dataSourceName := "data.aviatrix_device_interfaces.foo"
	wanPrimaryInterface := terraform.Output(t, terraformOptions, "wan_primary_interface")
	wanPrimaryInterfacePublicIP := terraform.Output(t, terraformOptions, "wan_primary_interface_public_ip")

	assert.NotEmpty(t, wanPrimaryInterface)
	assert.NotEmpty(t, wanPrimaryInterfacePublicIP)
	assert.Equal(t, dataSourceName, terraformOptions.StatePath)

}
func TestAccDataSourceAviatrixDeviceInterfaces_basic(t *testing.T) {
	rName := acctest.RandString(5)
	resourceName := "data.aviatrix_device_interfaces.foo"

	skipAcc := os.Getenv("SKIP_DATA_DEVICE_INTERFACES")
	if skipAcc == "yes" {
		t.Skip("Skipping Data Source Device Interfaces tests as SKIP_DATA_DEVICE_INTERFACES is set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDeviceInterfacesConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceAviatrixDeviceInterfaces(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "wan_interfaces.0.wan_primary_interface"),
					resource.TestCheckResourceAttrSet(resourceName, "wan_interfaces.0.wan_primary_interface_public_ip"),
				),
			},
		},
	})
}

func testAccDataSourceDeviceInterfacesConfigBasic(rName string) string {
	return fmt.Sprintf(`
data "aviatrix_device_interfaces" "foo" {
	device_name = "%s"
}
	`, os.Getenv("CLOUDN_DEVICE_NAME"))
}

func testAccDataSourceAviatrixDeviceInterfaces(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("root module has no data source called %s", name)
		}

		return nil
	}
}
