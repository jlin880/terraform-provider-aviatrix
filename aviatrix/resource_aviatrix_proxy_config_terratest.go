package aviatrix

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/random"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAviatrixProxyConfig_basic(t *testing.T) {
    terraformOptions := &terraform.Options{
        TerraformDir: "./",
        Vars: map[string]interface{}{
            "account_name":       fmt.Sprintf("tfa-%s", acctest.RandString(5)),
            "cloud_type":         "1",
            "aws_account_number": os.Getenv("AWS_ACCOUNT_NUMBER"),
            "aws_iam":            "false",
            "aws_access_key":     os.Getenv("AWS_ACCESS_KEY"),
            "aws_secret_key":     os.Getenv("AWS_SECRET_KEY"),
            "http_proxy":         "proxy.aviatrixtest.com:3128",
            "https_proxy":        "proxy.aviatrixtest.com:3129",
        },
    }

    defer terraform.Destroy(t, terraformOptions)

    terraform.InitAndApply(t, terraformOptions)

    proxyConfigID := terraform.Output(t, terraformOptions, "proxy_config_id")

    // Verify that the proxy config exists
    aviatrixClient := getAviatrixClientFromTerraformOptions(t, terraformOptions)
    proxyConfig, err := aviatrixClient.GetProxyConfig()
    require.NoError(t, err)
    require.NotNil(t, proxyConfig)
    require.Equal(t, proxyConfigID, proxyConfig.ProxyConfigID)
}

func getAviatrixClientFromTerraformOptions(t *testing.T, terraformOptions *terraform.Options) *goaviatrix.Client {
    accountName := terraformOptions.Vars["account_name"].(string)

    aviatrixClient, err := goaviatrix.NewClientFromEnv(accountName)
    require.NoError(t, err)

    return aviatrixClient
}

