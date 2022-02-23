package resolvers

import (
	"errors"
	"gorm.io/gorm"
	"server/entity"
)

func GetAllBulletin() ([]entity.Bulletin, error) {
	var allBulletin []entity.Bulletin
	result := db.Find(&allBulletin)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if result.Error != nil {
		return nil, result.Error
	}
	return allBulletin, nil
}
func AddBulletin(title string, content string, topping bool) error {
	newBulletin := entity.Bulletin{Title: title, Content: content, Topping: topping}
	result := db.Create(&newBulletin)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func FindBulletinByTitle(title string) (*entity.Bulletin, error) {
	var bulletin entity.Bulletin
	result := db.Where(map[string]interface{}{"Title": title}).First(&bulletin)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if result.Error != nil {
		return nil, result.Error
	}
	return &bulletin, nil
}
