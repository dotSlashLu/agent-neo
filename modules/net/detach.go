package net

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/dotSlashLu/agent-neo/lib"
	llib "github.com/dotSlashLu/agent-neo/lib/libvirt"
	"github.com/libvirt/libvirt-go"
	"strings"
)

type netType int16

const (
	netTypeIntra netType = iota
	netTypeInter
)

// proto:
// 	UUID llib.UUID
// 	MAC  [17]byte
// 	VLAN int16
// 	Type netType
//
// Right now the device's bridge is hardcoded as ovsbr0 / ovsbr1 according to
// the netType.
// Better find the device by parsing the VM's definition XML to find the
// device by its MAC.
func (m *Module) detach(recv []byte) ([]byte, error) {
	type paramsProto struct {
		UUID llib.UUID
		MAC  [17]byte
		VLAN int16
		Type netType
	}
	params := paramsProto{}
	binary.Read(bytes.NewReader(recv), m.Config.Endianness_, &params)
	uuid := lib.TrimBuf(params.UUID[:])
	mac := lib.TrimBuf(params.MAC[:])

	vlanXML := ""
	if params.VLAN != 0 {
		vlanXML = fmt.Sprintf(`
			<vlan>
				<tag id="%s" />
			</vlan>
		`, params.VLAN)
	}
	bridge := "ovsbr0"
	if params.Type == netTypeInter {
		bridge = "ovsbr1"
	}
	dev := strings.Replace(string(mac), ":", "", -1)
	ifXML := fmt.Sprintf(`
		<interface type='bridge'>
			<mac address='%s'/>
			<source bridge='%s'/>
			<virtualport type='openvswitch'>
			</virtualport>
			<target dev='v%s'/>
			<model type='virtio'/>
			%s
		</interface>
	`, string(mac), bridge, dev, vlanXML)
	conn, err := llib.Connect()
	if err != nil {
		return lib.RespError(err)
	}
	dom, err := conn.LookupDomainByUUIDString(string(uuid))
	if err != nil {
		return lib.RespError(err)
	}
	flags := libvirt.DOMAIN_DEVICE_MODIFY_CONFIG |
		libvirt.DOMAIN_DEVICE_MODIFY_LIVE
	if err = dom.DetachDeviceFlags(ifXML, flags); err != nil {
		return lib.RespError(err)
	}
	return lib.RespOk("")
}
