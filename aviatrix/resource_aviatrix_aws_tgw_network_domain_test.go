package test

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/avianto/go-metro"
	"github.com/avianto/go-metro/execution"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAccAviatrixAwsTgwNetworkDomain_basic(t *testing.T) {
	rName := strings.ToLower(metro.RandomString(5))
	tgwName := strings.ToLower(metro.RandomString(5))
	awsSideAsNumber := "64512"
	ndName := strings.ToLower(metro.RandomString(5))
	resourceName := fmt.Sprintf("aviatrix_aws_tgw_network_domain.%s", ndName)

	terraformOptions := &terraform.Options{
		TerraformDir: "../examples/aws-tgw-network-domain",
		Vars: map[string]interface{}{
			"prefix":            rName,
			"tgw_name":          tgwName,
			"aws_side_as_number": awsSideAsNumber,
			"nd_name":           ndName,
			"aws_account_number": os.Getenv("AWS_ACCOUNT_NUMBER"),
			"aws_access_key":     os.Getenv("AWS_ACCESS_KEY"),
			"aws_secret_key":     os.Getenv("AWS_SECRET_KEY"),
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	// Check if the network domain exists
	client := getAviatrixClient(t)
	timeout := 5 * time.Minute
	interval := 2 * time.Second

	err := metro.Retry(
		func() error {
			nd := getSecurityDomain(t, client, tgwName, ndName)
			if nd == nil {
				return fmt.Errorf("network domain %s not found", ndName)
			}

			return nil
		},
		timeout,
		interval,
	)

	assert.NoError(t, err)

	// Check the attributes of the resource
	nd := getSecurityDomain(t, client, tgwName, ndName)

	assert.Equal(t, nd.Name, ndName)
	assert.Equal(t, nd.AwsTgwName, tgwName)
}


func testAccAwsTgwNetworkDomainBasic(rName string, tgwName string, awsSideAsNumber string, ndName string) string {
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
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		awsSideAsNumber, tgwName, ndName)
}

func testAccCheckAwsTgwNetworkDomainExists(resourceName string, tgwName string, ndName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource %s not found", resourceName)
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		nd := &goaviatrix.SecurityDomain{
			Name:       ndName,
			AwsTgwName: tgwName,
		}

		_, err := client.GetSecurityDomain(nd)
		if err == goaviatrix.ErrNotFound {
			return fmt.Errorf("network domain %s not found", ndName)
		}

		return nil
	}
}

func testAccCheckAwsTgwNetworkDomainDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_aws_tgw_network_domain" {
			continue
		}

		if rs.Primary.Attributes["name"] == "Default_Domain" || rs.Primary.Attributes["name"] == "Shared_Service_Domain" ||
			rs.Primary.Attributes["name"] == "Aviatrix_Edge_Domain" {
			continue
		}

		awsTgw := &goaviatrix.AWSTgw{
			Name: rs.Primary.Attributes["tgw_name"],
		}

		_, err := client.ListTgwDetails(awsTgw)

		if err != goaviatrix.ErrNotFound {
			nd := &goaviatrix.SecurityDomain{
				Name:       rs.Primary.Attributes["name"],
				AwsTgwName: rs.Primary.Attributes["tgw_name"],
			}

			_, err := client.GetSecurityDomain(nd)
			if err != goaviatrix.ErrNotFound {
				return fmt.Errorf("network domain still exists: %v", err)
			}
		} else {
			break
		}
	}

	return nil
}
