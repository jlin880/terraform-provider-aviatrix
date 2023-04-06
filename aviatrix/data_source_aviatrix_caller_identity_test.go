package aviatrix

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceAviatrixCallerIdentity_basic(t *testing.T) {
	t.Parallel()

	testCid := os.Getenv("AVIATRIX_CID")
	skipIdentity := os.Getenv("SKIP_DATA_CALLER_IDENTITY")
	if skipIdentity == "true" {
		t.Skip("Skipping Data Source Caller Identity test as SKIP_DATA_CALLER_IDENTITY is set")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "./",
		Vars: map[string]interface{}{
			"cid": testCid,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	output := terraform.Output(t, terraformOptions, "version")
	if !strings.Contains(output, ".") {
		t.Fatalf("Expected version to contain '.' but got %s", output)
	}
}

func TestMain(m *testing.M) {
	testCid := os.Getenv("AVIATRIX_CID")
	if testCid == "" {
		fmt.Println("Missing environment variable AVIATRIX_CID")
		os.Exit(1)
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "./",
		Vars: map[string]interface{}{
			"cid": testCid,
		},
	}

	terraform.InitAndApply(m, terraformOptions)

	exitVal := m.Run()

	defer terraform.Destroy(m, terraformOptions)

	os.Exit(exitVal)
}

func testAccDataSourceAviatrixCallerIdentityConfigBasic(rName string) string {
	return `
data "aviatrix_caller_identity" "foo" {
}
	`
}

func testAccDataSourceAviatrixCallerIdentity(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("root module has no data source called %s", name)
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)
		client.CID = rs.Primary.Attributes["cid"]

		version, _, err := client.GetCurrentVersion()
		if err != nil {
			return fmt.Errorf("valid CID was not returned. Get version API gave the following Error: %v", err)
		}
		if !strings.Contains(version, ".") {
			return fmt.Errorf("valid CID was not returned. Get version API gave the wrong version")
		}

		return nil
	}
}
