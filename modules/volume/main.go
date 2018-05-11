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
	switch fn {
	case "create":
		return m.create(params)
	case "attach":
		return m.attach(params)
	case "detach":
		return m.detach(params)
	default:
		panic("method not implemented")
	}
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
