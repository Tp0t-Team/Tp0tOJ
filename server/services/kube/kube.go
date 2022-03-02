package kube

import (
	"context"
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
	"strconv"
)

//var kubeConfig *rest.Config
var clientSet *kubernetes.Clientset

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
func K8sPodAlloc(replicaId uint64, containerName string, imgLabel string, portConfigs []corev1.ContainerPort, servicePorts []corev1.ServicePort) {
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
		log.Panicln(err)
	}
	var list *corev1.PodList
	list, err = clientSet.CoreV1().Pods(corev1.NamespaceDefault).List(context.TODO(), metav1.ListOptions{
		LabelSelector: labels.SelectorFromSet(map[string]string{"app": id}).String(),
	})
	if err != nil {
		return
	}
	deployment.Spec.Template.Spec.NodeName = list.Items[0].Spec.NodeName
	_, err = clientSet.AppsV1().Deployments(corev1.NamespaceDefault).Update(context.TODO(), deployment, metav1.UpdateOptions{})
	if err != nil {
		return
	}
	_, err = clientSet.CoreV1().Services(corev1.NamespaceDefault).Create(context.TODO(), service, metav1.CreateOptions{})
	if err != nil {
		log.Panicln(err)
	}

}

func K8sPodList() {

}

func K8sPodDestroy() {

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

func ImgBuild() {

}

func ImgDelete() {

}

func ImgStatus() {

}
