package configure

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
	"server/utils"
)

const WriteUpPath string = "./writeup"

func Parse(config *utils.Config) {
	prefix, _ := os.Getwd()
	f, err := ioutil.ReadFile(prefix + "/resources/config.yaml")
	if err != nil {
		log.Panicln(err.Error())
	}
	err = yaml.Unmarshal(f, config)
	if err != nil {
		log.Panicf("Unmarshal: %v", err)
	}

}

var Configure *utils.Config

func init() {
	Configure = new(utils.Config)
	Parse(Configure)
}
