package server

import (
	"testing"

	"github.com/colonyos/colonies/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestSubmitWorkflowSpecSecurity(t *testing.T) {
	env, client, server, _, done := setupTestEnv1(t)

	// The setup looks like this:
	//   executor1 is member of colony1
	//   executor2 is member of colony2

	diamond := generateDiamondtWorkflowSpec(env.colony1ID)

	_, err := client.SubmitWorkflowSpec(diamond, env.executor2PrvKey)
	assert.NotNil(t, err)
	_, err = client.SubmitWorkflowSpec(diamond, env.colony1PrvKey)
	assert.NotNil(t, err)
	_, err = client.SubmitWorkflowSpec(diamond, env.colony2PrvKey)
	assert.NotNil(t, err)
	_, err = client.SubmitWorkflowSpec(diamond, env.executor1PrvKey)
	assert.Nil(t, err)

	server.Shutdown()
	<-done
}

func TestGetProcessGraphSecurity(t *testing.T) {
	env, client, server, _, done := setupTestEnv1(t)

	// The setup looks like this:
	//   executor1 is member of colony1
	//   executor2 is member of colony2

	diamond := generateDiamondtWorkflowSpec(env.colony1ID)
	graph, err := client.SubmitWorkflowSpec(diamond, env.executor1PrvKey)
	assert.Nil(t, err)

	_, err = client.GetProcessGraph(graph.ID, env.executor2PrvKey)
	assert.NotNil(t, err)
	_, err = client.GetProcessGraph(graph.ID, env.colony1PrvKey)
	assert.NotNil(t, err)
	_, err = client.GetProcessGraph(graph.ID, env.colony2PrvKey)
	assert.NotNil(t, err)
	_, err = client.GetProcessGraph(graph.ID, env.executor1PrvKey)
	assert.Nil(t, err)

	server.Shutdown()
	<-done
}

func TestGetProcessGraphsSecurity(t *testing.T) {
	env, client, server, _, done := setupTestEnv1(t)

	// The setup looks like this:
	//   executor1 is member of colony1
	//   executor2 is member of colony2

	_, err := client.GetWaitingProcessGraphs(env.colony1ID, 100, env.executor2PrvKey)
	assert.NotNil(t, err)
	_, err = client.GetWaitingProcessGraphs(env.colony1ID, 100, env.colony1PrvKey)
	assert.NotNil(t, err)
	_, err = client.GetWaitingProcessGraphs(env.colony1ID, 100, env.colony2PrvKey)
	assert.NotNil(t, err)
	_, err = client.GetWaitingProcessGraphs(env.colony1ID, 100, env.executor1PrvKey)
	assert.Nil(t, err)

	server.Shutdown()
	<-done
}

func TestDeleteProcessGraphSecurity(t *testing.T) {
	env, client, server, _, done := setupTestEnv1(t)

	// The setup looks like this:
	//   executor1 is member of colony1
	//   executor2 is member of colony2

	diamond := generateDiamondtWorkflowSpec(env.colony1ID)
	graph, err := client.SubmitWorkflowSpec(diamond, env.executor1PrvKey)
	assert.Nil(t, err)

	err = client.DeleteProcessGraph(graph.ID, env.executor2PrvKey)
	assert.NotNil(t, err)
	err = client.DeleteProcessGraph(graph.ID, env.colony1PrvKey)
	assert.NotNil(t, err)
	err = client.DeleteProcessGraph(graph.ID, env.colony2PrvKey)
	assert.NotNil(t, err)
	err = client.DeleteProcessGraph(graph.ID, env.executor1PrvKey)
	assert.Nil(t, err)

	server.Shutdown()
	<-done
}

func TestDeleteAllProcessGraphsSecurity(t *testing.T) {
	env, client, server, _, done := setupTestEnv1(t)

	// The setup looks like this:
	//   executor1 is member of colony1
	//   executor2 is member of colony2

	diamond := generateDiamondtWorkflowSpec(env.colony1ID)
	_, err := client.SubmitWorkflowSpec(diamond, env.executor1PrvKey)
	assert.Nil(t, err)

	err = client.DeleteAllProcessGraphs(env.colony1ID, env.executor2PrvKey)
	assert.NotNil(t, err)
	err = client.DeleteAllProcessGraphs(env.colony1ID, env.colony1PrvKey)
	assert.Nil(t, err)
	err = client.DeleteAllProcessGraphs(env.colony1ID, env.colony2PrvKey)
	assert.NotNil(t, err)
	err = client.DeleteAllProcessGraphs(env.colony1ID, env.executor1PrvKey)
	assert.NotNil(t, err)

	server.Shutdown()
	<-done
}

func TestAddChildSecurity(t *testing.T) {
	env, client, server, _, done := setupTestEnv1(t)

	executor, executor3PrvKey, err := utils.CreateTestExecutorWithKey(env.colony2ID)
	assert.Nil(t, err)
	executor3, err := client.AddExecutor(executor, env.colony2PrvKey)
	assert.Nil(t, err)
	err = client.ApproveExecutor(executor3.ID, env.colony2PrvKey)
	assert.Nil(t, err)

	// The setup looks like this:
	//   executor1 is member of colony1
	//   executor2 is member of colony2
	//   executor3 is member of colony2

	diamond := generateDiamondtWorkflowSpec(env.colony2ID)
	processGraph, err := client.SubmitWorkflowSpec(diamond, env.executor2PrvKey)
	assert.Nil(t, err)

	parentProcessID := processGraph.Roots[0]

	childFunctionSpec := utils.CreateTestFunctionSpec(env.colony2ID)
	childFunctionSpec.NodeName = "task5"

	_, err = client.AddChild(processGraph.ID, parentProcessID, "", childFunctionSpec, false, env.executor1PrvKey)
	assert.NotNil(t, err) // Error, executor1 not member of member of colony2

	_, err = client.AddChild(processGraph.ID, parentProcessID, "", childFunctionSpec, false, env.colony1PrvKey)
	assert.NotNil(t, err) // Error, invalid prvkey

	_, err = client.AddChild(processGraph.ID, parentProcessID, "", childFunctionSpec, false, env.colony2PrvKey)
	assert.NotNil(t, err) // Error, invalid prvkey

	_, err = client.AddChild(processGraph.ID, parentProcessID, "", childFunctionSpec, false, env.executor2PrvKey)
	assert.NotNil(t, err) // Error, process must be running

	// Assign task1 to executor2
	_, err = client.Assign(env.colony2ID, -1, env.executor2PrvKey)
	assert.Nil(t, err)

	// Now, we should be able to add a child since we got assigned task1
	_, err = client.AddChild(processGraph.ID, parentProcessID, "", childFunctionSpec, false, executor3PrvKey)
	assert.NotNil(t, err) // Error, process is not assigned to executor3

	// But, executor2 should be able to add a child since process is assigned to executor2
	_, err = client.AddChild(processGraph.ID, parentProcessID, "", childFunctionSpec, false, env.executor2PrvKey)
	assert.Nil(t, err)

	server.Shutdown()
	<-done
}
