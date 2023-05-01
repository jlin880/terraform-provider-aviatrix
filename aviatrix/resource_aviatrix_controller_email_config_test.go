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
	"github.com/gruntwork-io/terratest/modules/terraform"
)

func TestAccAviatrixControllerEmailConfig_basic(t *testing.T) {
	resourceName := "aviatrix_controller_email_config.test"
	rName := random.UniqueId() + "@test.com"
	skipAcc := os.Getenv("SKIP_CONTROLLER_EMAIL_CONFIG")
	if skipAcc == "yes" {
		t.Skip("Skipping Controller Email Config test as SKIP_CONTROLLER_CERT_DOMAIN_CONFIG is set")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../examples/controller_email_config/",
		Vars: map[string]interface{}{
			"admin_alert_email":                   rName,
			"critical_alert_email":                rName,
			"security_event_email":                rName,
			"status_change_email":                 rName,
			"status_change_notification_interval": 20,
		},
	}

	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApply(t, terraformOptions)

	// Import the resource state into Terraform state
	importedResourceOptions := terraform.ImportState{
		TerraformState: terraform.GetState(t, terraformOptions),
		ResourceAddr:   resourceName,
	}
	terraform.Import(t, &importedResourceOptions)

	// Check that the email configuration exists and is correct
	assert.NoError(t, checkControllerEmailConfigExists(t, terraformOptions, rName))
}

func checkControllerEmailConfigExists(t *testing.T, terraformOptions *terraform.Options, rName string) error {
	resourceID := strings.Replace(terraformOptions.Vars["controller_ip"].(string), ".", "-", -1)
	client := goaviatrix.NewClient(terraformOptions.Vars["controller_ip"].(string), terraformOptions.Vars["aviatrix_account_email"].(string), terraformOptions.Vars["aviatrix_account_password"].(string))

	emailConfig, err := client.GetNotificationEmails(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get controller email notification settings: %v", err)
	}

	expectedEmailConfig := &goaviatrix.NotificationEmails{
		AdminAlertEmail:                   rName,
		CriticalAlertEmail:                rName,
		SecurityEventEmail:                rName,
		StatusChangeEmail:                 rName,
		StatusChangeNotificationInterval: 20,
	}

	if emailConfig.AdminAlertEmail != expectedEmailConfig.AdminAlertEmail {
		return fmt.Errorf("admin alert email mismatch")
	}
	if emailConfig.CriticalAlertEmail != expectedEmailConfig.CriticalAlertEmail {
		return fmt.Errorf("critical alert email mismatch")
	}
	if emailConfig.SecurityEventEmail != expectedEmailConfig.SecurityEventEmail {
		return fmt.Errorf("security event email mismatch")
	}
	if emailConfig.StatusChangeEmail != expectedEmailConfig.StatusChangeEmail {
		return fmt.Errorf("status change email mismatch")
	}
	if emailConfig.StatusChangeNotificationInterval != expectedEmailConfig.StatusChangeNotificationInterval {
		return fmt.Errorf("status change notification interval mismatch")
	}
	if resourceID != expectedEmailConfig.StatusChangeEmail {
		return fmt.Errorf("controller email config ID not found")
	}

	return nil
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
