package resolvers

import (
	"errors"
	"gorm.io/gorm"
	"server/entity"
)

func FindResetTokenByUserId(userId uint64) (*entity.ResetToken, error) {
	var resetToken entity.ResetToken
	result := db.Where(map[string]interface{}{"UserId": userId}).First(&resetToken)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if result.Error != nil {
		return nil, result.Error
	}
	return &resetToken, nil
}

func FindResetTokenByToken(token string) (*entity.ResetToken, error) {
	var resetToken entity.ResetToken
	result := db.Where(map[string]interface{}{"Token": token}).First(&resetToken)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if result.Error != nil {
		return nil, result.Error
	}
	return &resetToken, nil
}
func ResetPassword(token string, password string) error {
	return db.Transaction(func(tx *gorm.DB) error {
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
}
