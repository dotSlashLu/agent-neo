package net

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/beevik/etree"
	"github.com/dotSlashLu/agent-neo/lib"
	llib "github.com/dotSlashLu/agent-neo/lib/libvirt"
	"github.com/libvirt/libvirt-go"
	"log"
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
	domXML, err := dom.GetXMLDesc(libvirt.DOMAIN_XML_SECURE)
	if err != nil {
		return lib.RespError(err)
	}
	doc := etree.NewDocument()
	if err = doc.ReadFromString(domXML); err != nil {
		return lib.RespError(err)
	}
	domain := doc.SelectElement("domain")
	devices := domain.SelectElement("devices")
	ifaces := devices.SelectElements("interface")
	var iface *etree.Element
	for _, iface = range ifaces {
		macElem := iface.SelectElement("mac")
		if macElem.SelectAttrValue("address", "") == string(mac) {
			break
		}
	}
	if iface == nil {
		return lib.RespError(errors.New("can't find iface"))
	}
	newDoc := etree.NewDocument()
	newDoc.SetRoot(iface)
	ifaceXML, _ := newDoc.WriteToString()
	// ifaceXML := iface.Text()
	log.Println("ifacexml", ifaceXML)
	flags := libvirt.DOMAIN_DEVICE_MODIFY_CONFIG |
		libvirt.DOMAIN_DEVICE_MODIFY_LIVE
	if err = dom.DetachDeviceFlags(ifaceXML, flags); err != nil {
		return lib.RespError(err)
	}
	return lib.RespOk("")
}

