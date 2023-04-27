package test

import (
    "fmt"
    "os"
    "testing"

    "github.com/gruntwork-io/terratest/modules/acctest"
    "github.com/gruntwork-io/terratest/modules/random"
    "github.com/gruntwork-io/terratest/modules/terraform"
    "github.com/stretchr/testify/assert"

    "github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
)

func preAzurePeerCheck(t *testing.T, msgCommon string) {
    vNet1 := os.Getenv("AZURE_VNET_ID")
    if vNet1 == "" {
        t.Fatal("Environment variable AZURE_VNET_ID is not set" + msgCommon)
    }
    vNet2 := os.Getenv("AZURE_VNET_ID2")
    if vNet2 == "" {
        t.Fatal("Environment variable AZURE_VNET_ID2 is not set" + msgCommon)
    }

    region1 := os.Getenv("AZURE_REGION")
    if region1 == "" {
        t.Fatal("Environment variable AZURE_REGION is not set" + msgCommon)
    }
    region2 := os.Getenv("AZURE_REGION2")
    if region2 == "" {
        t.Fatal("Environment variable AZURE_REGION2 is not set" + msgCommon)
    }
}

func TestAccAviatrixAzurePeer_basic(t *testing.T) {
    var azurePeer goaviatrix.AzurePeer
    vNet1 := os.Getenv("AZURE_VNET_ID")
    vNet2 := os.Getenv("AZURE_VNET_ID2")
    region1 := os.Getenv("AZURE_REGION")
    region2 := os.Getenv("AZURE_REGION2")

    randomSuffix := acctest.RandString(6)
    resourceName := fmt.Sprintf("aviatrix_azure_peer.test_azure_peer_%s", randomSuffix)

    skipAcc := os.Getenv("SKIP_AZURE_PEER")
    if skipAcc == "yes" {
        t.Skip("Skipping Aviatrix Azure peering tests as SKIP_AZURE_PEER is set")
    }
    msgCommon := ". Set SKIP_AZURE_PEER to yes to skip Azure peer tests"

    terraformOptions := &terraform.Options{
        TerraformDir: "./",
        Vars: map[string]interface{}{
            "account_name":        fmt.Sprintf("tf-testing-%s", randomSuffix),
            "cloud_type":          8,
            "arm_subscription_id": os.Getenv("ARM_SUBSCRIPTION_ID"),
            "arm_directory_id":    os.Getenv("ARM_DIRECTORY_ID"),
            "arm_application_id":  os.Getenv("ARM_APPLICATION_ID"),
            "arm_application_key": os.Getenv("ARM_APPLICATION_KEY"),
            "vnet_name_resource_group1": vNet1,
            "vnet_name_resource_group2": vNet2,
            "vnet_reg1": region1,
            "vnet_reg2": region2,
        },
    }

    defer terraform.Destroy(t, terraformOptions)

    terraform.InitAndApply(t, terraformOptions)

    preAzurePeerCheck(t, msgCommon)

    // Verify the Azure peer exists
    peerExists := func() error {
        _, err := goaviatrix.NewClient().GetAzurePeer(&goaviatrix.AzurePeer{
            VNet1: vNet1,
            VNet2: vNet2,
        })
        return err
    }
    assert.Eventually(t, peerExists, "30m", "1m")
}
