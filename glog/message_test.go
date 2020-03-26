package glog

import (
	"encoding/json"
	"testing"
)

func TestMarshal(t *testing.T) {
	m := &message{
		logLevel: infoLevel,
		Message:  "test message",
		Parameter: map[string]interface{}{
			"name": "swkwon",
			"age":  40,
		},
	}

	b, e := json.Marshal(m)
	if e != nil {
		t.Error(e)
	}

	t.Log(string(b))
}

func TestJson(t *testing.T) {
	type Msg struct {
		Message string `json:"message,omitempty"`
		Name    string `json:"name"`
		Age     int    `json:"age"`
	}

	m := &Msg{
		Name: "swkwon",
	}

	b, _ := json.Marshal(m)
	t.Log(string(b))
}
