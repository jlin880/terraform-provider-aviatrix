package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/assert"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/terraform_testing"
)

func TestAccAviatrixSamlEndpoint(t *testing.T) {
	t.Parallel()

	var samlEndpoint goaviatrix.SamlEndpoint
	idpMetadata := "https://idp.example.com/metadata"
	idpMetadataType := "PingFederate"
	rName := fmt.Sprintf("test-%s", random.UniqueId())

	resourceName := "aviatrix_saml_endpoint.foo"
	terraformOptions := &terraform.Options{
		TerraformDir: ".",
		Vars: map[string]interface{}{
			"endpoint_name":     rName,
			"idp_metadata":      idpMetadata,
			"idp_metadata_type": idpMetadataType,
		},
	}

	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApply(t, terraformOptions)

	checkResourceAttrs(t, terraformOptions, resourceName, rName, idpMetadata, idpMetadataType)

	state := terraform_testing.LoadStateFromFile(t, terraformOptions.StateFile)

	resourceChecks := []resource.TestCheckFunc{
		tesAccCheckSamlEndpointExists(resourceName, &samlEndpoint),
	}

	for _, resourceCheck := range resourceChecks {
		err := resourceCheck(state)
		assert.NoError(t, err)
	}
}

func testAccSamlEndpointConfigBasic(rName string, idpMetadata string, idpMetadataType string) string {
	return fmt.Sprintf(`
resource "aviatrix_saml_endpoint" "foo" {
	endpoint_name     = "%s"
	idp_metadata_type = "%s"
	idp_metadata      = "%s"
}
	`, rName, idpMetadataType, idpMetadata)
}

func tesAccCheckSamlEndpointExists(n string, samlEndpoint *goaviatrix.SamlEndpoint) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("aviatrix SAML endpoint not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no aviatrix SAML endpoint ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		foundSamlEndpoint := &goaviatrix.SamlEndpoint{
			EndPointName: rs.Primary.Attributes["endpoint_name"],
		}

		_, err := client.GetSamlEndpoint(foundSamlEndpoint)
		if err != nil {
			return err
		}
		if foundSamlEndpoint.EndPointName != rs.Primary.Attributes["endpoint_name"] {
			return fmt.Errorf("endpoint_name not found in created attributes")
		}

		*samlEndpoint = *foundSamlEndpoint

		return nil
	}
}

func testAccCheckSamlEndpointDestroy(t *testing.T, terraformOptions *terraform.Options) {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, resourceName := range terraform.GetResourceNames(t, terraformOptions) {
		if resourceName != "aviatrix_saml_endpoint.foo" {
			continue
		}

		resourceInstanceState := terraform.InstanceState{
			ID: resourceName,
			Attributes: map[string]string{
				"endpoint_name": terraform.Output(t, terraformOptions, "endpoint_name"),
			},
		}

		foundSamlEndpoint := &goaviatrix.SamlEndpoint{
			EndPointName: resourceInstanceState.Attributes["endpoint_name"],
		}

		_, err := client.GetSamlEndpoint(foundSamlEndpoint)
		if err != goaviatrix.ErrNotFound {
			t.Errorf("aviatrix Saml Endpoint still exists")
		}
	}
}
