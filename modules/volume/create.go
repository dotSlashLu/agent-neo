package volume

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"github.com/dotSlashLu/agent-neo/lib"
)

type paramsT struct {
	UUID [36]byte // vm uuid
	Name [32]byte // random str
	Size int32    // size in MB
}

/*
   create a volume backend
   signature:
   struct {
       UUID   [36]byte // vm uuid
       Name   [32]byte // random str
       Size   int32    // size in MB
   }
   python struct fmt: 36s 32s i
*/
func (m *Module) create(params []byte) ([]byte, error) {
	fmt.Println("params", params)
	p := paramsT{}
	err := binary.Read(bytes.NewReader(params), m.Config.Endianness_, &p)
	if err != nil {
		respError(err)
	}
	log.Println("parsed params", p)

	imgName := fmt.Sprintf("/data/kvm_img/%s/%s.qcow2", lib.TrimBuf(p.UUID[:]),
		lib.TrimBuf(p.Name[:]))
	cmdStr := fmt.Sprintf("create -f qcow2 %s %dM", imgName, p.Size)
	log.Println("cmd params", cmdStr)
	cmd := exec.Command("qemu-img", strings.Split(cmdStr, " ")...)
	var stderr, stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		return respError(err)
	}
	log.Println("stdout:", stdout.String(), "\nstderr:", stderr.String())

	type resp struct {
		Status string `json:"status"`
	}
	ret, _ := json.Marshal(resp{"ok"})
	return ret, nil
}
