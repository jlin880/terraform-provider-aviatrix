package aviatrix_test

import (
	"fmt"
	"math/rand"
	"os"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestTerraformAviatrixVPNProfile_basic(t *testing.T) {
	t.Parallel()

	profileName := fmt.Sprintf("tf-%s", random.UniqueId())

	skipAcc := os.Getenv("SKIP_VPN_PROFILE")
	if skipAcc == "yes" {
		t.Skip("Skipping VPN Profile test as SKIP_VPN_PROFILE is set")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../examples/aviatrix_vpn_profile/",
		Vars: map[string]interface{}{
			"profile_name":       profileName,
			"aws_account_number": os.Getenv("AWS_ACCOUNT_NUMBER"),
			"aws_access_key":     os.Getenv("AWS_ACCESS_KEY"),
			"aws_secret_key":     os.Getenv("AWS_SECRET_KEY"),
			"aws_vpc_id":         os.Getenv("AWS_VPC_ID"),
			"aws_region":         os.Getenv("AWS_REGION"),
			"aws_subnet":         os.Getenv("AWS_SUBNET"),
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	client := goaviatrix.NewClient(os.Getenv("AVIATRIX_API_KEY"), os.Getenv("AVIATRIX_API_SECRET"), os.Getenv("AVIATRIX_CONTROLLER_IP"))

	vpnProfile, err := client.GetProfile(&goaviatrix.Profile{Name: profileName})
	assert.NoError(t, err)
	assert.Equal(t, vpnProfile.Name, profileName)
	assert.Equal(t, vpnProfile.BaseRule, "allow_all")
	assert.Equal(t, vpnProfile.Users[0], fmt.Sprintf("tfu-%s", profileName))
	assert.Equal(t, vpnProfile.Policy[0].Action, "deny")
	assert.Equal(t, vpnProfile.Policy[0].Proto, "tcp")
	assert.Equal(t, vpnProfile.Policy[0].Port, "443")
	assert.Equal(t, vpnProfile.Policy[0].Target, "10.0.0.0/32")
}


func RandomString(n int) string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}

func TestAccCheckVPNProfileExists(t *testing.T) {
	t.Parallel()

	profileName := fmt.Sprintf("tf-%s", random.UniqueId())

	skipAcc := os.Getenv("SKIP_VPN_PROFILE")
	if skipAcc == "yes" {
		t.Skip("Skipping VPN Profile test as SKIP_VPN_PROFILE is set")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../examples/aviatrix_vpn_profile/",
		Vars: map[string]interface{}{
			"profile_name":       profileName,
			"aws_account_number": os.Getenv("AWS_ACCOUNT_NUMBER"),
			"aws_access_key":     os.Getenv("AWS_ACCESS_KEY"),
			"aws_secret_key":     os.Getenv("AWS_SECRET_KEY"),
			"aws_vpc_id":         os.Getenv("AWS_VPC_ID"),
			"aws_region":         os.Getenv("AWS_REGION"),
			"aws_subnet":         os.Getenv("AWS_SUBNET"),
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	vpnProfile := &goaviatrix.Profile{}

	err := resource.Retry(t, "Check if VPN profile exists", 3, 10*time.Second, func() (string, error) {
		client := goaviatrix.NewClient(os.Getenv("AVIATRIX_API_KEY"), os.Getenv("AVIATRIX_API_SECRET"), os.Getenv("AVIATRIX_CONTROLLER_IP"))

		err := client.GetProfile(vpnProfile)
		if err != nil {
			return "", err
		}

		if vpnProfile.Name != profileName {
			return "", fmt.Errorf("unexpected VPN profile name: got %s, expected %s", vpnProfile.Name, profileName)
		}

		return "VPN profile exists", nil
	})

	if err != nil {
		t.Fatalf("Error checking VPN profile exists: %v", err)
	}
}

func testAccCheckVPNProfileDestroy(t *testing.T, terraformOptions *terraform.Options) {
	t.Parallel()

	skipAcc := os.Getenv("SKIP_VPN_PROFILE")
	if skipAcc == "yes" {
		t.Skip("Skipping VPN Profile test as SKIP_VPN_PROFILE is set")
	}

	vpnProfileName := terraform.Output(t, terraformOptions, "profile_name")
	if vpnProfileName == "" {
		t.Fatal("VPN profile name not found in Terraform output")
	}

	err := terraform.Destroy(t, terraformOptions)
	if err != nil {
		t.Fatalf("Error destroying Terraform: %v", err)
	}

	client := goaviatrix.NewClient(os.Getenv("AVIATRIX_API_KEY"), os.Getenv("AVIATRIX_API_SECRET"), os.Getenv("AVIATRIX_CONTROLLER_IP"))

	foundVPNProfile := &goaviatrix.Profile{Name: vpnProfileName}
	err = client.GetProfile(foundVPNProfile)

	assert.EqualError(t, err, goaviatrix.ErrNotFound.Error(), "VPN profile still exists")
}
