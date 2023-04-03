package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/goaviatrix"
)

func TestAccAviatrixDeviceInterfaceConfig_basic(t *testing.T) {
	t.Parallel()

	skipIfEnvSet(t, "SKIP_DEVICE_INTERFACE_CONFIG")

	deviceName := os.Getenv("CLOUDN_DEVICE_NAME")
	uniqueID := random.UniqueId()

	terraformOptions := &terraform.Options{
		TerraformDir: "./fixtures/device_interface_config",
		Vars: map[string]interface{}{
			"device_name": deviceName,
			"resource_name": fmt.Sprintf("test_device_interface_config_%s", uniqueID),
		},
	}

	//terraform init and apply
	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApply(t, terraformOptions)

	//Get the instance of Client struct to perform tests
	client := aviatrix.NewClient()

	resourceName := fmt.Sprintf("aviatrix_device_interface_config.test_device_interface_config_%s", uniqueID)

	//Read the resource attributes into struct using aviatrix API
	device, err := client.GetDevice(&goaviatrix.Device{Name: deviceName})
	if err != nil {
		t.Fatal(err)
	}

	//Assert that the expected device and primary interface values are returned
	assert.Equal(t, deviceName, device.Name)
	assert.Equal(t, terraformOptions.Vars["wan_primary_interface"].(string), device.PrimaryInterface)

	//Import the state into Terraform
	importOpts := &terraform.ImportOptions{
		TerraformAddr: resourceName,
		ID:            resourceName,
	}

	//terraform import
	terraform.Import(t, terraformOptions, *importOpts)

	//terraform refresh to verify state
	terraform.Refresh(t, terraformOptions)

	//Verify that the imported state matches the terraform state
	assert.NoError(t, terraform.OutputStruct(resourceName, &device))
}

func skipIfEnvSet(t *testing.T, envVar string) {
	if os.Getenv(envVar) == "yes" {
		t.Skip(fmt.Sprintf("Skipping test as %s is set", envVar))
	}
}


func testAccDeviceInterfaceConfigBasic() string {
	return fmt.Sprintf(`
data "aviatrix_device_interfaces" "test" {
	device_name = "%s"
}

resource "aviatrix_device_interface_config" "test_device_interface_config" {
	device_name                     = data.aviatrix_device_interfaces.test.device_name
	wan_primary_interface           = "eth0"
	wan_primary_interface_public_ip = data.aviatrix_device_interfaces.test.wan_interfaces[0].wan_primary_interface_public_ip
}
`, os.Getenv("CLOUDN_DEVICE_NAME"))
}

func testAccCheckDeviceInterfaceConfigExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("ONE device_interface_config Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no device_interface_config ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		device := &goaviatrix.Device{Name: rs.Primary.Attributes["device_name"]}

		device, err := client.GetDevice(device)
		if err != nil {
			return err
		}

		if device.Name != rs.Primary.ID ||
			device.PrimaryInterface != rs.Primary.Attributes["wan_primary_interface"] {
			return fmt.Errorf("device_interface_config not found")
		}

		return nil
	}
}
