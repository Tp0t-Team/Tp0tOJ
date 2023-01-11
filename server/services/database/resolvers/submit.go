package resolvers

import (
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
	"regexp"
	"server/entity"
	"server/services/sse"
	"server/services/types"
	"server/utils"
	"strings"
	"time"
)

func CheckSubmitCorrectByUserIdAndChallengeId(userId uint64, challengeId uint64, outsideTX *gorm.DB) bool {
	if outsideTX == nil {
		outsideTX = db
	}
	result := outsideTX.Where(map[string]interface{}{"user_id": userId, "challenge_id": challengeId, "correct": true, "available": true}).First(&entity.Submit{})
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false
	} else if result.Error != nil {
		log.Println(result.Error)
		return false
	}
	return true
}

func FindAllSubmitByChallengeId(challengeId uint64) []entity.Submit {
	var submits []entity.Submit
	result := db.Preload("User").Preload("Challenge").Where(map[string]interface{}{"challenge_id": challengeId}).Find(&submits)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return []entity.Submit{}
	} else if result.Error != nil {
		log.Println(result.Error)
		return nil
	}
	return submits
}

func FindSubmitCorrectByChallengeId(challengeId uint64) []entity.Submit {
	var submits []entity.Submit
	result := db.Preload("User").Preload("Challenge").Where(map[string]interface{}{"challenge_id": challengeId, "correct": true, "available": true}).Find(&submits)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return []entity.Submit{}
	} else if result.Error != nil {
		log.Println(result.Error)
		return nil
	}
	return submits
}

func FindNotAdminSubmitCorrectByChallengeId(challengeId uint64) []entity.Submit {
	var submits []entity.Submit
	result := db.Preload("User").Preload("Challenge").Where(map[string]interface{}{"challenge_id": challengeId, "correct": true, "available": true}).Find(&submits)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return []entity.Submit{}
	} else if result.Error != nil {
		log.Println(result.Error)
		return nil
	}
	var results []entity.Submit
	for _, submit := range submits {
		if !CheckAdminByUserId(submit.UserId) {
			results = append(results, submit)
		}
	}
	return results
}

func FindSubmitCorrectByUserId(userId uint64) []entity.Submit {
	var submits []entity.Submit
	result := db.Preload("User").Preload("Challenge").Where(map[string]interface{}{"user_id": userId, "correct": true, "available": true}).Find(&submits)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return []entity.Submit{}
	} else if result.Error != nil {
		log.Println(result.Error)
		return nil
	}
	return submits
}

func FindSubmitCorrectSorted() []entity.Submit {
	var submits []entity.Submit
	result := db.Preload("User").Preload("Challenge").Where(map[string]interface{}{"correct": true, "available": true}).Order("submit_time").Find(&submits)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return []entity.Submit{}
	} else if result.Error != nil {
		log.Println(result.Error)
		return nil
	}
	return submits
}

var BloodName = []string{"first", "second", "third"}

func AddSubmit(userId uint64, challengeId uint64, flag string, submitTime time.Time, setBlood bool) (success bool, isCorrect bool) {
	submitCache := false
	deleteReplica := new(uint64)
	// blood infos
	bloodIndex := -1
	challengeName := ""
	var correct = false
	err := db.Transaction(func(tx *gorm.DB) error {
		if CheckSubmitCorrectByUserIdAndChallengeId(userId, challengeId, tx) {
			return errors.New("already finish this challenge")
		}
		alloc, err := FindReplicaAllocByUserIdAndChallengeId(userId, challengeId, tx)
		if err != nil {
			return err
		}
		if alloc == nil {
			return errors.New("no alloc exists")
		}
		var submits []entity.Submit
		result := tx.Preload("Challenge").Where(map[string]interface{}{"challenge_id": challengeId, "correct": true, "available": true}).Find(&submits)
		if result.Error != nil {
			return result.Error
		}
		var config types.ChallengeConfig
		err = json.Unmarshal([]byte(alloc.Replica.Challenge.Configuration), &config)
		if err != nil {
			return err
		}
		var newSubmit entity.Submit
		if alloc.Replica.FlagType == types.Multiple {
			flags := strings.Split(alloc.Replica.Flag, "\n")
			for _, i := range flags {
				i = strings.TrimSpace(i)
				if flag == i {
					correct = true
					break
				}
			}
		} else if alloc.Replica.FlagType == types.Regexp {
			re, _ := regexp.Compile(alloc.Replica.Flag)
			if re != nil {
				correct = re.MatchString(flag)
			}
		} else {
			correct = flag == alloc.Replica.Flag
		}
		newSubmit = entity.Submit{
			UserId:      userId,
			ChallengeId: challengeId,
			SubmitTime:  submitTime,
			Flag:        flag,
			Correct:     correct,
			Available:   alloc.Replica.Challenge.State == "enabled",
		}
		tx.Create(&newSubmit)
		if newSubmit.Correct {
			submitCache = true
			BehaviorComplete(challengeId, userId, flag, submitTime, tx)
		}
		if correct && alloc.Replica.Challenge.State == "enabled" && setBlood {
			var challenge entity.Challenge
			result := tx.First(&challenge, []uint64{challengeId})
			if result.Error != nil {
				return err
			}
			if challenge.FirstBloodId == nil {
				challenge.FirstBloodId = &userId
				bloodIndex = 0
			} else if challenge.SecondBloodId == nil {
				challenge.SecondBloodId = &userId
				bloodIndex = 1
			} else if challenge.ThirdBloodId == nil {
				challenge.ThirdBloodId = &userId
				bloodIndex = 2
			}
			tx.Save(&challenge)
			challengeName = challenge.Name
		}

		if !alloc.Replica.Singleton {
			*deleteReplica = alloc.ReplicaId
		} else {
			deleteReplica = nil
		}

		return nil
	})
	if err != nil {
		log.Println(err)
		return false, false
	}
	if bloodIndex >= 0 {
		user, err := FindUser(userId)
		if err != nil {
			log.Println(err)
		} else {
			title := fmt.Sprintf("%s %s blood", challengeName, BloodName[bloodIndex])
			info := fmt.Sprintf("Congratulations! %s get the %s blood of %s", user.Name, BloodName[bloodIndex], challengeName)
			ok := AddBulletin(title, info, BloodName[bloodIndex], false)
			if !ok {
				log.Println("blood info bulletin create failed")
			}
			sse.PublishMessage(sse.Message{
				Title: title,
				Info:  info,
			})
		}
	}
	if submitCache {
		err := utils.Cache.Submit(userId, challengeId, submitTime)
		if err != nil {
			log.Println(err)
			return false, false
		}
		if deleteReplica != nil {
			if !DeleteReplicaById(*deleteReplica, nil) {
				log.Println("Delete replica failed.")
			}
		}
	}
	return true, correct

}

func DeleteSubmitsByChallengeId(challengeId uint64, outsideTX *gorm.DB) bool {
	if outsideTX == nil {
		outsideTX = db
	}
	err := outsideTX.Transaction(func(tx *gorm.DB) error {
		var submits []entity.Submit
		submits = FindAllSubmitByChallengeId(challengeId)
		if submits == nil {
			return errors.New("ChallengeId not find during Delete submits")
		}
		for _, submit := range submits {
			tx.Delete(&submit)
		}
		return nil
	})
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}
