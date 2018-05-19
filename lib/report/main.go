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
