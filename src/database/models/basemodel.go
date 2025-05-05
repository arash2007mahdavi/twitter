package models

import (
	"database/sql"
	"time"
)

type BaseModel struct {
	Id int `gorm:"primarykey"`

	CreatedBy  int           `gorm:"not null"`
	ModifiedBy sql.NullInt64 `gorm:"null"`
	DeletedBy  sql.NullInt64 `gorm:"null"`

	CreatedAt  time.Time    `gorm:"type:TIMWSTAMP with time zone;not null"`
	ModifiedAt sql.NullTime `gorm:"type:TIMESTAMP with time zone;null"`
	DeletedAt  sql.NullTime `gorm:"type:TIMESTAMP with time zone;null"`
}