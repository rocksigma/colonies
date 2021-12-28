package core

import (
	"colonies/pkg/crypto"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

const (
	WAITING int = 0
	RUNNING     = 1
	SUCCESS     = 2
	FAILED      = 3
)

type Process struct {
	ID                string       `json:"processid"`
	AssignedRuntimeID string       `json:"assignedruntimeid"`
	IsAssigned        bool         `json:"isassigned"`
	Status            int          `json:"status"`
	SubmissionTime    time.Time    `json:"submissiontime"`
	StartTime         time.Time    `json:"starttime"`
	EndTime           time.Time    `json:"endtime"`
	Deadline          time.Time    `json:"deadline"`
	Retries           int          `json:"retries"`
	Attributes        []*Attribute `json:"attributes"`
	ProcessSpec       *ProcessSpec `json:"spec"`
}

func CreateProcess(processSpec *ProcessSpec) *Process {
	uuid := uuid.New()
	id := crypto.GenerateHashFromString(uuid.String()).String()

	var attributes []*Attribute

	process := &Process{ID: id,
		Status:      WAITING,
		IsAssigned:  false,
		Attributes:  attributes,
		ProcessSpec: processSpec,
	}

	return process
}

func CreateProcessFromDB(processSpec *ProcessSpec,
	id string,
	assignedRuntimeID string,
	isAssigned bool,
	status int,
	submissionTime time.Time,
	startTime time.Time,
	endTime time.Time,
	deadline time.Time,
	retries int,
	attributes []*Attribute) *Process {
	return &Process{ID: id,
		AssignedRuntimeID: assignedRuntimeID,
		IsAssigned:        isAssigned,
		Status:            status,
		SubmissionTime:    submissionTime,
		StartTime:         startTime,
		EndTime:           endTime,
		Deadline:          deadline,
		Retries:           retries,
		Attributes:        attributes,
		ProcessSpec:       processSpec,
	}
}

func ConvertJSONToProcess(jsonString string) (*Process, error) {
	var process *Process
	err := json.Unmarshal([]byte(jsonString), &process)
	if err != nil {
		return nil, err
	}

	return process, nil
}

func ConvertProcessArrayToJSON(processes []*Process) (string, error) {
	jsonBytes, err := json.MarshalIndent(processes, "", "    ")
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}

func ConvertJSONToProcessArray(jsonString string) ([]*Process, error) {
	var processes []*Process
	err := json.Unmarshal([]byte(jsonString), &processes)
	if err != nil {
		return processes, err
	}

	return processes, nil
}

func IsProcessArraysEqual(processes1 []*Process, processes2 []*Process) bool {
	counter := 0
	for _, process1 := range processes1 {
		for _, process2 := range processes2 {
			if process1.Equals(process2) {
				counter++
			}
		}
	}

	if counter == len(processes1) && counter == len(processes2) {
		return true
	}

	return false
}

func (process *Process) Equals(process2 *Process) bool {
	same := true
	if process.ID != process2.ID &&
		process.AssignedRuntimeID != process2.AssignedRuntimeID &&
		process.Status != process2.Status &&
		process.IsAssigned != process2.IsAssigned &&
		process.SubmissionTime != process2.SubmissionTime &&
		process.StartTime != process2.StartTime &&
		process.EndTime != process2.EndTime &&
		process.Deadline != process2.Deadline &&
		process.Retries != process2.Retries {
		same = false
	}

	if !IsAttributeArraysEqual(process.Attributes, process2.Attributes) {
		same = false
	}

	if !process.ProcessSpec.Equals(process2.ProcessSpec) {
		same = false
	}

	return same
}

func (process *Process) Assign() {
	process.IsAssigned = true
}

func (process *Process) Unassign() {
	process.IsAssigned = false
}

func (process *Process) SetStatus(status int) {
	process.Status = status
}

func (process *Process) SetAssignedRuntimeID(runtimeID string) {
	process.AssignedRuntimeID = runtimeID
	process.IsAssigned = true
}

func (process *Process) SetAttributes(attributes []*Attribute) {
	process.Attributes = attributes
}

func (process *Process) SetSubmissionTime(submissionTime time.Time) {
	process.SubmissionTime = submissionTime
}

func (process *Process) SetStartTime(startTime time.Time) {
	process.StartTime = startTime
}

func (process *Process) SetEndTime(endTime time.Time) {
	process.EndTime = endTime
}

func (process *Process) WaitingTime() time.Duration {
	return process.StartTime.Sub(process.SubmissionTime)
}

func (process *Process) ProcessingTime() time.Duration {
	return process.EndTime.Sub(process.StartTime)
}

func (process *Process) ToJSON() (string, error) {
	jsonBytes, err := json.MarshalIndent(process, "", "    ")
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}
