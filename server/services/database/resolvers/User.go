package resolvers

import (
	"errors"
	"gorm.io/gorm"
	"server/entity"
	"time"
)

// AddUser support role[admin|member|team] state[banned|disabled|normal]
func AddUser(name string, password string, mail string, role string, state string) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		checkResult := tx.Where(map[string]interface{}{"Mail": mail}).First(&entity.User{})
		if errors.Is(checkResult.Error, gorm.ErrRecordNotFound) {
			newUser := entity.User{Name: name, Password: password, Mail: mail, Role: role, State: state, JoinTime: time.Now(), Score: 0}
			result := tx.Create(&newUser)
			if result.Error != nil {
				return result.Error
			}
			return nil
		} else if checkResult.Error != nil {
			return checkResult.Error
		} else {
			return errors.New("exists")
		}
	})
	if err != nil {
		return err
	}
	return nil
}

func FindUserByMail(mail string) (*entity.User, error) {
	var user entity.User
	result := db.Where(map[string]interface{}{"Mail": mail}).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func FindUser(id uint64) (*entity.User, error) {
	var user entity.User
	result := db.Find(&user, []uint64{id})
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func CheckAdminByUserId(userId uint64) (bool, error) {
	var users entity.User
	result := db.Where(map[string]interface{}{"UserId": userId, "Role": "admin"}).First(&users)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, nil
	} else if result.Error != nil {
		return false, result.Error
	}
	return true, nil
}

func UpdateUserInfo(userId uint64, name string, role string, mail string, state string) error {
	var user entity.User
	result := db.Where(map[string]interface{}{"UserId": userId}).First(&user)
	if result.Error != nil {
		return result.Error
	}
	user.Name = name
	user.Role = role
	user.Mail = mail
	user.State = state
	db.Save(&user)
	return nil
}

func FindAllUser() ([]entity.User, error) {
	var users []entity.User
	result := db.Where(map[string]interface{}{}).Find(&users)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func CheckMailExistence(mail string) (bool, error) {
	result := db.Where(map[string]interface{}{"Mail": mail}).First(&entity.User{})
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, nil
	} else if result.Error != nil {
		return false, result.Error
	}
	return true, nil
}
