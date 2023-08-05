package aviatrix_test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAviatrixClient(t *testing.T) {
	// Set up the Terraform options for the Aviatrix provider
	terraformOptions := &terraform.Options{
		TerraformDir: "<path_to_terraform_module>", // Update with the path to your Terraform module
	}

	// Run Terraform apply and defer destroy
	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApply(t, terraformOptions)

	// Get the output values from Terraform
	username := terraform.Output(t, terraformOptions, "username")
	password := terraform.Output(t, terraformOptions, "password")
	controllerIP := terraform.Output(t, terraformOptions, "controller_ip")

	// Create the Aviatrix client configuration
	config := &aviatrix.Config{
		Username:     username,
		Password:     password,
		ControllerIP: controllerIP,
		VerifyCert:   true,
		PathToCACert: "<path_to_ca_cert>", // Update with the path to your CA certificate
	}

	// Create the Aviatrix client
	client, err := config.Client()

	// Assert that the Aviatrix client is not nil and no error occurred
	assert.NotNil(t, client)
	assert.NoError(t, err)
}
