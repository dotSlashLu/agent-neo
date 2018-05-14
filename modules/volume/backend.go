package volume

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dotSlashLu/agent-neo/lib"
	llib "github.com/dotSlashLu/agent-neo/lib/libvirt"
	"log"
	"os"
	"os/exec"
	"strings"
)

/*
   create a volume backend
   signature:
   struct {
       UUID   [36]byte // vm uuid
       Name   [32]byte // random str
       Size   int32    // size in MB
   }
*/
func (m *Module) create(recv []byte) ([]byte, error) {
	type paramsProto struct {
		UUID llib.UUID // vm uuid
		Name [32]byte  // random str
		Size int32     // size in MB
	}
	p := paramsProto{}
	err := binary.Read(bytes.NewReader(recv), m.Config.Endianness_, &p)
	if err != nil {
		return respError(err)
	}
	log.Println("parsed params", p)
	imgDir := fmt.Sprintf("/data/kvm_img/%s", lib.TrimBuf(p.UUID[:]))
	if _, err := os.Stat(imgDir); os.IsNotExist(err) {
		return respError(errors.New(imgDir + " doesn't exist"))
	}
	imgName := fmt.Sprintf("%s/%s.qcow2", imgDir, lib.TrimBuf(p.Name[:]))
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

func (m *Module) delete(recv []byte) ([]byte, error) {
	type paramsProto struct {
		UUID llib.UUID // vm uuid
		Name [32]byte  // random str
	}
	params := paramsProto{}
	err := binary.Read(bytes.NewReader(recv), m.Config.Endianness_, &params)
	if err != nil {
		return respError(err)
	}
	uuid := lib.TrimBuf(params.UUID[:])
	name := lib.TrimBuf(params.Name[:])
	imgName := fmt.Sprintf("/data/kvm_img/%s/%s.qcow2", uuid, name)

	if _, err := os.Stat(imgName); os.IsNotExist(err) {
		goto SUCC
	}
	if err := os.Remove(imgName); err != nil {
		return respError(err)
	}

SUCC:
	resp := struct {
		Status string `json:"status"`
	}{"ok"}
	ret, _ := json.Marshal(resp)
	return ret, nil
}
