package resolvers

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"log"
	"server/entity"
	"strconv"
	"time"
)

func FindResetTokenByUserId(userId uint64) *entity.ResetToken {
	var resetToken entity.ResetToken
	result := db.Where(map[string]interface{}{"user_id": userId}).First(&resetToken)
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
	result := db.Where(map[string]interface{}{"token": token}).First(&resetToken)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	} else if result.Error != nil {
		log.Println(result.Error)
		return nil
	}
	return &resetToken
}

func makeToken() (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return id.String() + "-" + strconv.FormatInt(time.Now().UnixMilli(), 16), nil
}

func AddResetToken(userId uint64) *entity.ResetToken {
	var resetToken entity.ResetToken
	result := db.Where(map[string]interface{}{"user_id": userId}).First(&resetToken)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		token, err := makeToken()
		if err != nil {
			log.Println(result.Error)
			return nil
		}
		resetToken = entity.ResetToken{
			Token:  token,
			UserId: userId,
		}
		db.Create(&resetToken)
		return &resetToken
	} else if result.Error != nil {
		log.Println(result.Error)
		return nil
	}
	if time.Now().Sub(resetToken.UpdatedAt) < 5*time.Minute {
		return nil
	}
	token, err := makeToken()
	if err != nil {
		log.Println(result.Error)
		return nil
	}
	resetToken.Token = token
	db.Save(&resetToken)
	return &resetToken
}

func ResetPassword(token string, password string) bool {
	//this error may should be handled by caller
	err := db.Transaction(func(tx *gorm.DB) error {
		var resetToken entity.ResetToken
		result := tx.Where(map[string]interface{}{"token": token}).First(&resetToken)
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
