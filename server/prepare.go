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
	"encoding/json"
	"encoding/pem"
	"flag"
	"fmt"
	"github.com/dustin/go-humanize"
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"math/big"
	rd "math/rand"
	"net"
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

type TransportProgress struct {
	Done  uint64
	Total uint64
}

func (progress *TransportProgress) Write(p []byte) (int, error) {
	n := len(p)
	progress.Done += uint64(n)
	progress.PrintProgress()
	return n, nil
}

func (progress *TransportProgress) PrintProgress() {
	fmt.Printf("\r%s", strings.Repeat(" ", 35))
	fmt.Printf("\r%s / %s", humanize.Bytes(progress.Done), humanize.Bytes(progress.Total))
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
	copyCmd.Stderr = os.Stderr
	copyCmd.Stdin = os.Stdin
	copyCmd.Stdout = os.Stdout
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
	} else if !os.IsNotExist(err) {
		fmt.Println(err)
		os.Exit(1)
	}
	log.Println("download k3s install script...")
	client := &http.Client{}
	k3sReq, _ := http.NewRequest("GET", "https://rancher-mirror.oss-cn-beijing.aliyuncs.com/k3s/k3s-install.sh", nil)
	k3sReq.Header.Set("Accept-Encoding", "*")
	k3sRes, err := client.Do(k3sReq)
	if err != nil {
		fmt.Println("china mirror error, fallback to default")
		k3sReq, _ = http.NewRequest("GET", "https://get.k3s.io/", nil)
		k3sReq.Header.Set("Accept-Encoding", "*")
		k3sRes, err = client.Do(k3sReq)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	k3sInstallSH, err := os.Create("k3s-install.sh")
	if err != nil {
		return
	}
	_, err = io.Copy(k3sInstallSH, io.TeeReader(k3sRes.Body, &TransportProgress{Done: 0, Total: uint64(k3sRes.ContentLength)}))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println()
	err = k3sInstallSH.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = os.Chmod("k3s-install.sh", 0755)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	log.Println("install k3s master...")
	k3sInstall := exec.Command("bash", "-c", fmt.Sprintf("sudo INSTALL_K3S_MIRROR=%s ./k3s-install.sh --write-kubeconfig-mode 644 --node-external-ip %s --node-name %s", os.Getenv("INSTALL_K3S_MIRROR"), masterIP, masterIP))
	k3sInstall.Stderr = os.Stderr
	k3sInstall.Stdin = os.Stdin
	k3sInstall.Stdout = os.Stdout
	//k3sPipe, err := k3sInstall.StdinPipe()
	//if err != nil {
	//	fmt.Println(err)
	//	os.Exit(1)
	//}
	//_, err = io.Copy(k3sPipe, k3sRes.Body)
	//if err != nil {
	//	fmt.Println(err)
	//	os.Exit(1)
	//}
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
	chmodCmd := exec.Command("bash", "-c", "sudo chmod 766 resources/k3s.yaml")
	chmodCmd.Stderr = os.Stderr
	chmodCmd.Stdin = os.Stdin
	chmodCmd.Stdout = os.Stdout
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
	} else if !os.IsNotExist(err) {
		fmt.Println(err)
		os.Exit(1)
	}
	mkdirCmd := exec.Command("bash", "-c", "sudo mkdir /etc/rancher/k3s/OJRegistry")
	mkdirCmd.Stderr = os.Stderr
	mkdirCmd.Stdin = os.Stdin
	mkdirCmd.Stdout = os.Stdout
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
	binaryName := fmt.Sprintf("OJ_%s_%s", runtime.GOOS, runtime.GOARCH)
	_, err := os.Stat(binaryName)
	if err == nil {
		return
	} else if !os.IsNotExist(err) {
		fmt.Println(err)
		os.Exit(1)
	}
	log.Println("get latest release info...")
	client := &http.Client{}
	releaseInfoReq, _ := http.NewRequest("GET", "https://api.github.com/repos/Tp0t-Team/Tp0tOJ/releases/latest", nil)
	releaseInfoReq.Header.Set("Accept-Encoding", "*")
	releaseInfoRes, err := client.Do(releaseInfoReq)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	releaseInfoData := bytes.Buffer{}
	_, err = io.Copy(&releaseInfoData, io.TeeReader(releaseInfoRes.Body, &TransportProgress{Done: 0, Total: uint64(releaseInfoRes.ContentLength)}))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println()
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
	log.Println("download latest release...")
	binaryReq, _ := http.NewRequest("GET", fmt.Sprintf("https://github.com/Tp0t-Team/Tp0tOJ/releases/download/%s/%s", releaseInfo.TagName, binaryName), nil)
	binaryReq.Header.Set("Accept-Encoding", "*")
	binaryRes, err := client.Do(binaryReq)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	binary, err := os.Create(binaryName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	_, err = io.Copy(binary, io.TeeReader(binaryRes.Body, &TransportProgress{Done: 0, Total: uint64(binaryRes.ContentLength)}))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println()
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
	} else if !os.IsNotExist(err) {
		fmt.Println(err)
		os.Exit(1)
	}

	//extSubjectAltName := pkix.Extension{}
	//extSubjectAltName.Id = asn1.ObjectIdentifier{2, 5, 29, 17}
	//extSubjectAltName.Critical = false
	//extSubjectAltName.Value = []byte(fmt.Sprintf("IP:%s", masterIP))

	CACert := &x509.Certificate{
		SerialNumber:          big.NewInt(rd.Int63()),
		BasicConstraintsValid: true,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(1, 0, 0),
		Subject: pkix.Name{
			CommonName: "ca-" + masterIP,
		},
		Issuer: pkix.Name{
			CommonName: "ca-" + masterIP,
		},
		IsCA: true,
		//ExtKeyUsage:    []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		//KeyUsage:       x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature,
		EmailAddresses: []string{},
		//ExtraExtensions: []pkix.Extension{extSubjectAltName},
	}
	Cert := &x509.Certificate{
		SerialNumber:          big.NewInt(rd.Int63()),
		BasicConstraintsValid: true,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(1, 0, 0),
		Subject: pkix.Name{
			CommonName: "ca-" + masterIP,
		},
		Issuer: pkix.Name{
			CommonName: "ca-" + masterIP,
		},
		IsCA:           false,
		ExtKeyUsage:    []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:       x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature,
		EmailAddresses: []string{},
		//ExtraExtensions: []pkix.Extension{extSubjectAltName},
		IPAddresses: []net.IP{net.ParseIP(masterIP)},
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
	err = pem.Encode(caFile, &pem.Block{Type: "CERTIFICATE", Bytes: buf})
	//_, err = caFile.Write(buf)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = caFile.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//CACert, _ = x509.ParseCertificate(buf)
	buf = nil
	buf, err = x509.CreateCertificate(rand.Reader, Cert, CACert, &Key.PublicKey, Key)
	certFile, err := os.Create("resources/tls.crt")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = pem.Encode(certFile, &pem.Block{Type: "CERTIFICATE", Bytes: buf})
	//_, err = certFile.Write(buf)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = certFile.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	buf = nil
	buf = x509.MarshalPKCS1PrivateKey(Key)
	keyFile, err := os.Create("resources/tls.key")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = pem.Encode(keyFile, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: buf})
	//_, err = keyFile.Write(buf)
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
	} else if !os.IsNotExist(err) {
		fmt.Println(err)
		os.Exit(1)
	}
	newConfig := utils.Config{
		Server: utils.Server{
			Host:     "127.0.0.1",
			Username: "Tp0t",
			Password: "admin",
			Mail:     "admin@example.com",
			Port:     0,
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
			PortAllocEnd:   31000,
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
	} else if !os.IsNotExist(err) {
		fmt.Println(err)
		os.Exit(1)
	}
	err = sudoCopy("/var/lib/rancher/k3s/server/node-token", "node-token")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	tokenData := bytes.Buffer{}
	readCmd := exec.Command("bash", "-c", "sudo chmod 0777 node-token")
	readCmd.Stderr = os.Stderr
	readCmd.Stdin = os.Stdin
	readCmd.Stdout = os.Stdout
	//tokenPipe, err := readCmd.StdoutPipe()
	//if err != nil {
	//	fmt.Println(err)
	//	os.Exit(1)
	//}
	err = readCmd.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	tokenFile, err := os.Open("node-token")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	_, err = io.Copy(&tokenData, tokenFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = tokenFile.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = os.Remove("node-token")
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

const RegistryConfigPath = "resources/docker-registry"

func PrepareRegistry(masterIP string, registryUsername string, registryPassword string) {
	_, err := os.Stat(fmt.Sprintf("/etc/docker/certs.d/%s:5000/ca.crt", masterIP))
	if err == nil {
		return
	} else if !os.IsNotExist(err) {
		fmt.Println(err)
		os.Exit(1)
	}

	err = os.MkdirAll(RegistryConfigPath, 0755)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	log.Println("make htpasswd file...")
	htpasswdCmdOut := bytes.Buffer{}
	htpasswdCmd := exec.Command("docker", "run", "--rm", "--entrypoint", "htpasswd", "httpd:alpine", "-Bbn", registryUsername, registryPassword)
	//pipe, err := htpasswdCmd.StdoutPipe()
	htpasswdCmd.Stdout = &htpasswdCmdOut
	//if err != nil {
	//	fmt.Println(err)
	//	os.Exit(1)
	//}
	err = htpasswdCmd.Run()
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
	_, err = io.Copy(passwdFile, &htpasswdCmdOut)
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

	mkdirCmd := exec.Command("bash", "-c", fmt.Sprintf("sudo mkdir -p /etc/docker/certs.d/%s:5000", masterIP))
	mkdirCmd.Stderr = os.Stderr
	mkdirCmd.Stdin = os.Stdin
	mkdirCmd.Stdout = os.Stdout
	err = mkdirCmd.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = sudoCopy("resources/ca.crt", fmt.Sprintf("/etc/docker/certs.d/%s:5000/ca.crt", masterIP))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func GenerateStartScript() {
	_, err := os.Stat("start.sh")
	if err == nil {
		return
	} else if !os.IsNotExist(err) {
		fmt.Println(err)
		os.Exit(1)
	}

	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fileName := fmt.Sprintf("./OJ_%s_%s\n", runtime.GOOS, runtime.GOARCH)
	tempString := "if test -z \"$(docker ps | grep oj_registry_instance)\"; then\n" +
		"\tdocker run -d --net=host --restart=always --name oj_registry_instance " +
		"-v %s/data:/var/lib/registry " +
		"-v %s/auth:/auth " +
		"-e \"REGISTRY_AUTH=htpasswd\" " +
		"-e \"REGISTRY_AUTH_HTPASSWD_REALM=Registry Realm\" " +
		"-e REGISTRY_AUTH_HTPASSWD_PATH=/auth/htpasswd " +
		"-v %s/certs:/certs " +
		"-e REGISTRY_HTTP_TLS_CERTIFICATE=/certs/tls.crt " +
		"-e REGISTRY_HTTP_TLS_KEY=/certs/tls.key " +
		"-e REGISTRY_STORAGE_DELETE_ENABLED=true " +
		"registry\n" +
		"fi\n" +
		fileName
	startSH, err := os.Create("start.sh")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	path := pwd + "/" + RegistryConfigPath
	_, err = startSH.Write([]byte(fmt.Sprintf(tempString, path, path, path)))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = startSH.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = os.Chmod("start.sh", 0755)
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
	if _, err := os.Stat("resources"); os.IsNotExist(err) {
		err = os.MkdirAll("resources", 0755)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	registryUsername := fmt.Sprintf("%02x", md5.Sum([]byte(strconv.FormatInt(rd.Int63(), 10))))[:8]
	registryPassword := fmt.Sprintf("%02x", md5.Sum([]byte(strconv.FormatInt(rd.Int63(), 10))))[:8]

	log.Println(" - install K3S:")
	InstallK3S(*masterIP)
	log.Println(" - create cert:")
	CreateCert(*masterIP)
	log.Println(" - config registry:")
	ConfigK3SRegistry(*masterIP)
	log.Println(" - download binary:")
	DownloadBinary()
	log.Println(" - generate default config:")
	CreateDefaultConfig(*masterIP, registryUsername, registryPassword)
	log.Println(" - generate agent script:")
	GenerateAgentScript(*masterIP)
	log.Println(" - prepare registry:")
	PrepareRegistry(*masterIP, registryUsername, registryPassword)
	log.Println(" - generate start script:")
	GenerateStartScript()

	if err := os.Remove("registries-config.yaml"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("\nCongratulations! Environment prepare success.\n" +
		"You can start the platform use: `./start.sh`\n" +
		"You can easily Add more agent machine by using agent-install.sh\n" +
		"For example, on the agent machine run:\n" +
		"    ./agent-install.sh <ip of agent machine>")
}
