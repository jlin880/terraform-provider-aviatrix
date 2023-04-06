package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAccDataSourceAviatrixFirewall_basic(t *testing.T) {
	skipAcc := os.Getenv("SKIP_DATA_FIREWALL")
	if skipAcc == "yes" {
		t.Skip("Skipping data source firewall tests as 'SKIP_DATA_FIREWALL' is set")
	}

	testAccDir := "./testdata/aviatrix_firewall"
	rName := RandomString(5)
	terraformOptions := &terraform.Options{
		TerraformDir: testAccDir,
		Vars: map[string]interface{}{
			"rand": rName,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	// Check the firewall data source
	expectedGwName := fmt.Sprintf("test-gw-%s", rName)
	expectedBasePolicy := "allow-all"
	expectedBaseLogEnabled := true

	actualGwName := terraform.Output(t, terraformOptions, "aviatrix_firewall.test.gw_name")
	actualBasePolicy := terraform.Output(t, terraformOptions, "aviatrix_firewall.test.base_policy")
	actualBaseLogEnabled := terraform.Output(t, terraformOptions, "aviatrix_firewall.test.base_log_enabled")

	assert.Equal(t, expectedGwName, actualGwName)
	assert.Equal(t, expectedBasePolicy, actualBasePolicy)
	assert.Equal(t, expectedBaseLogEnabled, actualBaseLogEnabled)

	// Check the firewall policy data source
	expectedSrcIP := "10.15.0.224/32"
	expectedDstIP := "10.12.0.172/32"
	expectedProtocol := "tcp"
	expectedPort := "0:65535"
	expectedDescription := "This is policy no.1"
	expectedAction := "allow"
	expectedLogEnabled := false

	actualSrcIP := terraform.Output(t, terraformOptions, "aviatrix_firewall.test_policy.0.src_ip")
	actualDstIP := terraform.Output(t, terraformOptions, "aviatrix_firewall.test_policy.0.dst_ip")
	actualProtocol := terraform.Output(t, terraformOptions, "aviatrix_firewall.test_policy.0.protocol")
	actualPort := terraform.Output(t, terraformOptions, "aviatrix_firewall.test_policy.0.port")
	actualDescription := terraform.Output(t, terraformOptions, "aviatrix_firewall.test_policy.0.description")
	actualAction := terraform.Output(t, terraformOptions, "aviatrix_firewall.test_policy.0.action")
	actualLogEnabled := terraform.Output(t, terraformOptions, "aviatrix_firewall.test_policy.0.log_enabled")

	assert.Equal(t, expectedSrcIP, actualSrcIP)
	assert.Equal(t, expectedDstIP, actualDstIP)
	assert.Equal(t, expectedProtocol, actualProtocol)
	assert.Equal(t, expectedPort, actualPort)
	assert.Equal(t, expectedDescription, actualDescription)
	assert.Equal(t, expectedAction, actualAction)
	assert.Equal(t, expectedLogEnabled, actualLogEnabled)
}

func TestTerraformDataSourceAviatrixFirewall_basic(t *testing.T) {
	t.Parallel()

	// Set the skip flag.
	skipAcc := os.Getenv("SKIP_DATA_FIREWALL")
	if skipAcc == "yes" {
		t.Skip("Skipping data source firewall tests as 'SKIP_DATA_FIREWALL' is set")
	}

	// Define the resource name and a message to display in the pre-checks.
	rName := randomUniqueID()
	resourceName := "data.aviatrix_firewall.test"
	msg := ". Set 'SKIP_DATA_FIREWALL' to 'yes' to skip data source firewall tests"

	// Set up Terraform options.
	terraformOptions := &terraform.Options{
		TerraformDir: "./",
		Vars: map[string]interface{}{
			"account_name":       fmt.Sprintf("tfa-%s", rName),
			"aws_account_number": os.Getenv("AWS_ACCOUNT_NUMBER"),
			"aws_access_key":     os.Getenv("AWS_ACCESS_KEY"),
			"aws_secret_key":     os.Getenv("AWS_SECRET_KEY"),
			"region":             os.Getenv("AWS_REGION"),
			"cidr":               "10.0.0.0/16",
			"gw_name":            fmt.Sprintf("test-gw-%s", rName),
			"base_policy":        "allow-all",
			"base_log_enabled":   true,
			"src_ip_1":           "10.15.0.224/32",
			"log_enabled_1":      false,
			"dst_ip_1":           "10.12.0.172/32",
			"action_1":           "allow",
			"port_1":             "0:65535",
			"description_1":      "This is policy no.1",
			"src_ip_2":           "10.15.1.224/32",
			"log_enabled_2":      true,
			"dst_ip_2":           "10.12.1.172/32",
			"action_2":           "deny",
			"port_2":             "0:65535",
			"description_2":      "This is policy no.2",
		},
	}

	// At the end of the test, destroy the Terraform resources.
	defer terraform.Destroy(t, terraformOptions)

	// Create the Terraform resources.
	terraform.InitAndApply(t, terraformOptions)

	// Run the tests.
	resource := terraform.Output(t, terraformOptions, "gw_name")
	if resource != fmt.Sprintf("test-gw-%s", rName) {
		t.Fatalf("Error: expected output gw_name to be 'test-gw-%s' but got '%s'", rName, resource)
	}

	basePolicy := terraform.Output(t, terraformOptions, "base_policy")
	if basePolicy != "allow-all" {
		t.Fatalf("Error: expected output base_policy to be 'allow-all' but got '%s'", basePolicy)
	}

	baseLogEnabled := terraform.Output(t, terraformOptions, "base_log_enabled")
	if baseLogEnabled != "true" {
		t.Fatalf("Error: expected output base_log_enabled to be 'true' but got '%s'", baseLogEnabled)
	}
}
func TestAccDataSourceAviatrixFirewall(t *testing.T) {
	skipAcc := os.Getenv("SKIP_DATA_FIREWALL")
	if skipAcc == "yes" {
		t.Skip("Skipping data source firewall tests as 'SKIP_DATA_FIREWALL' is set")
	}

	rName := fmt.Sprintf("terratest-%s", RandomString(6))
	resourceName := "data.aviatrix_firewall.test"

	terraformOptions := &terraform.Options{
		TerraformDir: "../examples/firewall",
		Vars: map[string]interface{}{
			"account_name":       fmt.Sprintf("tfa-%s", rName),
			"aws_account_number": os.Getenv("AWS_ACCOUNT_NUMBER"),
			"aws_access_key":     os.Getenv("AWS_ACCESS_KEY"),
			"aws_secret_key":     os.Getenv("AWS_SECRET_KEY"),
			"vpc_name":           fmt.Sprintf("tfv-%s", rName),
			"vpc_region":         os.Getenv("AWS_REGION"),
			"gw_name":            fmt.Sprintf("test-gw-%s", rName),
			"gw_size":            "t2.micro",
			"fw_policy1_src_ip":  "10.15.0.224/32",
			"fw_policy1_dst_ip":  "10.12.0.172/32",
			"fw_policy2_src_ip":  "10.15.1.224/32",
			"fw_policy2_dst_ip":  "10.12.1.172/32",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	data := terraform.OutputAll(t, terraformOptions)
	assert.Equal(t, data[resourceName+".gw_name"].Value, fmt.Sprintf("test-gw-%s", rName))
	assert.Equal(t, data[resourceName+".base_policy"].Value, "allow-all")
	assert.Equal(t, data[resourceName+".base_log_enabled"].Value, "true")
}
