import (
	"testing"
	"github.com/stretchr/testify/require"
)

func TestAccAviatrixControllerPrivateOob_basic(t *testing.T) {
	rName := acctest.RandString(5)

	skipAcc := os.Getenv("SKIP_CONTROLLER_PRIVATE_OOB")
	if skipAcc == "yes" {
		t.Skip("Skipping Controller Private OOB test as SKIP_CONTROLLER_PRIVATE_OOB is set")
	}
	msgCommon := ". Set SKIP_CONTROLLER_PRIVATE_OOB to yes to skip Controller Private OOB tests"
	resourceName := "aviatrix_controller_private_oob.test_private_oob"

	terraformOptions := &terraform.Options{
		TerraformDir: "./",
		Vars: map[string]interface{}{
			"account_name":         "tfa-" + rName,
			"aws_account_number":   os.Getenv("AWS_ACCOUNT_NUMBER"),
			"aws_access_key":       os.Getenv("AWS_ACCESS_KEY"),
			"aws_secret_key":       os.Getenv("AWS_SECRET_KEY"),
			"enable_private_oob":   true,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	require.NoError(t, testAccCheckControllerPrivateOobExists(resourceName))
}

func testAccCheckControllerPrivateOobExists(n string) error {
	client := getApiClientFromEnvVar()

	privateOobState, err := client.GetPrivateOobState()
	if err != nil {
		return fmt.Errorf("could not retrieve controller private oob status")
	}

	if !privateOobState {
		return fmt.Errorf("controller private oob is not enabled")
	}

	return nil
}

func testAccCheckControllerPrivateOobDestroy(t *testing.T, client *goaviatrix.Client) {
	maxRetries := 20
	sleepBetweenRetries := 10 * time.Second

	err := retry.DoWithRetry(context.Background(), "Check Controller Private OOB Destroy", maxRetries, sleepBetweenRetries, func() (string, error) {
		privateOobState, err := client.GetPrivateOobState()
		if err != nil {
			return "", err
		}
		if privateOobState {
			return "", fmt.Errorf("Controller Private OOB is still enabled")
		}
		return "Controller Private OOB destroyed successfully", nil
	})
	if err != nil {
		t.Fatalf("Failed to destroy Controller Private OOB: %v", err)
	}
}