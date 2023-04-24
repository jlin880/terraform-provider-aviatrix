package aviatrix_test

import (
    "context"
    "fmt"
    "os"
    "testing"

    "github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
    "github.com/gruntwork-io/terratest/modules/acctest"
    "github.com/gruntwork-io/terratest/modules/random"
    "github.com/gruntwork-io/terratest/modules/terraform"
    "github.com/stretchr/testify/assert"
)

func TestAccAviatrixAwsTgwConnectPeer_basic(t *testing.T) {
    if os.Getenv("SKIP_AWS_TGW_CONNECT_PEER") == "yes" {
        t.Skip("Skipping Branch Router test as SKIP_AWS_TGW_CONNECT_PEER is set")
    }

    rName := random.UniqueId()
    terraformOptions := &terraform.Options{
        TerraformDir: ".",
        Vars: map[string]interface{}{
            "region":           os.Getenv("AWS_REGION"),
            "resource_name":    "test-aws-tgw-connect-peer",
            "connection_name":  fmt.Sprintf("aws-tgw-connect-%s", rName),
            "connect_peer_name": fmt.Sprintf("connect-peer-%s", rName),
            "peer_as_number":   "65001",
            "peer_gre_address": "172.31.1.11",
            "bgp_inside_cidrs": []string{"169.254.6.0/29"},
            "tgw_gre_address":  "10.0.0.32",
        },
    }

    defer terraform.Destroy(t, terraformOptions)

    terraform.InitAndApply(t, terraformOptions)

    checkAwsTgwConnectPeerAttributes := func() error {
        awsTgwConnectPeer, err := getAwsTgwConnectPeer(t, terraformOptions, rName)
        if err != nil {
            return err
        }
        assert.Equal(t, awsTgwConnectPeer.ConnectionName, terraformOptions.Vars["connection_name"].(string))
        assert.Equal(t, awsTgwConnectPeer.ConnectPeerName, terraformOptions.Vars["connect_peer_name"].(string))
        assert.Equal(t, awsTgwConnectPeer.PeerAsNumber, terraformOptions.Vars["peer_as_number"].(string))
        assert.Equal(t, awsTgwConnectPeer.PeerGreAddress, terraformOptions.Vars["peer_gre_address"].(string))
        assert.Equal(t, awsTgwConnectPeer.BgpInsideCidrs, terraformOptions.Vars["bgp_inside_cidrs"].([]string))
        assert.Equal(t, awsTgwConnectPeer.TgwGreAddress, terraformOptions.Vars["tgw_gre_address"].(string))

        return nil
    }

    assert.NoError(t, checkAwsTgwConnectPeerAttributes())
}
func TestAccAwsTgwConnectPeerBasic(t *testing.T) {
	t.Parallel()

	terraformOptions := &terraform.Options{
		TerraformDir: "../path/to/terraform/config",
		Vars: map[string]interface{}{
			"region": os.Getenv("AWS_REGION"),
			"name":   "test",
		},
	}

	// Run terraform init, apply, and destroy
	terraform.InitAndApply(t, terraformOptions)
	defer terraform.Destroy(t, terraformOptions)

	// Get the connection name from the output
	connectionName := terraform.Output(t, terraformOptions, "connection_name")

	// Create an AWS session
	awsSession, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create a new Aviatrix client
	client, err := goaviatrix.NewClientWithAWS(awsSession)
	if err != nil {
		t.Fatal(err)
	}

	// Get the TGW connect peer
	awsTgwConnectPeer := &goaviatrix.AwsTgwConnectPeer{
		ConnectionName:  connectionName,
		TgwName:         fmt.Sprintf("aws-tgw-%s", terraformOptions.Vars["name"].(string)),
		ConnectPeerName: fmt.Sprintf("connect-peer-%s", terraformOptions.Vars["name"].(string)),
	}
	foundAwsTgwConnectPeer, err := client.GetTGWConnectPeer(context.Background(), awsTgwConnectPeer)
	if err != nil {
		t.Fatal(err)
	}

	// Assert that the TGW connect peer exists
	assert.Equal(t, awsTgwConnectPeer.ConnectionName, foundAwsTgwConnectPeer.ConnectionName)
	assert.Equal(t, awsTgwConnectPeer.TgwName, foundAwsTgwConnectPeer.TgwName)
	assert.Equal(t, awsTgwConnectPeer.ConnectPeerName, foundAwsTgwConnectPeer.ConnectPeerName)
}

func TestAccAwsTgwConnectPeerDestroy(t *testing.T) {
	t.Parallel()

	terraformOptions := &terraform.Options{
		TerraformDir: "../path/to/terraform/config",
		Vars: map[string]interface{}{
			"region": os.Getenv("AWS_REGION"),
			"name":   "test",
		},
	}

	// Run terraform init, apply, and destroy
	terraform.InitAndApply(t, terraformOptions)
	defer terraform.Destroy(t, terraformOptions)

	// Create an AWS session
	awsSession, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create a new Aviatrix client
	client, err := goaviatrix.NewClientWithAWS(awsSession)
	if err != nil {
		t.Fatal(err)
	}

func testAccCheckAwsTgwConnectPeerExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("aws_tgw_connect_peer Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no aws_tgw_connect_peer ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		foundAwsTgwConnectPeer := &goaviatrix.AwsTgwConnectPeer{
			ConnectionName:  rs.Primary.Attributes["connection_name"],
			TgwName:         rs.Primary.Attributes["tgw_name"],
			ConnectPeerName: rs.Primary.Attributes["connect_peer_name"],
		}

		_, err := client.GetTGWConnectPeer(context.Background(), foundAwsTgwConnectPeer)
		if err != nil {
			return err
		}
		if foundAwsTgwConnectPeer.ID() != rs.Primary.ID {
			return fmt.Errorf("aws_tgw_connect_peer not found")
		}

		return nil
	}
}

func testAccCheckAwsTgwConnectPeerDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_aws_tgw_connect_peer" {
			continue
		}
		foundAwsTgwConnectPeer := &goaviatrix.AwsTgwConnectPeer{
			ConnectionName:  rs.Primary.Attributes["connection_name"],
			TgwName:         rs.Primary.Attributes["tgw_name"],
			ConnectPeerName: rs.Primary.Attributes["connect_peer_name"],
		}
		_, err := client.GetTGWConnectPeer(context.Background(), foundAwsTgwConnectPeer)
		if err == nil {
			return fmt.Errorf("aws_tgw_connect_peer still exists")
		}
	}

	return nil
}

func getAwsTgwConnectPeerID(resourceName string) string {
	return fmt.Sprintf("tgw:%s:conn:%s:peer:%s", testAccProvider.Meta().(*goaviatrix.Client).AccountName, resourceName, resourceName)
}