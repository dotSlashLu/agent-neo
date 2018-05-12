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
       UUID   [36]byte // vm uuid
       Name   [32]byte // random str
       Target [3]byte  // vdb? vdc?
       Slot   [4]byte  // 0x007++
   }
   python struct fmt: 36s 32s 3s 4s

   This can be tested in the command line with

	   virsh attach-disk 23 --source /data/kvm_img/test.qcow2 \
	   --target vde --persistent --driver qemu --subdriver qcow2 \
	   --live --print-xml
*/
func (m *Module) attach(recv []byte) ([]byte, error) {
	type paramsProto struct {
		UUID   [36]byte // vm uuid
		Name   [32]byte // random str
		Target [3]byte  // vdb? vdc?
		Slot   [4]byte  // 0x007++
	}

	p := paramsProto{}
	err := binary.Read(bytes.NewReader(recv), m.Config.Endianness_, &p)
	if err != nil {
		fmt.Println("error parsing params", err.Error())
		return []byte{}, err
	}
	fmt.Println(p)
	xmlStr := getDeviceXML(p.UUID[:], p.Name[:], p.Target[:])
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

	type resp struct {
		Status string `json:"status"`
	}
	ret, err := json.Marshal(resp{"ok"})
	if err != nil {
		return []byte{}, err
	}
	return ret, nil
}

func getDeviceXML(uuid []byte, fileName []byte, target []byte) string {
	deviceXMLTemplate := `
        <disk type='file' device='disk'>
            <driver name='qemu' type='qcow2' cache='none'/>
            <source file='{{.filePath}}'/>
            <target dev='{{.target}}' bus='virtio'/>
        </disk>
    `
	filePath := fmt.Sprintf("/data/kvm_img/%s/%s.qcow2", uuid, fileName)
	var deviceXMLBuf bytes.Buffer
	t := template.Must(template.New("device").Parse(deviceXMLTemplate))
	templateVals := map[string]string{
		"filePath": filePath,
		"target":   string(bytes.Trim(target[:], "\x00")),
	}
	writer := bufio.NewWriter(&deviceXMLBuf)
	t.Execute(writer, templateVals)
	writer.Flush()
	ret := deviceXMLBuf.String()
	fmt.Println("gen xml", ret)
	return ret
}
