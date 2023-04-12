package aviatrix_test

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/gruntwork-io/terratest/modules/terraform/resource"
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

func TestAccDataSourceAviatrixDeviceInterfaces(t *testing.T) {
	t.Parallel()

	rName := random.UniqueId()
	resourceName := "data.aviatrix_device_interfaces.foo"

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

	err := resource.Retry(t, 5, 5*time.Second, func() *resource.RetryError {
		_, err := terraform.Provider().(*schema.Provider).Meta().(*aviatrix.Client).GetDeviceInterfaces(os.Getenv("CLOUDN_DEVICE_NAME"))

		if err != nil {
			if strings.Contains(err.Error(), "failed to authenticate") {
				return resource.RetryableError(fmt.Errorf("authentication failed, retrying: %s", err))
			}
			return resource.NonRetryableError(fmt.Errorf("failed to get device interfaces: %s", err))
		}
		return nil
	})

	if err != nil {
		t.Fatal(err)
	}

	err = testAccDataSourceAviatrixDeviceInterfaces(resourceName)(terraformOptions.State)
	assert.NoError(t, err)
}

func testAccDataSourceAviatrixDeviceInterfaces(name string) terraform.ResourceCheck {
	return terraform.ResourceCheck{
		Name: name,
		Exists: true,
		ExpectedOutput: "wan_primary_interface,wan_primary_interface_public_ip",
	}
}
