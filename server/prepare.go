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
	"io"
	"math/big"
	rd "math/rand"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"
)

type ReleaseInfo struct {
	TagName string `json:"tag_name"`
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

	// TODO: make default config.yaml
	// TODO: copy k3s.yaml
	// TODO: generate agent node install script
}
