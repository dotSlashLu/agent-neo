package power

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/dotSlashLu/agent-neo/lib"
	llib "github.com/dotSlashLu/agent-neo/lib/libvirt"
)

type Module struct {
	Name   string
	Config *lib.Config
}

var PowerModule = &Module{"power", nil}

func New(config *lib.Config) *Module {
	PowerModule.Config = config
	return PowerModule
}

func (m *Module) Call(fn string, params []byte) ([]byte, error) {
	return m.power(params)
}

/*
	proto
		UUID [36]byte
		Op   [10]byte
	Op: suspend, resume
*/
func (m *Module) power(recv []byte) ([]byte, error) {
	type paramsProto struct {
		UUID llib.UUID
		Op   [10]byte
	}
	params := paramsProto{}
	binary.Read(bytes.NewReader(recv), m.Config.Endianness_, &params)
	uuid := lib.TrimBuf(params.UUID[:])
	op := lib.TrimBuf(params.Op[:])

	conn, err := llib.Connect()
	if err != nil {
		return lib.RespError(err)
	}
	defer func() {
		conn.Close()
	}()
	dom, err := conn.LookupDomainByUUIDString(string(uuid))
	if err != nil {
		return lib.RespError(err)
	}

	var opMethod func() error
	switch string(op) {
	case "suspend":
		opMethod = dom.Suspend
	case "resume":
		opMethod = dom.Resume
	default:
		panic(fmt.Sprintf("%s op not implemented", op))
	}
	if err := opMethod(); err != nil {
		return lib.RespError(err)
	}
	ret, _ := json.Marshal(struct {
		Status string `json:"status"`
	}{"ok"})
	return ret, nil
}
