package misc

import (
	"github.com/dotSlashLu/agent-neo/lib"
	llib "github.com/dotSlashLu/agent-neo/lib"
	"encoding/xml"
	"encoding/binary"
)

func (m *Module) clone(recv []byte) ([]byte, error) {
	type paramsProto struct {
		srcUUID llib.UUID
		dstUUID llib.UUID
		dstip 	[15]byte
		dstmac	[17]byte
	}
	params := paramsProto{}
	binary.Read(bytes.NewReader(recv), m.config.Endianness_, &params)
	return nil, nil
}
