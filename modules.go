package main

import (
	"github.com/dotSlashLu/agent-neo/modules"
	misc "github.com/dotSlashLu/agent-neo/modules/misc"
	net "github.com/dotSlashLu/agent-neo/modules/net"
	power "github.com/dotSlashLu/agent-neo/modules/power"
	volume "github.com/dotSlashLu/agent-neo/modules/volume"
)

var registeredModules = map[string]module.CallableModule{
	"volume": volume.New(config),
	"power":  power.New(config),
	"net":    net.New(config),
	"misc":   misc.New(config),
}
