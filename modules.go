package main

import (
	"github.com/dotSlashLu/agent-neo/modules"
	volume "github.com/dotSlashLu/agent-neo/modules/volume"
)

var registeredModules = map[string]module.CallableModule{
	"volume": volume.New(config),
}
