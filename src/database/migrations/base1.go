package migrations

import (
	"twitter/src/database"
	"twitter/src/database/models"

	"gorm.io/gorm"
)

func Up1() {
	db := database.GetDB()
	
	tables := []interface{}{}

	user := &models.User{}

	checkTable(db, user, &tables)

	err := db.Migrator().CreateTable(tables...)
	if err != nil {
		panic(err)
	}
}

func checkTable(db *gorm.DB, table interface{}, tables *[]interface{}) {
	if !db.Migrator().HasTable(table) {
		*tables = append(*tables, table)
	}
}