package aviatrix

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
)

func TestAccAviatrixVpcMigrateState(t *testing.T) {
	t.Parallel()

	// Create a new instance state with schema version 0
	is := &terraform.InstanceState{
		Version: 0,
		Attributes: map[string]string{
			"aviatrix_firenet_vpc": "true",
		},
	}

	// Migrate the state to the latest version
	migratedState, err := resourceAviatrixVpcMigrateState(1, is, nil)
	if err != nil {
		t.Fatalf("Failed to migrate state: %s", err)
	}

	// Verify that the state has been migrated
	aviatrixFirenetVpc := migratedState.Attributes["aviatrix_firenet_vpc"]
	if aviatrixFirenetVpc != "false" {
		t.Errorf("Expected aviatrix_firenet_vpc to be 'false', got: %s", aviatrixFirenetVpc)
	}
}
