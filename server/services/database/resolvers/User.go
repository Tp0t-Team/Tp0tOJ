package resolvers

import (
	"errors"
	"gorm.io/gorm"
	"k8s.io/apimachinery/pkg/util/json"
	"log"
	"server/entity"
	"server/services/types"
	"time"
)

// AddUser support role[admin|member|team] state[banned|disabled|normal]
func AddUser(name string, password string, mail string, role string, state string) bool {
	err := db.Transaction(func(tx *gorm.DB) error {
		checkResult := tx.Where(map[string]interface{}{"Mail": mail}).First(&entity.User{})
		if errors.Is(checkResult.Error, gorm.ErrRecordNotFound) {
			newUser := entity.User{Name: name, Password: password, Mail: mail, Role: role, State: state, JoinTime: time.Now()}
			result := tx.Create(&newUser)
			if result.Error != nil {
				return result.Error
			}
			challenges := FindAllChallenges()
			if challenges == nil {
				return errors.New("error occurred getting challenges")
			}
			for _, challenge := range challenges {
				var config types.ChallengeConfig
				err := json.Unmarshal([]byte(challenge.Configuration), &config)
				if err != nil {
					log.Println(err)
					return errors.New("challenge Config Error")
				}
				//Alloc all replicas for the singleton and enabled challenge
				if challenge.State == "enabled" && config.Singleton {
					replicas := FindReplicaByChallengeId(challenge.ChallengeId)
					if replicas == nil || len(replicas) != 1 {
						return errors.New("found more than one or none replica for singleton challenge")
					}
					log.Println(newUser.UserId)
					ok := AddReplicaAlloc(replicas[0].ReplicaId, newUser.UserId, tx)
					if !ok {
						return errors.New("add replicaAlloc error")
					}
				}

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
	result := db.Where(map[string]interface{}{"Mail": mail}).First(&user)
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
	result := db.Where(map[string]interface{}{"UserId": userId, "Role": "admin"}).First(&users)
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
	result := db.Where(map[string]interface{}{"UserId": userId}).First(&user)
	if result.Error != nil {
		log.Println(result.Error)
		return false
	}
	user.Name = name
	user.Role = role
	user.Mail = mail
	user.State = state
	result = db.Save(&user)
	if result.Error != nil {
		log.Println(result.Error)
		return false
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
	result := db.Where(map[string]interface{}{"Mail": mail}).First(&entity.User{})
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false
	} else if result.Error != nil {
		log.Println(result.Error)
		return false
	}
	return true
}
