package volume

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"text/template"

	llib "github.com/dotSlashLu/agent-neo/lib/libvirt"
	"github.com/libvirt/libvirt-go"
)

/*
   attach a volume backend to vm
   signature:
   struct {
       UUID   [68]byte // vm uuid
       Name   [32]byte // random str
       Target [3]byte  // vdb? vdc?
       Slot   [4]byte  // 0x007++
   }
   68s 32s 3s 4s

   This can be tested in the command line with

	   virsh attach-disk 23 --source /data/kvm_img/test.qcow2 \
	   --target vde --persistent --driver qemu --subdriver qcow2 \
	   --live --print-xml
*/
func (m *Module) attach(params []byte) (string, error) {
	type paramsT struct {
		UUID   [68]byte // vm uuid
		Name   [32]byte // random str
		Target [3]byte  // vdb? vdc?
		Slot   [4]byte  // 0x007++
	}
	type resp struct {
		Status string `json:"status"`
	}

	p := paramsT{}
	err := binary.Read(bytes.NewReader(params), binary.LittleEndian, &p)
	if err != nil {
		fmt.Println("error parsing params", err.Error())
		return "", err
	}
	fmt.Println(p)
	imgName := fmt.Sprintf("/data/kvm_img/%s/%s.qcow2", p.UUID, p.Name)
	xmlStr := defineDevice(imgName, p.Target[:], p.Slot[:])
	conn, err := llib.Connect()
	defer llib.ConnectClose(conn)
	if err != nil {
		return respError(err)
	}
	dom, err := conn.LookupDomainByUUID(p.UUID[:])
	if err != nil {
		return respError(err)
	}
	err = dom.AttachDeviceFlags(xmlStr, libvirt.DOMAIN_DEVICE_MODIFY_LIVE)
	if err != nil {
		return respError(err)
	}
	ret, err := json.Marshal(resp{"ok"})
	if err != nil {
		return "", err
	}
	return string(ret), nil
}

func defineDevice(imgName string, target []byte, slot []byte) string {
	deviceXMLTemplate := `
        <disk type='file' device='disk'>
            <driver name='qemu' type='qcow2' cache='none'/>
            <source file='{{.ImgName}}'/>
            <target dev='{{.Target}}' bus='virtio'/>
        </disk>
    `
	var deviceXMLBuf bytes.Buffer
	t := template.Must(template.New("device").Parse(deviceXMLTemplate))
	templateVals := map[string]string{
		"ImgName": imgName,
		"Target":  string(bytes.Trim(target[:], "\x00")),
		"Slot":    string(bytes.Trim(slot[:], "\x00")),
	}
	writer := bufio.NewWriter(&deviceXMLBuf)
	t.Execute(writer, templateVals)
	writer.Flush()
	ret := deviceXMLBuf.String()
	fmt.Println("gen xml", ret)
	return ret
}
