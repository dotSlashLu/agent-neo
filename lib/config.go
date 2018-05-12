package lib

import (
	"encoding/binary"
	"github.com/BurntSushi/toml"
	"io/ioutil"
)

type Config struct {
	Port        int    `toml:"port"`
	SSHPassword string `toml:"ssh_password"`
	ImageDir    string `toml:"image_dir"`
	Endianness  string `toml:"endianness"`
	Endianness_ binary.ByteOrder
}

func ParseConfig(filename string, c *Config) (err error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	_, err = toml.Decode(string(content), c)
	if err != nil {
		panic(err)
	}
	if c.Endianness == "big" {
		c.Endianness_ = binary.BigEndian
	} else {
		c.Endianness_ = binary.LittleEndian
	}
	return
}
