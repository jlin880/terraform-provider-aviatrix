package aviatrix_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAviatrixVGWConnMigrateState(t *testing.T) {
	t.Run("Migrate VGW Conn State v0 to v1", func(t *testing.T) {
		// Define the initial state (v0)
		initialState := terraform.InstanceState{
			ID:         "test-vgw-conn",
			Attributes: map[string]string{"enable_advertise_transit_cidr": "true"},
		}

		// Migrate the state
		migratedState, err := resourceAviatrixVGWConnMigrateState(0, &initialState, nil)
		assert.NoError(t, err)

		// Assert the migrated state (v1)
		expectedState := terraform.InstanceState{
			ID:         "test-vgw-conn",
			Attributes: map[string]string{"enable_advertise_transit_cidr": "false"},
		}
		assert.Equal(t, expectedState, *migratedState)
	})

	t.Run("Unsupported schema version", func(t *testing.T) {
		// Define the state with an unsupported schema version
		state := terraform.InstanceState{
			ID:         "test-vgw-conn",
			Attributes: map[string]string{},
		}

		// Migrate the state with an unsupported schema version
		_, err := resourceAviatrixVGWConnMigrateState(2, &state, nil)

		// Assert the error message
		expectedError := "unexpected schema version: 2"
		assert.EqualError(t, err, expectedError)
	})
}
