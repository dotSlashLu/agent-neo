package volume

import (
	"encoding/binary"
	"bytes"
	"encoding/json"

	llib "github.com/dotSlashLu/agent-neo/lib/libvirt"
)

func (m *Module) detach(recv []byte) ([]byte, error) {
	type paramProto struct {
		UUID   [32]byte // vm uuid
		Name   [32]byte // random str
		Target [3]byte  // vdb? vdc?
	}
	params := paramProto{}
	if err := binary.Read(bytes.NewReader(recv),
		binary.LittleEndian,
		&params); err != nil {
		return respError(err)
	}
	xml := getDeviceXML(params.UUID[:], params.Name[:], params.Target[:])
	conn, err := llib.Connect()
	defer func() {
		conn.Close()
	}()
	if err != nil {
		return respError(err)
	}
	dom, err := conn.LookupDomainByUUID(params.UUID[:])
	if err != nil {
		return respError(err)
	}
	dom.DetachDevice(xml)

	resp := struct{Status string `json:"status"`}{"ok"}
	ret, _ := json.Marshal(resp)
	return ret, nil
}
