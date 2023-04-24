package test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/assert"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAviatrixControllerEmailConfig_basic(t *testing.T) {
	resourceName := "aviatrix_controller_email_config.test"
	rName := random.UniqueId() + "@test.com"
	skipAcc := os.Getenv("SKIP_CONTROLLER_EMAIL_CONFIG")
	if skipAcc == "yes" {
		t.Skip("Skipping Controller Email Config test as SKIP_CONTROLLER_CERT_DOMAIN_CONFIG is set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProvidersVersionValidation,
		CheckDestroy: testAccCheckControllerEmailConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccControllerEmailConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControllerEmailConfigExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "admin_alert_email", rName),
					resource.TestCheckResourceAttr(resourceName, "critical_alert_email", rName),
					resource.TestCheckResourceAttr(resourceName, "security_event_email", rName),
					resource.TestCheckResourceAttr(resourceName, "status_change_email", rName),
					resource.TestCheckResourceAttr(resourceName, "status_change_notification_interval", "20"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccControllerEmailConfigBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_controller_email_config" "test" {
	admin_alert_email                   = "%s"
	critical_alert_email                = "%s"
	security_event_email                = "%s"
	status_change_email                 = "%s"
	status_change_notification_interval = 20
}
`, rName, rName, rName, rName)
}

func testAccCheckControllerEmailConfigExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("controller email config ID Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("controller email config ID is not set")
		}

		client := testAccProviderVersionValidation.Meta().(*goaviatrix.Client)

		emailConfig, err := client.GetNotificationEmails(context.Background())
		if err != nil {
			return fmt.Errorf("failed to get controller email notification settings: %v", err)
		}

		if emailConfig.AdminAlertEmail != rs.Primary.Attributes["admin_alert_email"] {
			return fmt.Errorf("admin alert email mismatch")
		}
		if emailConfig.CriticalAlertEmail != rs.Primary.Attributes["critical_alert_email"] {
			return fmt.Errorf("critical alert email mismatch")
		}
		if emailConfig.SecurityEventEmail != rs.Primary.Attributes["security_event_email"] {
			return fmt.Errorf("security event email mismatch")
		}
		if emailConfig.StatusChangeEmail != rs.Primary.Attributes["status_change_email"] {
			return fmt.Errorf("status change email mismatch")
		}
		if emailConfig.StatusChangeNotificationInterval != 20 {
			return fmt.Errorf("status change notification interval mismatch")
		}

func testAccControllerEmailConfigBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_controller_email_config" "test" {
	admin_alert_email                   = "%s"
	critical_alert_email                = "%s"
	security_event_email                = "%s"
	status_change_email                 = "%s"
	status_change_notification_interval = 20
}
`, rName, rName, rName, rName)
}

func testAccCheckControllerEmailConfigExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("controller email config ID Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("controller email config ID is not set")
		}

		client := testAccProviderVersionValidation.Meta().(*goaviatrix.Client)

		_, err := client.GetNotificationEmails(context.Background())
		if err != nil {
			return fmt.Errorf("failed to get controller email notification settings: %v", err)
		}

		if strings.Replace(client.ControllerIP, ".", "-", -1) != rs.Primary.ID {
			return fmt.Errorf("controller email config ID not found")
		}

		return nil
	}
}

func testAccCheckControllerEmailConfigDestroy(s *terraform.State) error {
	client := testAccProviderVersionValidation.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_controller_email_config" {
			continue
		}

		emailConfiguration, _ := client.GetNotificationEmails(context.Background())
		if emailConfiguration.StatusChangeNotificationInterval != 60 {
			return fmt.Errorf("controller email configured when it should be destroyed")
		}
	}

	return nil
}
