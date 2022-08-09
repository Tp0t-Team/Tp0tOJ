package resolvers

import (
	"errors"
	"gorm.io/gorm"
	"log"
	"server/entity"
	"time"
)

func EventStartGame(actionTime time.Time, outsideTX *gorm.DB) bool {
	if outsideTX == nil {
		outsideTX = db
	}
	event := entity.GameEvent{
		Time:   actionTime,
		Action: entity.ResumeEvent,
	}
	result := db.Create(&event)
	if result.Error != nil {
		log.Println(result.Error)
		return false
	}
	return true
}

func EventStopGame(actionTime time.Time, outsideTX *gorm.DB) bool {
	if outsideTX == nil {
		outsideTX = db
	}
	event := entity.GameEvent{
		Time:   actionTime,
		Action: entity.PauseEvent,
	}
	result := db.Create(&event)
	if result.Error != nil {
		log.Println(result.Error)
		return false
	}
	return true
}

func IsGameRunning(outsideTX *gorm.DB) bool {
	if outsideTX == nil {
		outsideTX = db
	}
	currentTime := time.Now()
	var currentEvent entity.GameEvent
	result := outsideTX.Where(outsideTX.Where(map[string]interface{}{"action": entity.PauseEvent}).Or(map[string]interface{}{"action": entity.ResumeEvent}))
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return true
	} else if result.Error != nil {
		log.Println(result.Error)
		return false
	}
	result = result.Where("time < ?", currentTime).Order("time desc").First(&currentEvent)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false
	} else if result.Error != nil {
		log.Println(result.Error)
		return false
	}
	if currentEvent.Action == entity.ResumeEvent {
		return true
	} else {
		return false
	}

}
