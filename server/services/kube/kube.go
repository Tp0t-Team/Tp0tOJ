package kube

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/heroku/docker-registry-client/registry"
	"io/ioutil"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/utils/pointer"
	"log"
	"os"
	"server/utils"
	"strconv"
)

//var kubeConfig *rest.Config
var clientSet *kubernetes.Clientset
var dockerClient *client.Client
var dockerPushAuth string
var registryClient *registry.Registry

type portAllocInfo struct {
	allocated map[int32]struct{}
	current   int32
}

var AutoPortSet map[string]*portAllocInfo

func init() {
	const config = "/etc/rancher/k3s/k3s.yaml"
	var kubeConfig *rest.Config
	var err error
	kubeConfig, err = clientcmd.BuildConfigFromFlags("", config)
	if err != nil {
		log.Panicln(err)
	}
	clientSet, err = kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		log.Panicln(err)
	}
	//AutoPortSet = map[string]*portAllocInfo{}
	autoPortSetLoad()

	dockerClient, err = client.NewClientWithOpts(client.WithTLSClientConfig("resources/ca.crt", "resources/tls.crt", "resources/tls.key"))
	if err != nil {
		log.Panicln(err)
	}

	authConfig := types.AuthConfig{
		Username:      utils.Configure.Kubernetes.Username,
		Password:      utils.Configure.Kubernetes.Password,
		ServerAddress: utils.Configure.Kubernetes.RegistryHost,
	}
	marshal, err := json.Marshal(authConfig)
	if err != nil {
		log.Panicln(err)
	}
	dockerPushAuth = base64.URLEncoding.EncodeToString(marshal)

	registryClient, err = registry.NewInsecure(utils.Configure.Kubernetes.RegistryHost, utils.Configure.Kubernetes.Username, utils.Configure.Kubernetes.Password)
	if err != nil {
		log.Panicln(err)
	}
}

type serialType struct {
	ports   []int32
	current int32
}

func autoPortSetSave() {
	temp := map[string]*serialType{}
	for key, info := range AutoPortSet {
		temp[key] = &serialType{ports: []int32{}, current: info.current}
		for port, _ := range info.allocated {
			temp[key].ports = append(temp[key].ports, port)
		}
	}
	marshal, err := json.Marshal(AutoPortSet)
	if err != nil {
		log.Panicln(err)
	}
	err = ioutil.WriteFile("resources/auto-port-alloc-cache.json", marshal, 0600)
	if err != nil {
		log.Panicln(err)
	}
}

func autoPortSetLoad() {
	AutoPortSet = map[string]*portAllocInfo{}
	if _, err := os.Stat("resources/auto-port-alloc-cache.json"); err != nil {
		return
	}
	file, err := ioutil.ReadFile("resources/auto-port-alloc-cache.json")
	if err != nil {
		log.Panicln(err)
	}
	var temp map[string]*serialType
	err = json.Unmarshal(file, &temp)
	if err != nil {
		log.Panicln(err)
	}
	for key, desc := range temp {
		AutoPortSet[key] = &portAllocInfo{allocated: map[int32]struct{}{}, current: desc.current}
		for _, port := range desc.ports {
			AutoPortSet[key].allocated[port] = struct{}{}
		}
	}
}

func ParseProtocol(name string) corev1.Protocol {
	if name == "TCP" {
		return corev1.ProtocolTCP
	} else if name == "UDP" {
		return corev1.ProtocolUDP
	} else {
		return corev1.ProtocolTCP
	}
}

func NewContainerPortConfig(protocol corev1.Protocol, containerPort int32) *corev1.ContainerPort {
	return &corev1.ContainerPort{Name: strconv.FormatInt(int64(containerPort), 10), Protocol: protocol, ContainerPort: containerPort}
}

//NewServicePortConfig
//portName is like 'ssh'\'ftp' etc. ,this can be anyone you like, but must be unique for same pod.
//externalPort is exposed port for container, and internalPort is service port inside of container, podPort is port you can visit it on the internet.
func NewServicePortConfig(portName string, protocol corev1.Protocol, externalPort int32, internalPort int32, podPort int32) *corev1.ServicePort {
	return &corev1.ServicePort{Name: portName, Protocol: protocol, Port: externalPort, TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: internalPort}, NodePort: podPort}
}

//K8sPodAlloc
func K8sPodAlloc(replicaId uint64, containerName string, imgLabel string, portConfigs []corev1.ContainerPort, servicePorts []corev1.ServicePort, flag string) bool {
	id := strconv.FormatUint(replicaId, 10) + containerName
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: id,
		}, Spec: appsv1.DeploymentSpec{
			Replicas: pointer.Int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": id,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": id,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  containerName,
							Image: imgLabel,
							Ports: portConfigs,
							Env: []corev1.EnvVar{
								{
									Name:  "FLAG",
									Value: flag,
								},
							},
						},
					},
				},
			},
		},
	}
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: id,
		},
		Spec: corev1.ServiceSpec{
			Ports: servicePorts,
			Selector: map[string]string{
				"app": id,
			},
			Type: corev1.ServiceTypeNodePort,
		},
	}
	var err error
	_, err = clientSet.AppsV1().Deployments(corev1.NamespaceDefault).Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		// don't panic & rollback
		log.Println(err)
		return false
	}
	var list *corev1.PodList
	list, err = clientSet.CoreV1().Pods(corev1.NamespaceDefault).List(context.TODO(), metav1.ListOptions{
		LabelSelector: labels.SelectorFromSet(map[string]string{"app": id}).String(),
	})
	if err != nil {
		// don't panic & rollback
		log.Println(err)
		return false
	}
	deployment.Spec.Template.Spec.NodeName = list.Items[0].Spec.NodeName
	_, err = clientSet.AppsV1().Deployments(corev1.NamespaceDefault).Update(context.TODO(), deployment, metav1.UpdateOptions{})
	if err != nil {
		log.Println(err)
		return false
	}
	// auto alloc port
	if _, ok := AutoPortSet[deployment.Spec.Template.Spec.NodeName]; !ok {
		AutoPortSet[deployment.Spec.Template.Spec.NodeName] = &portAllocInfo{
			allocated: map[int32]struct{}{},
			current:   utils.Configure.Kubernetes.PortAllocBegin,
		}
	}
	portMaxSize := int(utils.Configure.Kubernetes.PortAllocEnd - utils.Configure.Kubernetes.PortAllocBegin)
	host := deployment.Spec.Template.Spec.NodeName
	for index, _ := range service.Spec.Ports {
		if service.Spec.Ports[index].NodePort == 0 {
			if len(AutoPortSet[host].allocated) >= portMaxSize {
				log.Println(host + " has not enough ports to alloc")
				return false
			}
			for {
				if AutoPortSet[host].current >= utils.Configure.Kubernetes.PortAllocEnd {
					AutoPortSet[host].current = utils.Configure.Kubernetes.PortAllocBegin
				}
				if _, ok := AutoPortSet[host].allocated[AutoPortSet[host].current]; !ok {
					service.Spec.Ports[index].NodePort = AutoPortSet[host].current
					AutoPortSet[host].allocated[AutoPortSet[host].current] = struct{}{}
					AutoPortSet[host].current += 1
					break
				}
				AutoPortSet[host].current += 1
			}
		}
	}

	autoPortSetSave()

	_, err = clientSet.CoreV1().Services(corev1.NamespaceDefault).Create(context.TODO(), service, metav1.CreateOptions{})
	if err != nil {
		// don't panic & rollback
		log.Println(err)
		return false
	}
	return true
}

func K8sPodList() {

}

func K8sPodDestroy(replicaId uint64, containerName string) bool {
	id := strconv.FormatUint(replicaId, 10) + containerName
	var deployment *appsv1.Deployment
	var err error
	deployment, err = clientSet.AppsV1().Deployments(corev1.NamespaceDefault).Get(context.TODO(), id, metav1.GetOptions{})
	if err != nil {
		log.Println(err)
		return false
	}
	host := ""
	if deployment != nil {
		host = deployment.Spec.Template.Spec.NodeName
		err := clientSet.AppsV1().Deployments(corev1.NamespaceDefault).Delete(context.TODO(), id, metav1.DeleteOptions{})
		if err != nil {
			log.Println(err)
			return false
		}
	}
	var service *corev1.Service
	service, err = clientSet.CoreV1().Services(corev1.NamespaceDefault).Get(context.TODO(), id, metav1.GetOptions{})
	if err != nil {
		log.Println(err)
		return false
	}
	if service != nil {
		err := clientSet.CoreV1().Services(corev1.NamespaceDefault).Delete(context.TODO(), id, metav1.DeleteOptions{})
		if err != nil {
			log.Println(err)
			return false
		}
		if host != "" {
			for _, port := range service.Spec.Ports {
				if port.NodePort >= utils.Configure.Kubernetes.PortAllocBegin && port.NodePort < utils.Configure.Kubernetes.PortAllocEnd {
					delete(AutoPortSet[host].allocated, port.NodePort)
				}
			}
		}
	}
	return true
}

func K8sPodStatus() {

}

func K8sServiceGetUrls(replicaId uint64, containerName string) []string {
	id := strconv.FormatUint(replicaId, 10) + containerName
	var deployment *appsv1.Deployment
	var err error
	deployment, err = clientSet.AppsV1().Deployments(corev1.NamespaceDefault).Get(context.TODO(), id, metav1.GetOptions{})
	if err != nil {
		return nil
	}
	host := deployment.Spec.Template.Spec.NodeName
	var service *corev1.Service
	service, err = clientSet.CoreV1().Services(corev1.NamespaceDefault).Get(context.TODO(), id, metav1.GetOptions{})
	if err != nil {
		log.Println(err)
		return nil
	}
	if service == nil {
		return nil
	}
	result := []string{}
	for _, port := range service.Spec.Ports {
		url := host + ":" + strconv.FormatInt(int64(port.NodePort), 10)
		result = append(result, url)
	}
	return result
}

func DockerFileUpload() {

}

func ImgBuild(tarArchive string, imageName string) error {
	file, err := os.Open(tarArchive)
	if err != nil {
		return err
	}
	_, err = dockerClient.ImageBuild(context.TODO(), file, types.ImageBuildOptions{
		Dockerfile: "dockerfile",
		Tags:       []string{imageName},
		Remove:     true,
	})
	err = file.Close()
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	_, err = dockerClient.ImagePush(context.TODO(), imageName, types.ImagePushOptions{
		RegistryAuth: dockerPushAuth,
	})
	if err != nil {
		return err
	}
	_, err = dockerClient.ImageRemove(context.TODO(), imageName, types.ImageRemoveOptions{})
	if err != nil {
		return err
	}
	return nil
}

func ImgDelete(imageName string) error {
	digest, err := registryClient.ManifestDigest(imageName, "latest")
	if err != nil {
		return err
	}
	err = registryClient.DeleteManifest(imageName, digest)
	if err != nil {
		return err
	}
	return nil
}

func ImgStatus() {

}
