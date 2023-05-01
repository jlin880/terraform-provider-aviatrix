package aviatrix_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/gruntwork-io/terratest/modules/acctest"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAccAviatrixGeoVPN_basic(t *testing.T) {
	var geoVPN goaviatrix.GeoVPN

	rName := random.UniqueId()
	awsRegion := os.Getenv("AWS_REGION")
	awsVpcId := os.Getenv("AWS_VPC_ID")
	awsSubnet := os.Getenv("AWS_SUBNET")
	domainName := os.Getenv("DOMAIN_NAME")

	skipAcc := os.Getenv("SKIP_GEO_VPN")
	if skipAcc == "yes" {
		t.Skip("Skipping Geo VPN test as SKIP_GEO_VPN is set")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../../_examples/aviatrix_geo_vpn",
		Vars: map[string]interface{}{
			"prefix":     fmt.Sprintf("tfa-%s", rName),
			"aws_region": awsRegion,
			"vpc_id":     awsVpcId,
			"subnet":     awsSubnet,
			"domain":     domainName,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	resourceName := "aviatrix_geo_vpn.foo"

	assert.NoError(t, terraform.OutputStruct(resourceName, &geoVPN))

	assert.Equal(t, goaviatrix.AWS, geoVPN.CloudType)
	assert.Equal(t, "vpn", geoVPN.ServiceName)
	assert.Equal(t, domainName, geoVPN.DomainName)
	assert.Len(t, geoVPN.ELBDNSNames, 1)

	expectedElbName := fmt.Sprintf("%s-%s-%s-elb", "tfa", rName, awsRegion)
	assert.Contains(t, geoVPN.ELBDNSNames[0], expectedElbName)
}
