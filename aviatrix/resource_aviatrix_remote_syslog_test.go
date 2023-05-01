package test

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAviatrixRemoteSyslog_basic(t *testing.T) {
	if os.Getenv("SKIP_REMOTE_SYSLOG") == "yes" {
		t.Skip("Skipping remote syslog test as SKIP_REMOTE_SYSLOG is set")
	}

	rIndex := random.RandIntRange(0, 9)
	rName := random.UniqueId()
	resourceName := "aviatrix_remote_syslog.test_remote_syslog"

	terraformOptions := &terraform.Options{
		TerraformDir: "./",
		Vars: map[string]interface{}{
			"index":             rIndex,
			"name":              rName,
			"server":            "1.2.3.4",
			"port":              10,
			"protocol":          "TCP",
			"excluded_gateways": []string{"a", "b"},
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	checkRemoteSyslogExists(t, terraformOptions, resourceName, rIndex)
	checkResourceAttrs(t, terraformOptions, resourceName, rIndex, rName)
	checkRemoteSyslogExcludedGatewaysMatch(t, resourceName, rIndex, []string{"a", "b"})
}

func checkRemoteSyslogExists(t *testing.T, terraformOptions *terraform.Options, resourceName string, index int) {
	client := getAviatrixClient(t)

	err := terraform.WithRetryableErrors(t, terraformOptions, func() error {
		_, err := client.GetRemoteSyslogStatus(index)
		if err == goaviatrix.ErrNotFound {
			return fmt.Errorf("remote syslog %d not found", index)
		}

		return nil
	})

	assert.NoError(t, err)
}

func checkResourceAttrs(t *testing.T, terraformOptions *terraform.Options, resourceName string, index int, name string) {
	expectedResourceAttrs := map[string]string{
		"index": strconv.Itoa(index),
		"name":  name,
		"server": "1.2.3.4",
		"port": "10",
		"protocol": "TCP",
	}

	actualResourceAttrs := terraform.OutputAll(t, terraformOptions, resourceName)

	for attrName, expectedAttrValue := range expectedResourceAttrs {
		actualAttrValue := actualResourceAttrs[attrName].Value

		assert.Equal(t, expectedAttrValue, actualAttrValue, "Attribute %s does not match", attrName)
	}
}

func checkRemoteSyslogExcludedGatewaysMatch(t *testing.T, resourceName string, index int, input []string) {
	client := getAviatrixClient(t)

	err := terraform.WithRetryableErrors(t, terraformOptions, func() error {
		resp, _ := client.GetRemoteSyslogStatus(index)
		if !goaviatrix.Equivalent(resp.ExcludedGateways, input) {
			return fmt.Errorf("excluded gateways don't match with the input")
		}

		return nil
	})

	assert.NoError(t, err)
}

func getAviatrixClient(t *testing.T) *goaviatrix.Client {
	aviatrixAccessKey := os.Getenv("AVIATRIX_ACCESS_KEY")
	aviatrixSecretKey := os.Getenv("AVIATRIX_SECRET_KEY")

	if aviatrixAccessKey == "" || aviatrixSecretKey == "" {
		t.Fatal("Both AVIATRIX_ACCESS_KEY and AVIATRIX_SECRET_KEY environment variables must be set")
	}

	client, err := goaviatrix.NewClient(aviatrixAccessKey, aviatrixSecretKey, "", goaviatrix.DefaultRetryer)
	if err != nil {
		t.Fatalf("Failed to create Aviatrix client: %v", err)
	}

	return client
}
