package kube

var kubeConfig string

func init() {
	kubeConfig = "/etc/rancher/k3s/k3s.yaml"
}

func k8sPodAlloc() {

}

func k8sPodList() {

}

func k8sPodDestroy() {

}

func k8sPodStatus() {

}

func dockerFileUpload() {

}

func imgBuild() {

}

func imgDelete() {

}

func imgStatus() {

}

//
//func main() {
//
//	//把用户传递的命令行参数，解析为响应变量的值
//	flag.Parse()
//	//加载kubeconfig中的apiserver地址、证书配置等信息
//	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
//	if err != nil {
//		log.Panicln(err)
//	}
//
//	//NewForConfig为给定的配置创建一个新的Clientset（如下图所示包含所有的api-versions，这样做的目的是便于其它
//	//资源类型对这个Pod进行管理和控制？）。
//	clientset, err := kubernetes.NewForConfig(config)
//	if err != nil {
//		log.Panicln(err)
//	}
//	deploymentsClient := clientset.AppsV1().Deployments(corev1.NamespaceDefault)
//	deployment := &appsv1.Deployment{
//		ObjectMeta: metav1.ObjectMeta{
//			Name: "test",
//		}, Spec: appsv1.DeploymentSpec{
//			Replicas: pointer.Int32Ptr(1),
//			Selector: &metav1.LabelSelector{
//				MatchLabels: map[string]string{
//					"app": "demo",
//				},
//			},
//			Template: corev1.PodTemplateSpec{
//				ObjectMeta: metav1.ObjectMeta{
//					Labels: map[string]string{
//						"app": "demo",
//					},
//				},
//				Spec: corev1.PodSpec{
//					Containers: []corev1.Container{
//						{
//							Name:  "echo",
//							Image: "ealen/echo-server",
//							Ports: []corev1.ContainerPort{
//								{
//									Name:          "http",
//									Protocol:      corev1.ProtocolTCP,
//									ContainerPort: 80,
//								},
//							},
//						},
//					},
//				},
//			},
//		},
//	}
//	service := &corev1.Service{
//		ObjectMeta: metav1.ObjectMeta{
//			Name: "test",
//		},
//		Spec: corev1.ServiceSpec{
//			Ports: []corev1.ServicePort{
//				{
//					Port:     8080,
//					Protocol: corev1.ProtocolTCP,
//					Name:     "http",
//					TargetPort: intstr.IntOrString{
//						Type:   intstr.Int,
//						IntVal: 80,
//					},
//					NodePort: 30001,
//				},
//			},
//			Selector: map[string]string{
//				"app": "demo",
//			},
//			Type: corev1.ServiceTypeNodePort,
//		},
//	}
//	//ingress := networkv1.Ingress{
//	//    TypeMeta:   metav1.TypeMeta{},
//	//    ObjectMeta: metav1.ObjectMeta{},
//	//    Spec:       networkv1.IngressSpec{},
//	//    Status:     networkv1.IngressStatus{},
//	//}
//	_, err = deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
//	if err != nil {
//		log.Panicln(err)
//	}
//	list, err := clientset.CoreV1().Pods(corev1.NamespaceDefault).List(context.TODO(), metav1.ListOptions{
//		LabelSelector: labels.SelectorFromSet(map[string]string{"app": "demo"}).String(),
//	})
//	if err != nil {
//		return
//	}
//	deployment.Spec.Template.Spec.NodeName = list.Items[0].Spec.NodeName
//	_, err = deploymentsClient.Update(context.TODO(), deployment, metav1.UpdateOptions{})
//	if err != nil {
//		return
//	}
//	_, err = clientset.CoreV1().Services(corev1.NamespaceDefault).Create(context.TODO(), service, metav1.CreateOptions{})
//	if err != nil {
//		log.Panicln(err)
//	}
//}
