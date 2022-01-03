package rpc

import (
	"encoding/json"
)

const GetRuntimeMsgType = "getruntime"

type GetRuntimeMsg struct {
	RuntimeID string `json:"runtimeid"`
}

func CreateGetRuntimeMsg(runtimeID string) *GetRuntimeMsg {
	msg := &GetRuntimeMsg{}
	msg.RuntimeID = runtimeID

	return msg
}

func (msg *GetRuntimeMsg) ToJSON() (string, error) {
	jsonBytes, err := json.Marshal(msg)
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}

func (msg *GetRuntimeMsg) ToJSONIndent() (string, error) {
	jsonBytes, err := json.MarshalIndent(msg, "", "    ")
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}

func CreateGetRuntimeMsgFromJSON(jsonString string) (*GetRuntimeMsg, error) {
	var msg *GetRuntimeMsg

	err := json.Unmarshal([]byte(jsonString), &msg)
	if err != nil {
		return msg, err
	}

	return msg, nil
}