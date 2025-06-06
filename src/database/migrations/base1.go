package migrations

import (
	"twitter/src/database"
	"twitter/src/database/models"
)

func Up1() {
	db := database.GetDB()

	user := &models.User{}
	tweet := &models.Tweet{}
	comment := &models.Comment{}
	file := &models.File{}

	db.AutoMigrate(&user, &tweet, &comment, &file)
}