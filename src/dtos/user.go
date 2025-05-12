package dtos

type UserCreate struct {
	Username     string `json:"username" binding:"username,required"`
	Password     string `json:"password" binding:"password,required"`
	MobileNumber string `json:"mobile_number" binding:"mobile,required"`
}

type UserUpdate struct {
	Username     string `json:"username,omitempty" binding:"omitempty,username"`
	Firstname    string `json:"firstname,omitempty" binding:"omitempty,alpha,min=2,max=25"`
	Lastname     string `json:"lastname,omitempty" binding:"omitempty,alpha,min=3,max=35"`
	MobileNumber string `json:"mobile_number,omitempty" binding:"omitempty,mobile"`
	Password     string `json:"password,omitempty" binding:"omitempty,password"`
	Enabled      bool   `json:"enabled" binding:"omitempty"`
}

type UserResponse struct {
	Id           int    `json:"id"`
	Username     string `json:"username"`
	Firstname    string `json:"firstname,omitempty"`
	Lastname     string `json:"lastname,omitempty"`
	MobileNumber string `json:"mobile_number"`
}
