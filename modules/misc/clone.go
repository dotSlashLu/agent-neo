package misc

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/beevik/etree"
	"github.com/dotSlashLu/agent-neo/lib"
	llib "github.com/dotSlashLu/agent-neo/lib/libvirt"
	"github.com/libvirt/libvirt-go"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func updateXML(srcXML string, srcUUID []byte, dstUUID []byte,
	baseDir string, dstMac []byte, dstVLAN int16) (string, error) {
	doc := etree.NewDocument()
	if err := doc.ReadFromString(srcXML); err != nil {
		return "", err
	}
	domain := doc.SelectElement("domain")
	uuidElem := domain.SelectElement("uuid")
	uuidElem.SetText(string(dstUUID))
	nameElem := domain.SelectElement("name")
	nameElem.SetText(nameElem.Text() + "_clone")

	devicesElem := domain.SelectElement("devices")
	// copy disks
	diskElems := devicesElem.SelectElements("disk")
	for _, diskElem := range diskElems {
		sourceElem := diskElem.SelectElement("source")
		if sourceElem == nil {
			continue
		}
		file := sourceElem.SelectAttrValue("file", "")
		if err := copyFile(baseDir, string(dstUUID), file); err != nil {
			return "", err
		}
		newFile := strings.Replace(file, string(srcUUID), string(dstUUID), -1)
		sourceElem.CreateAttr("file", newFile)
	}

	// change interfaces
	// now only one interface is supported
	ifaceElem := devicesElem.SelectElement("interface")
	ifaceElem.SelectElement("mac").CreateAttr("address", string(dstMac))
	vlanElem := ifaceElem.SelectElement("vlan")
	if vlanElem != nil {
		vlan := strconv.Itoa(int(dstVLAN))
		vlanElem.SelectElement("tag").CreateAttr("id", vlan)
	}

	newXMLStr, _ := doc.WriteToString()
	return newXMLStr, nil
}

func copy(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

func copyFile(baseDir, newUUID, srcFile string) error {
	dir := baseDir + "/" + newUUID
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err = os.Mkdir(dir, os.ModePerm); err != nil {
			return err
		}
	}
	dstFile := fmt.Sprintf("%s/%s/%s", baseDir, newUUID, filepath.Base(srcFile))
	log.Printf("copy %s to %s", srcFile, dstFile)
	if err := copy(srcFile, dstFile); err != nil {
		return err
	}
	return nil
}

// proto
//	SrcUUID llib.UUID
//	DstUUID llib.UUID
//	// DstIP 	[15]byte
//	DstMac	[17]byte
//	DstVLAN int16
func (m *Module) clone(recv []byte) ([]byte, error) {
	type paramsProto struct {
		SrcUUID llib.UUID
		DstUUID llib.UUID
		// DstIP 	[15]byte
		DstMac  [17]byte
		DstVLAN int16
	}
	params := paramsProto{}
	binary.Read(bytes.NewReader(recv), m.Config.Endianness_, &params)
	// dstIP := lib.TrimBuf(params.DstIP[:])
	dstMac := lib.TrimBuf(params.DstMac[:])
	srcUUID := lib.TrimBuf(params.SrcUUID[:])
	dstUUID := lib.TrimBuf(params.DstUUID[:])

	conn, err := llib.Connect()
	if err != nil {
		return lib.RespError(err)
	}
	defer func() {
		conn.Close()
	}()
	dom, err := conn.LookupDomainByUUIDString(string(srcUUID))
	if err != nil {
		return lib.RespError(err)
	}
	domInfo, _ := dom.GetInfo()
	if domInfo.State == libvirt.DOMAIN_RUNNING {
		return lib.RespError(errors.New("domain is alive"))
	}
	domXML, err := dom.GetXMLDesc(libvirt.DOMAIN_XML_SECURE)
	if err != nil {
		return lib.RespError(err)
	}
	newXML, err := updateXML(domXML, srcUUID, dstUUID, m.Config.ImageDir,
		dstMac, params.DstVLAN)
	if err != nil {
		return lib.RespError(err)
	}
	_, err = conn.DomainDefineXML(newXML)
	if err != nil {
		return lib.RespError(err)
	}

	return lib.RespOk("")
}
