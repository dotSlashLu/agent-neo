package lib

import (
	"github.com/BurntSushi/toml"
	"io/ioutil"
)

type Config struct {
	Port        int    `toml:"port"`
	SSHPassword string `toml:"ssh_password"`
	ImageDir    string `toml:"image_dir"`
}

func ParseConfig(filename string, c *Config) (err error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	_, err = toml.Decode(string(content), c)
	return
}
