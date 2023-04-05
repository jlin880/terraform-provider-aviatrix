func TestAccAviatrixVPNUser_basic(t *testing.T) {
    var vpnUser goaviatrix.VPNUser

    rName := acctest.RandString(5)
    resourceName := "aviatrix_vpn_user.test_vpn_user"
    importStateVerifyIgnore := []string{"manage_user_attachment"}

    skipAcc := os.Getenv("SKIP_VPN_USER")
    if skipAcc == "yes" {
        t.Skip("Skipping VPN User test as SKIP_VPN_USER is set")
    }
    msg := ". Set SKIP_VPN_USER to yes to skip VPN User tests"

    resource.Test(t, resource.TestCase{
        PreCheck: func() {
            testAccPreCheck(t)
            preGatewayCheck(t, msg)
        },
        Providers:    testAccProviders,
        CheckDestroy: testAccCheckVPNUserDestroy,
        Steps: []resource.TestStep{
            {
                Config: testAccVPNUserConfigBasic(rName, "false", ""),
                Check: resource.ComposeTestCheckFunc(
                    testAccCheckVPNUserExists("aviatrix_vpn_user.test_vpn_user", &vpnUser),
                    resource.TestCheckResourceAttr(resourceName, "gw_name", fmt.Sprintf("tfl-%s", rName)),
                    resource.TestCheckResourceAttr(resourceName, "vpc_id", os.Getenv("AWS_VPC_ID")),
                    resource.TestCheckResourceAttr(resourceName, "user_email", "user@xyz.com"),
                    resource.TestCheckResourceAttr(resourceName, "user_name", fmt.Sprintf("tfu-%s", rName)),
                ),
            },
            {
                ResourceName:            resourceName,
                ImportState:             true,
                ImportStateVerify:       true,
                ImportStateVerifyIgnore: importStateVerifyIgnore,
            },
        },
    })
}

func testAccVPNUserConfigBasic(rName string, samlEnabled string, endpointName string) string {
    return fmt.Sprintf(`
resource "aviatrix_account" "test_account" {
    account_name       = "tfa-%s"
    cloud_type         = 1
    aws_account_number = "%s"
    aws_iam            = false
    aws_access_key     = "%s"
    aws_secret_key     = "%s"
}
resource "aviatrix_gateway" "test_gw" {
    cloud_type   = 1
    account_name = aviatrix_account.test_account.account_name
    gw_name      = "tfg-%s"
    vpc_id       = "%s"
    vpc_reg      = "%s"
    gw_size      = "t2.micro"
    subnet       = "%s"
    vpn_access   = true
    vpn_cidr     = "192.168.43.0/24"
    max_vpn_conn = "100" 
    enable_elb   = true
    elb_name     = "tfl-%s"
    saml_enabled = "%s"
}
resource "aviatrix_vpn_user" "test_vpn_user" {
    vpc_id        = aviatrix_gateway.test_gw.vpc_id
    gw_name       = aviatrix_gateway.test_gw.elb_name
    user_name     = "tfu-%s"
    user_email    = "user@xyz.com"
	saml_endpoint = "%s"
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		rName, os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_REGION"), os.Getenv("AWS_SUBNET"), rName,
		samlEnabled, rName, endpointName)
}

func testAccCheckVPNUserExists(n string, vpnUser *goaviatrix.VPNUser) error {
	client := goaviatrix.NewClient(
		os.Getenv("AVIATRIX_API_ENDPOINT"),
		os.Getenv("AVIATRIX_USERNAME"),
		os.Getenv("AVIATRIX_PASSWORD"),
		true,
	)

	rs, err := client.GetVPNUserByName(n)
	if err != nil {
		return err
	}

	if rs == nil {
		return fmt.Errorf("VPN User not found: %s", n)
	}

	*vpnUser = *rs

	return nil
}

func testAccCheckVPNUserDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_vpn_user" {
			continue
		}

		foundVPNUser := &goaviatrix.VPNUser{
			UserEmail: rs.Primary.Attributes["user_email"],
			VpcID:     rs.Primary.Attributes["vpc_id"],
			UserName:  rs.Primary.Attributes["user_name"],
			GwName:    rs.Primary.Attributes["gw_name"],
		}

		_, err := client.GetVPNUser(foundVPNUser)
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("VPN User still exists")
		}
	}

	return nil
}
