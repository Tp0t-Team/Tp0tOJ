package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"flag"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"math/big"
	rd "math/rand"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"time"
)

type ReleaseInfo struct {
	TagName string `json:"tag_name"`
}
type Config struct {
	Server     Server     `yaml:"server"`
	Email      Email      `yaml:"email"`
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

func main() {
	masterIP := flag.String("MasterIP", "", "master ip")
	serverName := flag.String("Name", "Tp0tOJ", "server name but seems useless")
	serverPort := flag.Int("Port", 8888, "backend server port")
	salt := flag.String("Salt", "", "Salt to encrypt password, default is some random thing")
	debug := flag.Bool("Debug", false, "debug mode")
	emailServerHost := flag.String("MailHost", "smtp.163.com", "mail host, is you don't set this, your user will have no chance to reset or modify password")
	emailUserName := flag.String("MailUser", "changeme", "mail server's username, is you don't set this, your user will have no chance to reset or modify password")
	emailPassword := flag.String("MailPass", "changeme", "mail server's password, is you don't set this, your user will have no chance to reset or modify password")
	firstBloodReward := flag.Float64("1stReward", 0.10, "this is the reward coefficient for FirstBlood winner")
	secondBloodReward := flag.Float64("2ndReward", 0.08, "this is the reward coefficient for SecondBlood winner")
	thirdBloodReward := flag.Float64("3thReward", 0.05, "this is the reward coefficient for ThirdBlood winner")
	halfLife := flag.Uint("HalfLife", 20, "half life is the solved counter for the challenge score to reduce half of basic score")
	//k8sEnable := flag.Bool("K8S", false, "choose this to generate default k8s config, AND REMEMBER TO MODIFY IT IN `resources` folder")
	flag.Parse()

	cmd := exec.Command("docker", "info")
	if err := cmd.Run(); err != nil {
		fmt.Println("You need an available docker environment first.")
		os.Exit(1)
	}
	if masterIP == nil || *masterIP == "" {
		fmt.Println("You need give your MasterIP.")
		os.Exit(1)
	}
	k3sRes, err := http.Get("http://rancher-mirror.cnrancher.com/k3s/k3s-install.sh")
	if err != nil {
		return
	}
	k3sInstall := exec.Command("bash", "-c", fmt.Sprintf("sudo sh -s - --node-ip %s --node-name %s", *masterIP, *masterIP))
	k3sPipe, err := k3sInstall.StdinPipe()
	if err != nil {
		return
	}
	_, err = io.Copy(k3sPipe, k3sRes.Body)
	if err != nil {
		return
	}
	releaseInfoRes, err := http.Get("https://api.github.com/repos/Tp0t-Team/Tp0tOJ/releases/latest")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	releaseInfoData := bytes.Buffer{}
	_, err = io.Copy(&releaseInfoData, releaseInfoRes.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	releaseInfo := ReleaseInfo{}
	err = json.Unmarshal(releaseInfoData.Bytes(), &releaseInfo)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if releaseInfo.TagName == "" {
		fmt.Println("no available release version.")
		os.Exit(1)
	}
	binaryName := fmt.Sprintf("OJ_%s_%s", runtime.GOOS, runtime.GOARCH)
	binaryRes, err := http.Get(fmt.Sprintf("https://github.com/Tp0t-Team/Tp0tOJ/releases/download/%s/%s", releaseInfo.TagName, binaryName))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	binary, err := os.Create("OJ")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	_, err = io.Copy(binary, binaryRes.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = binary.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = os.Chmod("OJ", 0755)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = os.MkdirAll("resources", 0755)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	CACert := &x509.Certificate{
		SerialNumber:          big.NewInt(rd.Int63()),
		BasicConstraintsValid: true,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(1, 0, 0),
		Subject:               pkix.Name{},
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature,
		EmailAddresses:        []string{},
	}
	Cert := &x509.Certificate{
		SerialNumber:          big.NewInt(rd.Int63()),
		BasicConstraintsValid: true,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(1, 0, 0),
		Subject:               pkix.Name{},
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature,
		EmailAddresses:        []string{},
	}
	Key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var buf []byte
	buf, err = x509.CreateCertificate(rand.Reader, CACert, CACert, &Key.PublicKey, Key)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	caFile, err := os.Create("resources/ca.crt")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	_, err = caFile.Write(buf)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = caFile.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	buf, err = x509.CreateCertificate(rand.Reader, Cert, CACert, &Key.PublicKey, Key)
	certFile, err := os.Create("resources/tls.crt")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	_, err = certFile.Write(buf)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = certFile.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	buf = x509.MarshalPKCS1PrivateKey(Key)
	keyFile, err := os.Create("resources/tls.key")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	_, err = keyFile.Write(buf)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = keyFile.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// make default config.yaml
	err = os.MkdirAll("resources", 0755)
	if err != nil {
		return
	}
	configFile, err := os.Create("resources/config.yaml")
	if err != nil {
		return
	}
	if *salt == "" {
		*salt = strconv.FormatInt(rd.Int63(), 10)
	}

	newConfig := Config{
		Server:     Server{Host: "127.0.0.1", Name: *serverName, Port: *serverPort, Salt: *salt, Debug: *debug},
		Email:      Email{Host: *emailServerHost, Username: *emailUserName, Password: *emailPassword},
		Challenge:  Challenge{FirstBloodReward: *firstBloodReward, SecondBloodReward: *secondBloodReward, ThirdBloodReward: *thirdBloodReward, HalfLife: int(*halfLife)},
		Kubernetes: Kubernetes{},
	}
	configData, err := yaml.Marshal(newConfig)
	if err != nil {
		return
	}
	_, err = configFile.Write(configData)
	if err != nil {
		return
	}

	// TODO: copy k3s.yaml
	// TODO: generate agent node install script
}
