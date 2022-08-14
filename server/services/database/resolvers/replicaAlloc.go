package resolvers

import (
	"encoding/json"
	"errors"
	"gorm.io/gorm"
	"log"
	"server/entity"
	"server/services/types"
	"time"
)

func FindReplicaAllocByUserId(userId uint64) []entity.ReplicaAlloc {
	var replicaAllocs []entity.ReplicaAlloc
	result := db.Preload("User").Preload("Replica").Preload("Replica.Challenge").Where(map[string]interface{}{"user_id": userId}).Find(&replicaAllocs)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return []entity.ReplicaAlloc{}
	} else if result.Error != nil {
		log.Println(result.Error)
		return nil
	}
	return replicaAllocs
}

func FindReplicaAllocByUserIdAndChallengeId(userId uint64, challengeId uint64, outsideTX *gorm.DB) (*entity.ReplicaAlloc, error) {
	if outsideTX == nil {
		outsideTX = db
	}
	var found []entity.ReplicaAlloc
	err := db.Transaction(func(tx *gorm.DB) error {
		var allocs []entity.ReplicaAlloc
		allocResult := tx.Preload("User").Preload("Replica").Preload("Replica.Challenge").Where(map[string]interface{}{"user_id": userId}).Find(&allocs)
		if errors.Is(allocResult.Error, gorm.ErrRecordNotFound) {
			return nil
		} else if allocResult.Error != nil {
			return allocResult.Error
		}
		for _, alloc := range allocs {
			if alloc.Replica.ChallengeId == challengeId {
				found = append(found, alloc)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if found == nil {
		return nil, nil
	}
	if len(found) != 1 {
		return nil, errors.New("repeat replica alloc")
	}
	return &found[0], nil
}
func FindReplicaAllocByUserIdAndReplicaId(userId uint64, replicaId uint64, outsideTX *gorm.DB) (*entity.ReplicaAlloc, error) {
	if outsideTX == nil {
		outsideTX = db
	}
	var found []entity.ReplicaAlloc
	err := db.Transaction(func(tx *gorm.DB) error {
		var allocs []entity.ReplicaAlloc
		allocResult := tx.Preload("User").Preload("Replica").Preload("Replica.Challenge").Where(map[string]interface{}{"user_id": userId}).Find(&allocs)
		if errors.Is(allocResult.Error, gorm.ErrRecordNotFound) {
			return nil
		} else if allocResult.Error != nil {
			return allocResult.Error
		}
		for _, alloc := range allocs {
			if alloc.Replica.ReplicaId == replicaId {
				found = append(found, alloc)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if found == nil {
		return nil, nil
	}
	if len(found) != 1 {
		return nil, errors.New("repeat replica alloc")
	}
	return &found[0], nil
}
func AddReplicaAlloc(replicaId uint64, userId uint64, outsideTX *gorm.DB) bool {
	if outsideTX == nil {
		outsideTX = db
	}
	found, err := FindReplicaAllocByUserIdAndReplicaId(userId, replicaId, outsideTX)
	if err != nil {
		log.Println(err)
		return false
	}
	if found != nil {
		log.Println("recreate a repeat replicaAlloc error")
		return false
	}
	replicaAlloc := entity.ReplicaAlloc{
		ReplicaId: replicaId,
		UserId:    userId,
	}
	result := outsideTX.Create(&replicaAlloc)
	if result.Error != nil {
		log.Println(result.Error)
		return false
	}
	foundedReplica, err := FindReplicaById(replicaId, outsideTX)
	if err != nil {
		log.Println(err)
	} else {
		BehaviorAllocReplica(foundedReplica.ChallengeId, userId, foundedReplica.Flag, time.Now(), outsideTX)
	}
	return true

}

func DeleteReplicaAllocByReplicaId(replicaId uint64, outsideTX *gorm.DB) bool {
	if outsideTX == nil {
		outsideTX = db
	}
	err := outsideTX.Transaction(func(tx *gorm.DB) error {
		var replicaAllocs []entity.ReplicaAlloc
		getResult := tx.Preload("User").Preload("Replica").Preload("Replica.Challenge").Where(map[string]interface{}{"replica_id": replicaId}).Find(&replicaAllocs)
		if getResult.Error != nil {
			return getResult.Error
		}
		for _, replicaAlloc := range replicaAllocs {
			tx.Delete(&replicaAlloc)
		}
		return nil
	})
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func AllocSingleton(challengeId uint64, userId uint64) error {
	reason := ""
	err := db.Transaction(func(tx *gorm.DB) error {
		var user *entity.User
		user, err := FindUserInTX(userId, tx)
		if err != nil {
			log.Println(err)
			reason = "Get User Info Error!"
			return errors.New("")
		}
		if user == nil {
			reason = "No such user."
			return errors.New("")
		}
		var challenge *entity.Challenge
		challenge, err = FindChallengeByIdInTX(challengeId, tx)
		if err != nil {
			log.Println(err)
			reason = "no such challenge"
			return errors.New("")
		}
		//Alloc replicas for the watching singleton and enabled challenge
		var config types.ChallengeConfig
		err = json.Unmarshal([]byte(challenge.Configuration), &config)
		if err != nil {
			log.Println(err)
			reason = "alloc static replica error"
			return errors.New("")
		}
		if challenge.State == "enabled" && config.Singleton {
			replicas := FindReplicaByChallengeIdInTX(challenge.ChallengeId, tx)
			log.Println("enable replica ", challenge.Name)
			if replicas == nil || len(replicas) != 1 {
				log.Println("found more than one or none replica for singleton challenge")
				reason = "alloc static replica error"
				return errors.New("")
			}
			replicaAlloc, err := FindReplicaAllocByUserIdAndChallengeId(user.UserId, challenge.ChallengeId, tx)
			if err != nil {
				log.Println(err)
				reason = "alloc static replica error"
				return errors.New("")
			}
			if replicaAlloc != nil {
				return nil
			}
			ok := AddReplicaAlloc(replicas[0].ReplicaId, user.UserId, nil)
			if !ok {
				log.Println("add replicaAlloc error")
				reason = "alloc static replica error"
				return errors.New("")
			}
		}
		return nil
	})
	if err != nil {
		if reason != "" {
			return errors.New(reason)
		} else {
			log.Println(err)
			return errors.New("alloc static replica failed")
		}
	}
	return nil
}
