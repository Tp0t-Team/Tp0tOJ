package resolvers

import (
	"errors"
	"gorm.io/gorm"
	"log"
	"server/entity"
)

func FindResetTokenByUserId(userId uint64) *entity.ResetToken {
	var resetToken entity.ResetToken
	result := db.Where(map[string]interface{}{"UserId": userId}).First(&resetToken)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	} else if result.Error != nil {
		log.Println(result.Error)
		return nil
	}
	return &resetToken
}

func FindResetTokenByToken(token string) *entity.ResetToken {
	var resetToken entity.ResetToken
	result := db.Where(map[string]interface{}{"Token": token}).First(&resetToken)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	} else if result.Error != nil {
		log.Println(result.Error)
		return nil
	}
	return &resetToken
}

func ResetPassword(token string, password string) bool {
	//this error may should be handled by caller
	err := db.Transaction(func(tx *gorm.DB) error {
		var resetToken entity.ResetToken
		result := tx.Where(map[string]interface{}{"Token": token}).First(&resetToken)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("invalid")
		} else if result.Error != nil {
			return result.Error
		}
		resetToken.User.Password = password
		tx.Save(&resetToken.User)
		tx.Delete(&resetToken)
		return nil
	})
	if err != nil {
		log.Println(err)
		return false
	}
	return true

}
