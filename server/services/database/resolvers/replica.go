package resolvers

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	corev1 "k8s.io/api/core/v1"
	"log"
	unsafeRand "math/rand"
	"server/entity"
	"server/services/kube"
	"server/services/types"
	"strconv"
)

func FindReplicaByChallengeId(challengeId uint64) []entity.Replica {
	var replicas []entity.Replica
	result := db.Where(map[string]interface{}{"challenge_id": challengeId}).Find(&replicas)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return []entity.Replica{}
	} else if result.Error != nil {
		log.Println(result.Error)
		return nil
	}
	return replicas
}

func AddReplica(challengeId uint64, outsideTX *gorm.DB) *entity.Replica {
	if outsideTX == nil {
		outsideTX = db
	}
	var newReplica entity.Replica
	err := outsideTX.Transaction(func(tx *gorm.DB) error {
		var challenge entity.Challenge
		challengeResult := tx.Where(map[string]interface{}{"challenge_id": challengeId}).First(&challenge)
		if challengeResult.Error != nil {
			return challengeResult.Error
		}
		if challenge.State == "disabled" {
			return errors.New("unable to create replica of disabled challenge")
		}
		var config types.ChallengeConfig
		err := json.Unmarshal([]byte(challenge.Configuration), &config)
		if err != nil {
			return err
		}
		var flag string
		if config.Flag.Dynamic {
			//dynamic flag should be generated in here,and transfer to starting Pod as environment value,
			//the entrypoint.sh of dockerfile need to clear the FLAG environment value before starting the service
			seed := make([]byte, 8)
			_, err := rand.Read(seed)
			if err != nil {
				return err
			}
			unsafeRand.Seed(int64(binary.BigEndian.Uint64(seed)))
			init := make([]byte, 16)
			_, err = unsafeRand.Read(init)
			if err != nil {
				return err
			}
			flag = fmt.Sprintf("%x", md5.Sum(init))
		} else {
			//strings.Split(config.Flag.Value,"\n")
			//TODO: for muti-flag typed challenge, may need some extra method to map flags to replica
			flag = config.Flag.Value
		}
		newReplica = entity.Replica{
			ChallengeId: challengeId,
			Status:      "disabled",
			Flag:        flag,
		}
		result := tx.Create(&newReplica)
		if result.Error != nil {
			return result.Error
		}
		if ok := EnableReplica(newReplica.ReplicaId, tx); !ok {
			return errors.New("enable replica failed")
		}
		return nil
	})
	if err != nil {
		log.Println(err)
		return nil
	}
	return &newReplica
}

func EnableReplica(replicaId uint64, outsideTX *gorm.DB) bool {
	if outsideTX == nil {
		outsideTX = db
	}
	err := outsideTX.Transaction(func(tx *gorm.DB) error {
		var replica entity.Replica
		getResult := tx.Where(map[string]interface{}{"replica_id": replicaId}).First(&replica)
		if getResult.Error != nil {
			return getResult.Error
		}
		if replica.Status == "enabled" {
			return nil
		}
		replica.Status = "enabled"

		var config types.ChallengeConfig
		err := json.Unmarshal([]byte(replica.Challenge.Configuration), &config)
		if err != nil {
			return err
		}
		createPodSuccess := true
		for _, node := range config.NodeConfig {
			ports := []corev1.ContainerPort{}
			for _, port := range node.Ports {
				ports = append(ports, *kube.NewContainerPortConfig(kube.ParseProtocol(port.Protocol), port.Port))
			}
			servicePorts := []corev1.ServicePort{}
			for _, port := range node.ServicePorts {
				servicePorts = append(servicePorts, *kube.NewServicePortConfig(port.Name, kube.ParseProtocol(port.Protocol), port.External, port.Internal, port.Pod))
			}
			ok := kube.K8sPodAlloc(replicaId, node.Name, node.Image, ports, servicePorts, replica.Flag)
			if !ok {
				createPodSuccess = false
				break
			}
		}
		if !createPodSuccess {
			for _, node := range config.NodeConfig {
				ok := kube.K8sPodDestroy(replicaId, node.Name)
				if !ok {
					// TODO: this a uncorrectable error
					return errors.New("alloc pod failed & rollback pod failed - replica: " + strconv.FormatUint(replicaId, 10))
				}
			}
			return errors.New("create pod failed")
		}

		tx.Save(&replica)
		return nil
	})
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func DisableReplica(replicaId uint64, outsideTX *gorm.DB) bool {
	if outsideTX == nil {
		outsideTX = db
	}
	err := outsideTX.Transaction(func(tx *gorm.DB) error {
		var replica entity.Replica
		getResult := tx.Where(map[string]interface{}{"replica_id": replicaId}).First(&replica)
		if getResult.Error != nil {
			return getResult.Error
		}
		if replica.Status == "disabled" {
			return nil
		}
		replica.Status = "disabled"
		tx.Save(&replica)

		var config types.ChallengeConfig
		err := json.Unmarshal([]byte(replica.Challenge.Configuration), &config)
		if err != nil {
			return err
		}
		for _, node := range config.NodeConfig {
			ok := kube.K8sPodDestroy(replicaId, node.Name)
			if !ok {
				// TODO: this a uncorrectable error
				return errors.New("destroy pod failed - replica: " + strconv.FormatUint(replicaId, 10))
			}
		}

		return nil
	})
	if err != nil {
		log.Println(err)
		return false
	}
	return true

}

func DeleteReplicaByChallengeId(challengeId uint64, outsideTX *gorm.DB) bool {
	if outsideTX == nil {
		outsideTX = db
	}
	err := outsideTX.Transaction(func(tx *gorm.DB) error {
		var replicas []entity.Replica
		getResult := tx.Where(map[string]interface{}{"challenge_id": challengeId}).Find(&replicas)
		if getResult.Error != nil {
			return getResult.Error
		}
		for _, replica := range replicas {
			ok := DeleteReplicaAllocByReplicaId(replica.ReplicaId, tx)
			if !ok {
				return errors.New("deleteReplicaAllocByReplicaId occurred error")
			}
			if ok := DisableReplica(replica.ReplicaId, tx); !ok {
				return errors.New("disable replica failed")
			}
			tx.Delete(&replica)
		}
		return nil
	})
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}
