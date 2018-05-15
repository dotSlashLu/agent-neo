package misc

import (
	"github.com/dotSlashLu/agent-neo/lib"
)

type Module struct {
	Name   string
	Config *lib.Config
}

var MiscModule = &Module{"misc", nil}

func New(config *lib.Config) *Module {
	MiscModule.Config = config
	return MiscModule
}

func (m *Module) Call(fn string, params []byte) ([]byte, error) {
	var method func([]byte) ([]byte, error)
	switch fn {
	case "clone":
		method = m.clone
	default:
		panic("method not implemented")
	}
	return method(params)
}
