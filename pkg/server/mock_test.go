package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/colonyos/colonies/pkg/cluster"
	"github.com/colonyos/colonies/pkg/core"
	"github.com/colonyos/colonies/pkg/rpc"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupFakeServer() (*ColoniesServer, *controllerMock, *validatorMock, *dbMock, *gin.Context, *httptest.ResponseRecorder) {
	server := &ColoniesServer{}
	validatorMock := &validatorMock{}
	server.validator = validatorMock
	controllerMock := &controllerMock{}
	server.controller = controllerMock
	dbMock := &dbMock{}
	server.db = dbMock
	ctx, w := getTestGinContext()

	return server, controllerMock, validatorMock, dbMock, ctx, w
}

func createFakeColoniesController() (*coloniesController, *dbMock) {
	node := cluster.Node{Name: "etcd", Host: "localhost", EtcdClientPort: 24100, EtcdPeerPort: 23100, RelayPort: 25100, APIPort: TESTPORT}
	clusterConfig := cluster.Config{}
	clusterConfig.AddNode(node)
	dbMock := &dbMock{}
	return createColoniesController(dbMock, node, clusterConfig, "/tmp/colonies/etcd", GENERATOR_TRIGGER_PERIOD, CRON_TRIGGER_PERIOD, false, -1, 500), dbMock
}

// controllerMock
type controllerMock struct {
	returnError string
	returnValue string
}

func (v *controllerMock) getCronPeriod() int {
	return -1
}

func (v *controllerMock) getGeneratorPeriod() int {
	return -1
}

func (v *controllerMock) getEtcdServer() *cluster.EtcdServer {
	return nil
}

func (v *controllerMock) getEventHandler() *eventHandler {
	return nil
}

func (v *controllerMock) getThisNode() cluster.Node {
	return cluster.Node{}
}

func (v *controllerMock) subscribeProcesses(executorID string, subscription *subscription) error {
	return nil
}

func (v *controllerMock) subscribeProcess(executorID string, subscription *subscription) error {
	return nil
}

func (v *controllerMock) getColonies() ([]*core.Colony, error) {
	return nil, nil
}

func (v *controllerMock) getColony(colonyID string) (*core.Colony, error) {
	return nil, nil
}

func (v *controllerMock) addColony(colony *core.Colony) (*core.Colony, error) {
	return nil, nil
}

func (v *controllerMock) deleteColony(colonyID string) error {
	return nil
}

func (v *controllerMock) renameColony(colonyID string, name string) error {
	return nil
}

func (v *controllerMock) addExecutor(executor *core.Executor, allowExecutorReregister bool) (*core.Executor, error) {
	return nil, nil
}

func (v *controllerMock) getExecutor(executorID string) (*core.Executor, error) {
	return nil, nil
}

func (v *controllerMock) getExecutorByColonyID(colonyID string) ([]*core.Executor, error) {
	return nil, nil
}

func (v *controllerMock) approveExecutor(executorID string) error {
	return nil
}

func (v *controllerMock) rejectExecutor(executorID string) error {
	return nil
}

func (v *controllerMock) deleteExecutor(executorID string) error {
	return nil
}

func (v *controllerMock) addProcessToDB(process *core.Process) (*core.Process, error) {
	return nil, nil
}

func (v *controllerMock) addProcess(process *core.Process) (*core.Process, error) {
	return nil, nil
}

func (v *controllerMock) addChild(processGraphID string, parentProcessID string, childProcessID string, process *core.Process, executorID string, insert bool) (*core.Process, error) {
	return nil, nil
}

func (v *controllerMock) getProcess(processID string) (*core.Process, error) {
	return nil, nil
}

func (v *controllerMock) findProcessHistory(colonyID string, executorID string, seconds int, state int) ([]*core.Process, error) {
	return nil, nil
}

func (v *controllerMock) findWaitingProcesses(colonyID string, executorType string, count int) ([]*core.Process, error) {
	return nil, nil
}

func (v *controllerMock) findRunningProcesses(colonyID string, executorType string, count int) ([]*core.Process, error) {
	return nil, nil
}

func (v *controllerMock) findSuccessfulProcesses(colonyID string, executorType string, count int) ([]*core.Process, error) {
	return nil, nil
}

func (v *controllerMock) findFailedProcesses(colonyID string, executorType string, count int) ([]*core.Process, error) {
	return nil, nil
}

func (v *controllerMock) updateProcessGraph(graph *core.ProcessGraph) error {
	return nil
}

func (v *controllerMock) createProcessGraph(workflowSpec *core.WorkflowSpec, args []interface{}, rootInput []interface{}) (*core.ProcessGraph, error) {
	return nil, nil
}

func (v *controllerMock) submitWorkflowSpec(workflowSpec *core.WorkflowSpec) (*core.ProcessGraph, error) {
	return nil, nil
}

func (v *controllerMock) getProcessGraphByID(processGraphID string) (*core.ProcessGraph, error) {
	return nil, nil
}

func (v *controllerMock) findWaitingProcessGraphs(colonyID string, count int) ([]*core.ProcessGraph, error) {
	return nil, nil
}

func (v *controllerMock) findRunningProcessGraphs(colonyID string, count int) ([]*core.ProcessGraph, error) {
	return nil, nil
}

func (v *controllerMock) findSuccessfulProcessGraphs(colonyID string, count int) ([]*core.ProcessGraph, error) {
	return nil, nil
}

func (v *controllerMock) findFailedProcessGraphs(colonyID string, count int) ([]*core.ProcessGraph, error) {
	return nil, nil
}

func (v *controllerMock) deleteProcess(processID string) error {
	return nil
}

func (v *controllerMock) deleteAllProcesses(colonyID string, state int) error {
	return nil
}

func (v *controllerMock) deleteProcessGraph(processID string) error {
	return nil
}

func (v *controllerMock) deleteAllProcessGraphs(colonyID string, state int) error {
	return nil
}

func (v *controllerMock) closeSuccessful(processID string, executorID string, output []interface{}) error {
	return nil
}

func (v *controllerMock) notifyChildren(process *core.Process) error {
	return nil
}

func (v *controllerMock) closeFailed(processID string, errs []string) error {
	return nil
}

func (v *controllerMock) handleDefunctProcessgraph(processGraphID string, processID string, err error) error {
	return nil
}

func (v *controllerMock) assign(executorID string, colonyID string) (*core.Process, error) {
	return nil, nil
}

func (v *controllerMock) unassignExecutor(processID string) error {
	return nil
}

func (v *controllerMock) resetProcess(processID string) error {
	return nil
}

func (v *controllerMock) getColonyStatistics(colonyID string) (*core.Statistics, error) {
	return nil, nil
}

func (v *controllerMock) getStatistics() (*core.Statistics, error) {
	return nil, nil
}

func (v *controllerMock) addAttribute(attribute *core.Attribute) (*core.Attribute, error) {
	return nil, nil
}

func (v *controllerMock) getAttribute(attributeID string) (*core.Attribute, error) {
	return nil, nil
}

func (v *controllerMock) addFunction(function *core.Function) (*core.Function, error) {
	return nil, nil
}

func (v *controllerMock) getFunctionsByExecutorID(executorID string) ([]*core.Function, error) {
	return nil, nil
}

func (v *controllerMock) getFunctionsByColonyID(colonyID string) ([]*core.Function, error) {
	return nil, nil
}

func (v *controllerMock) getFunctionByID(functionID string) (*core.Function, error) {
	return nil, nil
}

func (v *controllerMock) deleteFunction(functionID string) error {
	return nil
}

func (v *controllerMock) addGenerator(generator *core.Generator) (*core.Generator, error) {
	if v.returnError == "addGenerator" {
		return nil, errors.New("error")
	}

	return nil, nil
}

func (v *controllerMock) getGenerator(generatorID string) (*core.Generator, error) {
	if v.returnError == "getGenerator" {
		return nil, errors.New("error")
	}

	if v.returnValue == "getGenerator" {
		return &core.Generator{}, nil
	}

	return nil, nil
}

func (v *controllerMock) resolveGenerator(generatorName string) (*core.Generator, error) {
	return nil, nil
}

func (v *controllerMock) getGenerators(colonyID string, count int) ([]*core.Generator, error) {
	return nil, nil
}

func (v *controllerMock) packGenerator(generatorID string, colonyID, arg string) error {
	return nil
}

func (v *controllerMock) generatorTriggerLoop() {
}

func (v *controllerMock) triggerGenerators() {
}

func (v *controllerMock) submitWorkflow(generator *core.Generator, counter int) {
}

func (v *controllerMock) addCron(cron *core.Cron) (*core.Cron, error) {
	return nil, nil
}

func (v *controllerMock) deleteGenerator(generatorID string) error {
	return nil
}

func (v *controllerMock) getCron(cronID string) (*core.Cron, error) {
	return nil, nil
}

func (v *controllerMock) getCrons(colonyID string, count int) ([]*core.Cron, error) {
	return nil, nil
}

func (v *controllerMock) runCron(cronID string) (*core.Cron, error) {
	return nil, nil
}

func (v *controllerMock) deleteCron(cronID string) error {
	return nil
}

func (v *controllerMock) calcNextRun(cron *core.Cron) time.Time {
	return time.Time{}
}

func (v *controllerMock) startCron(cron *core.Cron) {
}

func (v *controllerMock) triggerCrons() {
}

func (v *controllerMock) cronTriggerLoop() {
}

func (v *controllerMock) resetDatabase() error {
	return nil
}

func (v *controllerMock) stop() {
}

func (v *controllerMock) isLeader() bool {
	return false
}

func (v *controllerMock) tryBecomeLeader() bool {
	return false
}

func (v *controllerMock) timeoutLoop() {
}

func (v *controllerMock) blockingCmdQueueWorker() {
}

func (v *controllerMock) retentionWorker() {
}

func (v *controllerMock) cmdQueueWorker() {
}

// validatorMock
type validatorMock struct {
}

func (v *validatorMock) RequireServerOwner(recoveredID string, serverID string) error {
	return nil
}

func (v *validatorMock) RequireColonyOwner(recoveredID string, colonyID string) error {
	return nil
}

func (v *validatorMock) RequireExecutorMembership(recoveredID string, colonyID string, approved bool) error {
	return nil
}

type dbMock struct {
	returnError string
	returnValue string
}

func (db *dbMock) Close() {
}

func (db *dbMock) Initialize() error {
	return nil
}

func (db *dbMock) Drop() error {
	return nil
}

func (db *dbMock) AddColony(colony *core.Colony) error {
	if db.returnError == "AddColony" {
		return errors.New("error")
	}
	return nil
}

func (db *dbMock) GetColonies() ([]*core.Colony, error) {
	if db.returnError == "GetColonies" {
		return nil, errors.New("error")
	}

	return nil, nil
}

func (db *dbMock) GetColonyByID(id string) (*core.Colony, error) {
	if db.returnError == "GetColonyByID" {
		return nil, errors.New("error")
	}

	if db.returnValue == "GetColonyByID" {
		return &core.Colony{}, nil
	}

	return nil, nil
}

func (db *dbMock) RenameColony(id string, name string) error {
	if db.returnError == "RenameColony" {
		return errors.New("error")
	}

	return nil
}

func (db *dbMock) DeleteColonyByID(colonyID string) error {
	if db.returnError == "DeleteColonyByID" {
		return errors.New("error")
	}

	return nil
}

func (db *dbMock) CountColonies() (int, error) {
	return -1, nil
}

func (db *dbMock) AddExecutor(executor *core.Executor) error {
	if db.returnError == "AddExecutor" {
		return errors.New("error")
	}

	return nil
}

func (db *dbMock) AddOrReplaceExecutor(executor *core.Executor) error {
	return nil
}

func (db *dbMock) GetExecutors() ([]*core.Executor, error) {
	return nil, nil
}

func (db *dbMock) GetExecutorByID(executorID string) (*core.Executor, error) {
	if db.returnError == "GetExecutorByID" {
		return nil, errors.New("error")
	}

	if db.returnValue == "GetExecutorByID" {
		return &core.Executor{}, nil
	}

	return nil, nil
}

func (db *dbMock) GetExecutorsByColonyID(colonyID string) ([]*core.Executor, error) {
	if db.returnError == "GetExecutorByColonyID" {
		return nil, errors.New("error")
	}

	return nil, nil
}

func (db *dbMock) GetExecutorByName(colonyID string, executorName string) (*core.Executor, error) {
	if db.returnError == "GetExecutorByName" {
		return nil, errors.New("error")
	}

	if db.returnValue == "GetExecutorByName" {
		return &core.Executor{}, nil
	}

	return nil, nil
}

func (db *dbMock) ApproveExecutor(executor *core.Executor) error {
	if db.returnError == "ApproveExecutor" {
		return errors.New("error")
	}

	return nil
}

func (db *dbMock) RejectExecutor(executor *core.Executor) error {
	if db.returnError == "RejectExecutor" {
		return errors.New("error")
	}

	return nil
}

func (db *dbMock) MarkAlive(executor *core.Executor) error {
	return nil
}

func (db *dbMock) DeleteExecutorByID(executorID string) error {
	if db.returnError == "DeleteExecutorByID" {
		return errors.New("error")
	}

	return nil
}

func (db *dbMock) DeleteExecutorsByColonyID(colonyID string) error {
	return nil
}

func (db *dbMock) CountExecutors() (int, error) {
	return -1, nil
}

func (db *dbMock) CountExecutorsByColonyID(colonyID string) (int, error) {
	return -1, nil
}

func (db *dbMock) AddFunction(function *core.Function) error {
	return nil
}

func (db *dbMock) GetFunctionByID(functionID string) (*core.Function, error) {
	return nil, nil
}

func (db *dbMock) GetFunctionsByExecutorID(executorID string) ([]*core.Function, error) {
	return nil, nil
}

func (db *dbMock) GetFunctionsByColonyID(colonyID string) ([]*core.Function, error) {
	return nil, nil
}

func (db *dbMock) GetFunctionsByExecutorIDAndName(executorID string, name string) (*core.Function, error) {
	return nil, nil
}

func (db *dbMock) UpdateFunctionStats(executorID string, name string, counter int, minWaitTime float64, maxWaitTime float64, minExecTime float64, maxExecTime float64, avgWaitTime float64, avgExecTime float64) error {
	return nil
}

func (db *dbMock) DeleteFunctionByID(functionID string) error {
	return nil
}

func (db *dbMock) DeleteFunctionByName(executorID string, name string) error {
	return nil
}

func (db *dbMock) DeleteFunctionsByExecutorID(executorID string) error {
	return nil
}

func (db *dbMock) DeleteFunctionsByColonyID(colonyID string) error {
	return nil
}

func (db *dbMock) DeleteFunctions() error {
	return nil
}

func (db *dbMock) AddProcess(process *core.Process) error {
	if db.returnError == "AddProcess" {
		return errors.New("error")
	}

	return nil
}

func (db *dbMock) GetProcesses() ([]*core.Process, error) {
	return nil, nil
}

func (db *dbMock) GetProcessByID(processID string) (*core.Process, error) {
	if db.returnError == "GetProcessByID" {
		return nil, errors.New("error")
	}
	return nil, nil
}

func (db *dbMock) FindProcessesByColonyID(colonyID string, seconds int, state int) ([]*core.Process, error) {
	return nil, nil
}

func (db *dbMock) FindProcessesByExecutorID(colonyID string, executorID string, seconds int, state int) ([]*core.Process, error) {
	return nil, nil
}

func (db *dbMock) FindWaitingProcesses(colonyID string, executorType string, count int) ([]*core.Process, error) {
	return nil, nil
}

func (db *dbMock) FindRunningProcesses(colonyID string, executorType string, count int) ([]*core.Process, error) {
	return nil, nil
}

func (db *dbMock) FindSuccessfulProcesses(colonyID string, executorType string, count int) ([]*core.Process, error) {
	return nil, nil
}

func (db *dbMock) FindFailedProcesses(colonyID string, executorType string, count int) ([]*core.Process, error) {
	return nil, nil
}

func (db *dbMock) FindAllRunningProcesses() ([]*core.Process, error) {
	return nil, nil
}

func (db *dbMock) FindAllWaitingProcesses() ([]*core.Process, error) {
	return nil, nil
}

func (db *dbMock) FindUnassignedProcesses(colonyID string, executorID string, executorType string, count int) ([]*core.Process, error) {
	return nil, nil
}

func (db *dbMock) DeleteProcessByID(processID string) error {
	return nil
}

func (db *dbMock) DeleteAllProcesses() error {
	return nil
}

func (db *dbMock) DeleteAllWaitingProcessesByColonyID(colonyID string) error {
	return nil
}

func (db *dbMock) DeleteAllRunningProcessesByColonyID(colonyID string) error {
	return nil
}

func (db *dbMock) DeleteAllSuccessfulProcessesByColonyID(colonyID string) error {
	return nil
}

func (db *dbMock) DeleteAllFailedProcessesByColonyID(colonyID string) error {
	return nil
}

func (db *dbMock) DeleteAllProcessesByColonyID(colonyID string) error {
	return nil
}

func (db *dbMock) DeleteAllProcessesByProcessGraphID(processGraphID string) error {
	return nil
}

func (db *dbMock) DeleteAllProcessesInProcessGraphsByColonyID(colonyID string) error {
	return nil
}

func (db *dbMock) ResetProcess(process *core.Process) error {
	return nil
}

func (db *dbMock) SetInput(processID string, output []interface{}) error {
	return nil
}

func (db *dbMock) SetOutput(processID string, output []interface{}) error {
	return nil
}

func (db *dbMock) SetErrors(processID string, errs []string) error {
	return nil
}

func (db *dbMock) SetProcessState(processID string, state int) error {
	return nil
}

func (db *dbMock) SetParents(processID string, parents []string) error {
	return nil
}

func (db *dbMock) SetChildren(processID string, children []string) error {
	return nil
}

func (db *dbMock) SetWaitForParents(processID string, waitingForParent bool) error {
	return nil
}

func (db *dbMock) Assign(executorID string, process *core.Process) error {
	return nil
}

func (db *dbMock) Unassign(process *core.Process) error {
	return nil
}

func (db *dbMock) MarkSuccessful(processID string) (float64, float64, error) {
	return -1.0, -1.0, nil
}

func (db *dbMock) MarkFailed(processID string, errs []string) error {
	return nil
}

func (db *dbMock) CountProcesses() (int, error) {
	return -1, nil
}

func (db *dbMock) CountWaitingProcesses() (int, error) {
	return -1, nil
}

func (db *dbMock) CountRunningProcesses() (int, error) {
	return -1, nil
}

func (db *dbMock) CountSuccessfulProcesses() (int, error) {
	return -1, nil
}

func (db *dbMock) CountFailedProcesses() (int, error) {
	return -1, nil
}

func (db *dbMock) CountWaitingProcessesByColonyID(colonyID string) (int, error) {
	return -1, nil
}

func (db *dbMock) CountRunningProcessesByColonyID(colonyID string) (int, error) {
	return -1, nil
}

func (db *dbMock) CountSuccessfulProcessesByColonyID(colonyID string) (int, error) {
	return -1, nil
}

func (db *dbMock) CountFailedProcessesByColonyID(colonyID string) (int, error) {
	return -1, nil
}

func (db *dbMock) AddAttribute(attribute core.Attribute) error {
	return nil
}

func (db *dbMock) AddAttributes(attribute []core.Attribute) error {
	return nil
}

func (db *dbMock) GetAttributeByID(attributeID string) (core.Attribute, error) {
	return core.Attribute{}, nil
}

func (db *dbMock) GetAttributesByColonyID(colonyID string) ([]core.Attribute, error) {
	return nil, nil
}

func (db *dbMock) GetAttribute(targetID string, key string, attributeType int) (core.Attribute, error) {
	return core.Attribute{}, nil
}

func (db *dbMock) GetAttributes(targetID string) ([]core.Attribute, error) {
	return nil, nil
}

func (db *dbMock) GetAttributesByType(targetID string, attributeType int) ([]core.Attribute, error) {
	return nil, nil
}

func (db *dbMock) UpdateAttribute(attribute core.Attribute) error {
	return nil
}

func (db *dbMock) DeleteAttributeByID(attributeID string) error {
	return nil
}

func (db *dbMock) DeleteAllAttributesByColonyID(colonyID string) error {
	return nil
}

func (db *dbMock) DeleteAllAttributesByColonyIDWithState(colonyID string, state int) error {
	return nil
}

func (db *dbMock) DeleteAllAttributesByProcessGraphID(processGraphID string) error {
	return nil
}

func (db *dbMock) DeleteAllAttributesInProcessGraphsByColonyID(colonyID string) error {
	return nil
}

func (db *dbMock) DeleteAllAttributesInProcessGraphsByColonyIDWithState(colonyID string, state int) error {
	return nil
}

func (db *dbMock) DeleteAttributesByTargetID(targetID string, attributeType int) error {
	return nil
}

func (db *dbMock) DeleteAllAttributesByTargetID(targetID string) error {
	return nil
}

func (db *dbMock) DeleteAllAttributes() error {
	return nil
}

func (db *dbMock) AddProcessGraph(processGraph *core.ProcessGraph) error {
	return nil
}

func (db *dbMock) GetProcessGraphByID(processGraphID string) (*core.ProcessGraph, error) {
	return nil, nil
}

func (db *dbMock) SetProcessGraphState(processGraphID string, state int) error {
	return nil
}

func (db *dbMock) FindWaitingProcessGraphs(colonyID string, count int) ([]*core.ProcessGraph, error) {
	return nil, nil
}

func (db *dbMock) FindRunningProcessGraphs(colonyID string, count int) ([]*core.ProcessGraph, error) {
	return nil, nil
}

func (db *dbMock) FindSuccessfulProcessGraphs(colonyID string, count int) ([]*core.ProcessGraph, error) {
	return nil, nil
}

func (db *dbMock) FindFailedProcessGraphs(colonyID string, count int) ([]*core.ProcessGraph, error) {
	return nil, nil
}

func (db *dbMock) DeleteProcessGraphByID(processGraphID string) error {
	return nil
}

func (db *dbMock) DeleteAllProcessGraphsByColonyID(colonyID string) error {
	return nil
}

func (db *dbMock) DeleteAllWaitingProcessGraphsByColonyID(colonyID string) error {
	return nil
}

func (db *dbMock) DeleteAllRunningProcessGraphsByColonyID(colonyID string) error {
	return nil
}

func (db *dbMock) DeleteAllSuccessfulProcessGraphsByColonyID(colonyID string) error {
	return nil
}

func (db *dbMock) DeleteAllFailedProcessGraphsByColonyID(colonyID string) error {
	return nil
}

func (db *dbMock) CountWaitingProcessGraphs() (int, error) {
	return -1, nil
}

func (db *dbMock) CountRunningProcessGraphs() (int, error) {
	return -1, nil
}

func (db *dbMock) CountSuccessfulProcessGraphs() (int, error) {
	return -1, nil
}

func (db *dbMock) CountFailedProcessGraphs() (int, error) {
	return -1, nil
}

func (db *dbMock) CountWaitingProcessGraphsByColonyID(colonyID string) (int, error) {
	return -1, nil
}

func (db *dbMock) CountRunningProcessGraphsByColonyID(colonyID string) (int, error) {
	return -1, nil
}

func (db *dbMock) CountSuccessfulProcessGraphsByColonyID(colonyID string) (int, error) {
	return -1, nil
}

func (db *dbMock) CountFailedProcessGraphsByColonyID(colonyID string) (int, error) {
	return -1, nil
}

func (db *dbMock) AddGenerator(generator *core.Generator) error {
	return nil
}

func (db *dbMock) SetGeneratorLastRun(generatorID string) error {
	return nil
}

func (db *dbMock) SetGeneratorFirstPack(generatorID string) error {
	return nil
}

func (db *dbMock) GetGeneratorByID(generatorID string) (*core.Generator, error) {
	return nil, nil
}

func (db *dbMock) GetGeneratorByName(name string) (*core.Generator, error) {
	return nil, nil
}

func (db *dbMock) FindGeneratorsByColonyID(colonyID string, count int) ([]*core.Generator, error) {
	return nil, nil
}

func (db *dbMock) FindAllGenerators() ([]*core.Generator, error) {
	return nil, nil
}

func (db *dbMock) DeleteGeneratorByID(generatorID string) error {
	return nil
}

func (db *dbMock) DeleteAllGeneratorsByColonyID(colonyID string) error {
	return nil
}

func (db *dbMock) AddGeneratorArg(generatorArg *core.GeneratorArg) error {
	return nil
}

func (db *dbMock) GetGeneratorArgs(generatorID string, count int) ([]*core.GeneratorArg, error) {
	return nil, nil
}

func (db *dbMock) CountGeneratorArgs(generatorID string) (int, error) {
	if db.returnError == "CountGeneratorArgs" {
		return -1, errors.New("error")

	}
	return -1, nil
}

func (db *dbMock) DeleteGeneratorArgByID(generatorArgsID string) error {
	return nil
}

func (db *dbMock) DeleteAllGeneratorArgsByGeneratorID(generatorID string) error {
	return nil
}

func (db *dbMock) DeleteAllGeneratorArgsByColonyID(generatorID string) error {
	return nil
}

func (db *dbMock) AddCron(cron *core.Cron) error {
	return nil
}

func (db *dbMock) UpdateCron(cronID string, nextRun time.Time, lastRun time.Time, lastProcessGraphID string) error {
	return nil
}

func (db *dbMock) GetCronByID(cronID string) (*core.Cron, error) {

	return nil, nil
}

func (db *dbMock) FindCronsByColonyID(colonyID string, count int) ([]*core.Cron, error) {
	return nil, nil
}

func (db *dbMock) FindAllCrons() ([]*core.Cron, error) {
	return nil, nil
}

func (db *dbMock) DeleteCronByID(cronID string) error {
	return nil
}

func (db *dbMock) DeleteAllCronsByColonyID(colonyID string) error {
	return nil
}

func (db *dbMock) Lock(timeout int) error {

	return nil
}

func (db *dbMock) Unlock() error {
	return nil
}

func (db *dbMock) ApplyRetentionPolicy(retentionPeriod int64) error {

	return nil
}

// gin mockups
func getTestGinContext() (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = &http.Request{
		Header: make(http.Header),
	}

	return ctx, w
}

func assertRPCError(t *testing.T, body string) {
	rpcReplyMsg, err := rpc.CreateRPCReplyMsgFromJSON(body)
	assert.Nil(t, err)
	assert.True(t, rpcReplyMsg.Error)
}
