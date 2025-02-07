package server

import (
	"testing"

	"github.com/colonyos/colonies/pkg/core"
	"github.com/colonyos/colonies/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestAddExecutor(t *testing.T) {
	client, server, serverPrvKey, done := prepareTests(t)

	colony, colonyPrvKey, err := utils.CreateTestColonyWithKey()
	_, err = client.AddColony(colony, serverPrvKey)
	assert.Nil(t, err)

	executor, executorPrvKey, err := utils.CreateTestExecutorWithKey(colony.ID)
	assert.Nil(t, err)
	addedExecutor, err := client.AddExecutor(executor, colonyPrvKey)
	assert.Nil(t, err)
	assert.True(t, executor.Equals(addedExecutor))
	err = client.ApproveExecutor(executor.ID, colonyPrvKey)
	assert.Nil(t, err)

	// Just to make the comparison below work, the state will change after it has been approved
	addedExecutor.State = core.APPROVED

	executorFromServer, err := client.GetExecutor(executor.ID, executorPrvKey)
	assert.Nil(t, err)
	assert.NotNil(t, executorFromServer)
	assert.True(t, addedExecutor.Equals(executorFromServer))

	server.Shutdown()
	<-done
}

func TestAddExecutorReRegister(t *testing.T) {
	client, server, serverPrvKey, done := prepareTests(t)

	colony, colonyPrvKey, err := utils.CreateTestColonyWithKey()
	_, err = client.AddColony(colony, serverPrvKey)
	assert.Nil(t, err)

	executor, _, err := utils.CreateTestExecutorWithKey(colony.ID)
	assert.Nil(t, err)
	_, err = client.AddExecutor(executor, colonyPrvKey)
	assert.Nil(t, err)
	assert.Nil(t, err)

	_, err = client.AddExecutor(executor, colonyPrvKey)
	assert.NotNil(t, err)

	server.allowExecutorReregister = true
	_, err = client.AddExecutor(executor, colonyPrvKey)
	assert.Nil(t, err)

	server.Shutdown()
	<-done
}

func TestGetExecutors(t *testing.T) {
	client, server, serverPrvKey, done := prepareTests(t)

	colony, colonyPrvKey, err := utils.CreateTestColonyWithKey()
	assert.Nil(t, err)
	_, err = client.AddColony(colony, serverPrvKey)
	assert.Nil(t, err)

	executor1, executor1PrvKey, err := utils.CreateTestExecutorWithKey(colony.ID)
	_, err = client.AddExecutor(executor1, colonyPrvKey)
	assert.Nil(t, err)
	err = client.ApproveExecutor(executor1.ID, colonyPrvKey)
	assert.Nil(t, err)

	executor2, _, err := utils.CreateTestExecutorWithKey(colony.ID)
	_, err = client.AddExecutor(executor2, colonyPrvKey)
	assert.Nil(t, err)
	err = client.ApproveExecutor(executor2.ID, colonyPrvKey)
	assert.Nil(t, err)

	// Just to make the comparison below work, the state will change after it has been approved
	executor1.State = core.APPROVED
	executor2.State = core.APPROVED

	var executors []*core.Executor
	executors = append(executors, executor1)
	executors = append(executors, executor2)

	executorsFromServer, err := client.GetExecutors(colony.ID, executor1PrvKey)
	assert.Nil(t, err)
	assert.True(t, core.IsExecutorArraysEqual(executors, executorsFromServer))

	server.Shutdown()
	<-done
}

func TestApproveRejectExecutor(t *testing.T) {
	env, client, server, _, done := setupTestEnv2(t)

	// Add an approved eecutor to use for the test below
	approvedExecutor, approvedExecutorPrvKey, err := utils.CreateTestExecutorWithKey(env.colonyID)
	assert.Nil(t, err)
	_, err = client.AddExecutor(approvedExecutor, env.colonyPrvKey)
	assert.Nil(t, err)
	err = client.ApproveExecutor(approvedExecutor.ID, env.colonyPrvKey)
	assert.Nil(t, err)

	testExecutor, _, err := utils.CreateTestExecutorWithKey(env.colonyID)
	assert.Nil(t, err)
	_, err = client.AddExecutor(testExecutor, env.colonyPrvKey)
	assert.Nil(t, err)

	eecutorFromServer, err := client.GetExecutor(testExecutor.ID, approvedExecutorPrvKey)
	assert.Nil(t, err)
	assert.False(t, eecutorFromServer.IsApproved())

	err = client.ApproveExecutor(testExecutor.ID, env.colonyPrvKey)
	assert.Nil(t, err)

	eecutorFromServer, err = client.GetExecutor(testExecutor.ID, approvedExecutorPrvKey)
	assert.Nil(t, err)
	assert.True(t, eecutorFromServer.IsApproved())

	err = client.RejectExecutor(testExecutor.ID, env.colonyPrvKey)
	assert.Nil(t, err)

	eecutorFromServer, err = client.GetExecutor(testExecutor.ID, approvedExecutorPrvKey)
	assert.Nil(t, err)
	assert.False(t, eecutorFromServer.IsApproved())

	server.Shutdown()
	<-done
}

func TestDeleteExecutor(t *testing.T) {
	env, client, server, _, done := setupTestEnv2(t)

	executor, executorPrvKey, err := utils.CreateTestExecutorWithKey(env.colonyID)
	assert.Nil(t, err)
	_, err = client.AddExecutor(executor, env.colonyPrvKey)
	assert.Nil(t, err)
	err = client.ApproveExecutor(executor.ID, env.colonyPrvKey)
	assert.Nil(t, err)

	// Try to get it
	executorFromServer, err := client.GetExecutor(executor.ID, executorPrvKey)
	assert.Nil(t, err)
	assert.NotNil(t, executorFromServer)
	assert.True(t, executor.ID == executorFromServer.ID)

	// Now delete it
	err = client.DeleteExecutor(executor.ID, env.colonyPrvKey)
	assert.Nil(t, err)

	// Try to get it again, it should be gone
	executorFromServer, err = client.GetExecutor(executor.ID, executorPrvKey)
	assert.NotNil(t, err)
	assert.Nil(t, executorFromServer)

	server.Shutdown()
	<-done
}
