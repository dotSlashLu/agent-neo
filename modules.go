package main

import (
	"github.com/dotSlashLu/agent-neo/modules"
	volume  "github.com/dotSlashLu/agent-neo/modules/volume"
	power	"github.com/dotSlashLu/agent-neo/modules/power"
)

var registeredModules = map[string]module.CallableModule{
	"volume": volume.New(config),
	"power": power.New(config),
}
