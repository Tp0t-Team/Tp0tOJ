package kube

import corev1 "k8s.io/api/core/v1"

var DeletingReplicas = map[string]interface{}{}

type AllocTask struct {
	ReplicaId     uint64
	ContainerName string
	ImgLabel      string
	ServicePorts  []corev1.ServicePort
	Flag          string
}

type DestroyTask struct {
	ReplicaId     uint64
	ContainerName string
}

type Task struct {
	Tasks []interface{}
	CB    func(status bool)
}

var TaskQ chan Task

func init() {
	TaskQ = make(chan Task)
	go func() {
		for {
			desc := <-TaskQ
			success := true
			for _, item := range desc.Tasks {
				if cased, ok := item.(*AllocTask); ok {
					if !K8sPodAlloc(cased.ReplicaId, cased.ContainerName, cased.ImgLabel, cased.ServicePorts, cased.Flag) {
						success = false
						break
					}
				} else if cased, ok := item.(*DestroyTask); ok {
					if !K8sPodDestroy(cased.ReplicaId, cased.ContainerName) {
						success = false
						break
					}
				}
			}
			if desc.CB != nil {
				desc.CB(success)
			}
		}
	}()
}
