package rpc

import (
	"testing"

	"github.com/colonyos/colonies/pkg/core"
	"github.com/stretchr/testify/assert"
)

func TestRPCGetProcessesMsg(t *testing.T) {
	msg := CreateGetProcessesMsg(core.GenerateRandomID(), 1, 2, "test_executor_type")
	jsonString, err := msg.ToJSON()
	assert.Nil(t, err)

	msg2, err := CreateGetProcessesMsgFromJSON(jsonString + "error")
	assert.NotNil(t, err)

	msg2, err = CreateGetProcessesMsgFromJSON(jsonString)
	assert.Nil(t, err)

	assert.True(t, msg.Equals(msg2))
}

func TestRPCGetProcessesMsgIndent(t *testing.T) {
	msg := CreateGetProcessesMsg(core.GenerateRandomID(), 1, 2, "")
	jsonString, err := msg.ToJSONIndent()
	assert.Nil(t, err)

	msg2, err := CreateGetProcessesMsgFromJSON(jsonString + "error")
	assert.NotNil(t, err)

	msg2, err = CreateGetProcessesMsgFromJSON(jsonString)
	assert.Nil(t, err)

	assert.True(t, msg.Equals(msg2))
}

func TestRPCGetProcessesMsgEquals(t *testing.T) {
	msg := CreateGetProcessesMsg(core.GenerateRandomID(), 1, 2, "test_executor_type")
	assert.True(t, msg.Equals(msg))
	assert.False(t, msg.Equals(nil))
}
