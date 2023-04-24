package test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/assert"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

func TestAviatrixControllerCertDomainConfig_basic(t *testing.T) {
	skipAcc := os.Getenv("SKIP_CONTROLLER_CERT_DOMAIN_CONFIG")
	if skipAcc == "yes" {
		t.Skip("Skipping Controller Cert Domain Config test as SKIP_CONTROLLER_CERT_DOMAIN_CONFIG is set")
	}

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/controller_cert_domain_config/",
		Vars: map[string]interface{}{
			"cert_domain": random.UniqueId() + ".com",
		},
	})

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	controllerCertDomainResource := "aviatrix_controller_cert_domain_config.test"

	terraform.OutputRequired(t, terraformOptions, "resource_id")

	resourceID := terraform.Output(t, terraformOptions, "resource_id")

	client := goaviatrix.NewClientFromEnvParams()

	certDomainConfig, err := client.GetCertDomain(context.Background())
	assert.Nil(t, err)

	assert.Equal(t, terraformOptions.Vars["cert_domain"], certDomainConfig.CertDomain)

	assert.Equal(t, strings.Replace(client.ControllerIP, ".", "-", -1), resourceID)
}
func testAccCheckControllerCertDomainConfigDestroy(t *testing.T, state *terraform.State) error {
	client := testAccProviderVersionValidation.Meta().(*goaviatrix.Client)
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "aviatrix_controller_cert_domain_config" {
			continue
		}
	
		certDomainConfig, _ := client.GetCertDomain(context.Background())
		if !certDomainConfig.IsDefault {
			return fmt.Errorf("controller cert domain configured when it should be destroyed")
		}
	}
	
	return nil
}
