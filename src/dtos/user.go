package dtos

type UserCreate struct {
	Username     string `json:"username" binding:"username,required"`
	Password     string `json:"password" binding:"password,required"`
	MobileNumber string `json:"mobileNumber" binding:"mobile,required"`
}

type UserUpdate struct {
	Username     string `json:"username" binding:"omitempty,username"`
	Firstname    string `json:"firstname" binding:"omitempty,alpha,min=2,max=25"`
	Lastname     string `json:"lastname" binding:"omitempty,alpha,min=3,max=35"`
	MobileNumber string `json:"mobileNumber" binding:"omitempty,mobile"`
	Password     string `json:"password" binding:"omitempty,password"`
	Enabled      bool   `json:"enabled" binding:"omitempty,"`
}

type UserResponse struct {
	Id           int    `json:"id"`
	Username     string `json:"username"`
	Firstname    string `json:"firstname,omitempty"`
	Lastname     string `json:"lastname,omitempty"`
	MobileNumber string `json:"mobileNumber"`
	Password     string `json:"password"`
	Enabled      bool   `json:"enabled"`
}
