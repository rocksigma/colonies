package server

import (
	"testing"

	"github.com/colonyos/colonies/pkg/core"
	"github.com/stretchr/testify/assert"
)

func TestVerifyWorkflowSpec(t *testing.T) {
	colonyID := core.GenerateRandomID()

	argsif := make([]interface{}, 1)
	argsif[0] = "arg1"

	funcSpec1 := core.FunctionSpec{
		NodeName:    "gen_task1",
		FuncName:    "gen_test_func",
		Args:        argsif,
		MaxWaitTime: -1,
		MaxExecTime: 2,
		MaxRetries:  10,
		Conditions:  core.Conditions{ColonyID: colonyID, ExecutorType: "bemisexecutor"},
		Env:         make(map[string]string)}

	funcSpec2 := core.FunctionSpec{
		NodeName:    "gen_task2",
		FuncName:    "gen_test_func",
		Args:        argsif,
		MaxWaitTime: -1,
		MaxExecTime: 2,
		MaxRetries:  30,
		Conditions:  core.Conditions{ColonyID: colonyID, ExecutorType: "bemisexecutor"},
		Env:         make(map[string]string)}

	workflowSpec := core.CreateWorkflowSpec(colonyID)
	funcSpec2.AddDependency("task1")
	workflowSpec.AddFunctionSpec(&funcSpec1)
	workflowSpec.AddFunctionSpec(&funcSpec2)

	err := VerifyWorkflowSpec(workflowSpec) // Should not work
	assert.NotNil(t, err)

	funcSpec1 = core.FunctionSpec{
		NodeName:    "gen_task1",
		FuncName:    "gen_test_func",
		Args:        argsif,
		MaxWaitTime: -1,
		MaxExecTime: 2,
		MaxRetries:  10,
		Conditions:  core.Conditions{ColonyID: colonyID, ExecutorType: "bemisexecutor"},
		Env:         make(map[string]string)}

	funcSpec2 = core.FunctionSpec{
		NodeName:    "gen_task2",
		FuncName:    "gen_test_func",
		Args:        argsif,
		MaxWaitTime: -1,
		MaxExecTime: 2,
		MaxRetries:  30,
		Conditions:  core.Conditions{ColonyID: colonyID, ExecutorType: "bemisexecutor"},
		Env:         make(map[string]string)}

	workflowSpec = core.CreateWorkflowSpec(colonyID)
	funcSpec2.AddDependency("gen_task1") // Should work
	workflowSpec.AddFunctionSpec(&funcSpec1)
	workflowSpec.AddFunctionSpec(&funcSpec2)

	err = VerifyWorkflowSpec(workflowSpec) // Should not work
	assert.Nil(t, err)
}
