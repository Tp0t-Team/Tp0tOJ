package resolvers

import (
	"gorm.io/gorm"
	"log"
	"server/entity"
	"server/utils/configure"
	"time"
)

func BehaviorAllocReplica(challenge uint64, user uint64, flag string, actionTime time.Time, outsideTX *gorm.DB) {
	if !configure.Configure.Server.BehaviorLog {
		return
	}
	if outsideTX == nil {
		outsideTX = db
	}
	behavior := entity.Behavior{
		ActionTime:  actionTime,
		Action:      entity.ActionAllocReplica,
		ChallengeId: challenge,
		UserId:      user,
		Content:     flag,
	}
	result := db.Create(&behavior)
	if result.Error != nil {
		log.Println(result.Error)
	}
}

func BehaviorSubmit(challenge uint64, user uint64, flag string, actionTime time.Time, outsideTX *gorm.DB) {
	if !configure.Configure.Server.BehaviorLog {
		return
	}
	if outsideTX == nil {
		outsideTX = db
	}
	behavior := entity.Behavior{
		ActionTime:  actionTime,
		Action:      entity.ActionSubmit,
		ChallengeId: challenge,
		UserId:      user,
		Content:     flag,
	}
	result := outsideTX.Create(&behavior)
	if result.Error != nil {
		log.Println(result.Error)
	}
}

func BehaviorComplete(challenge uint64, user uint64, flag string, actionTime time.Time, outsideTX *gorm.DB) {
	if !configure.Configure.Server.BehaviorLog {
		return
	}
	if outsideTX == nil {
		outsideTX = db
	}
	behavior := entity.Behavior{
		ActionTime:  actionTime,
		Action:      entity.ActionComplete,
		ChallengeId: challenge,
		UserId:      user,
		Content:     flag,
	}
	result := outsideTX.Create(&behavior)
	if result.Error != nil {
		log.Println(result.Error)
	}
}

func BehaviorLogin(user uint64, ip string, actionTime time.Time, outsideTX *gorm.DB) {
	if !configure.Configure.Server.BehaviorLog {
		return
	}
	if outsideTX == nil {
		outsideTX = db
	}
	behavior := entity.Behavior{
		ActionTime:  actionTime,
		Action:      entity.ActionLogin,
		ChallengeId: 0,
		UserId:      user,
		Content:     ip,
	}
	result := outsideTX.Create(&behavior)
	if result.Error != nil {
		log.Println(result.Error)
	}
}

func BehaviorWatchDescription(challenge uint64, user uint64, actionTime time.Time, outsideTX *gorm.DB) {
	if !configure.Configure.Server.BehaviorLog {
		return
	}
	if outsideTX == nil {
		outsideTX = db
	}
	behavior := entity.Behavior{
		ActionTime:  actionTime,
		Action:      entity.ActionWatchDescription,
		ChallengeId: challenge,
		UserId:      user,
		Content:     "",
	}
	result := outsideTX.Create(&behavior)
	if result.Error != nil {
		log.Println(result.Error)
	}
}
