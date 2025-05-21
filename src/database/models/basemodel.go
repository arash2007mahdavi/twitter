package models

import (
	"database/sql"
	"time"
)

type BaseModel struct {
	Id int `json:"id,omitempty" gorm:"primaryKey"`

	CreatedBy  int           `json:"created_by,omitempty" gorm:"not null"`
	ModifiedBy sql.NullInt64 `json:"modified_by,omitempty" gorm:"null"`
	DeletedBy  sql.NullInt64 `json:"deleted_by,omitempty" gorm:"null"`

	CreatedAt  time.Time    `json:"created_at,omitempty" gorm:"type:TIMESTAMP with time zone;not null"`
	ModifiedAt sql.NullTime `json:"modified_at,omitempty" gorm:"type:TIMESTAMP with time zone;null"`
	DeletedAt  sql.NullTime `json:"deleted_at,omitempty" gorm:"type:TIMESTAMP with time zone;null"`
}