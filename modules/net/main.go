package net

import (
	"github.com/dotSlashLu/agent-neo/lib"
)

type Module struct {
	Name   string
	Config *lib.Config
}

var NetModule = &Module{"net", nil}

func New(config *lib.Config) *Module {
	NetModule.Config = config
	return NetModule
}

func (m *Module) Call(fn string, params []byte) ([]byte, error) {
	var method func([]byte) ([]byte, error)
	switch fn {
	case "detach":
		method = m.detach
	default:
		panic("method not implemented")
	}
	return method(params)
}
