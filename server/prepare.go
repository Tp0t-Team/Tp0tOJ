package main

import (
	"archive/tar"
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	_ "embed"
	"encoding/asn1"
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
	"server/utils"
	"strconv"
	"strings"
	"time"
)

//go:embed unpack.sh
var unpackShell string

type ReleaseInfo struct {
	TagName string `json:"tag_name"`
}

type K3sRegistryMirrorItem struct {
	Endpoint []string `yaml:"endpoint"`
}
type K3sRegistryTLSConfigItem struct {
	CAFile   string `yaml:"ca_file"`
	CertFile string `yaml:"cert_file"`
	KeyFile  string `yaml:"key_file"`
}
type K3sRegistryConfigItem struct {
	TLS K3sRegistryTLSConfigItem `yaml:"tls"`
}
type K3sRegistry struct {
	Mirrors map[string]K3sRegistryMirrorItem `yaml:"mirrors"`
	Configs map[string]K3sRegistryConfigItem `yaml:"configs"`
}

func sudoCopy(src string, dst string) error {
	copyCmd := exec.Command("bash", "-c", fmt.Sprintf("sudo cp %s %s", src, dst))
	err := copyCmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func InstallK3S(masterIP string) {
	_, err := os.Stat("resources/k3s.yaml")
	if err == nil {
		return
	} else if err != os.ErrNotExist {
		fmt.Println(err)
		os.Exit(1)
	}
	k3sRes, err := http.Get("http://rancher-mirror.cnrancher.com/k3s/k3s-install.sh")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	k3sInstall := exec.Command("bash", "-c", fmt.Sprintf("sudo sh -s - --node-external-ip %s --node-name %s", masterIP, masterIP))
	k3sPipe, err := k3sInstall.StdinPipe()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	_, err = io.Copy(k3sPipe, k3sRes.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = k3sInstall.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = sudoCopy("/etc/rancher/k3s/k3s.yaml", "resources/k3s.yaml")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	chmodCmd := exec.Command("bash", "-c", "sudo chmod 0555 resources/k3s.yaml")
	err = chmodCmd.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func ConfigK3SRegistry(masterIP string) {
	_, err := os.Stat("/etc/rancher/k3s/registries.yaml")
	if err == nil {
		return
	} else if err != os.ErrNotExist {
		fmt.Println(err)
		os.Exit(1)
	}
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	mkdirCmd := exec.Command("bash", "-c", "sudo mkdir /etc/rancher/k3s/OJRegistry")
	err = mkdirCmd.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = sudoCopy("resources/ca.crt", " /etc/rancher/k3s/OJRegistry/ca.crt")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = sudoCopy("resources/tls.crt", " /etc/rancher/k3s/OJRegistry/tls.crt")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = sudoCopy("resources/tls.key", " /etc/rancher/k3s/OJRegistry/tls.key")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	registryHost := fmt.Sprintf("%s:5000", masterIP)
	registries := K3sRegistry{
		Mirrors: map[string]K3sRegistryMirrorItem{
			registryHost: {Endpoint: []string{"https://" + registryHost}},
		},
		Configs: map[string]K3sRegistryConfigItem{
			registryHost: {
				TLS: K3sRegistryTLSConfigItem{
					CAFile:   "/etc/rancher/k3s/OJRegistry/ca.crt",
					CertFile: "/etc/rancher/k3s/OJRegistry/tls.crt",
					KeyFile:  "/etc/rancher/k3s/OJRegistry/tls.key",
				},
			},
		},
	}
	registryFile, err := os.Create("registries-config.yaml")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	registryData, err := yaml.Marshal(registries)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	_, err = registryFile.Write(registryData)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = registryFile.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = sudoCopy("registries-config.yaml", "/etc/rancher/k3s/registries.yaml")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func DownloadBinary() {
	_, err := os.Stat("OJ")
	if err == nil {
		return
	} else if err != os.ErrNotExist {
		fmt.Println(err)
		os.Exit(1)
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
}

func CreateCert(masterIP string) {
	_, err := os.Stat("resources/tls.key")
	if err == nil {
		return
	} else if err != os.ErrNotExist {
		fmt.Println(err)
		os.Exit(1)
	}

	extSubjectAltName := pkix.Extension{}
	extSubjectAltName.Id = asn1.ObjectIdentifier{2, 5, 29, 17}
	extSubjectAltName.Critical = false
	extSubjectAltName.Value = []byte(fmt.Sprintf("IP:%s", masterIP))

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
		ExtraExtensions:       []pkix.Extension{extSubjectAltName},
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
		ExtraExtensions:       []pkix.Extension{extSubjectAltName},
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
}

func CreateDefaultConfig(masterIP string, registryUsername string, registryPassword string) {
	_, err := os.Stat("resources/config.yaml")
	if err == nil {
		return
	} else if err != os.ErrNotExist {
		fmt.Println(err)
		os.Exit(1)
	}
	newConfig := utils.Config{
		Server: utils.Server{
			Host:     "127.0.0.1",
			Username: "Tp0t",
			Password: "admin",
			Mail:     "admin@example.com",
			Port:     8080,
			Salt:     strconv.FormatInt(rd.Int63(), 10),
			Debug:    false,
		},
		Email: utils.Email{
			Host:     "smtp.example.com",
			Username: "exampleUsername",
			Password: "examplePassword",
		},
		Challenge: utils.Challenge{
			FirstBloodReward:  0.10,
			SecondBloodReward: 0.08,
			ThirdBloodReward:  0.05,
			HalfLife:          20,
		},
		Kubernetes: utils.Kubernetes{
			PortAllocBegin: 30000,
			PortAllocEnd:   40000,
			Username:       registryUsername,
			Password:       registryPassword,
			RegistryHost:   fmt.Sprintf("%s:5000", masterIP),
		},
	}
	configFile, err := os.Create("resources/config.yaml")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	configData, err := yaml.Marshal(newConfig)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	_, err = configFile.Write(configData)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = configFile.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func TarAddFile(file string, archiveFile string, tarArchive *tar.Writer) error {
	configInfo, err := os.Stat(file)
	if err != nil {
		return err
	}
	configHeader, err := tar.FileInfoHeader(configInfo, "")
	if err != nil {
		return err
	}
	configHeader.Name = archiveFile
	err = tarArchive.WriteHeader(configHeader)
	if err != nil {
		return err
	}
	configFile, err := os.Open(file)
	if err != nil {
		return err
	}
	_, err = io.Copy(tarArchive, configFile)
	if err != nil {
		return err
	}
	err = configFile.Close()
	if err != nil {
		return err
	}
	return nil
}

func GenerateAgentScript(masterIP string) {
	_, err := os.Stat("agent-install.sh")
	if err == nil {
		return
	} else if err != os.ErrNotExist {
		fmt.Println(err)
		os.Exit(1)
	}
	tokenData := bytes.Buffer{}
	readCmd := exec.Command("bash", "-c", "sudo cat /var/lib/rancher/k3s/server/node-token")
	tokenPipe, err := readCmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = readCmd.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	_, err = io.Copy(&tokenData, tokenPipe)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	token := strings.TrimSpace(string(tokenData.Bytes()))
	k3sCmdSting := fmt.Sprintf("curl -sfL http://rancher-mirror.cnrancher.com/k3s/k3s-install.sh | K3S_URL=https://%s:6443/ K3S_TOKEN=%s sh -s - --node-external-ip $1 --node-name $1\n", masterIP, token)

	file, err := os.Create("agent-install.sh")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	_, err = file.Write([]byte(unpackShell))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	_, err = file.Write([]byte(k3sCmdSting))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	cmdBlock := "sudo cp registries-config.yaml /etc/rancher/k3s/registries.yaml\n" +
		"sudo cp -r OJRegistry /etc/rancher/k3s\n" +
		"rm registries-config.yaml\n" +
		"rm -r OJRegistry\n" +
		"exit 0\n"
	_, err = file.Write([]byte(cmdBlock))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	_, err = file.Write([]byte("__ARCHIVE_BELOW__\n"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	tarBuffer := bytes.Buffer{}
	tarArchive := tar.NewWriter(&tarBuffer)
	err = TarAddFile("registries-config.yaml", "registries-config.yaml", tarArchive)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = TarAddFile("resources/ca.crt", "OJRegistry/ca.crt", tarArchive)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = TarAddFile("resources/tls.crt", "OJRegistry/tls.crt", tarArchive)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = TarAddFile("resources/tls.key", "OJRegistry/tls.key", tarArchive)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = tarArchive.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	_, err = io.Copy(file, &tarBuffer)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = file.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = os.Chmod("agent-install.sh", 0755)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

const RegistryConfigPath = "/opt/docker-registry" // TODO:

func StartRegistry(masterIP string, registryUsername string, registryPassword string) {
	_, err := os.Stat(fmt.Sprintf("/etc/docker/certs.d/%s:5000/ca.crt", masterIP))
	if err == nil {
		return
	} else if err != os.ErrNotExist {
		fmt.Println(err)
		os.Exit(1)
	}

	htpasswdCmd := exec.Command("docker", "run", "--rm", "--entrypoint", "htpasswd", "httpd:alpine", "-Bbn", registryUsername, registryPassword)
	err = htpasswdCmd.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	pipe, err := htpasswdCmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = os.MkdirAll(RegistryConfigPath+"/auth", 0755)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	passwdFile, err := os.Create(RegistryConfigPath + "/auth/htpasswd")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	_, err = io.Copy(passwdFile, pipe)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = passwdFile.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = os.MkdirAll(RegistryConfigPath+"/certs", 0755)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	src, err := os.Open("resources/tls.key")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	dst, err := os.Create(RegistryConfigPath + "/certs/tls.key")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	_, err = io.Copy(dst, src)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = src.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = dst.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	src, err = os.Open("resources/tls.crt")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	dst, err = os.Create(RegistryConfigPath + "/certs/tls.crt")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	_, err = io.Copy(dst, src)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = src.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = dst.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = os.MkdirAll(RegistryConfigPath+"/data", 0755)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// TODO: start docker here?
	//docker run -d -p 5000:5000 --restart=always --name registry \
	//-v $RegistryPath/data:/var/lib/registry \
	//-v $RegistryPath/auth:/auth \
	//-e "REGISTRY_AUTH=htpasswd" \
	//-e "REGISTRY_AUTH_HTPASSWD_REALM=Registry Realm" \
	//-e REGISTRY_AUTH_HTPASSWD_PATH=/auth/htpasswd \
	//-v $RegistryPath/certs:/certs \
	//-e REGISTRY_HTTP_TLS_CERTIFICATE=/certs/tls.crt \
	//-e REGISTRY_HTTP_TLS_KEY=/certs/tls.key \
	//registry

	err = sudoCopy("resources/ca.crt", fmt.Sprintf("/etc/docker/certs.d/%s:5000/ca.crt", masterIP))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	masterIP := flag.String("MasterIP", "", "master ip")

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
	if _, err := os.Stat("resources"); err == os.ErrNotExist {
		err = os.MkdirAll("resources", 0755)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	registryUsername := fmt.Sprintf("%02x", md5.Sum([]byte(strconv.FormatInt(rd.Int63(), 10))))[:8]
	registryPassword := fmt.Sprintf("%02x", md5.Sum([]byte(strconv.FormatInt(rd.Int63(), 10))))[:8]

	InstallK3S(*masterIP)
	CreateCert(*masterIP)
	ConfigK3SRegistry(*masterIP)
	DownloadBinary()
	CreateDefaultConfig(*masterIP, registryUsername, registryPassword)
	GenerateAgentScript(*masterIP)
	StartRegistry(*masterIP, registryUsername, registryPassword)

	if err := os.Remove("registries.yaml"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Congratulations! Environment prepare success.\n" +
		"You can start the platform use: `./OJ`\n" +
		"You can easily Add more agent machine by using agent-install.sh\n" +
		"For example, on the agent machine run:\n" +
		"    ./agent-install.sh <ip of agent machine>")
}
