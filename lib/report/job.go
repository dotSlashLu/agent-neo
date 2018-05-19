package report

import (
	"bytes"
	"encoding/binary"
	"github.com/dotSlashLu/agent-neo/lib"
)

type jobMsg struct {
	jobID  [20]byte
	msgLen int32
	msg    []byte
}

func ReportJob(jobID []byte, msg string) {
	msg := jobMsg{
		jobID,
		len(msg),
		[]byte(msg),
	}
	buf = new(bytes.Buffer)
	if err := binary.Write(buf, lib.Cfg.Endianness_, msg); err != nil {
		return err
	}
	return report(ReportTypeJob, buf.Bytes())
}
