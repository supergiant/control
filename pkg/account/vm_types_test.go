package account

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMachineTypes(t *testing.T) {
	require.NotNil(t, awsMachines.Sizes(), "aws sizes should be set")
}
