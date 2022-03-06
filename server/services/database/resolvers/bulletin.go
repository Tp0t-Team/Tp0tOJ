package resolvers

import (
	"errors"
	"gorm.io/gorm"
	"log"
	"server/entity"
)

func GetAllBulletin() []entity.Bulletin {
	var allBulletin []entity.Bulletin
	result := db.Find(&allBulletin)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return []entity.Bulletin{}
	} else if result.Error != nil {
		log.Println(result.Error)
		return nil
	}
	return allBulletin
}

func AddBulletin(title string, content string, topping bool) bool {
	newBulletin := entity.Bulletin{Title: title, Content: content, Topping: topping}
	result := db.Create(&newBulletin)
	if result.Error != nil {
		log.Println(result.Error)
		return false
	}
	return true
}

func FindBulletinByTitle(title string) (*entity.Bulletin, error) {
	var bulletin entity.Bulletin
	result := db.Where(map[string]interface{}{"Title": title}).First(&bulletin)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if result.Error != nil {
		//log.Println(result.Error)
		return nil, result.Error
	}
	return &bulletin, nil
}
