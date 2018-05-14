package volume

import (
	"bytes"
	"encoding/binary"
	"encoding/json"

	"github.com/dotSlashLu/agent-neo/lib"
	llib "github.com/dotSlashLu/agent-neo/lib/libvirt"
)

func (m *Module) detach(recv []byte) ([]byte, error) {
	type paramProto struct {
		UUID   llib.UUID // vm uuid
		Name   [32]byte  // random str
		Target [3]byte   // vdb? vdc?
	}
	params := paramProto{}
	if err := binary.Read(bytes.NewReader(recv),
		m.Config.Endianness_,
		&params); err != nil {
		return respError(err)
	}
	uuid := lib.TrimBuf(params.UUID[:])
	name := lib.TrimBuf(params.Name[:])
	target := lib.TrimBuf(params.Target[:])

	xml := getDeviceXML(uuid, name, target)
	conn, err := llib.Connect()
	if err != nil {
		return respError(err)
	}
	defer func() {
		conn.Close()
	}()
	dom, err := conn.LookupDomainByUUIDString(string(uuid))
	if err != nil {
		return respError(err)
	}
	dom.DetachDevice(xml)

	resp := struct {
		Status string `json:"status"`
	}{"ok"}
	ret, _ := json.Marshal(resp)
	return ret, nil
}
