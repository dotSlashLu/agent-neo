package libvirt

import "github.com/libvirt/libvirt-go"

type UUID [32]byte

func Connect() (*libvirt.Connect, error) {
	return libvirt.NewConnect("qemu:///")
}

func ConnectClose(c *libvirt.Connect) {
	c.Close()
}
