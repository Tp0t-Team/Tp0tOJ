package load

import (
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"golang.org/x/text/unicode/norm"
	"gorm.io/gorm"
	"log"
	"os"
	"regexp"
	"server/entity"
	"server/services/database"
	"server/utils/configure"
	"time"
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

func addUser(info Item) (uint64, error) {
	user := entity.User{Name: info.name, Password: "", Mail: info.mail, Role: "member", State: "normal", JoinTime: time.Now()}
	err := database.DataBase.Transaction(func(tx *gorm.DB) error {
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
		return 0, err
	}
	return user.UserId, nil
}

func Run(args []string) {
	cli := flag.NewFlagSet("load", flag.ExitOnError)

	cli.Usage = func() {
		fmt.Println("Usage: load <file>")
		fmt.Println("  <file> is a csv file with [mail,username] format and no header")
	}
	err := cli.Parse(args)
	if err != nil {
		log.Panicln(err)
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
	for _, item := range items {
		var id uint64
		id, err = addUser(item)
		if err != nil {
			break
		}
		added = append(added, id)
	}
	if err != nil {
		database.DataBase.Delete(&entity.User{}, added)
		log.Panicln(err)
	}
	// TODO: maybe need send mail here?
	log.Println("Load success.")
}
