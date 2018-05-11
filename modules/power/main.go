package power

import (
	"fmt"
	"bytes"
	"encoding/json"
	"encoding/binary"
	"github.com/dotSlashLu/agent-neo/lib"
	llib "github.com/dotSlashLu/agent-neo/lib/libvirt"
)

type Module struct {
	Name   string
	Config *lib.Config
}

var PowerModule = &Module{"volume", nil}

func New(config *lib.Config) *Module {
	PowerModule.Config = config
	return PowerModule
}

func (m *Module) Call(fn string, params []byte) ([]byte, error) {
	return power(params)
}

func power(recv []byte) ([]byte, error) {
	type paramsProto struct {
		Op		[10]byte
		UUID	llib.UUID
	}
	params := paramsProto{}
	binary.Read(bytes.NewReader(recv), binary.LittleEndian, &params)

	conn, err := llib.Connect()
	if err != nil {
		return respError(err)
	}
	defer func() {
		conn.Close()
	}()
	dom, err := conn.LookupDomainByUUID(params.UUID[:])
	if err != nil {
		return respError(err)
	}

	resp := struct {
		Status string `json:"status"`
	}{"ok"}
	switch string(params.Op[:]) {
	case "suspend":
		if err := dom.Suspend(); err != nil {
			return respError(err)
		}
	case "resume":
		if err := dom.Resume(); err != nil {
			return respError(err)
		}
	default:
		panic(fmt.Sprintf("%s op not implemented", params.Op))
	}
	ret, _ := json.Marshal(resp)
	return ret, nil
}

func respError(e error) ([]byte, error) {
	type resp struct {
		Status string `json:"status"`
		Error  string `json:"error"`
	}
	str, err := json.Marshal(resp{"error", e.Error()})
	if err != nil {
		panic(err)
	}
	return str, nil
}
