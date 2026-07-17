package User

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Phone    int    `json:"phone"`
	Email    string `json:"email" gorm:"index"`
	Password string `json:"password"`
	Name     string `json:"name"`
}
