package aviatrix

import (
	"fmt"
	"math/rand"
	"os"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestTerraformAviatrixVPNProfile_basic(t *testing.T) {
	t.Parallel()

	rName := fmt.Sprintf("tf-%s", RandomString(5))
	resourceName := fmt.Sprintf("aviatrix_vpn_profile.test_vpn_profile")

	skipAcc := os.Getenv("SKIP_VPN_PROFILE")
	if skipAcc == "yes" {
		t.Skip("Skipping VPN Profile test as SKIP_VPN_PROFILE is set")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../examples/aviatrix_vpn_profile/",
		Vars: map[string]interface{}{
			"res_name":           rName,
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

	vpnProfile, err := client.GetProfile(&goaviatrix.Profile{Name: fmt.Sprintf("tfp-%s", rName)})
	assert.NoError(t, err)
	assert.Equal(t, vpnProfile.Name, fmt.Sprintf("tfp-%s", rName))
	assert.Equal(t, vpnProfile.BaseRule, "allow_all")
	assert.Equal(t, vpnProfile.Users[0], fmt.Sprintf("tfu-%s", rName))
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


func testAccCheckVPNProfileExists(n string, vpnProfile *goaviatrix.Profile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("VPN Profile Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no VPN Profile ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		foundVPNProfile := &goaviatrix.Profile{
			Name: rs.Primary.Attributes["name"],
		}

		foundVPNProfile2, err := client.GetProfile(foundVPNProfile)
		if err != nil {
			return err
		}
		if foundVPNProfile2.Name != rs.Primary.ID {
			return fmt.Errorf("VPN profile not found")
		}

		*vpnProfile = *foundVPNProfile
		return nil
	}
}

func testAccCheckVPNProfileDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_vpn_profile" {
			continue
		}

		foundVPNProfile := &goaviatrix.Profile{
			Name: rs.Primary.Attributes["name"],
		}

		_, err := client.GetProfile(foundVPNProfile)
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("VPN Profile still exists")
		}
	}

	return nil
}
