package core

import (
	. "colonies/pkg/utils"
	"testing"
)

func TestCreateColony(t *testing.T) {
	colony, err := CreateColony(GenerateRandomID(), "test_colony_name")
	CheckError(t, err)

	if colony.Name() != "test_colony_name" {
		Fatal(t, "invalid colony name")
	}

	if len(colony.ID()) != 64 {
		Fatal(t, "invalid id")
	}
}
