package report

import (
	"bytes"
	"encoding/binary"
	"github.com/dotSlashLu/agent-neo/lib"
	llib "github.com/dotSlashLu/agent-neo/lib/libvirt"
)

type VMStat int

const (
	VMStatRunning VMStat = iota
	VMStatShutdown
	VMStatPaused
	VMStatCreating
)

type vMStatMessage struct {
	uuid llib.UUID
	stat VMStat
}

func ReportVMStat(uuid []byte, status VMStat) error {
	msg := vMStatMessage{
		uuid,
		status,
	}
	buf = new(bytes.Buffer)
	if err := binary.Write(buf, lib.Cfg.Endianness_, msg); err != nil {
		return err
	}
	return report(ReportTypeVMStat, buf.Bytes())
}
