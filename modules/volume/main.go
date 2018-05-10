package volume

import (
	"encoding/binary"
	"fmt"
	"os/exec"
	// "github.com/dotSlashLu/agent-neo/modules"
	"bufio"
	"encoding/json"
	// "strings"
	"bytes"
	"text/template"
)

type Module struct {
	Name string
}

var VolumeModule Module

func (m Module) Call(fn string, params []byte) (string, error) {
	switch fn {
	case "echo":
		return echo(params)
	case "createVolume":
		return createVolume(params)
	default:
		panic("method not implemented")
	}
}

func echo(params []byte) (string, error) {
	type resp struct {
		Status string `json:"status"`
	}

	content := string(params)
	res, err := json.Marshal(resp{content})
	if err != nil {
		return "", err
	}
	return string(res), nil
}

func createVolume(params []byte) (string, error) {
	// 68s32s3s4si
	type paramsT struct {
		UUID   [68]byte // vm uuid
		Name   [32]byte // random str
		Target [3]byte  // vdb? vdc?
		Slot   [4]byte  // 0x007++
		Size   int32    // size in MB
	}
	type resp struct {
		Status string `json:"status"`
	}

	p := paramsT{}
	if err := binary.Read(bytes.NewReader(params), binary.LittleEndian, &p); err != nil {
		fmt.Println("error parsing params", err.Error())
		return "", err
	}
	fmt.Println(p)
	imgName := fmt.Sprintf("/data/kvm_img/%s/%s.qcow2", p.UUID, p.Name)
	cmd := exec.Command("qemu-img", "create", "-f", "qcow2", imgName)
	cmd.Run()
	deviceXMLTemplate := `
		<disk type='file' device='disk'>                                                              
			<driver name='qemu' type='qcow2' cache='none'/>                                             
		  	<source file='{{.ImgName}}'/>
		  	<target dev='{{.Target}}' bus='virtio'/>                                                            
		  	<address type='pci' domain='0x0000' bus='0x00' 
		  		slot='{{.Slot}}' function='0x0'/>                 
		</disk>                                                                                       
	`
	var deviceXMLBuf bytes.Buffer
	t := template.Must(template.New("device").Parse(deviceXMLTemplate))
	templateVals := map[string]string{
		"ImgName": imgName,
		"Target":  string(bytes.Trim(p.Target[:], "\x00")),
		"Slot":    string(bytes.Trim(p.Slot[:], "\x00")),
	}
	writer := bufio.NewWriter(&deviceXMLBuf)
	t.Execute(writer, templateVals)
	writer.Flush()
	fmt.Println("gen xml", deviceXMLBuf.String())
	ret, err := json.Marshal(resp{"ok"})
	if err != nil {
		return "", err
	}
	return string(ret), nil
}

func attachVolume(params []byte) (string, error) {
	/*#
	virsh attach-disk {vm-name} \
	--source /var/lib/libvirt/images/{img-name-here} \
	--target vdb \
	--persistent*/
	return "", nil
}

func init() {
	VolumeModule = Module{"volume"}
}
