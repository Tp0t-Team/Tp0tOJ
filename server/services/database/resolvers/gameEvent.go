package resolvers

import (
	"errors"
	"gorm.io/gorm"
	"log"
	"server/entity"
	"time"
)

func AddEvent(eventAction int, eventTime time.Time) bool {
	event := entity.GameEvent{
		Time:   eventTime,
		Action: uint64(eventAction),
	}
	result := db.Create(&event)
	if result.Error != nil {
		log.Println(result.Error)
		return false
	}
	return true
}

func UpdateEvent(eventId uint64, eventTime time.Time) bool {

	var event entity.GameEvent

	result := db.Where(map[string]interface{}{"event_id": eventId}).First(&event)
	if result.Error != nil {
		log.Println(errors.New("Update Event error:\n" + result.Error.Error()))
		return false
	}
	event.Time = eventTime
	result = db.Save(&event)
	if result.Error != nil {
		log.Println(result.Error)
		return false
	}
	return true
}
func DeleteEvent(eventId uint64) bool {

	var event entity.GameEvent

	result := db.Where(map[string]interface{}{"event_id": eventId}).First(&event)
	if result.Error != nil {
		log.Println(errors.New("Update Event error:\n" + result.Error.Error()))
		return false
	}
	result = db.Delete(&event)
	if result.Error != nil {
		log.Println(result.Error)
		return false
	}
	return true
}

func GetAllEvents() []entity.GameEvent {

	var events []entity.GameEvent

	result := db.Where(map[string]interface{}{}).Find(&events)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return []entity.GameEvent{}
	} else if result.Error != nil {
		log.Println(result.Error)
		return nil
	}
	return events
}

func IsGameRunning(outsideTX *gorm.DB) bool {
	if outsideTX == nil {
		outsideTX = db
	}
	currentTime := time.Now()
	var currentEvent entity.GameEvent
	var events []entity.GameEvent
	result := outsideTX.Where(map[string]interface{}{"action": entity.PauseEvent}).Or(map[string]interface{}{"action": entity.ResumeEvent}).Find(&events)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return true
	} else if result.Error != nil {
		log.Println(result.Error)
		return false
	}
	if len(events) == 0 {
		return true
	}
	result = outsideTX.Where(
		outsideTX.Where(map[string]interface{}{"action": entity.PauseEvent}).Or(map[string]interface{}{"action": entity.ResumeEvent}),
	).Where("time <= ?", currentTime).Order("time desc").First(&currentEvent)
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
