package main

import (
	"github.com/dotSlashLu/agent-neo/modules"
	power "github.com/dotSlashLu/agent-neo/modules/power"
	volume "github.com/dotSlashLu/agent-neo/modules/volume"
	net "github.com/dotSlashLu/agent-neo/modules/net"
)

var registeredModules = map[string]module.CallableModule{
	"volume": volume.New(config),
	"power":  power.New(config),
	"net": net.New(config),
}
