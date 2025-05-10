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
	user_following := &models.UserFollowings{}
	user_tweet := &models.UserTweet{}
	user_comment := &models.UserComment{}
	tweet_comment := &models.TweetComment{}
	tweet := &models.Tweet{}
	comment := &models.Comment{}

	checkTable(db, user, &tables)
	checkTable(db, tweet, &tables)
	checkTable(db, comment, &tables)
	checkTable(db, user_followers, &tables)
	checkTable(db, user_following, &tables)
	checkTable(db, user_tweet, &tables)
	checkTable(db, user_comment, &tables)
	checkTable(db, tweet_comment, &tables)

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