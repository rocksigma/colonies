package postgresql

import (
	"testing"

	"github.com/colonyos/colonies/pkg/core"
	"github.com/colonyos/colonies/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestGeneratorClosedDB(t *testing.T) {
	db, err := PrepareTests()
	assert.Nil(t, err)

	db.Close()

	generator := utils.FakeGenerator(t, core.GenerateRandomID())
	generator.ID = core.GenerateRandomID()
	err = db.AddGenerator(generator)
	assert.NotNil(t, err)

	err = db.SetGeneratorLastRun("invalid_id")
	assert.NotNil(t, err)

	err = db.SetGeneratorFirstPack("invalid_id")
	assert.NotNil(t, err)

	_, err = db.GetGeneratorByID("invalid_id")
	assert.NotNil(t, err)

	_, err = db.GetGeneratorByName("invalid_name")
	assert.NotNil(t, err)

	_, err = db.FindGeneratorsByColonyID("invalid_id", 100)
	assert.NotNil(t, err)

	_, err = db.FindAllGenerators()
	assert.NotNil(t, err)

	err = db.DeleteGeneratorByID("invalid_id")
	assert.NotNil(t, err)

	err = db.DeleteAllGeneratorsByColonyID("invalid_id")
	assert.NotNil(t, err)
}

func TestAddGenerator(t *testing.T) {
	db, err := PrepareTests()
	assert.Nil(t, err)

	defer db.Close()

	generator := utils.FakeGenerator(t, core.GenerateRandomID())
	generator.ID = core.GenerateRandomID()
	err = db.AddGenerator(generator)
	assert.Nil(t, err)
}

func TestGetGeneratorByID(t *testing.T) {
	db, err := PrepareTests()
	assert.Nil(t, err)

	defer db.Close()

	generator := utils.FakeGenerator(t, core.GenerateRandomID())
	generator.ID = core.GenerateRandomID()
	err = db.AddGenerator(generator)
	assert.Nil(t, err)

	generatorFromDB, err := db.GetGeneratorByID("invalid_id")
	assert.Nil(t, err)
	assert.Nil(t, generatorFromDB)

	generatorFromDB, err = db.GetGeneratorByID(generator.ID)
	assert.Nil(t, err)
	assert.True(t, generator.Equals(generatorFromDB))
}

func TestGetGeneratorByName(t *testing.T) {
	db, err := PrepareTests()
	assert.Nil(t, err)

	defer db.Close()

	generator := utils.FakeGenerator(t, core.GenerateRandomID())
	generator.ID = core.GenerateRandomID()
	generator.Name = "test_name"
	err = db.AddGenerator(generator)
	assert.Nil(t, err)

	generatorFromDB, err := db.GetGeneratorByName("invalid_name")
	assert.Nil(t, err)
	assert.Nil(t, generatorFromDB)

	generatorFromDB, err = db.GetGeneratorByName("test_name")
	assert.Nil(t, err)
	assert.True(t, generator.Equals(generatorFromDB))
}

func TestSetGeneratorLastRun(t *testing.T) {
	db, err := PrepareTests()
	assert.Nil(t, err)

	defer db.Close()

	generator := utils.FakeGenerator(t, core.GenerateRandomID())
	generator.ID = core.GenerateRandomID()
	err = db.AddGenerator(generator)
	assert.Nil(t, err)

	generatorFromDB, err := db.GetGeneratorByID(generator.ID)
	assert.Nil(t, err)
	assert.True(t, generator.Equals(generatorFromDB))

	lastRun := generatorFromDB.LastRun.Unix()

	err = db.SetGeneratorLastRun(generator.ID)
	assert.Nil(t, err)

	generatorFromDB, err = db.GetGeneratorByID(generator.ID)
	assert.Nil(t, err)

	assert.Greater(t, generatorFromDB.LastRun.Unix(), lastRun)
}

func TestSetGeneratorFirstPack(t *testing.T) {
	db, err := PrepareTests()
	assert.Nil(t, err)

	defer db.Close()

	generator := utils.FakeGenerator(t, core.GenerateRandomID())
	generator.ID = core.GenerateRandomID()
	err = db.AddGenerator(generator)
	assert.Nil(t, err)

	generatorFromDB, err := db.GetGeneratorByID(generator.ID)
	assert.Nil(t, err)
	assert.True(t, generator.Equals(generatorFromDB))

	err = db.SetGeneratorFirstPack(generator.ID)
	assert.Nil(t, err)

	generatorFromDB, err = db.GetGeneratorByID(generator.ID)
	assert.Nil(t, err)

	assert.True(t, generatorFromDB.FirstPack.Unix() > 0)
}

func TestFindGeneratorsByColonyID(t *testing.T) {
	db, err := PrepareTests()
	assert.Nil(t, err)

	defer db.Close()

	colonyID := core.GenerateRandomID()
	generator1 := utils.FakeGenerator(t, colonyID)
	generator1.ID = core.GenerateRandomID()
	err = db.AddGenerator(generator1)
	assert.Nil(t, err)

	generator2 := utils.FakeGenerator(t, colonyID)
	generator2.ID = core.GenerateRandomID()
	err = db.AddGenerator(generator2)
	assert.Nil(t, err)

	generatorsFromDB, err := db.FindGeneratorsByColonyID(colonyID, 100)
	assert.Nil(t, err)
	assert.Len(t, generatorsFromDB, 2)

	count := 0
	for _, generator := range generatorsFromDB {
		if generator.ID == generator1.ID {
			count++
		}
		if generator.ID == generator2.ID {
			count++
		}
	}
	assert.True(t, count == 2)
}

func TestFindAllGenerators(t *testing.T) {
	db, err := PrepareTests()
	assert.Nil(t, err)

	defer db.Close()

	colonyID1 := core.GenerateRandomID()
	generator1 := utils.FakeGenerator(t, colonyID1)
	generator1.ID = core.GenerateRandomID()
	err = db.AddGenerator(generator1)
	assert.Nil(t, err)

	colonyID2 := core.GenerateRandomID()
	generator2 := utils.FakeGenerator(t, colonyID2)
	generator2.ID = core.GenerateRandomID()
	err = db.AddGenerator(generator2)
	assert.Nil(t, err)

	generatorsFromDB, err := db.FindAllGenerators()
	assert.Nil(t, err)
	assert.Len(t, generatorsFromDB, 2)
}

func TestDeleteGeneratorByID(t *testing.T) {
	db, err := PrepareTests()
	assert.Nil(t, err)

	defer db.Close()

	colonyID := core.GenerateRandomID()
	generator1 := utils.FakeGenerator(t, colonyID)
	generator1.ID = core.GenerateRandomID()
	err = db.AddGenerator(generator1)
	assert.Nil(t, err)

	generator2 := utils.FakeGenerator(t, colonyID)
	generator2.ID = core.GenerateRandomID()
	err = db.AddGenerator(generator2)
	assert.Nil(t, err)

	generatorFromDB, err := db.GetGeneratorByID(generator1.ID)
	assert.Nil(t, err)
	assert.NotNil(t, generatorFromDB)

	generatorArg := core.CreateGeneratorArg(generator1.ID, colonyID, "arg")
	err = db.AddGeneratorArg(generatorArg)
	assert.Nil(t, err)

	count, err := db.CountGeneratorArgs(generator1.ID)
	assert.Nil(t, err)
	assert.Equal(t, count, 1)

	err = db.DeleteGeneratorByID(generator1.ID)
	assert.Nil(t, err)

	generatorFromDB, err = db.GetGeneratorByID(generator1.ID)
	assert.Nil(t, err)
	assert.Nil(t, generatorFromDB)

	generatorFromDB, err = db.GetGeneratorByID(generator2.ID)
	assert.Nil(t, err)
	assert.NotNil(t, generatorFromDB)

	count, err = db.CountGeneratorArgs(generator1.ID)
	assert.Nil(t, err)
	assert.Equal(t, count, 0)
}

func TestDeleteAllGeneratorsByColonyID(t *testing.T) {
	db, err := PrepareTests()
	assert.Nil(t, err)

	defer db.Close()

	colonyID1 := core.GenerateRandomID()
	generator1 := utils.FakeGenerator(t, colonyID1)
	generator1.ID = core.GenerateRandomID()
	err = db.AddGenerator(generator1)
	assert.Nil(t, err)

	generator2 := utils.FakeGenerator(t, colonyID1)
	generator2.ID = core.GenerateRandomID()
	err = db.AddGenerator(generator2)
	assert.Nil(t, err)

	colonyID2 := core.GenerateRandomID()
	generator3 := utils.FakeGenerator(t, colonyID2)
	err = db.AddGenerator(generator3)
	assert.Nil(t, err)

	generatorArg := core.CreateGeneratorArg(generator1.ID, colonyID1, "arg")
	err = db.AddGeneratorArg(generatorArg)
	assert.Nil(t, err)
	generatorArg = core.CreateGeneratorArg(generator2.ID, colonyID1, "arg")
	err = db.AddGeneratorArg(generatorArg)
	assert.Nil(t, err)
	generatorArg = core.CreateGeneratorArg(generator3.ID, colonyID2, "arg")
	err = db.AddGeneratorArg(generatorArg)
	assert.Nil(t, err)

	count, err := db.CountGeneratorArgs(generator1.ID)
	assert.Nil(t, err)
	assert.Equal(t, count, 1)

	generatorFromDB, err := db.GetGeneratorByID(generator1.ID)
	assert.Nil(t, err)
	assert.NotNil(t, generatorFromDB)

	err = db.DeleteAllGeneratorsByColonyID(colonyID1)
	assert.Nil(t, err)

	generatorFromDB, err = db.GetGeneratorByID(generator1.ID)
	assert.Nil(t, err)
	assert.Nil(t, generatorFromDB)

	generatorFromDB, err = db.GetGeneratorByID(generator2.ID)
	assert.Nil(t, err)
	assert.Nil(t, generatorFromDB)

	generatorFromDB, err = db.GetGeneratorByID(generator3.ID)
	assert.Nil(t, err)
	assert.NotNil(t, generatorFromDB)

	count, err = db.CountGeneratorArgs(generator1.ID)
	assert.Nil(t, err)
	assert.Equal(t, count, 0)

	count, err = db.CountGeneratorArgs(generator2.ID)
	assert.Nil(t, err)
	assert.Equal(t, count, 0)

	count, err = db.CountGeneratorArgs(generator3.ID)
	assert.Nil(t, err)
	assert.Equal(t, count, 1)
}
