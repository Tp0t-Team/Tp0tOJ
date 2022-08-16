package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"server/entity"
	"server/services/database"
	"server/utils/configure"
	"strconv"
	"time"
)

var timeZone *time.Location

func init() {
	var err error
	timeZone, err = time.LoadLocation("UTC")
	if err != nil {
		log.Panicln(err)
	}
}

type BulletinJSON struct {
	BulletinId  string `json:"bulletinId"`
	Content     string `json:"content"`
	Topping     bool   `json:"topping"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
	Title       string `json:"title"`
	PublishTime string `json:"publishTime"`
}

func newBulletinJSON(input entity.Bulletin) BulletinJSON {
	return BulletinJSON{
		BulletinId:  strconv.FormatUint(input.BulletinId, 10),
		Content:     input.Content,
		Topping:     input.Topping,
		CreatedAt:   input.CreatedAt.In(timeZone).Format(time.RFC3339),
		UpdatedAt:   input.UpdatedAt.In(timeZone).Format(time.RFC3339),
		Title:       "",
		PublishTime: input.PublishTime.In(timeZone).Format(time.RFC3339),
	}
}

type ChallengeJSON struct {
	ChallengeId   string `json:"challengeId"`
	Name          string `json:"name"`
	Configuration string `json:"configuration"`
	FirstBloodId  string `json:"firstBloodId"`
	SecondBloodId string `json:"secondBloodId"`
	ThirdBloodId  string `json:"thirdBloodId"`
	State         string `json:"state"`
	CreatedAt     string `json:"createdAt"`
	UpdatedAt     string `json:"updatedAt"`
}

func optionalUint(v *uint64) string {
	if v == nil {
		return ""
	} else {
		return strconv.FormatUint(*v, 10)
	}
}

func newChallengeJSON(input entity.Challenge) ChallengeJSON {
	return ChallengeJSON{
		ChallengeId:   strconv.FormatUint(input.ChallengeId, 10),
		Name:          input.Name,
		Configuration: input.Configuration,
		FirstBloodId:  optionalUint(input.FirstBloodId),
		SecondBloodId: optionalUint(input.SecondBloodId),
		ThirdBloodId:  optionalUint(input.ThirdBloodId),
		State:         input.State,
		CreatedAt:     input.CreatedAt.In(timeZone).Format(time.RFC3339),
		UpdatedAt:     input.UpdatedAt.In(timeZone).Format(time.RFC3339),
	}
}

type ReplicaJSON struct {
	ReplicaId   string `json:"replicaId"`
	ChallengeId string `json:"challengeId"`
	Singleton   bool
	Status      string `json:"status"`
	Flag        string `json:"flag"`
	FlagType    string
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

const (
	Single = iota
	Multiple
	Regexp
	Dynamic
)

func newReplicaJSON(input entity.Replica) ReplicaJSON {
	var flagType string
	switch input.FlagType {
	case Single:
		flagType = "Single"
		break
	case Multiple:
		flagType = "Multiple"
		break
	case Regexp:
		flagType = "Regexp"
		break
	case Dynamic:
		flagType = "Dynamic"
		break
	}

	return ReplicaJSON{
		ReplicaId:   strconv.FormatUint(input.ReplicaId, 10),
		ChallengeId: strconv.FormatUint(input.ChallengeId, 10),
		Singleton:   input.Singleton,
		Status:      input.Status,
		Flag:        input.Flag,
		FlagType:    flagType,
		CreatedAt:   input.CreatedAt.In(timeZone).Format(time.RFC3339),
		UpdatedAt:   input.UpdatedAt.In(timeZone).Format(time.RFC3339),
	}
}

type ReplicaAllocJSON struct {
	ReplicaAllocId string `json:"replicaAllocId"`
	UserId         string `json:"userId"`
	ReplicaId      string `json:"replicaId"`
}

func newReplicaAllocJSON(input entity.ReplicaAlloc) ReplicaAllocJSON {
	return ReplicaAllocJSON{
		ReplicaAllocId: strconv.FormatUint(input.ReplicaAllocId, 10),
		UserId:         strconv.FormatUint(input.UserId, 10),
		ReplicaId:      strconv.FormatUint(input.ReplicaId, 10),
	}
}

type ResetTokenJSON struct {
	TokenId   string `json:"tokenId"`
	Token     string `json:"token"`
	UserId    string `json:"userId"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

func newResetTokenJSON(input entity.ResetToken) ResetTokenJSON {
	return ResetTokenJSON{
		TokenId:   strconv.FormatUint(input.TokenId, 10),
		Token:     input.Token,
		UserId:    strconv.FormatUint(input.UserId, 10),
		CreatedAt: input.CreatedAt.In(timeZone).Format(time.RFC3339),
		UpdatedAt: input.UpdatedAt.In(timeZone).Format(time.RFC3339),
	}
}

type SubmitJSON struct {
	SubmitId    string `json:"submitId"`
	UserId      string `json:"userId"`
	SubmitTime  string `json:"submitTime"`
	Flag        string `json:"flag"`
	Correct     bool   `json:"correct"`
	Available   bool   `json:"available"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
	ChallengeId string `json:"challengeId"`
}

func newSubmitJSON(input entity.Submit) SubmitJSON {
	return SubmitJSON{
		SubmitId:    strconv.FormatUint(input.SubmitId, 10),
		UserId:      strconv.FormatUint(input.UserId, 10),
		SubmitTime:  input.SubmitTime.In(timeZone).Format(time.RFC3339),
		Flag:        input.Flag,
		Correct:     false,
		Available:   false,
		CreatedAt:   input.CreatedAt.In(timeZone).Format(time.RFC3339),
		UpdatedAt:   input.UpdatedAt.In(timeZone).Format(time.RFC3339),
		ChallengeId: strconv.FormatUint(input.ChallengeId, 10),
	}
}

type UserJSON struct {
	UserId    string `json:"userId"`
	Name      string `json:"name"`
	Password  string `json:"password"`
	State     string `json:"state"`
	Mail      string `json:"mail"`
	JoinTime  string `json:"joinTime"`
	Role      string `json:"role"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

func newUserJSON(input entity.User) UserJSON {
	return UserJSON{
		UserId:    strconv.FormatUint(input.UserId, 10),
		Name:      input.Name,
		Password:  input.Password,
		State:     input.State,
		Mail:      input.Mail,
		JoinTime:  input.JoinTime.In(timeZone).Format(time.RFC3339),
		Role:      input.Role,
		CreatedAt: input.CreatedAt.In(timeZone).Format(time.RFC3339),
		UpdatedAt: input.UpdatedAt.In(timeZone).Format(time.RFC3339),
	}
}

type BehaviorJSON struct {
	Id          string `json:"id"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
	ActionTime  string `json:"actionTime"`
	Action      string `json:"action"`
	ChallengeId string `json:"challengeId"`
	UserId      string `json:"userId"`
	Content     string `json:"content"`
}

const (
	ActionAllocReplica = iota
	ActionSubmit
	ActionComplete
	ActionLogin
	ActionWatchDescription
	ActionUploadWP
)

func newBehaviorJSON(input entity.Behavior) BehaviorJSON {
	var action string
	switch input.Action {
	case ActionAllocReplica:
		action = "AllocReplica"
		break
	case ActionSubmit:
		action = "Submit"
		break
	case ActionComplete:
		action = "Complete"
		break
	case ActionLogin:
		action = "Login"
		break
	case ActionWatchDescription:
		action = "WatchDescription"
		break
	case ActionUploadWP:
		action = "WatchDescription"
		break
	}

	return BehaviorJSON{
		Id:          strconv.FormatUint(input.Id, 10),
		CreatedAt:   input.CreatedAt.In(timeZone).Format(time.RFC3339),
		UpdatedAt:   input.UpdatedAt.In(timeZone).Format(time.RFC3339),
		ActionTime:  input.ActionTime.In(timeZone).Format(time.RFC3339),
		Action:      action,
		ChallengeId: strconv.FormatUint(input.ChallengeId, 10),
		UserId:      strconv.FormatUint(input.UserId, 10),
		Content:     input.Content,
	}
}

type GameEventJSON struct {
	EventId   string `json:"eventId"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	Time      string `json:"time"`
	Action    string `json:"action"`
}

const (
	PauseEvent = iota
	ResumeEvent
)

func newGameEventJSON(input entity.GameEvent) GameEventJSON {

	var action string
	switch input.Action {
	case PauseEvent:
		action = "PauseGame"
		break
	case ResumeEvent:
		action = "ResumeGame"
		break
	}
	return GameEventJSON{
		EventId:   strconv.FormatUint(input.EventId, 10),
		CreatedAt: input.CreatedAt.In(timeZone).Format(time.RFC3339),
		UpdatedAt: input.UpdatedAt.In(timeZone).Format(time.RFC3339),
		Time:      input.Time.In(timeZone).Format(time.RFC3339),
		Action:    action,
	}
}

func main() {
	dir := flag.String("dir", "./data", "export directory")
	flag.Parse()

	if _, err := os.Stat(*dir); os.IsNotExist(err) {
		err := os.MkdirAll(*dir, 0777)
		if err != nil {
			log.Panicln(err)
		}
	}

	database.InitDB(configure.Configure.Database.Dsn)

	var bulletins []entity.Bulletin
	result := database.DataBase.Find(&bulletins)
	if result.Error != nil {
		log.Println(result.Error)
	} else {
		var formated = []BulletinJSON{}
		for _, item := range bulletins {
			formated = append(formated, newBulletinJSON(item))
		}
		data, err := json.Marshal(formated)
		if err != nil {
			log.Println(err)
		}
		err = os.WriteFile(*dir+"/bulletins.json", data, 0666)
		if err != nil {
			log.Println(err)
		}
	}

	var challenges []entity.Challenge
	result = database.DataBase.Find(&challenges)
	if result.Error != nil {
		log.Println(result.Error)
	} else {
		var formated = []ChallengeJSON{}
		for _, item := range challenges {
			formated = append(formated, newChallengeJSON(item))
		}
		data, err := json.Marshal(formated)
		if err != nil {
			log.Println(err)
		}
		err = os.WriteFile(*dir+"/challenges.json", data, 0666)
		if err != nil {
			log.Println(err)
		}
	}

	var replicas []entity.Replica
	result = database.DataBase.Find(&replicas)
	if result.Error != nil {
		log.Println(result.Error)
	} else {
		var formated = []ReplicaJSON{}
		for _, item := range replicas {
			formated = append(formated, newReplicaJSON(item))
		}
		data, err := json.Marshal(formated)
		if err != nil {
			log.Println(err)
		}
		err = os.WriteFile(*dir+"/replicas.json", data, 0666)
		if err != nil {
			log.Println(err)
		}
	}

	var allocs []entity.ReplicaAlloc
	result = database.DataBase.Find(&allocs)
	if result.Error != nil {
		log.Println(result.Error)
	} else {
		var formated = []ReplicaAllocJSON{}
		for _, item := range allocs {
			formated = append(formated, newReplicaAllocJSON(item))
		}
		data, err := json.Marshal(formated)
		if err != nil {
			log.Println(err)
		}
		err = os.WriteFile(*dir+"/allocs.json", data, 0666)
		if err != nil {
			log.Println(err)
		}
	}

	var tokens []entity.ResetToken
	result = database.DataBase.Find(&tokens)
	if result.Error != nil {
		log.Println(result.Error)
	} else {
		var formated = []ResetTokenJSON{}
		for _, item := range tokens {
			formated = append(formated, newResetTokenJSON(item))
		}
		data, err := json.Marshal(formated)
		if err != nil {
			log.Println(err)
		}
		err = os.WriteFile(*dir+"/tokens.json", data, 0666)
		if err != nil {
			log.Println(err)
		}
	}

	var submits []entity.Submit
	result = database.DataBase.Find(&submits)
	if result.Error != nil {
		log.Println(result.Error)
	} else {
		var formated = []SubmitJSON{}
		for _, item := range submits {
			formated = append(formated, newSubmitJSON(item))
		}
		data, err := json.Marshal(formated)
		if err != nil {
			log.Println(err)
		}
		err = os.WriteFile(*dir+"/submits.json", data, 0666)
		if err != nil {
			log.Println(err)
		}
	}

	var users []entity.User
	result = database.DataBase.Find(&users)
	if result.Error != nil {
		log.Println(result.Error)
	} else {
		var formated = []UserJSON{}
		for _, item := range users {
			formated = append(formated, newUserJSON(item))
		}
		data, err := json.Marshal(formated)
		if err != nil {
			log.Println(err)
		}
		err = os.WriteFile(*dir+"/users.json", data, 0666)
		if err != nil {
			log.Println(err)
		}
	}

	var behaviors []entity.Behavior
	result = database.DataBase.Find(&behaviors)
	if result.Error != nil {
		log.Println(result.Error)
	} else {
		var formated = []BehaviorJSON{}
		for _, item := range behaviors {
			formated = append(formated, newBehaviorJSON(item))
		}
		data, err := json.Marshal(formated)
		if err != nil {
			log.Println(err)
		}
		err = os.WriteFile(*dir+"/behaviors.json", data, 0666)
		if err != nil {
			log.Println(err)
		}
	}

	var events []entity.GameEvent
	result = database.DataBase.Find(&events)
	if result.Error != nil {
		log.Println(result.Error)
	} else {
		var formated = []GameEventJSON{}
		for _, item := range events {
			formated = append(formated, newGameEventJSON(item))
		}
		data, err := json.Marshal(formated)
		if err != nil {
			log.Println(err)
		}
		err = os.WriteFile(*dir+"/events.json", data, 0666)
		if err != nil {
			log.Println(err)
		}
	}

	log.Println("finish.")
}
