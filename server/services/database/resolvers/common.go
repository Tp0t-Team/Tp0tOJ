package resolvers

import (
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB(DB *gorm.DB) {
	db = DB
}
