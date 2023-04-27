package aviatrix_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/acctest"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

func TestAccDataSourceAviatrixNetworkDomains_basic(t *testing.T) {
	rName := random.UniqueId()
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	tgwName := random.UniqueId()[:5] + random.RandomStringFromSet(5, charset)
	awsSideAsNumber := "64512"
	ndName := random.UniqueId()[:5] + random.RandomStringFromSet(5, charset)
	resourceName := "data.aviatrix_network_domains.test"

	skipAcc := os.Getenv("SKIP_DATA_NETWORK_DOMAINS")
	if skipAcc == "yes" {
		t.Skip("Skipping Data Source All Network Domains tests as SKIP_DATA_NETWORK_DOMAINS is set")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../examples/data-sources/network_domains",
		Vars: map[string]interface{}{
			"rName":             rName,
			"tgwName":           tgwName,
			"awsSideAsNumber":   awsSideAsNumber,
			"ndName":            ndName,
			"awsAccountNumber":  os.Getenv("AWS_ACCOUNT_NUMBER"),
			"awsAccessKey":      os.Getenv("AWS_ACCESS_KEY"),
			"awsSecretKey":      os.Getenv("AWS_SECRET_KEY"),
			"aviatrixVersion":   os.Getenv("AVIATRIX_VERSION"),
			"controllerAccount": os.Getenv("AVIATRIX_CONTROLLER_ACCOUNT"),
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	terraform.Output(t, terraformOptions, "aws_tgw_network_domain_id")

	testResourceExists(t, resourceName, tgwName, ndName)
}

func TestAccDataSourceAviatrixNetworkDomains_basic(t *testing.T) {
	t.Parallel()

	rName := randomUniqueName()
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	tgwName := acctest.RandStringFromCharSet(5, charset) + acctest.RandString(5)
	awsSideAsNumber := "64512"
	ndName := acctest.RandStringFromCharSet(5, charset) + acctest.RandString(5)
	resourceName := "data.aviatrix_network_domains.test"

	skipAcc := os.Getenv("SKIP_DATA_NETWORK_DOMAINS")
	if skipAcc == "yes" {
		t.Skip("Skipping Data Source All Network Domains tests as SKIP_DATA_NETWORK_DOMAINS is set")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../examples/data-sources/network_domains",
		Vars: map[string]interface{}{
			"account_name":             fmt.Sprintf("tfa-%s", rName),
			"aws_account_number":       os.Getenv("AWS_ACCOUNT_NUMBER"),
			"aws_access_key":           os.Getenv("AWS_ACCESS_KEY"),
			"aws_secret_key":           os.Getenv("AWS_SECRET_KEY"),
			"aws_side_as_number":       awsSideAsNumber,
			"aws_region":               "us-west-1",
			"tgw_name":                 tgwName,
			"network_domain_name":      ndName,
			"enable_controller_vpc_dns": false,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	resource.Check(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAviatrixNetworkDomainsConfigBasic(rName, tgwName, awsSideAsNumber, ndName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsTgwNetworkDomainExists("aviatrix_aws_tgw_network_domain.test", tgwName, ndName),
					resource.TestCheckResourceAttr(resourceName, "network_domains.0.tgw_name", tgwName),
					resource.TestCheckResourceAttr(resourceName, "network_domains.3.account", fmt.Sprintf("tfa-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "network_domains.3.name", ndName),
					resource.TestCheckResourceAttr(resourceName, "network_domains.2.cloud_type", "AWS"),
					resource.TestCheckResourceAttr(resourceName, "network_domains.1.region", "us-west-1"),
				),
			},
		},
	})
}

func testAccDataSourceAviatrixNetworkDomainsConfigBasic(rName string, tgwName string, awsSideAsNumber string, ndName string) string {
	return fmt.Sprintf(`
variable "account_name" {}
variable "aws_account_number" {}
variable "aws_access_key" {}
variable "aws_secret_key" {}
variable "aws_side_as_number" {}
variable "aws_region" {}
variable "tgw_name" {}
variable "network_domain_name" {}
variable "enable_controller_vpc_dns" {}

provider "aviatrix" {
  account_name = var.account_name
  aws_account_number = var.aws_account_number
  aws_access_key = var.aws_access_key
  aws_secret_key = var.aws_secret_key
}

provider "aws" {
  region = var.aws_region
}

resource "aviatrix_account" "test" {
	account_name       = var.account_name
	cloud_type         = 1
	aws_account_number = var.aws_account_number
	aws_iam            = false
	aws_access_key     = var.aws_access_key
	aws_secret_key     = var


func testAccDataSourceAviatrixNetworkDomainsConfigBasic(rName string, tgwName string, awsSideAsNumber string, ndName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test" {
	account_name       = "tfa-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}
resource "aviatrix_aws_tgw" "test" {
	account_name       = aviatrix_account.test.account_name
	aws_side_as_number = "%s"
	region             = "us-west-1"
	tgw_name           = "%s"
}
resource "aviatrix_aws_tgw_network_domain" "Default_Domain" {
	name     = "Default_Domain"
	tgw_name = aviatrix_aws_tgw.test.tgw_name
}
resource "aviatrix_aws_tgw_network_domain" "Shared_Service_Domain" {
	name     = "Shared_Service_Domain"
	tgw_name = aviatrix_aws_tgw.test.tgw_name
}
resource "aviatrix_aws_tgw_network_domain" "Aviatrix_Edge_Domain" {
	name     = "Aviatrix_Edge_Domain"
	tgw_name = aviatrix_aws_tgw.test.tgw_name
}
resource "aviatrix_aws_tgw_peering_domain_conn" "default_nd_conn1" {
	tgw_name1    = aviatrix_aws_tgw.test.tgw_name
	domain_name1 = aviatrix_aws_tgw_network_domain.Aviatrix_Edge_Domain.name
	tgw_name2    = aviatrix_aws_tgw.test.tgw_name
	domain_name2 = aviatrix_aws_tgw_network_domain.Default_Domain.name
}
resource "aviatrix_aws_tgw_peering_domain_conn" "default_nd_conn2" {
	tgw_name1    = aviatrix_aws_tgw.test.tgw_name
	domain_name1 = aviatrix_aws_tgw_network_domain.Aviatrix_Edge_Domain.name
	tgw_name2    = aviatrix_aws_tgw.test.tgw_name
	domain_name2 = aviatrix_aws_tgw_network_domain.Shared_Service_Domain.name
}
resource "aviatrix_aws_tgw_peering_domain_conn" "default_nd_conn3" {
	tgw_name1    = aviatrix_aws_tgw.test.tgw_name
	domain_name1 = aviatrix_aws_tgw_network_domain.Default_Domain.name
	tgw_name2    = aviatrix_aws_tgw.test.tgw_name
	domain_name2 = aviatrix_aws_tgw_network_domain.Shared_Service_Domain.name
}
resource "aviatrix_aws_tgw_network_domain" "test" {
	name       = "%s"
	tgw_name   = aviatrix_aws_tgw.test.tgw_name
	depends_on = [
    	aviatrix_aws_tgw_network_domain.Default_Domain,
    	aviatrix_aws_tgw_network_domain.Shared_Service_Domain,
    	aviatrix_aws_tgw_network_domain.Aviatrix_Edge_Domain
  ]
}
data "aviatrix_network_domains" "test"{
	depends_on = [
    	aviatrix_aws_tgw_network_domain.Default_Domain,
    	aviatrix_aws_tgw_network_domain.Shared_Service_Domain,
    	aviatrix_aws_tgw_network_domain.Aviatrix_Edge_Domain,
        aviatrix_aws_tgw_network_domain.test
  ]
}
`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		awsSideAsNumber, tgwName, ndName)
}
