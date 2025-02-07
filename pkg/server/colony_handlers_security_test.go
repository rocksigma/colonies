package server

import (
	"testing"

	"github.com/colonyos/colonies/pkg/core"
	"github.com/colonyos/colonies/pkg/security/crypto"
	"github.com/stretchr/testify/assert"
)

func TestAddColonySecurity(t *testing.T) {
	client, server, serverPrvKey, done := prepareTests(t)

	crypto := crypto.CreateCrypto()

	privateKey, err := crypto.GeneratePrivateKey()
	assert.Nil(t, err)

	colonyID, err := crypto.GenerateID(privateKey)
	assert.Nil(t, err)

	colony := core.CreateColony(colonyID, "test_colony_name")

	_, err = client.AddColony(colony, "invalid_api_key")
	assert.NotNilf(t, err, "it should be possible to create a colony without correct api key")

	_, err = client.AddColony(colony, serverPrvKey)
	assert.Nil(t, err)

	server.Shutdown()
	<-done
}

func TestDeleteColonySecurity(t *testing.T) {
	client, server, serverPrvKey, done := prepareTests(t)

	crypto := crypto.CreateCrypto()

	privateKey, err := crypto.GeneratePrivateKey()
	assert.Nil(t, err)

	invalidPrivateKey, err := crypto.GeneratePrivateKey()
	assert.Nil(t, err)

	colonyID, err := crypto.GenerateID(privateKey)
	assert.Nil(t, err)

	colony := core.CreateColony(colonyID, "test_colony_name")

	_, err = client.AddColony(colony, serverPrvKey)
	assert.Nil(t, err)

	err = client.DeleteColony(colonyID, invalidPrivateKey)
	assert.NotNil(t, err)

	err = client.DeleteColony(colonyID, serverPrvKey)
	assert.Nil(t, err)

	server.Shutdown()
	<-done
}

func TestRenameColonySecurity(t *testing.T) {
	env, client, server, serverPrvKey, done := setupTestEnv1(t)

	// The setup looks like this:
	//   executor1 is member of colony1
	//   executor2 is member of colony2

	err := client.RenameColony(env.colony1ID, "test_new_name", serverPrvKey)
	assert.NotNil(t, err)

	err = client.RenameColony(env.colony1ID, "test_new_name", env.executor1PrvKey)
	assert.NotNil(t, err)

	err = client.RenameColony(env.colony1ID, "test_new_name", env.executor2PrvKey)
	assert.NotNil(t, err)

	err = client.RenameColony(env.colony1ID, "test_new_name", env.colony2PrvKey)
	assert.NotNil(t, err)

	err = client.RenameColony(env.colony1ID, "test_new_name", env.colony1PrvKey)
	assert.Nil(t, err)

	server.Shutdown()
	<-done
}

func TestGetColoniesSecurity(t *testing.T) {
	_, client, server, serverPrvKey, done := setupTestEnv1(t)

	// The setup looks like this:
	//   executor1 is member of colony1
	//   executor2 is member of colony2

	// Now, try to get colonies info using an invalid api
	_, err := client.GetColonies(core.GenerateRandomID())
	assert.NotNil(t, err) // Should not work

	// Now, try to get colonies info using an invalid api
	_, err = client.GetColonies(serverPrvKey)
	assert.Nil(t, err) // Should work

	server.Shutdown()
	<-done
}

func TestGetColonyByIDSecurity(t *testing.T) {
	env, client, server, _, done := setupTestEnv1(t)

	// The setup looks like this:
	//   executor1 is member of colony1
	//   executor2 is member of colony2

	// Now, try to get colony1 info using colony2 credentials
	_, err := client.GetColonyByID(env.colony1ID, env.colony2PrvKey)
	assert.NotNil(t, err) // Should not work

	// Now, try to get colony1 info using colony1 credentials
	_, err = client.GetColonyByID(env.colony1ID, env.colony1PrvKey)
	assert.NotNil(t, err) // Should work, cannot use colony1PrvKey as credential

	// Now, try to get colony1 info using executor1 credentials
	_, err = client.GetColonyByID(env.colony1ID, env.executor1PrvKey)
	assert.Nil(t, err) // Should work

	// Now, try to get colony1 info using executor1 credentials
	_, err = client.GetColonyByID(env.colony1ID, env.executor2PrvKey)
	assert.NotNil(t, err) // Should not work

	server.Shutdown()
	<-done
}

func TestGetProcessStatSecurity(t *testing.T) {
	env, client, server, _, done := setupTestEnv1(t)

	// The setup looks like this:
	//   executor1 is member of colony1
	//   executor2 is member of colony2

	_, err := client.ColonyStatistics(env.colony1ID, env.executor2PrvKey)
	assert.NotNil(t, err) // Should not work

	_, err = client.ColonyStatistics(env.colony2ID, env.executor1PrvKey)
	assert.NotNil(t, err) // Should not work

	_, err = client.ColonyStatistics(env.colony1ID, env.executor1PrvKey)
	assert.Nil(t, err) // Should work

	_, err = client.ColonyStatistics(env.colony2ID, env.executor2PrvKey)
	assert.Nil(t, err) // Should work

	_, err = client.ColonyStatistics(env.colony1ID, env.colony1PrvKey)
	assert.Nil(t, err) // Should work

	_, err = client.ColonyStatistics(env.colony2ID, env.colony2PrvKey)
	assert.Nil(t, err) // Should work

	_, err = client.ColonyStatistics(env.colony1ID, env.colony2PrvKey)
	assert.NotNil(t, err) // Should not work

	server.Shutdown()
	<-done
}
