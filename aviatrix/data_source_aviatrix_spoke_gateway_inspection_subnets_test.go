package test

import (
	"fmt"
	"testing"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

func TestAccDataSourceAviatrixSpokeGatewayInspectionSubnets_basic(t *testing.T) {
	t.Parallel()

	// Specify the subscription ID, directory ID, application ID, and application key for the Azure provider
	subscriptionID := "<AZURE_SUBSCRIPTION_ID>"
	directoryID := "<AZURE_DIRECTORY_ID>"
	applicationID := "<AZURE_APPLICATION_ID>"
	applicationKey := "<AZURE_APPLICATION_KEY>"

	// Construct the Terraform options with the path to Terraform code
	terraformOptions := &terraform.Options{
		TerraformDir: "./",
		Vars: map[string]interface{}{
			"subscription_id": subscriptionID,
			"directory_id":    directoryID,
			"application_id":  applicationID,
			"application_key": applicationKey,
		},
	}

	// Clean up resources after the test is done
	defer terraform.Destroy(t, terraformOptions)

	// Deploy the Terraform code
	terraform.InitAndApply(t, terraformOptions)

	// Get the output variables
	gatewayName := terraform.Output(t, terraformOptions, "gateway_name")
	expectedSubnets := []string{"18.9.16.0/20~~test-vpc-Public-subnet-1", "18.9.32.0/20~~test-vpc-Private-subnet-1", "18.9.48.0/20~~test-vpc-Public-subnet-2", "18.9.64.0/20~~test-vpc-Private-subnet-2"}

	// Verify that the spoke gateway inspection subnets match the expected subnets
	actualSubnets, err := GetSpokeGatewayInspectionSubnets(t, gatewayName)
	if err != nil {
		t.Fatalf("Failed to get spoke gateway inspection subnets: %v", err)
	}
	if !Equivalent(actualSubnets, expectedSubnets) {
		t.Fatalf("Spoke gateway inspection subnets do not match the expected subnets. Expected: %v, but got: %v", expectedSubnets, actualSubnets)
	}
}

// GetSpokeGatewayInspectionSubnets returns the inspection subnets for a spoke gateway with the given name
func GetSpokeGatewayInspectionSubnets(t *testing.T, gatewayName string) ([]string, error) {
	client := GetAviatrixClient(t)

	subnets, err := client.GetSubnetsForInspection(gatewayName)
	if err != nil {
		return nil, fmt.Errorf("Failed to get spoke gateway inspection subnets for gateway %s: %v", gatewayName, err)
	}

	return subnets, nil
}

// GetAviatrixClient returns an Aviatrix API client
func GetAviatrixClient(t *testing.T) *goaviatrix.Client {
	aviatrixAccessKey := "<AVIATRIX_ACCESS_KEY>"
	aviatrixSecretKey := "<AVIATRIX_SECRET_KEY>"
	aviatrixAPIEndpoint := "<AVIATRIX_API_ENDPOINT>"
	client, err := goaviatrix.NewClient(aviatrixAccessKey, aviatrixSecretKey, aviatrixAPIEndpoint)
	if err != nil {
		t.Fatalf("Failed to create Aviatrix API client: %v", err)
	}
	return client
}

// Equivalent returns true if the two slices have the same elements, regardless of order
