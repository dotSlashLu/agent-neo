package report

import (
	"bytes"
	"encoding/binary"
	"github.com/dotSlashLu/agent-neo/lib"
	"net"
)

type VMStat int
const (
	VMStatRunning	VMStat = iota
	VMStatShutdown
	VMStatPaused
	VMStatCreating
)
type reportMessage struct {
	jobID 	[20]byte
	msgLen	int32
	msg		[]byte
}

func Report(jobID []byte, msg string) (ok bool) {
	report := reportMessage{
		jobID
		len(msg)
		[]byte(msg)
	}
	conn, err := net.Dial("tcp", "localhost:10240")
	if err != nil {
		panic(err)
	}
	buf := new(Bytes.Buffer)
	binary.Write(buf, lib.Cfg.Endianness_, report)
	if err = lib.SendAll(conn, buf); err != nil {
		panic(err)
	}
	return true
}

func SetVMStat(uuid []byte, status) {
}
