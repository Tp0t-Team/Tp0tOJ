package resolvers

import (
	"encoding/json"
	"errors"
	"gorm.io/gorm"
	"log"
	"server/entity"
	"server/services/types"
)

func FindReplicaByChallengeId(challengeId uint64) []entity.Replica {
	var replicas []entity.Replica
	result := db.Where(map[string]interface{}{"ChallengeId": challengeId}).Find(&replicas)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return []entity.Replica{}
	} else if result.Error != nil {
		log.Println(result.Error)
		return nil
	}
	return replicas
}

func AddReplica(challengeId uint64) bool {
	err := db.Transaction(func(tx *gorm.DB) error {
		var challenge entity.Challenge
		challengeResult := tx.Where(map[string]interface{}{"ChallengeId": challengeId}).First(&challenge)
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
			// TODO: generate dynamic flag
		} else {
			flag = config.Flag.Value
		}
		newReplica := entity.Replica{
			ChallengeId: challengeId,
			Status:      "disabled",
			Flag:        flag,
		}
		result := tx.Create(&newReplica)
		if result.Error != nil {
			return result.Error
		}
		return nil
	})
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func EnableReplica(replicaId uint64, outsideTX *gorm.DB) bool {
	if outsideTX == nil {
		outsideTX = db
	}
	err := outsideTX.Transaction(func(tx *gorm.DB) error {
		var replica entity.Replica
		getResult := tx.Where(map[string]interface{}{"ReplicaId": replicaId}).First(&replica)
		if getResult.Error != nil {
			return getResult.Error
		}
		replica.Status = "enabled"
		db.Save(&replica)
		// TODO: K8sPodAlloc
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
		getResult := tx.Where(map[string]interface{}{"ReplicaId": replicaId}).First(&replica)
		if getResult.Error != nil {
			return getResult.Error
		}
		replica.Status = "disabled"
		// TODO: K8sPodDestroy
		db.Save(&replica)
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
		getResult := tx.Where(map[string]interface{}{"ChallengeId": challengeId}).Find(&replicas)
		if getResult.Error != nil {
			return getResult.Error
		}
		for _, replica := range replicas {
			// TODO: delete replicaAlloc
			ok := DeleteReplicaAllocByReplicaId(replica.ReplicaId, tx)
			if !ok {
				return errors.New("deleteReplicaAllocByReplicaId occurred error")
			}
			DisableReplica(replica.ReplicaId, tx)
			db.Delete(&replica)
		}
		return nil
	})
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}
