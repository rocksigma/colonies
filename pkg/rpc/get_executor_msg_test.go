package rpc

import (
	"testing"

	"github.com/colonyos/colonies/pkg/core"
	"github.com/stretchr/testify/assert"
)

func TestRPCGetExecutorMsg(t *testing.T) {
	msg := CreateGetExecutorMsg(core.GenerateRandomID())
	jsonString, err := msg.ToJSON()
	assert.Nil(t, err)

	msg2, err := CreateGetExecutorMsgFromJSON(jsonString + "error")
	assert.NotNil(t, err)

	msg2, err = CreateGetExecutorMsgFromJSON(jsonString)
	assert.Nil(t, err)

	assert.True(t, msg.Equals(msg2))
}

func TestRPCGetExecutorMsgIndent(t *testing.T) {
	msg := CreateGetExecutorMsg(core.GenerateRandomID())
	jsonString, err := msg.ToJSONIndent()
	assert.Nil(t, err)

	msg2, err := CreateGetExecutorMsgFromJSON(jsonString + "error")
	assert.NotNil(t, err)

	msg2, err = CreateGetExecutorMsgFromJSON(jsonString)
	assert.Nil(t, err)

	assert.True(t, msg.Equals(msg2))
}

func TestRPCGetExecutorMsgEquals(t *testing.T) {
	msg := CreateGetExecutorMsg(core.GenerateRandomID())
	assert.True(t, msg.Equals(msg))
	assert.False(t, msg.Equals(nil))
}
