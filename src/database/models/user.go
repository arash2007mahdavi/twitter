package models

type User struct {
	BaseModel
	Username     string `gorm:"type:string;size:30;not null;unique"`
	Firstname    string `gorm:"type:string;size:20;null"`
	Lastname     string `gorm:"type:string;size:40;null"`
	MobileNumber string `json:"mobile_number" gorm:"type:string;size:11;not null;unique"`
	Password     string `gorm:"type:string;size:50;not null"`
	Enabled      bool   `gorm:"type:bool;default:true"`
}
