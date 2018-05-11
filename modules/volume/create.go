package volume

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strconv"
)

/*
   create a volume backend
   signature:
   struct {
       UUID   [68]byte // vm uuid
       Name   [32]byte // random str
       Size   int32    // size in MB
   }
   68s 32s i
*/
func (m *Module) create(params []byte) (string, error) {
	type paramsT struct {
		UUID [68]byte // vm uuid
		Name [32]byte // random str
		Size int32    // size in MB
	}
	type resp struct {
		Status string `json:"status"`
	}

	p := paramsT{}
	err := binary.Read(bytes.NewReader(params), binary.LittleEndian, &p)
	if err != nil {
		respError(err)
	}
	imgName := fmt.Sprintf("/data/kvm_img/%s/%s.qcow2", p.UUID, p.Name)
	cmd := exec.Command("qemu-img", "create", "-f", "qcow2",
		imgName, strconv.Itoa(int(p.Size)))
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Run(); err != nil {
		return respError(err)
	}
	errOut, _ := ioutil.ReadAll(stderr)
	if len(errOut) > 0 {
		return respError(errors.New(string(errOut)))
	}
	ret, _ := json.Marshal(resp{"ok"})
	return string(ret), nil
}
