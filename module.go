package main

import (
	"github.com/dotSlashLu/agent-neo/modules"
	volume "github.com/dotSlashLu/agent-neo/modules/volume"
)

type modules map[string]module.Module

func registerModule(ms modules, name string, m module.Module) {
	ms[name] = m
}

func registerModules() modules {
	ms := make(modules)
	registerModule(ms, "volume", volume.VolumeModule)
	return ms
}
