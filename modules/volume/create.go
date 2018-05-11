package volume

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"io/ioutil"
	"os/exec"
	// "strconv"
)

type paramsT struct {
	UUID [32]byte // vm uuid
	Name [32]byte // random str
	Size int32    // size in MB
}
/*
   create a volume backend
   signature:
   struct {
       UUID   [32]byte // vm uuid
       Name   [32]byte // random str
       Size   int32    // size in MB
   }
   python struct fmt: 68s 32s i
*/
func (m *Module) create(params []byte) ([]byte, error) {
	fmt.Println("params", params)
	p := paramsT{}
	err := binary.Read(bytes.NewReader(params), binary.LittleEndian, &p)
	if err != nil {
		respError(err)
	}
	log.Println("params", p)
	imgName := fmt.Sprintf("/data/kvm_img/%s/%s.qcow2", p.UUID, p.Name)
	// TODO fork/exec /usr/bin/qemu-img: invalid argument
	cmd := exec.Command("qemu-img", "create", "-f", "qcow2",
		imgName, fmt.Sprintf("%dM", p.Size))
	fmt.Println("execute qemu-img", fmt.Sprintf("create -f qcow2 %s %dM", imgName, p.Size))
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Run(); err != nil {
		return respError(err)
	}
	errOut, _ := ioutil.ReadAll(stderr)
	if len(errOut) > 0 {
		return respError(errors.New(string(errOut)))
	}
	type resp struct {
		Status string `json:"status"`
	}
	ret, _ := json.Marshal(resp{"ok"})
	return ret, nil
}
