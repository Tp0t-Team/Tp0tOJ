package load

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"golang.org/x/text/unicode/norm"
	"gorm.io/gorm"
	"io"
	"log"
	unsafeRand "math/rand"
	"os"
	"regexp"
	"server/entity"
	"server/services/database"
	"server/utils/configure"
	"server/utils/email"
	"strconv"
	"time"

	"github.com/cheggaaa/pb/v3"
)

var blankRegexp *regexp.Regexp

func init() {
	blankRegexp, _ = regexp.Compile("\\s")
}

func lengthLimit(data string, min int, max int) bool {
	return len([]rune(data)) >= min && len([]rune(data)) <= max
}

type Item struct {
	mail string
	name string
}

func (input *Item) checkPass() bool {
	input.name = blankRegexp.ReplaceAllString(input.name, "")
	input.mail = blankRegexp.ReplaceAllString(input.mail, "")
	input.name = norm.NFC.String(input.name)
	return lengthLimit(input.name, 1, 20) && lengthLimit(input.mail, 1, 50)
}

func makePassword() (string, error) {
	seed := make([]byte, 8)
	_, err := rand.Read(seed)
	if err != nil {
		return "", err
	}
	unsafeRand.Seed(int64(binary.BigEndian.Uint64(seed)))
	init := make([]byte, 8)
	_, err = unsafeRand.Read(init)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%02x", md5.Sum(init))[:8], nil
}

func passwordHash(password string) string {
	hash1 := sha256.New()
	_, err := io.WriteString(hash1, configure.Configure.Server.Salt+password)
	if err != nil {
		log.Panicln(err.Error())
	}
	hash2 := sha256.New()
	_, err = io.WriteString(hash2, configure.Configure.Server.Salt+fmt.Sprintf("%02x", hash1.Sum(nil)))
	if err != nil {
		log.Panicln(err.Error())
	}
	return fmt.Sprintf("%02x", hash2.Sum(nil))
}

func addUser(info Item) (uint64, *email.WelcomeInfo, error) {
	password, err := makePassword()
	if err != nil {
		return 0, nil, err
	}
	user := entity.User{Name: info.name, Password: passwordHash(password), Mail: info.mail, Role: "member", State: "normal", JoinTime: time.Now()}
	err = database.DataBase.Transaction(func(tx *gorm.DB) error {
		checkResult := tx.Where(map[string]interface{}{"mail": info.mail}).First(&entity.User{})
		if errors.Is(checkResult.Error, gorm.ErrRecordNotFound) {
			result := tx.Create(&user)
			if result.Error != nil {
				return result.Error
			}
		} else if checkResult.Error != nil {
			return checkResult.Error
		} else {
			return errors.New("exists")
		}
		return nil
	})
	if err != nil {
		return 0, nil, err
	}
	return user.UserId, &email.WelcomeInfo{
		Username: user.Name,
		Mail:     user.Mail,
		Password: password,
	}, nil
}

func sendPassword(info email.WelcomeInfo, url string, reward string, halfLife string) error {
	content, err := email.RenderWelcomeEmail(info)
	if err != nil {
		return err
	}
	info.Url = url
	info.Reward = reward
	info.HalfLife = halfLife
	err = email.SendMail(info.Mail, "Welcome to Tp0t OJ", content)
	return err
}

type NullableBoolValue struct {
	Value *bool
}

func (v NullableBoolValue) String() string {
	if v.Value == nil {
		return ""
	}
	if *v.Value {
		return "true"
	} else {
		return "false"
	}
}

func (v *NullableBoolValue) Set(s string) error {
	v.Value = nil
	if s == "" {
		return nil
	}
	var value bool
	if s == "true" {
		value = true
		v.Value = &value
	}
	if s == "false" {
		value = false
		v.Value = &value
	}
	if v.Value == nil {
		return errors.New("must be `true` or `false`")
	}
	return nil
}

func Run(args []string) {
	cli := flag.NewFlagSet("load", flag.ExitOnError)

	cli.Usage = func() {
		fmt.Println("Usage: load <file>")
		fmt.Println("  <file> is a csv file with [mail,username] format and no header")
		cli.PrintDefaults()
	}
	var welcome NullableBoolValue
	cli.Var(&welcome, "welcome", "auto send welcome emails for each user.")
	err := cli.Parse(args)
	if err != nil {
		log.Panicln(err)
	}

	if welcome.Value == nil {
		log.Panicln("You must specify whether to send welcome emails.")
	}

	if len(cli.Args()) != 1 {
		log.Panicln("Can only load one file.")
	}
	file := cli.Args()[0]

	if configure.LoadConfigError != nil {
		log.Printf("load config error: %s", configure.LoadConfigError.Error())
		os.Exit(1)
	}

	if stat, err := os.Stat(file); os.IsNotExist(err) || stat.IsDir() {
		log.Panicln("No such regular file.")
	}

	database.InitDB(configure.Configure.Database.Dsn)

	csvFile, err := os.Open(file)
	if err != nil {
		log.Panicln(err)
	}
	defer csvFile.Close()

	csvData, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		fmt.Println(err)
	}

	items := []Item{}
	for _, line := range csvData {
		if len(line) != 2 {
			log.Panicln("wrong csv format")
		}
		item := Item{
			mail: line[0],
			name: line[1],
		}
		if !item.checkPass() {
			log.Panicf("Invalid user info:\n%s,%s\n", line[0], line[1])
		}
		items = append(items)
	}

	err = nil
	added := []uint64{}
	addedInfo := []*email.WelcomeInfo{}
	for _, item := range items {
		id, info, err := addUser(item)
		if err != nil {
			break
		}
		added = append(added, id)
		addedInfo = append(addedInfo, info)
	}
	if err != nil {
		database.DataBase.Delete(&entity.User{}, added)
		log.Panicln(err)
	}
	fmt.Println("Users has been loaded into the database. Start sending welcome emails...")

	url := fmt.Sprintf("%s:%s/", configure.Configure.Server.Host, strconv.Itoa(configure.Configure.Server.Port))
	reward := fmt.Sprintf(
		"%g%%, %g%%, %g%%",
		configure.Configure.Challenge.FirstBloodReward*100,
		configure.Configure.Challenge.SecondBloodReward*100,
		configure.Configure.Challenge.ThirdBloodReward*100)
	halfLife := fmt.Sprintf("%d", configure.Configure.Challenge.HalfLife)

	const barTemplate = `{{with string . "prefix"}}{{.}} {{end}}[SENDING MAIL {{percent . }}] {{counters . }}{{with string . "suffix"}} {{.}}{{end}}`
	bar := pb.ProgressBarTemplate(barTemplate).Start(len(added))

	if *welcome.Value {
		for index, _ := range added {
			bar.Increment()
			err := sendPassword(*addedInfo[index], url, reward, halfLife)
			if err != nil {
				log.Println(err)
				log.Printf("[SEND MAIL ERROR] User: %s, Mail: %s, Password: %s\n", addedInfo[index].Username, addedInfo[index].Mail, addedInfo[index].Password)
			}
		}
	}

	log.Println("Load success.")
}
