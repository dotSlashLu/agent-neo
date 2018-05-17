// Report job status, vm status back to the remote controller interface.
package report

import (
	"bytes"
	"encoding/binary"
	"github.com/dotSlashLu/agent-neo/lib"
	llib "github.com/dotSlashLu/agent-neo/lib/libvirt"
	"net"
)

type ReportType int

const (
	ReportTypeJob ReportType = iota
	ReportTypeVMStat
)

type reportHeader struct {
	reportType ReportType
	bodyLen    int32
}

type reportJobMsg struct {
	jobID  [20]byte
	msgLen int32
	msg    []byte
}

// TODO
// Now two buffers are made for both header and body which is repeating
// and might hurt performance. Try to use only one buffer.
func report(reportType ReportType, body []byte) error {
	header := reportHeader{
		reportType,
		len(body),
	}
	conn, err := net.Dial("tcp", "localhost:10240")
	if err != nil {
		return err
	}
	buf := new(Bytes.Buffer)
	binary.Write(buf, lib.Cfg.Endianness_, header)
	buf.Write(body)
	if err = lib.SendAll(conn, buf.Bytes()); err != nil {
		return err
	}
}

func ReportJob(jobID []byte, msg string) {
	msg := reportJobMessage{
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

type VMStat int

const (
	VMStatRunning VMStat = iota
	VMStatShutdown
	VMStatPaused
	VMStatCreating
)

type reportVMStatMessage struct {
	uuid llib.UUID
	stat VMStat
}

func ReportVMStat(uuid []byte, status VMStat) error {
	msg := reportVMStatMessage{
		uuid,
		status,
	}
	buf = new(bytes.Buffer)
	if err := binary.Write(buf, lib.Cfg.Endianness_, msg); err != nil {
		return err
	}
	return report(ReportTypeVMStat, buf.Bytes())
}
