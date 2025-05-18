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
	user_followers := &models.UserFollowers{}
	tweet := &models.Tweet{}
	tweet_likes := &models.TweetLikes{}
	comment := &models.Comment{}
	comment_likes := &models.CommentLikes{}

	checkTable(db, user, &tables)
	checkTable(db, tweet, &tables)
	checkTable(db, comment, &tables)
	checkTable(db, comment_likes, &tables)
	checkTable(db, tweet_likes, &tables)
	checkTable(db, user_followers, &tables)

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