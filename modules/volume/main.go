package volume

import (
	"encoding/json"
	"github.com/dotSlashLu/agent-neo/lib"
)

type Module struct {
	Name   string
	Config *lib.Config
}

var VolumeModule = &Module{"volume", nil}

func New(config *lib.Config) *Module {
	VolumeModule.Config = config
	return VolumeModule
}

func (m *Module) Call(fn string, params []byte) ([]byte, error) {
	var method func([]byte) ([]byte, error)
	switch fn {
	case "create":
		method = m.create
	case "attach":
		method = m.attach
	case "detach":
		method = m.detach
	case "delete":
		method = m.delete
	default:
		panic("method not implemented")
	}
	return method(params)
}

func respError(e error) ([]byte, error) {
	type resp struct {
		Status string `json:"status"`
		Error  string `json:"error"`
	}
	str, err := json.Marshal(resp{"error", e.Error()})
	if err != nil {
		panic(err)
	}
	return str, nil
}
