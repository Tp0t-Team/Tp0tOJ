package utils

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
)

type Config struct {
	Server     Server     `yaml:"server"`
	Email      Email      `yaml:"email"`
	Redis      Redis      `yaml:"redis"`
	Challenge  Challenge  `yaml:"challenge"`
	Kubernetes Kubernetes `yaml:"kubernetes"`
}

type Server struct {
	Host  string `yaml:"host"`
	Name  string `yaml:"name"`
	Port  int    `yaml:"port"`
	Salt  string `yaml:"salt"`
	Debug bool   `yaml:"debug"`
}

type Email struct {
	Host     string `yaml:"host"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Redis struct {
	Password  string `yaml:"password"`
	MaxActive int    `yaml:"maxActive"`
	MaxIdle   int    `yaml:"maxIdle"`
	Database  int    `yaml:"database"`
	Host      string `yaml:"host"`
	Port      int    `yaml:"port"`
	MaxWait   string `yaml:"maxWait"`
	MinIdle   int    `yaml:"minIdle"`
	Timeout   int    `yaml:"timeout"`
}

type Challenge struct {
	SecondBloodReward float64 `yaml:"secondBloodReward"`
	ThirdBloodReward  float64 `yaml:"thirdBloodReward"`
	HalfLife          int     `yaml:"halfLife"`
	FirstBloodReward  float64 `yaml:"firstBloodReward"`
}

type Kubernetes struct {
	PortAllocBegin int32  `yaml:"portAllocBegin"`
	PortAllocEnd   int32  `yaml:"portAllocEnd"`
	Username       string `yaml:"username"`
	Password       string `yaml:"password"`
	RegistryHost   string `yaml:"registryHost"`
}

func Parse(config *Config) {
	f, err := ioutil.ReadFile("resources/config.yaml")
	if err != nil {
		log.Panicln(err.Error())
	}
	err = yaml.Unmarshal(f, config)
	if err != nil {
		log.Panicf("Unmarshal: %v", err)
	}

}

var Configure *Config

func init() {
	Configure = new(Config)
	Parse(Configure)
}
