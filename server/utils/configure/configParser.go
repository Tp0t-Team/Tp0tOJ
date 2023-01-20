package configure

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"server/utils"
)

const WriteUpPath string = "./writeup"

var LoadConfigError error = nil

func Parse(config *utils.Config) {
	prefix, _ := os.Getwd()
	f, err := ioutil.ReadFile(prefix + "/resources/config.yaml")
	if err != nil {
		LoadConfigError = err
		return
	}
	err = yaml.Unmarshal(f, config)
	if err != nil {
		LoadConfigError = err
	}

}

var Configure *utils.Config

func init() {
	Configure = new(utils.Config)
	Parse(Configure)
}
