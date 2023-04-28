package test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAccAviatrixCloudnRegistration_basic(t *testing.T) {
	t.Parallel()

	skipAcc := os.Getenv("SKIP_CLOUDN_REGISTRATION")
	if skipAcc == "yes" {
		t.Skip("Skipping Aviatrix CloudN Registration test as SKIP_CLOUDN_REGISTRATION is set")
	}

	cloudnName := fmt.Sprintf("cloudn-%s", random.UniqueId())
	localASNumber := "65707"

	terraformOptions := &terraform.Options{
		TerraformDir: "../examples/cloudn_registration",
		EnvVars: map[string]string{
			"AVIATRIX_USERNAME": os.Getenv("AVIATRIX_USERNAME"),
			"AVIATRIX_PASSWORD": os.Getenv("AVIATRIX_PASSWORD"),
			"CLOUDN_IP":         os.Getenv("CLOUDN_IP"),
			"CLOUDN_USERNAME":   os.Getenv("CLOUDN_USERNAME"),
			"CLOUDN_PASSWORD":   os.Getenv("CLOUDN_PASSWORD"),
		},
		Vars: map[string]interface{}{
			"cloudn_name":       cloudnName,
			"local_as_number":   localASNumber,
			"prepend_as_path":   []string{localASNumber},
			"aviatrix_version":  os.Getenv("AVIATRIX_VERSION"),
			"controller_domain": os.Getenv("CONTROLLER_DOMAIN"),
		},
	}

	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApply(t, terraformOptions)

	// Verify the cloudn registration exists
	client := goaviatrix.NewClient(os.Getenv("AVIATRIX_CONTROLLER_IP"), os.Getenv("AVIATRIX_USERNAME"), os.Getenv("AVIATRIX_PASSWORD"), os.Getenv("AVIATRIX_VERSION"))
	cloudnRegistration := &goaviatrix.CloudnRegistration{
		Name: cloudnName,
	}

	err := client.GetCloudnRegistration(context.Background(), cloudnRegistration)
	assert.NoError(t, err)
	assert.Equal(t, cloudnName, cloudnRegistration.Name)
	assert.Equal(t, os.Getenv("CLOUDN_IP"), cloudnRegistration.Address)
	assert.Equal(t, localASNumber, cloudnRegistration.LocalASNumber)

	// Verify the cloudn registration was imported correctly
	importedCloudnRegistrationOptions := &terraform.Options{
		TerraformDir: "../examples/cloudn_registration",
		EnvVars: map[string]string{
			"AVIATRIX_USERNAME": os.Getenv("AVIATRIX_USERNAME"),
			"AVIATRIX_PASSWORD": os.Getenv("AVIATRIX_PASSWORD"),
			"CLOUDN_IP":         os.Getenv("CLOUDN_IP"),
		},
	}

	importedCloudnRegistration := terraform.Import(t, importedCloudnRegistrationOptions, "aviatrix_cloudn_registration.test_cloudn_registration")
	assert.Equal(t, cloudnName, importedCloudnRegistration["name"])
	assert.Equal(t, os.Getenv("CLOUDN_IP"), importedCloudnRegistration["address"])
	assert.Equal(t, localASNumber, importedCloudnRegistration["local_as_number"])
}


func TestAccAviatrixCloudnRegistration_basic(t *testing.T) {
	skipAcc := os.Getenv("SKIP_CLOUDN_REGISTRATION")
	if skipAcc == "yes" {
		t.Skip("Skipping Aviatrix CloudN Registration test as SKIP_CLOUDN_REGISTRATION is set")
	}

	terraformOptions := createCloudnRegistrationTerraformOptions(t)
	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApply(t, terraformOptions)

	checkCloudnRegistrationResources(t, terraformOptions)

	importedTerraformOptions := importCloudnRegistrationTerraformOptions(t, terraformOptions)
	checkCloudnRegistrationResources(t, importedTerraformOptions)
}

func createCloudnRegistrationTerraformOptions(t *testing.T) *terraform.Options {
	rName := fmt.Sprintf("cloudn-%s", random.UniqueId())
	localASNumber := "65707"
	address := os.Getenv("CLOUDN_IP")
	username := os.Getenv("CLOUDN_USERNAME")
	password := os.Getenv("CLOUDN_PASSWORD")

	terraformOptions := &terraform.Options{
		TerraformDir: "../examples/cloudn_registration/",
		Vars: map[string]interface{}{
			"name":            rName,
			"address":         address,
			"username":        username,
			"password":        password,
			"local_as_number": localASNumber,
		},
	}

	return terraformOptions
}

func importCloudnRegistrationTerraformOptions(t *testing.T, terraformOptions *terraform.Options) *terraform.Options {
	importedTerraformOptions := &terraform.Options{
		TerraformDir: "../examples/cloudn_registration/",
		BackendConfig: map[string]interface{}{
			"key": terraformOptions.StateFile,
		},
	}

	return importedTerraformOptions
}

func checkCloudnRegistrationResources(t *testing.T, terraformOptions *terraform.Options) {
	name := terraformOptions.Vars["name"].(string)
	localASNumber := terraformOptions.Vars["local_as_number"].(string)

	client := goaviatrix.NewClient(os.Getenv("AVIATRIX_API_ENDPOINT"), os.Getenv("AVIATRIX_API_ACCOUNT_NAME"), os.Getenv("AVIATRIX_API_ACCESS_KEY"), os.Getenv("AVIATRIX_API_SECRET_KEY"))

	cloudnRegistration := &goaviatrix.CloudnRegistration{
		Name: name,
	}

	err := client.GetCloudnRegistration(context.Background(), cloudnRegistration)
	assert.Nil(t, err)
	assert.Equal(t, address, cloudnRegistration.Address)
	assert.Equal(t, name, cloudnRegistration.Name)
	assert.Equal(t, localASNumber, cloudnRegistration.LocalASNumber)
	assert.Equal(t, "", cloudnRegistration.PrependASPath)
}


func testAccAviatrixCloudnRegistrationPreCheck(t *testing.T) {
	requiredEnv := []string{
		"CLOUDN_IP",
		"CLOUDN_USERNAME",
		"CLOUDN_PASSWORD",
	}

	for _, v := range requiredEnv {
		if os.Getenv(v) == "" {
			t.Fatalf("%s must be set for aviatrix_cloudn_registration acceptance test", v)
		}
	}
}
