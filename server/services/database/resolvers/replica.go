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
	"server/utils/configure"
	"strconv"
	"strings"
)

func FindReplicaByChallengeId(challengeId uint64) []entity.Replica {
	var replicas []entity.Replica
	result := db.Preload("Challenge").Where(map[string]interface{}{"challenge_id": challengeId}).Find(&replicas)
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
			flag = fmt.Sprintf("%02x", md5.Sum(init))
		} else {
			//strings.Split(config.Flag.Value,"\n")
			//TODO: for muti-flag typed challenge, may need some extra method to map flags to replica
			flag = config.Flag.Value
		}
		newReplica = entity.Replica{
			ChallengeId: challengeId,
			Status:      "disabled",
			Flag:        flag,
			Singleton:   config.Singleton,
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
		getResult := tx.Preload("Challenge").Where(map[string]interface{}{"replica_id": replicaId}).First(&replica)
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
			log.Println(replica)
			return err
		}
		createPodSuccess := true
		for _, node := range config.NodeConfig {
			//ports := []corev1.ContainerPort{}
			//for _, port := range node.Ports {
			//	ports = append(ports, *kube.NewContainerPortConfig(kube.ParseProtocol(port.Protocol), port.Port))
			//}
			servicePorts := []corev1.ServicePort{}
			for _, port := range node.ServicePorts {
				servicePorts = append(servicePorts, *kube.NewServicePortConfig(port.Name, kube.ParseProtocol(port.Protocol), port.External, port.Internal, port.Pod))
			}
			//ok := kube.K8sPodAlloc(replicaId, node.Name, strings.ToLower(configure.Configure.Kubernetes.RegistryHost+"/"+node.Image), ports, servicePorts, replica.Flag)
			ok := kube.K8sPodAlloc(replicaId, node.Name, strings.ToLower(configure.Configure.Kubernetes.RegistryHost+"/"+node.Image), servicePorts, replica.Flag)

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
		getResult := tx.Preload("Challenge").Where(map[string]interface{}{"replica_id": replicaId}).First(&replica)
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
		getResult := tx.Preload("Challenge").Where(map[string]interface{}{"challenge_id": challengeId}).Find(&replicas)
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

//
//type replicaTimer struct {
//	clock map[uint64]*time.Timer
//	lock  sync.Mutex
//}

//func (t *replicaTimer) NewTimer(userId uint64) bool {
//	t.lock.Lock()
//	defer t.lock.Unlock()
//	if _, ok := t.clock[userId]; ok {
//		return false
//	}
//	t.clock[userId] = time.NewTimer(30 * time.Minute)
//	go func() {
//		select {
//		case <-t.clock[userId].C:
//			t.clock[userId].Stop()
//			err := db.Transaction(func(tx *gorm.DB) error {
//				var allocs []entity.ReplicaAlloc
//				tx.Preload("Replica").Where(map[string]interface{}{"user_id": userId}).First(&allocs)
//				for _, alloc := range allocs {
//					if !alloc.Replica.Singleton {
//						replicaId := alloc.ReplicaId
//						tx.Delete(&alloc)
//						if !DisableReplica(alloc.ReplicaId, tx) {
//							return errors.New("delete replica failed")
//						}
//						tx.Delete(&entity.Replica{ReplicaId: replicaId})
//					}
//				}
//				return nil
//			})
//			if err != nil {
//				log.Println(err)
//			}
//			delete(t.clock, userId)
//		case <-time.After(35 * time.Minute):
//			log.Println("some thing error at replicaTimer, timeout")
//			delete(t.clock, userId)
//		}
//	}()
//	return true
//}

//func (t *replicaTimer) DeleteTimer(userId uint64) bool {
//	t.lock.Lock()
//	defer t.lock.Unlock()
//	// need timer in the queue
//	if _, ok := t.clock[userId]; !ok {
//		return false
//	}
//	t.clock[userId].Stop()
//	delete(t.clock, userId)
//	return true
//}

//var ReplicaTimer = replicaTimer{clock: map[uint64]*time.Timer{}}

//func DeleteReplicaByUserId(userId uint64) bool {
//	err := db.Transaction(func(tx *gorm.DB) error {
//		var allocs []entity.ReplicaAlloc
//		tx.Preload("Replica").Where(map[string]interface{}{"user_id": userId}).First(&allocs)
//		for _, alloc := range allocs {
//			if !alloc.Replica.Singleton {
//				replicaId := alloc.ReplicaId
//				tx.Delete(&alloc)
//				if !DisableReplica(alloc.ReplicaId, tx) {
//					return errors.New("delete replica failed")
//				}
//				tx.Delete(&entity.Replica{ReplicaId: replicaId})
//			}
//		}
//		return nil
//	})
//	if err != nil {
//		log.Println(err)
//		return false
//	}
//	return true
//}

func StartReplicaForUser(userId uint64, challengeId uint64) bool {
	//ok := ReplicaTimer.NewTimer(userId)
	//if !ok {
	//	log.Println("new timer failed, cannot start replica")
	//	return false
	//}
	err := db.Transaction(func(tx *gorm.DB) error {
		challenge, err := FindChallengeById(challengeId)
		if err != nil {
			return err
		}
		var config types.ChallengeConfig
		err = json.Unmarshal([]byte(challenge.Configuration), &config)
		if err != nil {
			return err
		}
		if config.Singleton {
			return errors.New("cannot start replica for a singleton")
		}
		var allocs []entity.ReplicaAlloc
		tx.Preload("Replica").Where(map[string]interface{}{"user_id": userId}).First(&allocs)
		for _, alloc := range allocs {
			if !alloc.Replica.Singleton {
				replicaId := alloc.ReplicaId
				tx.Delete(&alloc)
				if !DisableReplica(alloc.ReplicaId, tx) {
					return errors.New("delete replica failed")
				}
				tx.Delete(&entity.Replica{ReplicaId: replicaId})
			}
		}
		newReplica := AddReplica(challengeId, tx)
		if newReplica == nil {
			return errors.New("create replica failed")
		}
		if !EnableReplica(newReplica.ReplicaId, tx) {
			return errors.New("enable replica failed")
		}
		if !AddReplicaAlloc(newReplica.ReplicaId, userId, tx) {
			return errors.New("alloc replica failed")
		}
		//ok := ReplicaTimer.NewTimer(userId)
		//if !ok {
		//	return errors.New("new timer failed, cannot start replica")
		//}
		return nil
	})
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}
