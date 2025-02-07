package server

import (
	"testing"
	"time"

	"github.com/colonyos/colonies/pkg/core"
	"github.com/colonyos/colonies/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestSubmitSecurity(t *testing.T) {
	env, client, server, _, done := setupTestEnv1(t)

	// The setup looks like this:
	//   executor1 is member of colony1
	//   executor2 is member of colony2

	funcSpec1 := utils.CreateTestFunctionSpec(env.colony1ID)
	_, err := client.Submit(funcSpec1, env.executor1PrvKey)
	assert.Nil(t, err)

	_, err = client.Submit(funcSpec1, env.executor2PrvKey)
	assert.NotNil(t, err) // Should not work, runtiume2 is not member of colony1

	funcSpec2 := utils.CreateTestFunctionSpec(env.colony2ID)
	_, err = client.Submit(funcSpec2, env.executor2PrvKey)
	assert.Nil(t, err)

	_, err = client.Submit(funcSpec2, env.executor1PrvKey)
	assert.NotNil(t, err) // Should not work, runtiume1 is not member of colony2

	server.Shutdown()
	<-done
}

func TestAssignSecurity(t *testing.T) {
	env, client, server, _, done := setupTestEnv1(t)

	// The setup looks like this:
	//   executor1 is member of colony1
	//   executor2 is member of colony2

	funcSpec1 := utils.CreateTestFunctionSpec(env.colony1ID)
	_, err := client.Submit(funcSpec1, env.executor1PrvKey)
	assert.Nil(t, err)

	time.Sleep(50 * time.Millisecond)

	funcSpec2 := utils.CreateTestFunctionSpec(env.colony2ID)
	_, err = client.Submit(funcSpec2, env.executor2PrvKey)
	assert.Nil(t, err)

	// Now try to assign a process from colony2 using executor1 credentials
	_, err = client.Assign(env.colony2ID, -1, env.executor1PrvKey)
	assert.NotNil(t, err) // Should not work

	// Now try to assign a process from colony2 using executor1 credentials
	_, err = client.Assign(env.colony1ID, -1, env.executor2PrvKey)
	assert.NotNil(t, err) // Should not work

	// Now try to assign a process from colony2 using executor1 credentials
	_, err = client.Assign(env.colony1ID, -1, env.executor1PrvKey)
	assert.Nil(t, err) // Should work

	// Now try to assign a process from colony2 using colony1 credentials
	_, err = client.Assign(env.colony1ID, -1, env.colony1PrvKey)
	assert.NotNil(t, err) // Should not work, only executors are allowed

	// Now try to assign a process from colony2 using colony1 credentials
	_, err = client.Assign(env.colony1ID, -1, env.colony2PrvKey)
	assert.NotNil(t, err) // Should not work, only executors are allowed, also invalid credentials are used

	server.Shutdown()
	<-done
}

func TestGetProcessHistForColonySecurity(t *testing.T) {
	env, client, server, _, done := setupTestEnv1(t)

	// The setup looks like this:
	//   executor1 is member of colony1
	//   executor2 is member of colony2

	numberOfRunningProcesses := 2
	for i := 0; i < numberOfRunningProcesses; i++ {
		funcSpec := utils.CreateTestFunctionSpec(env.colony1ID)
		_, err := client.Submit(funcSpec, env.executor1PrvKey)
		assert.Nil(t, err)
		_, err = client.Assign(env.colony1ID, -1, env.executor1PrvKey)
		assert.Nil(t, err)
	}

	_, err := client.GetProcessHistForColony(core.RUNNING, env.colony1ID, 60, env.executor2PrvKey)
	assert.NotNil(t, err) // Should not work

	_, err = client.GetProcessHistForColony(core.RUNNING, env.colony1ID, 60, env.colony2PrvKey)
	assert.NotNil(t, err) // Should not work

	_, err = client.GetProcessHistForColony(core.RUNNING, env.colony1ID, 60, env.executor1PrvKey)
	assert.Nil(t, err) // Should work

	_, err = client.GetProcessHistForColony(core.RUNNING, env.colony1ID, 60, env.colony1PrvKey)
	assert.Nil(t, err) // Should work

	server.Shutdown()
	<-done
}

func TestGetProcessHistForExecutorSecurity(t *testing.T) {
	env, client, server, _, done := setupTestEnv1(t)

	// The setup looks like this:
	//   executor1 is member of colony1
	//   executor2 is member of colony2

	numberOfRunningProcesses := 2
	for i := 0; i < numberOfRunningProcesses; i++ {
		funcSpec := utils.CreateTestFunctionSpec(env.colony1ID)
		_, err := client.Submit(funcSpec, env.executor1PrvKey)
		assert.Nil(t, err)
		_, err = client.Assign(env.colony1ID, -1, env.executor1PrvKey)
		assert.Nil(t, err)
	}

	_, err := client.GetProcessHistForExecutor(core.RUNNING, env.colony1ID, env.executor1ID, 60, env.executor2PrvKey)
	assert.NotNil(t, err) // Should not work

	_, err = client.GetProcessHistForExecutor(core.RUNNING, env.colony1ID, env.executor1ID, 60, env.colony2PrvKey)
	assert.NotNil(t, err) // Should not work

	_, err = client.GetProcessHistForExecutor(core.RUNNING, env.colony1ID, env.executor1ID, 60, env.executor1PrvKey)
	assert.Nil(t, err) // Should work

	_, err = client.GetProcessHistForExecutor(core.RUNNING, env.colony1ID, env.executor1ID, 60, env.colony1PrvKey)
	assert.Nil(t, err) // Should work

	server.Shutdown()
	<-done
}

func TestGetWaitingProcessesSecurity(t *testing.T) {
	env, client, server, _, done := setupTestEnv1(t)

	// The setup looks like this:
	//   executor1 is member of colony1
	//   executor2 is member of colony2

	numberOfRunningProcesses := 2
	for i := 0; i < numberOfRunningProcesses; i++ {
		funcSpec := utils.CreateTestFunctionSpec(env.colony1ID)
		_, err := client.Submit(funcSpec, env.executor1PrvKey)
		assert.Nil(t, err)
		_, err = client.Assign(env.colony1ID, -1, env.executor1PrvKey)
		assert.Nil(t, err)
	}

	_, err := client.GetRunningProcesses(env.colony1ID, "", numberOfRunningProcesses, env.executor2PrvKey)
	assert.NotNil(t, err) // Should not work

	_, err = client.GetRunningProcesses(env.colony1ID, "", numberOfRunningProcesses, env.executor1PrvKey)
	assert.Nil(t, err) // Should work

	server.Shutdown()
	<-done
}

func TestGetRunningProcessesSecurity(t *testing.T) {
	env, client, server, _, done := setupTestEnv1(t)

	// The setup looks like this:
	//   executor1 is member of colony1
	//   executor2 is member of colony2

	numberOfRunningProcesses := 2
	for i := 0; i < numberOfRunningProcesses; i++ {
		funcSpec := utils.CreateTestFunctionSpec(env.colony1ID)
		_, err := client.Submit(funcSpec, env.executor1PrvKey)
		assert.Nil(t, err)
	}

	_, err := client.GetWaitingProcesses(env.colony1ID, "", numberOfRunningProcesses, env.executor2PrvKey)
	assert.NotNil(t, err) // Should not work

	_, err = client.GetWaitingProcesses(env.colony1ID, "", numberOfRunningProcesses, env.executor1PrvKey)
	assert.Nil(t, err) // Should work

	server.Shutdown()
	<-done
}

func TestGetSuccessfulProcessesSecurity(t *testing.T) {
	env, client, server, _, done := setupTestEnv1(t)

	// The setup looks like this:
	//   executor1 is member of colony1
	//   executor2 is member of colony2

	numberOfRunningProcesses := 2
	for i := 0; i < numberOfRunningProcesses; i++ {
		funcSpec := utils.CreateTestFunctionSpec(env.colony1ID)
		_, err := client.Submit(funcSpec, env.executor1PrvKey)
		assert.Nil(t, err)
		processFromServer, err := client.Assign(env.colony1ID, -1, env.executor1PrvKey)
		assert.Nil(t, err)
		err = client.Close(processFromServer.ID, env.executor1PrvKey)
		assert.Nil(t, err)
	}

	_, err := client.GetSuccessfulProcesses(env.colony1ID, "", numberOfRunningProcesses, env.executor2PrvKey)
	assert.NotNil(t, err) // Should not work

	_, err = client.GetSuccessfulProcesses(env.colony1ID, "", numberOfRunningProcesses, env.executor1PrvKey)
	assert.Nil(t, err) // Should work

	server.Shutdown()
	<-done
}

func TestGetFailedProcessesSecurity(t *testing.T) {
	env, client, server, _, done := setupTestEnv1(t)

	// The setup looks like this:
	//   executor1 is member of colony1
	//   executor2 is member of colony2

	numberOfRunningProcesses := 2
	for i := 0; i < numberOfRunningProcesses; i++ {
		funcSpec := utils.CreateTestFunctionSpec(env.colony1ID)
		_, err := client.Submit(funcSpec, env.executor1PrvKey)
		assert.Nil(t, err)
		processFromServer, err := client.Assign(env.colony1ID, -1, env.executor1PrvKey)
		assert.Nil(t, err)
		err = client.Fail(processFromServer.ID, []string{"error"}, env.executor1PrvKey)
		assert.Nil(t, err)
	}

	_, err := client.GetFailedProcesses(env.colony1ID, "", numberOfRunningProcesses, env.executor2PrvKey)
	assert.NotNil(t, err) // Should not work

	_, err = client.GetFailedProcesses(env.colony1ID, "", numberOfRunningProcesses, env.executor1PrvKey)
	assert.Nil(t, err) // Should work

	server.Shutdown()
	<-done
}

func TestGetProcessSecurity(t *testing.T) {
	env, client, server, _, done := setupTestEnv1(t)

	// The setup looks like this:
	//   executor1 is member of colony1
	//   executor2 is member of colony2

	funcSpec := utils.CreateTestFunctionSpec(env.colony1ID)
	addedProcess, err := client.Submit(funcSpec, env.executor1PrvKey)
	assert.Nil(t, err)

	_, err = client.GetProcess(addedProcess.ID, env.executor2PrvKey)
	assert.NotNil(t, err) // Should not work

	_, err = client.GetProcess(addedProcess.ID, env.executor1PrvKey)
	assert.Nil(t, err) // Should work

	server.Shutdown()
	<-done
}

func TestDeleteProcessSecurity(t *testing.T) {
	env, client, server, _, done := setupTestEnv1(t)

	// The setup looks like this:
	//   executor1 is member of colony1
	//   executor2 is member of colony2

	funcSpec := utils.CreateTestFunctionSpec(env.colony1ID)
	addedProcess, err := client.Submit(funcSpec, env.executor1PrvKey)
	assert.Nil(t, err)

	err = client.DeleteProcess(addedProcess.ID, env.executor2PrvKey)
	assert.NotNil(t, err) // Should not work

	err = client.DeleteProcess(addedProcess.ID, env.colony1PrvKey)
	assert.NotNil(t, err) // Should not work

	err = client.DeleteProcess(addedProcess.ID, env.colony2PrvKey)
	assert.NotNil(t, err) // Should not work

	err = client.DeleteProcess(addedProcess.ID, env.executor1PrvKey)
	assert.Nil(t, err) // Should work

	server.Shutdown()
	<-done
}

func TestDeleteAllProcessSecurity(t *testing.T) {
	env, client, server, _, done := setupTestEnv1(t)

	// The setup looks like this:
	//   executor1 is member of colony1
	//   executor2 is member of colony2

	funcSpec := utils.CreateTestFunctionSpec(env.colony1ID)
	_, err := client.Submit(funcSpec, env.executor1PrvKey)
	assert.Nil(t, err)

	err = client.DeleteAllProcesses(env.colony1ID, env.executor2PrvKey)
	assert.NotNil(t, err) // Should not work

	err = client.DeleteAllProcesses(env.colony1ID, env.executor1PrvKey)
	assert.NotNil(t, err) // Should not work

	err = client.DeleteAllProcesses(env.colony1ID, env.colony2PrvKey)
	assert.NotNil(t, err) // Should not work

	err = client.DeleteAllProcesses(env.colony1ID, env.colony1PrvKey)
	assert.Nil(t, err) // Shoul dwork

	server.Shutdown()
	<-done
}

func TestCloseSuccessfulSecurity(t *testing.T) {
	env, client, server, _, done := setupTestEnv1(t)

	// The setup looks like this:
	//   executor1 is member of colony1
	//   executor2 is member of colony2

	funcSpec := utils.CreateTestFunctionSpec(env.colony1ID)
	_, err := client.Submit(funcSpec, env.executor1PrvKey)
	assert.Nil(t, err)
	processFromServer, err := client.Assign(env.colony1ID, -1, env.executor1PrvKey)
	assert.Nil(t, err)

	err = client.Close(processFromServer.ID, env.executor2PrvKey)
	assert.NotNil(t, err) // Should not work

	err = client.Close(processFromServer.ID, env.executor1PrvKey)
	assert.Nil(t, err) // Should work

	// Add another executor to colony1 and try to close the process statred by executor1, it should not be possible
	executor3, executor3PrvKey, err := utils.CreateTestExecutorWithKey(env.colony1ID)
	assert.Nil(t, err)
	_, err = client.AddExecutor(executor3, env.colony1PrvKey)
	assert.Nil(t, err)
	err = client.ApproveExecutor(executor3.ID, env.colony1PrvKey)
	assert.Nil(t, err)
	err = client.Close(processFromServer.ID, executor3PrvKey)
	assert.NotNil(t, err) // Should work

	server.Shutdown()
	<-done
}

func TestCloseFailedSecurity(t *testing.T) {
	env, client, server, _, done := setupTestEnv1(t)

	// The setup looks like this:
	//   executor1 is member of colony1
	//   executor2 is member of colony2

	funcSpec := utils.CreateTestFunctionSpec(env.colony1ID)
	_, err := client.Submit(funcSpec, env.executor1PrvKey)
	assert.Nil(t, err)
	processFromServer, err := client.Assign(env.colony1ID, -1, env.executor1PrvKey)
	assert.Nil(t, err)

	err = client.Fail(processFromServer.ID, []string{"error"}, env.executor2PrvKey)
	assert.NotNil(t, err) // Should not work

	err = client.Fail(processFromServer.ID, []string{"error"}, env.executor1PrvKey)
	assert.Nil(t, err) // Should work

	// Add another executor to colony1 and try to close the process started by executor1, it should not be possible
	executor3, executor3PrvKey, err := utils.CreateTestExecutorWithKey(env.colony1ID)
	assert.Nil(t, err)
	_, err = client.AddExecutor(executor3, env.colony1PrvKey)
	assert.Nil(t, err)
	err = client.ApproveExecutor(executor3.ID, env.colony1PrvKey)
	assert.Nil(t, err)
	err = client.Fail(processFromServer.ID, []string{"error"}, executor3PrvKey)
	assert.NotNil(t, err) // Should work

	server.Shutdown()
	<-done
}
