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

// func checkTable(db *gorm.DB, table interface{}, tables *[]interface{}) {
// 	if !db.Migrator().HasTable(table) {
// 		*tables = append(*tables, table)
// 	}
// }