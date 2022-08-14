package resolvers

import (
	"errors"
	"gorm.io/gorm"
	"log"
	"server/entity"
	"server/utils"
	"time"
)

// AddUser support role[admin|member|team] state[banned|disabled|normal]
func AddUser(name string, password string, mail string, role string, state string) bool {
	err := db.Transaction(func(tx *gorm.DB) error {
		checkResult := tx.Where(map[string]interface{}{"mail": mail}).First(&entity.User{})
		if errors.Is(checkResult.Error, gorm.ErrRecordNotFound) {
			newUser := entity.User{Name: name, Password: password, Mail: mail, Role: role, State: state, JoinTime: time.Now()}
			result := tx.Create(&newUser)
			if result.Error != nil {
				return result.Error
			}
			if role != "admin" {
				utils.Cache.AddUser(newUser.UserId)
			}
		} else if checkResult.Error != nil {
			return checkResult.Error
		} else {
			return errors.New("exists")
		}
		return nil
	})
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func FindUserByMail(mail string) (*entity.User, error) {
	var user entity.User
	result := db.Where(map[string]interface{}{"mail": mail}).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if result.Error != nil {
		//log.Println(result.Error)
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
		log.Println(result.Error)
		return nil, result.Error
	}
	return &user, nil
}

func CheckAdminByUserId(userId uint64) bool {
	var users entity.User
	result := db.Where(map[string]interface{}{"user_id": userId, "role": "admin"}).First(&users)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false
	} else if result.Error != nil {
		log.Println(result.Error)
		return false
	}
	return true
}

func UpdateUserInfo(userId uint64, name string, role string, mail string, state string) bool {
	var user entity.User
	result := db.Where(map[string]interface{}{"user_id": userId}).First(&user)
	if result.Error != nil {
		log.Println(result.Error)
		return false
	}
	needWarmUp := user.Role != role && (user.Role == "admin" || role == "admin")
	user.Name = name
	user.Role = role
	user.Mail = mail
	user.State = state
	//TODO: if we disable user, may need re calculate the blood info
	result = db.Save(&user)
	if result.Error != nil {
		log.Println(result.Error)
		return false
	}
	if needWarmUp {
		err := utils.Cache.WarmUp()
		if err != nil {
			log.Println(err)
			return false
		}
	}
	return true
}

func FindAllUser() []entity.User {
	var users []entity.User
	result := db.Where(map[string]interface{}{}).Find(&users)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return []entity.User{}
	} else if result.Error != nil {
		log.Println(result.Error)
		return nil
	}
	return users
}

func CheckMailExistence(mail string) bool {
	result := db.Where(map[string]interface{}{"mail": mail}).First(&entity.User{})
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false
	} else if result.Error != nil {
		log.Println(result.Error)
		return false
	}
	return true
}
