package models

import (
	"database/sql"
	"time"
)

type BaseModel struct {
	Id int `json:"id,omitempty" gorm:"primaryKey"`

	CreatedBy  int           `json:"created_by,omitempty" gorm:"not null"`
	ModifiedBy sql.NullInt64 `json:"modified_by,omitempty" gorm:"null" swaggertype:"string"`
	DeletedBy  sql.NullInt64 `json:"deleted_by,omitempty" gorm:"null" swaggertype:"string"`

	CreatedAt  time.Time    `json:"created_at,omitempty" gorm:"type:TIMESTAMP with time zone;not null"`
	ModifiedAt sql.NullTime `json:"modified_at,omitempty" gorm:"type:TIMESTAMP with time zone;null" swaggertype:"string"`
	DeletedAt  sql.NullTime `json:"deleted_at,omitempty" gorm:"type:TIMESTAMP with time zone;null" swaggertype:"string"`
}