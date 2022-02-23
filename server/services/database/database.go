package database

import (
	"errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
	"server/entity"
	"server/services/types"
	"time"
)

var db *gorm.DB

func init() {
	prefix, _ := os.Getwd()
	dbPath := prefix + "/resources/data.db"
	test, err := os.Lstat(dbPath)
	if os.IsExist(err) {
		_, err := os.Create(dbPath)
		if err != nil {
			log.Panicln(err, test)
			return
		}
	} else if err != nil {
		if err != nil {
			log.Panicln(err, test)
			return
		}
	}
	db, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Panicln("DB connect error", err.Error())
	}
	err = db.AutoMigrate(&entity.Bulletin{}, &entity.Challenge{}, &entity.Replica{}, &entity.ReplicaAlloc{}, &entity.ResetToken{}, &entity.Submit{}, &entity.User{})
	if err != nil {
		log.Panicln("DB connect error", err.Error())
		return
	}
}

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

//func CheckMailExistence(mail string) (bool, error) {
//	result := db.Where(map[string]interface{}{"Mail": mail}).First(&entity.User{})
//	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
//		return false, nil
//	} else if result.Error != nil {
//		return false, result.Error
//	}
//	return true, nil
//}

func FindChallengeByState(state string) ([]entity.Challenge, error) {
	var challenges []entity.Challenge
	result := db.Where(map[string]interface{}{"State": state}).Find(&challenges)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if result.Error != nil {
		return nil, result.Error
	}
	return challenges, nil
}
func FindChallengeById(id uint64) (*entity.Challenge, error) {
	var challenge entity.Challenge
	result := db.Where(map[string]interface{}{"ChallengeId": id}).First(&challenge)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if result.Error != nil {
		return nil, result.Error
	}
	return &challenge, nil
}

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

func FindReplicaAllocByUserId(userId uint64) ([]entity.ReplicaAlloc, error) {
	var replicaAllocs []entity.ReplicaAlloc
	result := db.Where(map[string]interface{}{"UserId": userId}).Find(&replicaAllocs)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if result.Error != nil {
		return nil, result.Error
	}
	return replicaAllocs, nil
}

func FindReplicaByChallengeId(challengeId uint64) ([]entity.Replica, error) {
	var replicas []entity.Replica
	result := db.Where(map[string]interface{}{"ChallengeId": challengeId}).Find(&replicas)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if result.Error != nil {
		return nil, result.Error
	}
	return replicas, nil
}

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

func CheckSubmitCorrectByUserIdAndChallengeId(userId uint64, challengeId uint64) (bool, error) {
	result := db.Where(map[string]interface{}{"UserId": userId, "ChallengeId": challengeId, "Correct": true}).Find(&entity.Submit{})
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, nil
	} else if result.Error != nil {
		return false, result.Error
	}
	return true, nil
}

func FindSubmitCorrectByChallengeId(challengeId uint64) ([]entity.Submit, error) {
	var submits []entity.Submit
	result := db.Where(map[string]interface{}{"ChallengeId": challengeId, "Correct": true}).Find(&submits)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if result.Error != nil {
		return nil, result.Error
	}
	return submits, nil
}

func FindSubmitCorrectByUserId(userId uint64) ([]entity.Submit, error) {
	var submits []entity.Submit
	result := db.Where(map[string]interface{}{"UserId": userId, "Correct": true}).Find(&submits)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if result.Error != nil {
		return nil, result.Error
	}
	return submits, nil
}

func FindSubmitCorrectSorted() ([]entity.Submit, error) {
	var submits []entity.Submit
	result := db.Where(map[string]interface{}{"Correct": true}).Order("SubmitTime").Find(&submits)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if result.Error != nil {
		return nil, result.Error
	}
	return submits, nil
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
	var users entity.User
	result := db.Where(map[string]interface{}{"UserId": userId}).First(&users)
	if result.Error != nil {
		return result.Error
	}
	users.Name = name
	users.Role = role
	users.Mail = mail
	users.State = state
	db.Save(&users)
}

func AddChallenge(input types.ChallengeMutateInput) error {
	//TODO:
}

func FindEnabledChallenges() ([]entity.Challenge, error) {
	var challenges []entity.Challenge
	result := db.Where(map[string]interface{}{"State": "enabled"}).Find(&challenges)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if result.Error != nil {
		return nil, result.Error
	}
	return challenges, nil
}
