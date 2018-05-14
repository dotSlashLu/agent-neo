package libvirt

import "github.com/libvirt/libvirt-go"

type UUID [36]byte

func Connect() (*libvirt.Connect, error) {
	return libvirt.NewConnect("qemu:///system")
}

func ConnectClose(c *libvirt.Connect) {
	c.Close()
}
